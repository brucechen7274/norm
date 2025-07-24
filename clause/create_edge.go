package clause

import (
	"fmt"
	"github.com/haysons/norm/resolver"
	"strconv"
)

type CreateEdge struct {
	IfNotExists bool
	Edge        *resolver.EdgeSchema
}

const CreateEdgeName = "CREATE_EDGE"

func (ce CreateEdge) Name() string {
	return CreateEdgeName
}

func (ce CreateEdge) MergeIn(clause *Clause) {
	clause.Expression = ce
}

func (ce CreateEdge) Build(nGQL Builder) error {
	nGQL.WriteString("CREATE EDGE ")
	if ce.IfNotExists {
		nGQL.WriteString("IF NOT EXISTS ")
	}
	nGQL.WriteString(ce.Edge.GetTypeName())
	ttlCols, ttlDuration, err := buildProps(ce.Edge.GetProps(), nGQL)
	if err != nil {
		return err
	}
	if len(ttlCols) > 1 {
		return fmt.Errorf("norm: %w, build create edge clause failed, must only one ttl col", ErrInvalidClauseParams)
	}
	if len(ttlCols) == 1 && ttlDuration != "" {
		nGQL.WriteString(" TTL_DURATION = ")
		nGQL.WriteString(ttlDuration)
		nGQL.WriteString(", TTL_COL = ")
		nGQL.WriteString(strconv.Quote(ttlCols[0]))
	}
	return nil
}
