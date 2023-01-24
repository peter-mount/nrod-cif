package cifretrieve

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2/log"
	"io"
	"net/http"
	"os"
	"path/filepath"
)

const (
	DOW = "montuewedthufrisatsun"
	URL = "https://publicdatafeeds.networkrail.co.uk/ntrod/CifFileAuthenticate?type=%s&day=%s"
	// URL_OLD is the old url in use before the end of Jan 2023
	URL_OLD = "https://datafeeds.networkrail.co.uk/ntrod/CifFileAuthenticate?type=%s&day=%s"
)

func (a *CIFRetriever) retrieveFull() error {
	return a.retrieve("CIF_ALL_FULL_DAILY", "toc-full.CIF.gz", "full")
}

func (a *CIFRetriever) retrieveUpdate(dow int) error {
	s := dow * 3
	return a.retrieve(
		"CIF_ALL_UPDATE_DAILY",
		fmt.Sprintf("toc-update-%s.CIF.gz", DOW[s:(s+3)]),
		"update")
}

func (a *CIFRetriever) retrieve(filetype string, day string, fname string) error {
	log.Printf("Retrieving %s", day)

	url := fmt.Sprintf(URL, filetype, day)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}

	req.SetBasicAuth(*a.username, *a.password)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("Request returned %d: %s", resp.StatusCode, resp.Status)
	}
	log.Printf("Retrieved %d bytes", resp.ContentLength)

	// Copy the body to a temporary file which is deleted when retrieve exits
	tempfile, err := os.CreateTemp("", "cif")
	if err != nil {
		return err
	}
	defer tempfile.Close()
	defer os.Remove(tempfile.Name())

	_, err = io.Copy(tempfile, resp.Body)
	if err != nil {
		return err
	}

	// Extract the CIF HD record to get extract timestamp, file type & file path
	cif, err := a.extractCifHeader(tempfile)
	if err != nil {
		return err
	}

	// Log it's not a CIF file but don't fail, just ignore it
	if cif == nil {
		log.Printf("%s is not a CIF file", day)
		return nil
	}

	// Add the cif to the list of new files
	a.files = append(a.files, cif)

	// Copy the file to the output directory, creating the directory as required

	err = os.MkdirAll(filepath.Dir(cif.path), 0755)
	if err != nil {
		return err
	}

	// Create but do not overwrite
	wfile, err := os.OpenFile(cif.path, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0644)
	if err != nil {
		// if we exist already then do nothing
		if os.IsExist(err) {
			log.Printf("%s already exists", day)
			return nil
		}
		return err
	}
	defer wfile.Close()

	log.Printf("Writing %s", cif.path)

	_, err = tempfile.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}

	_, err = io.Copy(wfile, tempfile)
	if err != nil {
		return err
	}

	// Set file time to that of the extract date
	return os.Chtimes(cif.path, cif.date, cif.date)
}
