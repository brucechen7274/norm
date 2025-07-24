package clause

type DropEdge struct {
	IfExists     bool
	EdgeTypeName string
}

const DropEdgeName = "DROP_EDGE"

func (de DropEdge) Name() string {
	return DropEdgeName
}

func (de DropEdge) MergeIn(clause *Clause) {
	clause.Expression = de
}

func (de DropEdge) Build(nGQL Builder) error {
	nGQL.WriteString("DROP EDGE ")
	if de.IfExists {
		nGQL.WriteString("IF EXISTS ")
	}
	nGQL.WriteString(de.EdgeTypeName)
	return nil
}
