package main

import (
	"github.com/haysons/norm"
	"log"
	"time"
)

type Woman struct {
	VID     string `norm:"vertex_id"`
	Name    string
	Age     int
	Married bool
	Salary  float64
}

func (t *Woman) VertexID() string {
	return t.VID
}

func (t *Woman) VertexTagName() string {
	return "woman"
}

type WomanUpdate struct {
	VID     string `norm:"vertex_id"`
	Name    string `norm:"prop:name;not_null;default:''"`
	Age     int    `norm:"prop:age;not_null;default:0"`
	Married bool   `norm:"prop:married;not_null;default:false"`
}

func (t *WomanUpdate) VertexID() string {
	return t.VID
}

func (t *WomanUpdate) VertexTagName() string {
	return "woman"
}

var db *norm.DB

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	conf := &norm.Config{
		Username:    "root",
		Password:    "nebula",
		SpaceName:   "test",
		Addresses:   []string{"127.0.0.1:9669"},
		ConnTimeout: 10 * time.Second,
	}
	var err error
	db, err = norm.Open(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migrator := db.Debug().Migrator()
	hasTag, err := migrator.HasVertexTag("woman")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("has vertex tag: %v", hasTag)

	if err = migrator.AutoMigrateVertexes(Woman{}); err != nil {
		log.Fatal(err)
	}
	womanProps, err := migrator.DescVertexTag("woman")
	if err != nil {
		log.Fatal(err)
	}
	for _, womanProp := range womanProps {
		log.Printf("woman prop: %+v\n", womanProp)
	}

	migrator = norm.NewMigrator(db.Debug())
	if err = migrator.AutoMigrateVertexes(WomanUpdate{}); err != nil {
		log.Fatal(err)
	}
	womanProps, err = migrator.DescVertexTag("woman")
	if err != nil {
		log.Fatal(err)
	}
	for _, womanProp := range womanProps {
		log.Printf("woman prop: %+v\n", womanProp)
	}

	if err = migrator.DropVertexTag("woman"); err != nil {
		log.Fatal(err)
	}
}
