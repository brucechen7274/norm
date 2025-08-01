package clause

import (
	"fmt"
	"github.com/haysons/norm/resolver"
)

type RebuildIndex struct {
	IndexType  resolver.IndexType
	IndexNames []string
}

const RebuildIndexName = "REBUILD_INDEX"

func (ri RebuildIndex) Name() string {
	return RebuildIndexName
}

func (ri RebuildIndex) MergeIn(clause *Clause) {
	clause.Expression = ri
}

func (ri RebuildIndex) Build(nGQL Builder) error {
	nGQL.WriteString("REBUILD ")
	switch ri.IndexType {
	case resolver.IndexTypeTag:
		nGQL.WriteString("TAG INDEX ")
	case resolver.IndexTypeEdge:
		nGQL.WriteString("EDGE INDEX ")
	default:
		return fmt.Errorf("norm: %w, build rebuild index clause failed, invalid target type %d", ErrInvalidClauseParams, ri.IndexType)
	}
	if len(ri.IndexNames) == 0 {
		return fmt.Errorf("norm: %w, build rebuild index clause failed, index names empty", ErrInvalidClauseParams)
	}
	for i, name := range ri.IndexNames {
		nGQL.WriteString(name)
		if i != len(ri.IndexNames)-1 {
			nGQL.WriteString(", ")
		}
	}
	return nil
}
