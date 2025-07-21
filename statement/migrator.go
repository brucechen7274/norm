package statement

import (
	"github.com/haysons/norm/clause"
	"github.com/haysons/norm/resolver"
	"reflect"
)

func (stmt *Statement) CreateVertexTags(vertex any, ifNotExist ...bool) *Statement {
	vertexType := reflect.TypeOf(vertex)
	vertexSchema, err := resolver.ParseVertex(vertexType)
	if err != nil {
		stmt.err = err
		return stmt
	}
	var notExistOpt bool
	if len(ifNotExist) > 0 {
		notExistOpt = ifNotExist[0]
	}
	for i, tag := range vertexSchema.GetTags() {
		if i > 0 {
			stmt.AddPart(NewPart())
		}
		stmt.AddClause(&clause.CreateTag{
			IfNotExist: notExistOpt,
			Tag:        tag,
		})
		stmt.SetPartType(PartTypeCreateTag)
	}
	return stmt
}

func (stmt *Statement) DropVertexTag(tagName string, ifExist ...bool) *Statement {
	if tagName == "" {
		return stmt
	}
	var existOpt bool
	if len(ifExist) > 0 {
		existOpt = ifExist[0]
	}
	stmt.AddClause(&clause.DropTag{
		TagName: tagName,
		IfExist: existOpt,
	})
	stmt.SetPartType(PartTypeDropTag)
	return stmt
}

func (stmt *Statement) AlterVertexTag(vertex any, op clause.AlterTagOperate, tagName ...string) *Statement {
	vertexType := reflect.TypeOf(vertex)
	vertexSchema, err := resolver.ParseVertex(vertexType)
	if err != nil {
		stmt.err = err
		return stmt
	}
	var tag *resolver.VertexTag
	if len(vertexSchema.GetTags()) > 1 && len(tagName) > 0 {
		tagNameOpt := tagName[0]
		for _, t := range vertexSchema.GetTags() {
			if t.TagName == tagNameOpt {
				tag = t
			}
		}
	} else {
		tag = vertexSchema.GetTags()[0]
	}
	stmt.AddClause(&clause.AlterTag{
		Tag:             tag,
		AlterTagOperate: op,
	})
	stmt.SetPartType(PartTypeAlterTag)
	return stmt
}
