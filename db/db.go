package db

import "github.com/globalsign/mgo"

const (
	host = "localhost:27017"
	dbs  = "go-blog"
)

var globalS *mgo.Session

func init() {
	dialInfo := &mgo.DialInfo{
		Addrs:    []string{host},
		Database: dbs,
	}
	s, err := mgo.DialWithInfo(dialInfo)
	if err != nil {
		panic(err)
	}
	globalS = s
}

func Connect(collection string) (*mgo.Session, *mgo.Collection) {
	s := globalS.Copy()
	c := s.DB("").C(collection)
	return s, c

}
