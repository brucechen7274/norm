package clause

import (
	"fmt"
	"github.com/haysons/norm/resolver"
	"sort"
	"strconv"
	"strings"
)

type CreateIndex struct {
	TargetType  IndexTarget
	IfNotExists bool
	IndexName   string
	TargetName  string
	Props       []*resolver.FieldIndex
}

const CreateIndexName = "CREATE_INDEX"

type IndexTarget int

const (
	IndexTargetTag IndexTarget = iota + 1
	IndexTargetEdge
)

func (ci CreateIndex) Name() string {
	return CreateIndexName
}

func (ci CreateIndex) MergeIn(clause *Clause) {
	clause.Expression = ci
}

func (ci CreateIndex) Build(nGQL Builder) error {
	nGQL.WriteString("CREATE ")
	switch ci.TargetType {
	case IndexTargetTag:
		nGQL.WriteString("TAG INDEX ")
	case IndexTargetEdge:
		nGQL.WriteString("EDGE INDEX ")
	default:
		return fmt.Errorf("norm: %w, build create index clause failed, invalid target type %d", ErrInvalidClauseParams, ci.TargetType)
	}
	if ci.IfNotExists {
		nGQL.WriteString("IF NOT EXISTS ")
	}
	nGQL.WriteString(ci.IndexName)
	nGQL.WriteString(" ON ")
	nGQL.WriteString(ci.TargetName)
	nGQL.WriteByte('(')
	sort.SliceStable(ci.Props, func(i, j int) bool {
		return ci.Props[i].Priority < ci.Props[j].Priority
	})
	for i, prop := range ci.Props {
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
		if i != len(ci.Props)-1 {
			nGQL.WriteString(", ")
		}
	}
	nGQL.WriteByte(')')
	return nil
}
