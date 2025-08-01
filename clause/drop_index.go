package clause

import (
	"fmt"
	"github.com/haysons/norm/resolver"
)

type DropIndex struct {
	IndexType resolver.IndexType
	IfExists  bool
	IndexName string
}

const DropIndexName = "DROP_INDEX"

func (di DropIndex) Name() string {
	return DropIndexName
}

func (di DropIndex) MergeIn(clause *Clause) {
	clause.Expression = di
}

func (di DropIndex) Build(nGQL Builder) error {
	nGQL.WriteString("DROP ")
	switch di.IndexType {
	case resolver.IndexTypeTag:
		nGQL.WriteString("TAG INDEX ")
	case resolver.IndexTypeEdge:
		nGQL.WriteString("EDGE INDEX ")
	default:
		return fmt.Errorf("norm: %w, build drop index clause failed, invalid target type %d", ErrInvalidClauseParams, di.IndexType)
	}
	if di.IfExists {
		nGQL.WriteString("IF EXISTS ")
	}
	nGQL.WriteString(di.IndexName)
	return nil
}
