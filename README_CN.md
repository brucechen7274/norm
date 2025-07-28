# norm

[English](README.md)

[![go report card](https://goreportcard.com/badge/haysons/norm)](https://goreportcard.com/report/github.com/haysons/norm)
[![MIT license](https://img.shields.io/badge/license-MIT-brightgreen.svg)](https://opensource.org/licenses/MIT)

## 🚀 介绍

**norm** 是一个轻量且开发者友好的 ORM 框架，专为 [nebula graph](https://nebula-graph.io) 设计。  
它旨在简化 Go 语言下的 nebula graph 开发体验，实现优雅且链式的 `nGQL` 查询构建，并支持无缝的结果映射。

无论你是在构建基于图的社交网络，还是知识图谱平台，`norm` 都能帮助你快速开发，同时保证代码的可读性和可维护性。

## 📦 安装

```bash
go get github.com/haysons/norm
```

## ✨ 特性

- 🔗 **链式 nGQL 构建器**：通过流畅的链式调用书写可读且优雅的查询语句。
- 📦 **基于结构体的映射**：查询结果可直接映射到 Go 结构体。
- 🧠 **智能解析**：轻松支持嵌套类型 — 顶点（vertex）、边（edge）、列表（list）、映射（map）、集合（set）等。
- 📚 **支持结构体内嵌**：最大化代码复用，同时保持代码清晰。
- 🔄 **自动迁移节点与边结构**：根据结构体定义自动创建或变更对应的 tag / edge schema。
- 🧪 **单元测试覆盖完善**：放心构建生产级应用。
- 💡 **开发者优先设计**：减少样板代码，提高开发效率。

## ⚡ 快速开始

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

📚 更多使用示例请参考 [example 目录](./example)。

## 🤝 贡献

欢迎社区的贡献！

- 🍴 Fork 本仓库
- 🔧 创建功能分支
- ✅ 提交 Pull Request

## 🙏 致谢

特别感谢以下项目对 `norm` 的启发和支持：

- [**gorm**](https://gorm.io)：广受喜爱的 Golang ORM，简单、强大且优雅。

## 📄 许可证

© 2024–至今 [@hayson](https://github.com/haysons)

基于 [MIT 许可证](./LICENSE) 发行
