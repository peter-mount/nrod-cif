package cifretrieve

import (
	"fmt"
	"os"
)

type CIFRetriever struct {
	basedir   *string `kernel:"flag,d,Directory to retrieve into"`
	username  *string `kernel:"flag,user,NROD username or use NRODUSER env"`
	password  *string `kernel:"flag,password,NROD password or use NRODPASS env"`
	getFull   *bool   `kernel:"flag,full,Retrieve weekly full CIF file"`
	getUpdate *bool   `kernel:"flag,update,Retrieve daily update CIF files"`
	output    *string `kernel:"flag,o,write list of updates to file"`
	// List of retrieved files
	files []*cif
}

func (a *CIFRetriever) Name() string {
	return "CIFRetriever"
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

	// _USR & _PSW forms used by Jenkins credentials in pipelines
	defaultParam(a.username, "NROD_USR")
	defaultParam(a.password, "NROD_PSW")

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
