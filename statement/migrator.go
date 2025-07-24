package statement

import (
	"github.com/haysons/norm/clause"
	"github.com/haysons/norm/resolver"
	"reflect"
)

// CreateVertexTags creates all tags associated with a vertex.
// If a vertex contains multiple tags, each tag will be created sequentially.
//
// Example:
//
//	type player struct {
//		VID  string `norm:"vertex_id"`
//		Name string
//		Age  int
//	}
//
//	func (t player) VertexID() string {
//		return t.VID
//	}
//
//	func (t player) VertexTagName() string {
//		return "player"
//	}
//
// Resulting nGQL:
// CREATE TAG IF NOT EXISTS player(name string, age int)
//
// Usage:
// stmt.CreateVertexTags(&player{}, true)
func (stmt *Statement) CreateVertexTags(vertex any, ifNotExists ...bool) *Statement {
	var notExistsOpt bool
	if len(ifNotExists) > 0 {
		notExistsOpt = ifNotExists[0]
	}

	switch v := vertex.(type) {
	case *resolver.VertexSchema:
		stmt.createVertexTags(v.GetTags(), notExistsOpt)
	case *resolver.VertexTag:
		stmt.createVertexTags([]*resolver.VertexTag{v}, notExistsOpt)
	default:
		vertexType := reflect.TypeOf(vertex)
		vertexSchema, err := resolver.ParseVertex(vertexType)
		if err != nil {
			stmt.err = err
			return stmt
		}
		stmt.createVertexTags(vertexSchema.GetTags(), notExistsOpt)
	}
	return stmt
}

func (stmt *Statement) createVertexTags(tags []*resolver.VertexTag, ifNotExists bool) {
	for i, tag := range tags {
		if i > 0 {
			stmt.AddPart(NewPart())
		}
		stmt.AddClause(&clause.CreateTag{
			IfNotExists: ifNotExists,
			Tag:         tag,
		})
		stmt.SetPartType(PartTypeCreateTag)
	}
}

// DropVertexTag drops a vertex tag by its name.
// If the tag name is empty, the operation is skipped.
// Optionally, you can specify whether to include the IF EXISTS clause.
func (stmt *Statement) DropVertexTag(tagName string, ifExists ...bool) *Statement {
	if tagName == "" {
		return stmt
	}
	var existsOpt bool
	if len(ifExists) > 0 {
		existsOpt = ifExists[0]
	}
	stmt.AddClause(&clause.DropTag{
		TagName:  tagName,
		IfExists: existsOpt,
	})
	stmt.SetPartType(PartTypeDropTag)
	return stmt
}

// AlterVertexTag modifies a tag of the given vertex using the specified operation.
// By default, it alters the first tag of the vertex. To target a specific tag,
// use clause.WithTagName in the options.
//
// Example:
//
//	type player struct {
//		VID  string `norm:"vertex_id"`
//		Name string
//		Age  int
//	}
//
//	func (t player) VertexID() string {
//		return t.VID
//	}
//
//	func (t player) VertexTagName() string {
//		return "player"
//	}
//
// Resulting nGQL:
// ALTER TAG player ADD (name string, age int)
//
// Usage:
// stmt.AlterVertexTag(&player{}, clause.AlterOperate{AddProps: []string{"name", "age"}})
func (stmt *Statement) AlterVertexTag(vertex any, op clause.AlterOperate, opts ...clause.Option) *Statement {
	alterOpts := new(clause.Options)
	for _, opt := range opts {
		opt(alterOpts)
	}
	switch v := vertex.(type) {
	case *resolver.VertexSchema:
		stmt.alterVertexTag(v.GetTags(), op, alterOpts.TagName)
	case *resolver.VertexTag:
		stmt.alterVertexTag([]*resolver.VertexTag{v}, op, alterOpts.TagName)
	default:
		vertexType := reflect.TypeOf(vertex)
		vertexSchema, err := resolver.ParseVertex(vertexType)
		if err != nil {
			stmt.err = err
			return stmt
		}
		stmt.alterVertexTag(vertexSchema.GetTags(), op, alterOpts.TagName)
	}
	return stmt
}

func (stmt *Statement) alterVertexTag(tags []*resolver.VertexTag, op clause.AlterOperate, tagName string) {
	var tag *resolver.VertexTag
	if len(tags) > 1 && tagName != "" {
		for _, t := range tags {
			if t.TagName == tagName {
				tag = t
				break
			}
		}
	} else {
		tag = tags[0]
	}
	stmt.AddClause(&clause.AlterTag{
		Tag:          tag,
		AlterOperate: op,
	})
	stmt.SetPartType(PartTypeAlterTag)
}

// CreateEdge creates an edge schema in the space.
//
//	type follow struct {
//		SrcID  string `norm:"edge_src_id"`
//		DstID  string `norm:"edge_dst_id"`
//		Degree int
//	}
//
//	func (e follow) EdgeTypeName() string {
//		return "follow"
//	}
//
// stmt.CreateEdge(&follow{}, true)
// CREATE EDGE IF NOT EXISTS follow(degree int)
func (stmt *Statement) CreateEdge(edge any, ifNotExists ...bool) *Statement {
	var notExistsOpt bool
	if len(ifNotExists) > 0 {
		notExistsOpt = ifNotExists[0]
	}
	var edgeSchema *resolver.EdgeSchema
	switch e := edge.(type) {
	case *resolver.EdgeSchema:
		edgeSchema = e
	default:
		edgeType := reflect.TypeOf(edge)
		var err error
		edgeSchema, err = resolver.ParseEdge(edgeType)
		if err != nil {
			stmt.err = err
			return stmt
		}
	}
	stmt.AddClause(&clause.CreateEdge{
		IfNotExists: notExistsOpt,
		Edge:        edgeSchema,
	})
	stmt.SetPartType(PartTypeCreateEdge)
	return stmt
}

// DropEdge drops an edge schema by its type name.
//
// stmt.DropEdge("e1", true)
// DROP EDGE IF EXISTS e1
func (stmt *Statement) DropEdge(edgeTypeName string, ifExists ...bool) *Statement {
	if edgeTypeName == "" {
		return stmt
	}
	var existsOpt bool
	if len(ifExists) > 0 {
		existsOpt = ifExists[0]
	}
	stmt.AddClause(&clause.DropEdge{
		EdgeTypeName: edgeTypeName,
		IfExists:     existsOpt,
	})
	stmt.SetPartType(PartTypeDropEdge)
	return stmt
}

// AlterEdge alters the definition of an existing edge type.
//
//	type e2 struct {
//		SrcID string `norm:"edge_src_id"`
//		DstID string `norm:"edge_dst_id"`
//		Rank  int    `norm:"edge_rank"`
//		Name  string `norm:"prop:name"`
//		Age   int    `norm:"prop:age"`
//	}
//
//	func (e *e2) EdgeTypeName() string {
//		return "e2"
//	}
//
// stmt.AlterEdge(&e2{}, clause.AlterOperate{AddProps: []string{"name", "age"}})
// ALTER EDGE e2 ADD (name string, age int)
func (stmt *Statement) AlterEdge(edge any, op clause.AlterOperate) *Statement {
	var edgeSchema *resolver.EdgeSchema
	switch e := edge.(type) {
	case *resolver.EdgeSchema:
		edgeSchema = e
	default:
		edgeType := reflect.TypeOf(edge)
		var err error
		edgeSchema, err = resolver.ParseEdge(edgeType)
		if err != nil {
			stmt.err = err
			return stmt
		}
	}
	stmt.AddClause(&clause.AlterEdge{
		Edge:         edgeSchema,
		AlterOperate: op,
	})
	stmt.SetPartType(PartTypeAlterEdge)
	return stmt
}
