# norm

[中文](README_CN.md)

[![go report card](https://goreportcard.com/badge/haysons/norm)](https://goreportcard.com/report/github.com/haysons/norm)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## 🚀 Introduction

**norm** is a lightweight, developer-friendly ORM framework designed specifically
for [nebula graph](https://nebula-graph.io).  
It aims to simplify the Go development experience by enabling elegant, chainable `nGQL` query construction and seamless
result mapping.

Whether you're building a graph-based social network or a knowledge graph platform, `norm` helps you move fast without
sacrificing readability or maintainability.

## 📦 Installation

```bash
go get github.com/haysons/norm
```

## ✨ Features

- 🔗 **Chainable nGQL builder**: Write readable, elegant queries with fluent chaining.
- 📦 **Struct-based mapping**: Map query results directly into Go structs.
- 🧠 **Smart parsing**: Supports nested types — vertex, edge, list, map, set — with ease.
- 📚 **Struct embedding support**: Maximize code reuse and maintain clarity.
- 🧪 **Fully unit tested**: Confidently build production-grade apps.
- 💡 **Developer-first design**: Less boilerplate, more productivity.

## ⚡ Quick Start

``` go
// Player vertex player
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

// Team vertex team
type Team struct {
    VID  string `norm:"vertex_id"`
    Name string `norm:"prop:name"`
}

func (t Team) VertexID() string {
    return t.VID
}

func (t Team) VertexTagName() string {
    return "team"
}

// Serve edge Serve
type Serve struct {
    SrcID     string `norm:"edge_src_id"`
    DstID     string `norm:"edge_dst_id"`
    Rank      int    `norm:"edge_rank"`
    StartYear int64  `norm:"prop:start_year"`
    EndYear   int64  `norm:"prop:end_year"`
}

func (s Serve) EdgeTypeName() string {
    return "serve"
}

func main() {
    // initialize the db object
    conf := &norm.Config{
        Username:    "root",
        Password:    "nebula",
        SpaceName:   "demo_basketballplayer",
        Addresses:   []string{"127.0.0.1:9669"},
    }
    db, err := norm.Open(conf)
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    
    // insert the player vertex
    player := &Player{
        VID:  "player1001",
        Name: "Kobe Bryant",
        Age:  33,
    }
    if err := db.InsertVertex(player).Exec(); err != nil {
        log.Fatalf("insert player failed: %v", err)
    }
    
    // insert the team vertex
    team := &Team{
        VID:  "team1001",
        Name: "Lakers",
    }
    if err := db.InsertVertex(team).Exec(); err != nil {
        log.Fatalf("insert team failed: %v", err)
    }
    
    // insert the serve edge
    serve := &Serve{
        SrcID:     "player1001",
        DstID:     "team1001",
        StartYear: time.Date(1996, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
        EndYear:   time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
    }
    if err := db.InsertEdge(serve).Exec(); err != nil {
        log.Fatalf("insert serve failed: %v", err)
    }

    // find the player vertex
    player = new(Player)
    err = db.
        Fetch("player", "player1001").
        Yield("vertex as v").
        FindCol("v", player)
    if err != nil {
        log.Fatalf("fetch player failed: %v", err)
    }
    log.Printf("player: %+v", player)
    
    // count the number of vertexes that the player vertex connects through different edges
    type edgeCnt struct {
        Edge string `norm:"col:e"`
        Cnt  int    `norm:"col:cnt"`
    }
    edgesCnt := make([]*edgeCnt, 0)
    err = db.Go().
        From("player1001").
        Over("*").
        Yield("type(edge) as t").
        GroupBy("$-.t").
        Yield("$-.t as e, count(*) as cnt").
        Find(&edgesCnt)
    if err != nil {
        log.Fatalf("get edge cnt failed: %v", err)
    }
    for _, c := range edgesCnt {
        log.Printf("edge cnt: %+v\n", c)
    }
}
```

## 🤝 Contributing

We welcome contributions from the community!

- 🍴 Fork the repo
- 🔧 Create a feature branch
- ✅ Submit a pull request

## 🙏 Acknowledgements

Special thanks to the following projects that inspired and supported `norm`:

- [**gorm**](https://gorm.io): The beloved ORM for Golang — simple, powerful, elegant.

## 📄 License

© 2024–NOW [@hayson](https://github.com/haysons)

Released under the [MIT License](./LICENSE)