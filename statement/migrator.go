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

func (stmt *Statement) DropTag(tagName string, ifExist ...bool) *Statement {
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
