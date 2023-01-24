package cifimport

import (
	"compress/gzip"
	"context"
	"fmt"
	"github.com/peter-mount/go-kernel/v2/log"
	"io"
	"os"
)

func (c *CIFImporter) importCIFTask(ctx context.Context) error {
	file := ctx.Value("file").(string)

	log.Printf("Parsing %s", file)

	f, err := os.Open(file)
	if err != nil {
		return err
	}
	defer f.Close()

	// gzip or plain
	var header [2]byte
	c1, err := io.ReadFull(f, header[:])
	if err != nil {
		return err
	}
	if c1 < 2 {
		return fmt.Errorf("")
	}
	_, err = f.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	reader := io.Reader(f)
	if header[0] == 0x1f && header[1] == 0x8b {
		reader, err = gzip.NewReader(f)
		if err != nil {
			return err
		}
	}

	skip, err := c.importCIF(reader)
	if err != nil {
		if skip {
			// Non fatal error so log it but don't kill the import
			log.Println(err)
		} else {
			return err
		}
	}

	return nil
}
