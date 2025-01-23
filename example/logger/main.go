package main

import (
	"github.com/haysons/nebulaorm"
	"github.com/haysons/nebulaorm/logger"
	"log"
	"time"
)

type Player struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name"`
	Age  int    `norm:"prop:age"`
}

func (p Player) VertexID() string {
	return p.VID
}

func (p Player) VertexTagName() string {
	return "player"
}

func main() {
	queryWithLoggerDebug()
	queryWithLoggerWarn()
}

func queryWithLoggerDebug() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conf := &nebulaorm.Config{
		Username:    "root",
		Password:    "nebula",
		SpaceName:   "demo_basketballplayer",
		Addresses:   []string{"127.0.0.1:9669"},
		ConnTimeout: 10 * time.Second,
	}
	db, err := nebulaorm.Open(conf, nebulaorm.WithLogger(logger.Default.LogMode(logger.DebugLevel)))
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	player := &Player{
		VID:  "player1001",
		Name: "Kobe Bryant",
		Age:  33,
	}
	if err := db.InsertVertex(player).Exec(); err != nil {
		log.Fatalf("insert player failed: %v", err)
	}

	player = new(Player)
	err = db.Fetch("player", "player1001").
		Yield("vertex as v").
		FindCol("v", player)
	if err != nil {
		log.Fatalf("fetch player failed: %v", err)
	}
	log.Printf("player: %+v", player)
}

func queryWithLoggerWarn() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conf := &nebulaorm.Config{
		Username:    "root",
		Password:    "nebula",
		SpaceName:   "demo_basketballplayer",
		Addresses:   []string{"127.0.0.1:9669"},
		ConnTimeout: 10 * time.Second,
	}
	db, err := nebulaorm.Open(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	player := &Player{
		VID:  "player1001",
		Name: "Kobe Bryant",
		Age:  33,
	}
	if err := db.Debug().InsertVertex(player).Exec(); err != nil {
		log.Fatalf("insert player failed: %v", err)
	}

	player = new(Player)
	err = db.Fetch("player", "player1001").
		Yield("vertex as v").
		FindCol("v", player)
	if err != nil {
		log.Fatalf("fetch player failed: %v", err)
	}
	log.Printf("player: %+v", player)
}
