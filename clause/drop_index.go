package clause

import "fmt"

type DropIndex struct {
	TargetType IndexTarget
	IfExists   bool
	IndexName  string
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
	switch di.TargetType {
	case IndexTargetTag:
		nGQL.WriteString("TAG INDEX ")
	case IndexTargetEdge:
		nGQL.WriteString("EDGE INDEX ")
	default:
		return fmt.Errorf("norm: %w, build drop index clause failed, invalid target type %d", ErrInvalidClauseParams, di.TargetType)
	}
	if di.IfExists {
		nGQL.WriteString("IF EXISTS ")
	}
	nGQL.WriteString(di.IndexName)
	return nil
}
