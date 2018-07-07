package client

import (
  "github.com/peter-mount/nrod-cif/cif"
)

func (c *CIFClient) GetSchedule( uid, date, stp string ) ( *cif.Response, error ) {
  res := &cif.Response{}

  if found, err := c.get( "/schedule/" + uid + "/" + date + "/" + stp, &res ); err != nil {
    return nil, err
  } else if found {
    return res, nil
  } else {
    return nil, nil
  }
}
