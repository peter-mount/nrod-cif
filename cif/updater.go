package cif

import (
  "compress/gzip"
  "log"
  "net/http"
//  "time"
)

// Handles the current status of an update
type Updater struct {
  cif          *CIF
  username      string
  password      string
  client       *http.Client
}

// NewUpdater creates a new Updater instance
func (cif *CIF) SetUpdater( user string, pass string ) {
  cif.Updater = &Updater{
    cif: cif,
    username: user,
    password: pass,
  }

  cif.Updater.client = &http.Client{
    CheckRedirect: cif.Updater.checkRedirect,
  }
}

func (u *Updater) checkRedirect( req *http.Request, via []*http.Request) error {
  log.Println( "Redirect", req.URL )
  // As this is a redirect to AWS S3 don't set the basic auth again
  //req.SetBasicAuth( u.username, u.password )
  return nil
}

func (u *Updater) do(req *http.Request) (*http.Response, error) {
  req.SetBasicAuth( u.username, u.password )
  resp, err := u.client.Do( req )
  return resp, err
}

// Update runs an update against a CIF
func (u *Updater) Update() {
  log.Println( "Begining update" )
  if err := u.update(); err != nil {
    log.Println( "Update failed:", err )
  } else {
    log.Println( "Update completed" );
  }
}

func (u *Updater) update() error {
  // Check for full update
  if err := u.updateDay( -1 ); err != nil {
    return err
  }

  return nil
}

// Days of week. The order is specified by time.Weekday.
// It's used to form the day attribute
const cifDay = "sunmontuewedthufrisat"

// updateDay retrieves the daily CIF and imports it.
// d is -1 for Full Import, 0 Sunday, 1 Monday... etc
func (u *Updater) updateDay( d int ) error {
  var tp string
  var dy string
  if d < 0 {
    tp = "CIF_ALL_FULL_DAILY"
    dy = "toc-full"
  } else {
    tp = "CIF_ALL_UPDATE_DAILY"
    dy = "toc-update-" + cifDay[ d*3 : (d*3)+3 ]
  }

  url := "https://datafeeds.networkrail.co.uk/ntrod/CifFileAuthenticate?type=" + tp + "&day=" + dy + ".CIF.gz"
  log.Println( "Retrieving", url )

  if req, err := http.NewRequest("GET", url, nil); err != nil {
    return err
  } else {
    if resp, err := u.do( req ); err != nil {
      return err
    } else {
      if reader, err := gzip.NewReader( resp.Body ); err != nil {
        return err
      } else {
        return u.cif.ImportCIF( reader )
      }
    }
  }
}
