package cifimport

import (
	"fmt"
	"github.com/peter-mount/go-kernel/v2/log"
	"time"
)

type HD struct {
	Id                    string // Record Identity, always "HD"
	FileMainframeIdentity string
	// The date that the most recent cif file imported was extracted from Network Rail
	DateOfExtract        time.Time
	CurrentFileReference string
	LastFileReference    string
	// Was the last import an update or a full import
	Update  bool
	Version string
	// The Start and End dates for schedules in the latest import.
	// You can be assured that there would be no schedules which are not contained
	// either fully or partially inside these dates to be present.
	UserStartDate time.Time
	UserEndDate   time.Time
}

// Parse HD record
// returns true if the file should be imported
func (c *CIFImporter) parseHD(l string) (bool, error) {
	var h *HD = &HD{}

	i := 0
	i = parseString(l, i, 2, &h.Id)
	i = parseString(l, i, 20, &h.FileMainframeIdentity)
	i = parseDDMMYY_HHMM(l, i, &h.DateOfExtract)
	i = parseString(l, i, 7, &h.CurrentFileReference)
	i = parseString(l, i, 7, &h.LastFileReference)

	var update string
	i = parseString(l, i, 1, &update)
	h.Update = update == "U"

	i = parseString(l, i, 1, &h.Version)
	i = parseDDMMYY(l, i, &h.UserStartDate)
	i = parseDDMMYY(l, i, &h.UserEndDate)

	log.Println(h.String())

	rows, err := c.tx.Query(
		"SELECT timetable.beginimport( $1, $2, $3, $4, $5, $6, $7 )",
		h.FileMainframeIdentity,
		h.DateOfExtract,
		h.CurrentFileReference,
		h.LastFileReference,
		h.Update,
		h.UserStartDate,
		h.UserEndDate,
	)
	if err != nil {
		return false, err
	}
	defer rows.Close()

	for rows.Next() {
		var doImport bool
		err := rows.Scan(&doImport)
		if err != nil {
			return false, err
		}

		if !doImport {
			log.Println("Skipping CIF import")
			return false, nil
		}
	}

	if h.Update {
		log.Println("Performing CIF Update")
	} else {
		log.Println("Performing Full import")
	}

	c.importhd = h
	return true, nil
}

// String returns a human readable version of the HD record.
func (h *HD) String() string {
	return fmt.Sprintf(
		"CIF %s Extracted %v Date Range %v - %v Update %v",
		h.FileMainframeIdentity,
		h.DateOfExtract.Format(HumanDateTime),
		h.UserStartDate.Format(HumanDate),
		h.UserEndDate.Format(HumanDate),
		h.Update)
}
