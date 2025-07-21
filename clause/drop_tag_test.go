package clause_test

import (
	"fmt"
	"github.com/haysons/norm/clause"
	"testing"
)

func TestDropTag(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.DropTag{TagName: "test"}},
			gqlWant: `DROP TAG test`,
		},
		{
			clauses: []clause.Interface{clause.DropTag{TagName: "test", IfExist: true}},
			gqlWant: `DROP TAG IF EXISTS test`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
