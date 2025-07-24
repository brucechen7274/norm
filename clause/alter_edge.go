package clause

import (
	"fmt"
	"github.com/haysons/norm/resolver"
	"slices"
)

type AlterEdge struct {
	Edge *resolver.EdgeSchema
	AlterOperate
}

const AlterEdgeName = "ALTER_EDGE"

func (ae AlterEdge) Name() string {
	return AlterEdgeName
}

func (ae AlterEdge) MergeIn(clause *Clause) {
	clause.Expression = ae
}

func (ae AlterEdge) Build(nGQL Builder) error {
	nGQL.WriteString("ALTER EDGE ")
	nGQL.WriteString(ae.Edge.GetTypeName())
	addProps := make([]*resolver.Prop, 0, len(ae.AddProps))
	changeProps := make([]*resolver.Prop, 0, len(ae.DropProps))
	ttlCols := make([]string, 0, 1)
	ttlDuration := ""
	for _, prop := range ae.Edge.GetProps() {
		if slices.Contains(ae.AddProps, prop.Name) {
			addProps = append(addProps, prop)
		}
		if slices.Contains(ae.ChangeProps, prop.Name) {
			changeProps = append(changeProps, prop)
		}
		if prop.TTL != "" {
			ttlCols = append(ttlCols, prop.Name)
			ttlDuration = prop.TTL
		}
	}
	if len(ttlCols) > 1 {
		return fmt.Errorf("norm: %w, build alter edge clause failed, must only one ttl col", ErrInvalidClauseParams)
	}
	if len(addProps) == 0 && len(ae.DropProps) == 0 && len(changeProps) == 0 && !ae.UpdateTTL {
		return fmt.Errorf("norm: %w, build alter edge clause failed, must has operate", ErrInvalidClauseParams)
	}
	return buildAlterProps(addProps, changeProps, ae.DropProps, ae.UpdateTTL, ttlCols, ttlDuration, nGQL)
}
