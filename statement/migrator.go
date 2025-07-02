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
