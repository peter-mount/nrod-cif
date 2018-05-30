package nrod

import (
  "github.com/hashicorp/consul/api"
  "log"
)

// NetworkRail contains the common configuration options
// needed when connecting to Network Rail.
type NetworkRail struct {
  // User details
  User struct {
    Username      string  `yaml:"username"`
    Password      string  `yaml:"password"`
  }                       `yaml:"user"`
  // ActiveMQ details
  ActiveMQ struct {
    Server        string  `yaml:"server"`
    Port          int     `yaml:"port"`
    ClientId      string  `yaml:"clientId"`
  }                       `yaml:"activeMQ"`
}


func readString( kv *api.KV, k string, v *string ) error {
  pair, _, err := kv.Get( k, nil )
  if err != nil {
    log.Printf( "Get %s failed", k )
    return err
  }
  if pair != nil {
    *v = string(pair.Value[:])
  }
  return nil
}

func (n *NetworkRail) ReadConsul( client *api.Client ) error {
  kv := client.KV()

  log.Println( "List values" )
  if kvp, _, err := kv.List( "dev/", nil ); err != nil {
    log.Println( "List failed" )
    return err
  } else {
    for i, p := range kvp {
      log.Printf( "%2d %16s %v", i, p.Key, p.Value )
    }
  }

  log.Println( "Get values" )
  if err := readString( kv, "dev/nrod/feed/user/username", &n.User.Username ); err != nil { return err }
  if err := readString( kv, "dev/nrod/feed/user/password", &n.User.Password ); err != nil { return err }

  return nil
}
