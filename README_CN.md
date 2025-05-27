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
- 🧪 **单元测试覆盖完善**：放心构建生产级应用。
- 💡 **开发者优先设计**：减少样板代码，提高开发效率。

## ⚡ 快速开始

``` go
// Player 节点
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

// Team 节点
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

// Serve 边
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
    // 初始化db对象
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
    
    // 写入player节点
    player := &Player{
        VID:  "player1001",
        Name: "Kobe Bryant",
        Age:  33,
    }
    if err := db.InsertVertex(player).Exec(); err != nil {
        log.Fatalf("insert player failed: %v", err)
    }
    
    // 写入team节点
    team := &Team{
        VID:  "team1001",
        Name: "Lakers",
    }
    if err := db.InsertVertex(team).Exec(); err != nil {
        log.Fatalf("insert team failed: %v", err)
    }
    
    // 写入serve边
    serve := &Serve{
        SrcID:     "player1001",
        DstID:     "team1001",
        StartYear: time.Date(1996, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
        EndYear:   time.Date(2012, 1, 1, 0, 0, 0, 0, time.Local).Unix(),
    }
    if err := db.InsertEdge(serve).Exec(); err != nil {
        log.Fatalf("insert serve failed: %v", err)
    }

    // 查询player节点
    player = new(Player)
    err = db.
        Fetch("player", "player1001").
        Yield("vertex as v").
        FindCol("v", player)
    if err != nil {
        log.Fatalf("fetch player failed: %v", err)
    }
    log.Printf("player: %+v", player)
    
    // 统计player节点通过不同边关联到的节点的数量
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
