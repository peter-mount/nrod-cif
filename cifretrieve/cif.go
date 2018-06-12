package cifretrieve

import (
  "compress/gzip"
  "fmt"
  "io"
  "log"
  "os"
  "time"
)

// Simple representation of a cif file
type cif struct {
  // Date of extract, used in sorting
  date    time.Time
  // true if a full import, false for an update
  full    bool
  // The path in the file system of the final cif
  path    string
}

// extractCifHeader extracts the CIF header line
func (a *CIFRetriever) extractCifHeader( f *os.File ) (*cif, error) {
  _, err := f.Seek( 0, io.SeekStart )
  if err != nil {
    return nil, err
  }

  gr, err := gzip.NewReader( f )
  if err != nil {
    return nil, err
  }

  var header [80]byte
  _, err = io.ReadFull( gr, header[:] )
  if err != nil {
    return nil, err
  }

  s:= string(header[:])
  if s[0:2] != "HD" {
    return nil, nil
  }

  dt := s[22:32]
  extractDate, err := time.Parse( "200601021504", "20" + dt[4:6] + dt[2:4] + dt[0:2] + dt[6:] )
  if err != nil {
    return nil, err
  }

  full := s[46] == 'F'

  fileType := "update"
  if full {
    fileType  = "full"
  }

  c := &cif{
    date: extractDate,
    full: full,
    path: fmt.Sprintf( "%s/%s-%s.cif.gz", *a.basedir, extractDate.Format( "2006/01/02" ), fileType ),
  }

  log.Printf( "CIF %s %v schedules between 20%s-%s-%s and 20%s-%s-%s",
    fileType,
    c.date,
    s[48:50],
    s[50:52],
    s[52:54],
    s[54:56],
    s[56:58],
    s[58:60] )

  return c, nil
}
