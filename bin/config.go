// Internal library used for the binary webservices
package bin

import (
  "github.com/peter-mount/golib/rest"
  "github.com/peter-mount/golib/statistics"
  "gopkg.in/robfig/cron.v2"
  "gopkg.in/yaml.v2"
  "io/ioutil"
  "path/filepath"
)

// Common configuration used to read config.yaml
type Config struct {
  // URL prefixes for lookups to the reference microservices
  Services struct {
    DarwinD3  string        `yaml:"darwind3"`
    Reference string        `yaml:"reference"`
    Timetable string        `yaml:"timetable"`
  }                         `yaml:"services"`

  Database struct {
    Path          string    `yaml:"path"`
    Cif           string    `yaml:"cif"`
    CifTimetable  string    `yaml:"cifTimetable"`
  }                         `yaml:"database"`

  NetworkRail struct {
    User struct {
      Username      string  `yaml:"username"`
      Password      string  `yaml:"password"`
    }                       `yaml:"user"`
    ActiveMQ struct {
      Server        string  `yaml:"server"`
      Port          int     `yaml:"port"`
      ClientId      string  `yaml:"clientId"`
    }                       `yaml:"activeMQ"`
  }                         `yaml:"networkrail"`

  Server struct {
    // Root context path, defaults to ""
    Context       string    `yaml:"context"`
    // The port to run on, defaults to 80
    Port          int       `yaml:"port"`
    // The permitted headers
    Headers     []string
    // The permitted Origins
    Origins     []string
    // The permitted methods
    Methods     []string
    // Web Server
    server       *rest.Server
    // Base Context
    Ctx          *rest.ServerContext
  }                         `yaml:"server"`

  Statistics struct {
    Log           bool      `yaml:"log"`
    Rest          string    `yaml:"rest"`
    Schedule      string    `yaml:"schedule"`
    statistics   *statistics.Statistics
  }                         `yaml:"statistics"`

  // Cron
  Cron         *cron.Cron
}

// ReadFile reads the provided file and imports yaml config
func (c *Config) readFile( configFile string ) error {
  if filename, err := filepath.Abs( configFile ); err != nil {
    return err
  } else if in, err := ioutil.ReadFile( filename ); err != nil {
    return err
  } else if err := yaml.Unmarshal( in, c ); err != nil {
    return err
  }
  return nil
}
