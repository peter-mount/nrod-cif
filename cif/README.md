# cif
--
    import "cif"

A GO library providing a database based on the Network Rail CIF Timetable feed.

## Usage

```go
const (
	// Import tiplocs only
	TIPLOC = 1
	// Import schedules only
	SCHEDULE = 1 << 1
	// The default mode used if nothing is set
	ALL = TIPLOC | SCHEDULE
)
```
Bitmasks for CIF.Mode used by CIF.Parse() & CIF.ParseFile() to determine what to
import. If not set then everything is imported

```go
const (
	DateTime      = "2006-01-02 15:04:05"
	Date          = "2006-01-02"
	HumanDateTime = "2006 Jan 02 15:04:05"
	HumanDate     = "2006 Jan 02"
	Time          = "15:04:05"
)
```

#### type CIF

```go
type CIF struct {
	// The mode the parser should use when importing NR CIF files.
	// This is a bit mask of TIPLOC or SCHEDULE. If not set then ALL is used.
	Mode int
}
```


#### func (*CIF) CRSHandler

```go
func (c *CIF) CRSHandler(w http.ResponseWriter, r *http.Request)
```
CRSHandler implements a net/http handler that implements a simple Rest service
to retrieve CRS/3Alpha records. The handler must have {id} set in the path for
this to work, where id would represent the CRS code.

For example:

router.HandleFunc( "/crs/{id}", db.CRSHandler ).Methods( "GET" )

where db is a pointer to an active CIF struct. When running this would allow GET
requests like /crs/MDE to return JSON representing that station.

#### func (*CIF) Close

```go
func (c *CIF) Close()
```
Close the database. If OpenDB() was used to open the db then that db is closed.
If UseDB() was used this simply detaches the CIF from that DB. The DB is not
closed()

#### func (*CIF) GetCRS

```go
func (c *CIF) GetCRS(tx *bolt.Tx, crs string) ([]*Tiploc, bool)
```
GetCRS retrieves an array of Tiploc records for the CRS/3Alpha code of a
station.

#### func (*CIF) GetHD

```go
func (c *CIF) GetHD() (*HD, error)
```
GetHD retrieves the latest HD record of the latest cif file imported into the
database.

#### func (*CIF) GetStanox

```go
func (c *CIF) GetStanox(tx *bolt.Tx, stanox int) ([]*Tiploc, bool)
```

#### func (*CIF) GetTiploc

```go
func (c *CIF) GetTiploc(tx *bolt.Tx, t string) (*Tiploc, bool)
```
GetTiploc retrieves a Tiploc from the cif database

tx An active readonly bolt.Tx

t The Tiploc to retrieve, 1..7 characters long, always upper case

Returns ( tiploc *Tiploc, exist bool )

If exist is true then tiploc will be a new Tiploc instance with the retrieved
data. If false then the tiploc is not in the database.

#### func (*CIF) ImportCIF

```go
func (c *CIF) ImportCIF(r io.Reader) error
```
ImportCIF imports a uncompressed CIF file retrieved from NetworkRail into the
cif database. If this file is a full export then the database will be cleared
first.

The CIF.Mode field determines how this import is performed. This field is a
bitmask so one or more options can be included. They are:


TIPLOC Import tiplocs

SCHEDULE Import schedules

ALL Import everything, the default and the same as TIPLOC | SCHEDULE

#### func (*CIF) ImportFile

```go
func (c *CIF) ImportFile(fname string) error
```
ImportFile imports a uncompressed CIF file retrieved from NetworkRail into the
cif database. If this file is a full export then the database will be cleared
first.

The CIF.Mode field determines how this import is performed. This field is a
bitmask so one or more options can be included. They are:


TIPLOC Import tiplocs

SCHEDULE Import schedules

ALL Import everything, the default and the same as TIPLOC | SCHEDULE

#### func (*CIF) ImportHandler

```go
func (c *CIF) ImportHandler(rw http.ResponseWriter, req *http.Request)
```
ImportHandler implements a net/http handler that implements a Rest service to
import an uncompressed CIF file from NetworkRail and import it into the cif
database.

For example:

router.HandleFunc( "/importCIF", db.ImportHandler ).Methods( "POST" )

Will define the path /importCIF to accept HTTP POST requests. You can then
submit a cif file to this endpoint to import a CIF file.

Example: To perform a full import, replacing all cif data in the database

curl -X POST --data-binary @toc-full.CIF http://localhost:8081/importCIF

To perform an update then simply submit an update cif file:

curl -X POST --data-binary @toc-update-sun.CIF http://localhost:8081/importCIF

BUG(peter-mount): The Rest service provided by ImportHandler is currently
unprotected so anyone can perform an import. We need to provide some means of
simple authentication to this handler.

#### func (*CIF) OpenDB

```go
func (c *CIF) OpenDB(dbFile string) error
```
OpenDB opens a CIF database.

#### func (*CIF) StanoxHandler

```go
func (c *CIF) StanoxHandler(w http.ResponseWriter, r *http.Request)
```
StanoxHandler implements a net/http handler that implements a simple Rest
service to retrieve stanox records. The handler must have {id} set in the path
for this to work, where id would represent the CRS code.

For example:

router.HandleFunc( "/stanox/{id}", db.StanoxHandler ).Methods( "GET" )

where db is a pointer to an active CIF struct. When running this would allow GET
requests like /stanox/89403 to return JSON representing that station.

#### func (*CIF) String

```go
func (c *CIF) String() string
```
String returns a human readable description of the latest CIF file imported into
this database.

#### func (*CIF) TiplocHandler

```go
func (c *CIF) TiplocHandler(w http.ResponseWriter, r *http.Request)
```
TiplocHandler implements a net/http handler that implements a simple Rest
service to retrieve Tiploc records. The handler must have {id} set in the path
for this to work, where id would represent the Tiploc code.

For example:

router.HandleFunc( "/tiploc/{id}", db.TiplocHandler ).Methods( "GET" )

where db is a pointer to an active CIF struct. When running this would allow GET
requests like /tiploc/MSTONEE to return JSON representing that station.

#### func (*CIF) UseDB

```go
func (c *CIF) UseDB(db *bolt.DB) error
```
UseDB Allows an already open database to be used with cif.

#### type HD

```go
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
```


#### func (*HD) Read

```go
func (h *HD) Read(c *codec.BinaryCodec)
```

#### func (*HD) String

```go
func (h *HD) String() string
```
String returns a human readable version of the HD record.

#### func (*HD) Write

```go
func (h *HD) Write(c *codec.BinaryCodec)
```

#### type Location

```go
type Location struct {
	// Type of location:
	Id string
	// Location including Suffix (for circular routes)
	// This is guaranteed to be unique per schedule, although for most purposes
	// like display you would use Tiploc
	Location string
	// Tiploc of this location. For some schedules like circular routes this can
	// appear more than once in a schedule.
	Tiploc string
	// Public Timetable
	Pta PublicTime
	Ptd PublicTime
	// Working Timetable
	Wta WorkingTime
	Wtd WorkingTime
	Wtp WorkingTime
	// Platform
	Platform string
	// Activity up to 6 codes
	Activity []string
	// The Line the train will take
	Line string
	// The Path the train will take
	Path string
	// Allowances at this location
	EngAllow  string
	PathAllow string
	PerfAllow string
}
```

A representation of a location within a schedule. There are three types of
location, defined by the Id field:

"LO" Origin, always the first location in a schedule

"LI" Intermediate: A stop or pass along the route

"LT" Destination: always the last lcoation in a schedule

For most purposes you would be interested in the Tiploc, Pta, Ptd and Platform
fields. Tiploc is the name of this location.

Pta & Ptd are the public timetable times, i.e. what is published to the general
public.

Pta is the arrival time and is valid for LI & LT entries only.

Ptd is the departue time and is valid for LO & LI entries only.

If either are not set then the train is not scheduled to stop at this location.

Wta, Wtd & Wtp are the working timetable, i.e. the actual timetable the service
runs to. Wta & Wtd are like Pta & Ptd but Wtp means the time the train is
scheduled to pass a location. If Wtp is set then Pta, Ptd, Wta & Wtp will not be
set.

#### func (*Location) Read

```go
func (l *Location) Read(c *codec.BinaryCodec)
```

#### func (*Location) Write

```go
func (l *Location) Write(c *codec.BinaryCodec)
```

#### type PublicTime

```go
type PublicTime struct {
	T int
}
```

Public Timetable time Note: 00:00 is not possible as in CIF that means no-time

#### func (*PublicTime) Get

```go
func (t *PublicTime) Get() int
```
Get returns the PublicTime in seconds of the day

#### func (*PublicTime) IsSet

```go
func (t *PublicTime) IsSet() bool
```
IsSet returns true if the time is set

#### func (PublicTime) Read

```go
func (t PublicTime) Read(c *codec.BinaryCodec)
```

#### func (*PublicTime) Set

```go
func (t *PublicTime) Set(v int)
```
Set sets the PublicTime in seconds of the day

#### func (*PublicTime) String

```go
func (t *PublicTime) String() string
```
String returns a PublicTime in HH:MM format or 5 blank spaces if it's not set.

#### func (PublicTime) Write

```go
func (t PublicTime) Write(c *codec.BinaryCodec)
```

#### type Schedule

```go
type Schedule struct {
	// The train UID
	TrainUID string
	// The date range the schedule is valid on
	RunsFrom time.Time
	RunsTo   time.Time
	// The day's of the week the service will run
	DaysRun    string
	BankHolRun string
	Status     string
	Category   string
	// The identity sometimes confusingly called the Headcode of the service.
	// This is the value you would see in the nrod-td feed
	TrainIdentity string
	// The headcode of this service. Don't confuse with TrainIdentity above
	Headcode                 int
	ServiceCode              int
	PortionId                string
	PowerType                string
	TimingLoad               string
	Speed                    int
	OperatingCharacteristics string
	SeatingClass             string
	Sleepers                 string
	Reservations             string
	CateringCode             string
	ServiceBranding          string
	// The STP Indicator
	STPIndicator string
	UICCode      int
	// The operator of this service
	ATOCCode            string
	ApplicableTimetable bool
	// LO, LI & LT entries
	Locations []*Location
	// The CIF extract this entry is from
	DateOfExtract time.Time
}
```

A train schedule

#### func (*Schedule) Equals

```go
func (s *Schedule) Equals(o *Schedule) bool
```
Equals returns true if two Schedule struts refer to the same schedule. This
checks the "primary key" for schedules which is TrainUID, RunsFrom &
STPIndicator

#### func (*Schedule) FullString

```go
func (s *Schedule) FullString() string
```
FullString is a debug function that returns a Schedule as a string in a human
readable format. Unlike String() this will contain everything about the
schedule.

#### func (*Schedule) Read

```go
func (s *Schedule) Read(c *codec.BinaryCodec)
```

#### func (*Schedule) String

```go
func (s *Schedule) String() string
```
String returns the "primary key" for schedules which is TrainUID, RunsFrom &
STPIndicator

#### func (*Schedule) Write

```go
func (s *Schedule) Write(c *codec.BinaryCodec)
```

#### type SimpleResponse

```go
type SimpleResponse struct {
	Status  int
	Message string
}
```


#### type Tiploc

```go
type Tiploc struct {
	// Tiploc key for this location
	Tiploc   string
	NLC      int
	NLCCheck string
	// Proper description for this location
	Desc string
	// Stannox code, 0 means none
	Stanox int
	// CRS code, "" for none. Codes starting with X or Z are usually not stations.
	CRS string
	// NLC description of the location
	NLCDesc string
	// The CIF extract this entry is from
	DateOfExtract time.Time
}
```

Tiploc represents a location on the rail network. This can be either a station,
a junction or a specific point along the line/

#### func (*Tiploc) Read

```go
func (t *Tiploc) Read(c *codec.BinaryCodec)
```

#### func (*Tiploc) String

```go
func (t *Tiploc) String() string
```
String returns a human readable version of a Tiploc

#### func (*Tiploc) Write

```go
func (t *Tiploc) Write(c *codec.BinaryCodec)
```

#### type WorkingTime

```go
type WorkingTime struct {
	T int
}
```

Working Timetable time. WorkingTime is similar to PublciTime, except we can have
seconds. In the Working Timetable, the seconds can be either 0 or 30.

#### func (*WorkingTime) Get

```go
func (t *WorkingTime) Get() int
```
Get returns the WorkingTime in seconds of the day

#### func (*WorkingTime) IsSet

```go
func (t *WorkingTime) IsSet() bool
```
IsSet returns true if the time is set

#### func (WorkingTime) Read

```go
func (t WorkingTime) Read(c *codec.BinaryCodec)
```

#### func (*WorkingTime) Set

```go
func (t *WorkingTime) Set(v int)
```
Set sets the WorkingTime in seconds of the day

#### func (*WorkingTime) String

```go
func (t *WorkingTime) String() string
```
String returns a PublicTime in HH:MM:SS format or 8 blank spaces if it's not
set.

#### func (WorkingTime) Write

```go
func (t WorkingTime) Write(c *codec.BinaryCodec)
```
