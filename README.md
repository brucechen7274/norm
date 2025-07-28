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
- 🔄 **Auto schema migration**: Automatically create or update vertex and edge schemas from structs.
- 🧪 **Fully unit tested**: Confidently build production-grade apps.
- 💡 **Developer-first design**: Less boilerplate, more productivity.

## ⚡ Quick Start

``` go
package main

import (
	"github.com/haysons/norm"
	"log"
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
	// init norm.DB
	conf := &norm.Config{
		Username:  "root",
		Password:  "nebula",
		SpaceName: "test",
		Addresses: []string{"127.0.0.1:9669"},
	}
	db, err := norm.Open(conf)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// migrate vertex player tags
	if err = db.Migrator().AutoMigrateVertexes(Player{}); err != nil {
		log.Fatalf("auto migrate vertex palyer failed: %v", err)
	}

	// insert the player vertex
	player := &Player{
		VID:  "player1001",
		Name: "Kobe Bryant",
		Age:  33,
	}
	if err := db.InsertVertex(player).Exec(); err != nil {
		log.Fatalf("insert vertex player failed: %v", err)
	}

	// find the player vertex
	player = new(Player)
	err = db.
		Fetch("player", "player1001").
		Yield("vertex as v").
		FindCol("v", player)
	if err != nil {
		log.Fatalf("fetch vertex player failed: %v", err)
	}
	log.Printf("player: %+v", player)
}
```

📚 See more usage patterns in the [example directory](./example).

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