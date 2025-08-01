package clause

import (
	"fmt"
	"github.com/haysons/norm/resolver"
	"sort"
	"strconv"
	"strings"
)

type CreateIndex struct {
	IfNotExists bool
	Index       *resolver.Index
}

const CreateIndexName = "CREATE_INDEX"

func (ci CreateIndex) Name() string {
	return CreateIndexName
}

func (ci CreateIndex) MergeIn(clause *Clause) {
	clause.Expression = ci
}

func (ci CreateIndex) Build(nGQL Builder) error {
	nGQL.WriteString("CREATE ")
	switch ci.Index.Type {
	case resolver.IndexTypeTag:
		nGQL.WriteString("TAG INDEX ")
	case resolver.IndexTypeEdge:
		nGQL.WriteString("EDGE INDEX ")
	default:
		return fmt.Errorf("norm: %w, build create index clause failed, invalid target type %d", ErrInvalidClauseParams, ci.Index.Type)
	}
	if ci.IfNotExists {
		nGQL.WriteString("IF NOT EXISTS ")
	}
	nGQL.WriteString(ci.Index.Name)
	nGQL.WriteString(" ON ")
	nGQL.WriteString(ci.Index.Target)
	nGQL.WriteByte('(')
	sort.SliceStable(ci.Index.Fields, func(i, j int) bool {
		return ci.Index.Fields[i].Priority < ci.Index.Fields[j].Priority
	})
	for i, prop := range ci.Index.Fields {
		nGQL.WriteString(prop.Prop)
		if strings.ToLower(prop.DataType) == "string" {
			if prop.Length == 0 {
				return fmt.Errorf("norm: %w, build create index clause failed, string property must has index length", ErrInvalidClauseParams)
			} else {
				nGQL.WriteByte('(')
				nGQL.WriteString(strconv.Itoa(prop.Length))
				nGQL.WriteByte(')')
			}
		}
		if i != len(ci.Index.Fields)-1 {
			nGQL.WriteString(", ")
		}
	}
	nGQL.WriteByte(')')
	return nil
}
