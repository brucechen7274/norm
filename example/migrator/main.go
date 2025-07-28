package main

import (
	"github.com/haysons/norm"
	"log"
	"time"
)

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

	migrateTags()

	migrateEdges()
}

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

func migrateTags() {
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

type Follow struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	P1    int
	P2    bool
}

func (e Follow) EdgeTypeName() string {
	return "follow"
}

type FollowUpdate struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	P1    int    `norm:"prop:p1;not_null;default:0"`
	P2    bool   `norm:"prop:p2;not_null;default:false"`
	P3    string `norm:"prop:p3;not_null;default:''"`
}

func (e FollowUpdate) EdgeTypeName() string {
	return "follow"
}

func migrateEdges() {
	migrator := db.Debug().Migrator()
	hasEdge, err := migrator.HasEdge("follow")
	if err != nil {
		log.Fatal(err)
	}
	log.Printf("has edge: %v", hasEdge)
	if err = migrator.AutoMigrateEdges(&Follow{}); err != nil {
		log.Fatal(err)
	}
	props, err := migrator.DescEdge("follow")
	if err != nil {
		log.Fatal(err)
	}
	for _, prop := range props {
		log.Printf("edge prop: %+v\n", prop)
	}

	migrator = norm.NewMigrator(db.Debug())
	if err = migrator.AutoMigrateEdges(&FollowUpdate{}); err != nil {
		log.Fatal(err)
	}
	props, err = migrator.DescEdge("follow")
	if err != nil {
		log.Fatal(err)
	}
	for _, prop := range props {
		log.Printf("edge prop: %+v\n", prop)
	}

	if err = migrator.DropEdge("follow"); err != nil {
		log.Fatal(err)
	}
}
