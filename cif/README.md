# cif
--
    import "github.com/peter-mount/nrod-cif/cif"

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

#### func  PublicTimeWrite

```go
func PublicTimeWrite(c *codec.BinaryCodec, t *PublicTime)
```
PublicTimeWrite is a workaround for writing null times. If the pointer is null
then a time is written where IsZero()==true

#### func  WorkingTimeWrite

```go
func WorkingTimeWrite(c *codec.BinaryCodec, t *WorkingTime)
```
WorkingTimeWrite is a workaround for writing null times. If the pointer is null
then a time is written where IsZero()==true

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
func (c *CIF) CRSHandler(r *rest.Rest) error
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

#### func (*CIF) GetSchedule

```go
func (c *CIF) GetSchedule(tx *bolt.Tx, uid string, date time.Time, stp string) *Schedule
```
GetSchedule returns a specific Schedule's for a specific TrainUID, startDate and
STPIndicator If no schedule exists for the required key then nil is returned

#### func (*CIF) GetSchedulesByUID

```go
func (c *CIF) GetSchedulesByUID(tx *bolt.Tx, uid string) []*Schedule
```
GetSchedulesByUID returns all Schedule's for a specific TrainUID. If no
schedules exist for the required TrainUID then the returned slice is empty.

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
func (c *CIF) ImportHandler(r *rest.Rest) error
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

#### func (*CIF) ResolveScheduleTiplocs

```go
func (c *CIF) ResolveScheduleTiplocs(tx *bolt.Tx, s *Schedule, r *Response)
```
ResolveScheduleTiplocs resolves the Tiploc's in a Schedule

#### func (*CIF) ScheduleHandler

```go
func (c *CIF) ScheduleHandler(r *rest.Rest) error
```
ScheduleHandler implements a net/http handler that implements a simple Rest
service to retrieve all schedules for a specific uid, date and STPIndicator The
handler must have {uid} set in the path for this to work.

For example:

router.HandleFunc( "/schedule/{uid}/{date}/{stp}", db.ScheduleHandler ).Methods(
"GET" )

where db is a pointer to an active CIF struct.

#### func (*CIF) ScheduleUIDHandler

```go
func (c *CIF) ScheduleUIDHandler(r *rest.Rest) error
```
ScheduleUIDHandler implements a net/http handler that implements a simple Rest
service to retrieve all schedules for a specific uid The handler must have {uid}
set in the path for this to work.

For example:

router.HandleFunc( "/schedule/{uid}", db.ScheduleUIDHandler ).Methods( "GET" )

where db is a pointer to an active CIF struct.

#### func (*CIF) StanoxHandler

```go
func (c *CIF) StanoxHandler(r *rest.Rest) error
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
func (c *CIF) TiplocHandler(r *rest.Rest) error
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
	Id string `json:"-" xml:"-"`
	// Location including Suffix (for circular routes)
	// This is guaranteed to be unique per schedule, although for most purposes
	// like display you would use Tiploc
	Location string `json:"-" xml:"-"`
	// Tiploc of this location. For some schedules like circular routes this can
	// appear more than once in a schedule.
	Tiploc string `json:"tpl" xml:"tpl,attr"`
	// Public Timetable
	Pta *PublicTime `json:"pta,omitempty" xml:"pta,attr,omitempty"`
	Ptd *PublicTime `json:"ptd,omitempty" xml:"ptd,attr,omitempty"`
	// Working Timetable
	Wta *WorkingTime `json:"wta,omitempty" xml:"wta,attr,omitempty"`
	Wtd *WorkingTime `json:"wtd,omitempty" xml:"wtd,attr,omitempty"`
	Wtp *WorkingTime `json:"wtp,omitempty" xml:"wtp,attr,omitempty"`
	// Platform
	Platform string `json:"plat,omitempty" xml:"plat,attr,omitempty"`
	// Activity up to 6 codes
	Activity []string `json:"activity,omitempty" xml:"activity,omitempty"`
	// The Line the train will take
	Line string `json:"line,omitempty" xml:"line,attr,omitempty"`
	// The Path the train will take
	Path string `json:"path,omitempty" xml:"path,attr,omitempty"`
	// Allowances at this location
	EngAllow  string `json:"engAllow,omitempty" xml:"engAllow,attr,omitempty"`
	PathAllow string `json:"pathAllow,omitempty" xml:"pathAllow,attr,omitempty"`
	PerfAllow string `json:"perfAllow,omitempty" xml:"perfAllow,attr,omitempty"`
}
```

A representation of a location within a schedule. There are three types of
location, defined by the Id field:

"LO" Origin, always the first location in a schedule

"LI" Intermediate: A stop or pass along the route

"LT" Destination: always the last lcoation in a schedule

For most purposes you would be interested in the Tiploc, Pta, Ptd and Platform
fields.

Tiploc is the name of this location.

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
BinaryCodec reader

#### func (*Location) Write

```go
func (l *Location) Write(c *codec.BinaryCodec)
```
BinaryCodec writer

#### type PublicTime

```go
type PublicTime struct {
}
```

Public Timetable time Note: 00:00 is not possible as in CIF that means no-time

#### func  PublicTimeRead

```go
func PublicTimeRead(c *codec.BinaryCodec) *PublicTime
```
PublicTimeRead is a workaround issue where a custom type cannot be omitempty in
JSON unless it's a nil So instead of using BinaryCodec.Read( v ), we call this &
set the return value in the struct as a pointer.

#### func (*PublicTime) Get

```go
func (t *PublicTime) Get() int
```
Get returns the PublicTime in seconds of the day

#### func (*PublicTime) IsZero

```go
func (t *PublicTime) IsZero() bool
```
IsZero returns true if the time is not present

#### func (*PublicTime) MarshalJSON

```go
func (t *PublicTime) MarshalJSON() ([]byte, error)
```
Custom JSON Marshaler. This will write null or the time as "HH:MM"

#### func (*PublicTime) MarshalXMLAttr

```go
func (t *PublicTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error)
```
Custom XML Marshaler.

#### func (*PublicTime) Read

```go
func (t *PublicTime) Read(c *codec.BinaryCodec)
```
BinaryCodec reader

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

#### func (*PublicTime) Write

```go
func (t *PublicTime) Write(c *codec.BinaryCodec)
```
BinaryCodec writer

#### type Response

```go
type Response struct {
	XMLName   xml.Name    `json:"-" xml:"response"`
	Status    int         `json:"status,omitempty" xml:"status,attr,omitempty"`
	Message   string      `json:"message,omitempty" xml:"message,attr,omitempty"`
	Schedules []*Schedule `json:"schedules,omitempty" xml:"schedules>schedule,omitempty"`
	//Tiploc     []*Tiploc              `json:"tiploc,omitempty" xml:"tiplocs>tiploc,omitempty"`
	Tiploc *TiplocMap `json:"tiploc,omitempty" xml:"tiplocs>tiploc,omitempty"`
	Self   string     `json:"self" xml:"self,attr,omitempty"`
}
```

Common struct used in forming all responses from rest endpoints. This makes the
responses similar in nature and reduces the amount of redundant code

#### func  NewResponse

```go
func NewResponse() *Response
```

#### func (*Response) AddTiploc

```go
func (r *Response) AddTiploc(t *Tiploc)
```
AddTiploc adds a Tiploc to the response

#### func (*Response) AddTiplocs

```go
func (r *Response) AddTiplocs(t []*Tiploc)
```
AddTiplocs adds an array of Tiploc's to the response

#### func (*Response) GetTiploc

```go
func (r *Response) GetTiploc(n string) (*Tiploc, bool)
```

#### func (*Response) TiplocsSetSelf

```go
func (r *Response) TiplocsSetSelf(rs *rest.Rest)
```
SetSelf sets the Self field to match this request

#### type Schedule

```go
type Schedule struct {
	XMLName xml.Name `json:"-" xml:"schedule"`
	// The train UID
	TrainUID string `json:"uid" xml:"uid,attr"`
	// The date range the schedule is valid on
	RunsFrom time.Time `json:"runsFrom" xml:"from,attr"`
	RunsTo   time.Time `json:"runsTo" xml:"to,attr"`
	// The day's of the week the service will run
	DaysRun    string `json:"daysRun" xml:"daysRun,attr"`
	BankHolRun string `json:"bankHolRun,omitempty" xml:"bankHolRun,attr,omitempty"`
	Status     string `json:"status" xml:"status,attr"`
	Category   string `json:"category" xml:"category,attr"`
	// The identity sometimes confusingly called the Headcode of the service.
	// This is the value you would see in the nrod-td feed
	TrainIdentity string `json:"trainIdentity,omitempty" xml:"trainIdentity,attr,omitempty"`
	// The headcode of this service. Don't confuse with TrainIdentity above
	Headcode                 int    `json:"headcode,omitempty" xml:"headcode,attr,omitempty"`
	ServiceCode              int    `json:"serviceCode,omitempty" xml:"serviceCode,attr,omitempty"`
	PortionId                string `json:"portionId,omitempty" xml:"portionId,attr,omitempty"`
	PowerType                string `json:"powerType,omitempty" xml:"powerType,attr,omitempty"`
	TimingLoad               string `json:"timingLoad,omitempty" xml:"timingLoad,attr,omitempty"`
	Speed                    int    `json:"speed,omitempty" xml:"speed,attr,omitempty"`
	OperatingCharacteristics string `json:",omitempty" xml:",omitempty"`
	SeatingClass             string `json:"seatingClass,omitempty" xml:"seatingClass,attr,omitempty"`
	Sleepers                 string `json:"sleepers,omitempty" xml:"sleepers,attr,omitempty"`
	Reservations             string `json:"reservations,omitempty" xml:"reservations,attr,omitempty"`
	CateringCode             string `json:"cateringCode,omitempty" xml:"cateringCode,attr,omitempty"`
	ServiceBranding          string `json:"branding,omitempty" xml:"branding,attr,omitempty"`
	// The STP Indicator
	STPIndicator string `json:"stp" xml:"stp,attr"`
	UICCode      int    `json:"uic,omitempty" xml:"uic,attr,omitempty"`
	// The operator of this service
	ATOCCode            string `json:"operator,omitempty" xml:"operator,attr,omitempty"`
	ApplicableTimetable bool   `json:"applicableTimetable" xml:"applicableTimetable,attr"`
	// LO, LI & LT entries
	Locations []*Location `json:"locations" xml:"location"`
	// The CIF extract this entry is from
	DateOfExtract time.Time `json:"dateOfExtract" xml:"dateOfExtract,attr"`
	// URL for this Schedule
	Self string `json:"self,omitempty" xml:"self,attr,omitempty"`
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
BinaryCodec reader

#### func (*Schedule) SetSelf

```go
func (s *Schedule) SetSelf(r *rest.Rest)
```
SetSelf sets the Schedule's Self field according to the inbound request. The
resulting URL should then refer back to the rest endpoint that would return this
Schedule.

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
BinaryCodec writer

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
	//XMLName         xml.Name  `xml:"tiploc"`
	// Tiploc key for this location
	Tiploc   string `json:"tiploc" xml:"tiploc,attr"`
	NLC      int    `json:"nlc" xml:"nlc,attr"`
	NLCCheck string `json:"nlcCheck" xml:"nlcCheck,attr"`
	// Proper description for this location
	Desc string `json:"desc" xml:"desc,attr,omitempty"`
	// Stannox code, 0 means none
	Stanox int `json:"stanox" xml:"stanox,attr,omitempty"`
	// CRS code, "" for none. Codes starting with X or Z are usually not stations.
	CRS string `json:"crs" xml:"crs,attr,omitempty"`
	// NLC description of the location
	NLCDesc string `json:"nlcDesc" xml:"nlcDesc,attr,omitempty"`
	// The CIF extract this entry is from
	DateOfExtract time.Time `json:"dateOfExtract" xml:"dateOfExtract,attr"`
	// Self (generated on rest only)
	Self string `json:"self,omitempty" xml:"self,attr,omitempty"`
}
```

Tiploc represents a location on the rail network. This can be either a station,
a junction or a specific point along the line/

#### func (*Tiploc) Read

```go
func (t *Tiploc) Read(c *codec.BinaryCodec)
```

#### func (*Tiploc) SetSelf

```go
func (t *Tiploc) SetSelf(r *rest.Rest)
```
SetSelf sets the Self field to match this request

#### func (*Tiploc) String

```go
func (t *Tiploc) String() string
```
String returns a human readable version of a Tiploc

#### func (*Tiploc) Write

```go
func (t *Tiploc) Write(c *codec.BinaryCodec)
```

#### type TiplocMap

```go
type TiplocMap struct {
}
```


#### func (*TiplocMap) MarshalJSON

```go
func (t *TiplocMap) MarshalJSON() ([]byte, error)
```

#### type WorkingTime

```go
type WorkingTime struct {
}
```

Working Timetable time. WorkingTime is similar to PublciTime, except we can have
seconds. In the Working Timetable, the seconds can be either 0 or 30.

#### func  WorkingTimeRead

```go
func WorkingTimeRead(c *codec.BinaryCodec) *WorkingTime
```
WorkingTimeRead is a workaround issue where a custom type cannot be omitempty in
JSON unless it's a nil So instead of using BinaryCodec.Read( v ), we call this &
set the return value in the struct as a pointer.

#### func (*WorkingTime) Get

```go
func (t *WorkingTime) Get() int
```
Get returns the WorkingTime in seconds of the day

#### func (*WorkingTime) IsZero

```go
func (t *WorkingTime) IsZero() bool
```
IsZero returns true if the time is not present

#### func (*WorkingTime) MarshalJSON

```go
func (t *WorkingTime) MarshalJSON() ([]byte, error)
```
Custom JSON Marshaler. This will write null or the time as "HH:MM:SS"

#### func (*WorkingTime) MarshalXMLAttr

```go
func (t *WorkingTime) MarshalXMLAttr(name xml.Name) (xml.Attr, error)
```
Custom XML Marshaler.

#### func (*WorkingTime) Read

```go
func (t *WorkingTime) Read(c *codec.BinaryCodec)
```
BinaryCodec reader

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

#### func (*WorkingTime) Write

```go
func (t *WorkingTime) Write(c *codec.BinaryCodec)
```
BinaryCodec writer
