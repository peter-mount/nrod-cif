package cifretrieve

import (
  "flag"
  "fmt"
  "github.com/peter-mount/golib/kernel"
  "os"
)

type CIFRetriever struct {
  // Base directory to retrieve into
  basedir *string
  // NROD Credentials
  username *string
  password *string
  // Retrieve full
  getFull *bool
  // Retrieve updates
  getUpdate *bool
  // File to write file update list to
  output *string
  // List of retrieved files
  files []*cif
}

func (a *CIFRetriever) Name() string {
  return "CIFRetriever"
}

func (a *CIFRetriever) Init(k *kernel.Kernel) error {
  a.basedir = flag.String("d", "", "Directory to retrieve into")
  a.username = flag.String("user", "", "username at NROD or use NRODUSER environment variable")
  a.password = flag.String("password", "", "password at NROD or use NRODPASS environment variable")
  a.getFull = flag.Bool("full", true, "Retrieve weekly full CIF file")
  a.getUpdate = flag.Bool("update", true, "Retrieve daily update CIF files")
  a.output = flag.String("o", "", "Write list of updates to file")

  return nil
}

// If *param is "" then use environment variable value
func defaultParam(param *string, envName string) {
  if *param == "" {
    *param = os.Getenv(envName)
  }
}

func (a *CIFRetriever) PostInit() error {
  if *a.basedir == "" {
    return fmt.Errorf("-d directory is required")
  }

  defaultParam(a.username, "NRODUSER")
  defaultParam(a.password, "NRODPASS")
  if *a.username == "" || *a.password == "" {
    return fmt.Errorf("NROD Username/Password required")
  }

  return nil
}

func (a *CIFRetriever) Start() error {

  if *a.getFull {
    if err := a.retrieveFull(); err != nil {
      return err
    }
  }

  if *a.getUpdate {
    for dow := 0; dow < 7; dow++ {
      if err := a.retrieveUpdate(dow); err != nil {
        return err
      }
    }
  }

  if *a.output != "" {
    if err := a.writeoutput(); err != nil {
      return err
    }
  }

  return nil
}
