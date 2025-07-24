package clause_test

import (
	"fmt"
	"github.com/haysons/norm/clause"
	"testing"
)

func TestDropEdge(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.DropEdge{EdgeTypeName: "e1"}},
			gqlWant: `DROP EDGE e1`,
		},
		{
			clauses: []clause.Interface{clause.DropEdge{EdgeTypeName: "e1", IfExists: true}},
			gqlWant: `DROP EDGE IF EXISTS e1`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
