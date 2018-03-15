package bin

import (
  "path/filepath"
)

// DbPath ensures the database name is set. If the name is not absolute then it's
// taken as being relative to the database path in config.
// s The required filename
// d The filename to use if s is ""
func (c *Config) DbPath( s *string, d string ) *Config {
  if *s == "" {
    *s = d
  }

  if (*s)[0] != '/' {
    *s = c.Database.Path + *s
  }

  return c
}

func (c *Config) initDb() error {

  if c.Database.Path == "" {
    c.Database.Path = "/database/"
  }

  if path, err := filepath.Abs( c.Database.Path ); err != nil {
    return err
  } else {
    c.Database.Path = path
  }

  if c.Database.Path[len(c.Database.Path)-1] != '/' {
    c.Database.Path = c.Database.Path + "/"
  }

  return nil
}
