// The rest server

package main

import (
  "fmt"
  "github.com/gorilla/mux"
  "log"
  "net/http"
  "time"
)

type Server struct {
  // Port to listen to
  Port    int
  // ===== Internal
  Router  *mux.Router
}

func (s *Server) Init() {
  // If not defined then use port 80
  if s.Port < 1 || s.Port > 65534 {
    s.Port = 8080
  }

  s.Router = mux.NewRouter()
}

func (s *Server) Start() error {
  log.Printf( "Listening on port %d\n", s.Port )
  return http.ListenAndServe( fmt.Sprintf( ":%d", s.Port ), s.Router )
}

// Sets Date header to a unix timestamp
func (s *Server) setUnixDate( w http.ResponseWriter, t int64 ) {
  var b []byte = appendTime( nil, time.Unix( t, int64(0) ) )
  w.Header().Add( "Date", string( b[:] ) )
}

// Sets response headers for a Json response.
// sc is != 0 then the status code
// t unix timestamp
// cache time in seconds to cache the response, <0 for no cache, 0 to ignore
func (s *Server) setJsonResponse( w http.ResponseWriter, sc int, t int64, cache int ) {
  var h http.Header = w.Header()

  h.Add( "Content-Type", "application/json" )

  if t > 0 {
    s.setUnixDate( w, t )
  }

  if cache < 0 {
    h.Add( "Cache-Control", "no-cache" )
  } else if cache > 0 {
    h.Add( "Cache-Control", fmt.Sprintf( "max-age=%d, s-maxage=%d", cache, cache ) )
  }

  // This must be last
  if sc > 0 {
    w.WriteHeader( sc )
  }
}

// appendTime is a non-allocating version of []byte(t.UTC().Format(TimeFormat))
func appendTime(b []byte, t time.Time) []byte {
	const days = "SunMonTueWedThuFriSat"
	const months = "JanFebMarAprMayJunJulAugSepOctNovDec"

	t = t.UTC()
	yy, mm, dd := t.Date()
	hh, mn, ss := t.Clock()
	day := days[3*t.Weekday():]
	mon := months[3*(mm-1):]

	return append(b,
		day[0], day[1], day[2], ',', ' ',
		byte('0'+dd/10), byte('0'+dd%10), ' ',
		mon[0], mon[1], mon[2], ' ',
		byte('0'+yy/1000), byte('0'+(yy/100)%10), byte('0'+(yy/10)%10), byte('0'+yy%10), ' ',
		byte('0'+hh/10), byte('0'+hh%10), ':',
		byte('0'+mn/10), byte('0'+mn%10), ':',
		byte('0'+ss/10), byte('0'+ss%10), ' ',
		'G', 'M', 'T')
}
