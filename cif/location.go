package cif

type Location struct {
  // LO,
  Id          string
  // Location including Suffix (for circular routes)
  Location    string
  // Tiploc of Location sans Suffix
  Tiploc      string
  // Public times in seconds of day
  Pta         PublicTime
  Ptd         PublicTime
  // Working times in seconds of day
  Wta         WorkingTime
  Wtd         WorkingTime
  Wtp         WorkingTime
  // Platform
  Platform    string
  // Activity up to 6 codes
  Activity  []string
  // Misc
  Line        string
  Path        string
  EngAllow    string
  PathAllow   string
  PerfAllow   string
}

func newLocation() *Location {
  var loc *Location = &Location{}
  loc.Pta.Set( -1 )
  loc.Ptd.Set( -1 )
  loc.Wta.Set( -1 )
  loc.Wtd.Set( -1 )
  loc.Wtp.Set( -1 )
  return loc
}

func (s *Schedule) appendLocation(l *Location) {
  s.Locations = append( s.Locations, l )
}
