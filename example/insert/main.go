package main

import (
	"github.com/haysons/norm"
	"log"
	"time"
)

// norm.DB is concurrency-safe and internally uses the connection pool provided by the nebula graph official SDK.
// Therefore, in general, only a single instance needs to be defined.
var db *norm.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conf := &norm.Config{
		Username:    "root",
		Password:    "nebula",
		SpaceName:   "demo_basketballplayer",
		Addresses:   []string{"127.0.0.1:9669"},
		ConnTimeout: 10 * time.Second,
	}
	var err error
	db, err = norm.Open(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	insertVertex()
	insertEdge()
}
