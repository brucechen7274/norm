package clause_test

import (
	"fmt"
	"github.com/haysons/norm/clause"
	"testing"
)

func TestDropIndex(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.DropIndex{TargetType: clause.IndexTargetTag, IndexName: "player_index_0"}},
			gqlWant: `DROP TAG INDEX player_index_0`,
		},
		{
			clauses: []clause.Interface{clause.DropIndex{TargetType: clause.IndexTargetTag, IfExists: true, IndexName: "player_index_0"}},
			gqlWant: `DROP TAG INDEX IF EXISTS player_index_0`,
		},
		{
			clauses: []clause.Interface{clause.DropIndex{TargetType: clause.IndexTargetEdge, IndexName: "follow_index"}},
			gqlWant: `DROP EDGE INDEX follow_index`,
		},
		{
			clauses: []clause.Interface{clause.DropIndex{TargetType: clause.IndexTargetEdge, IfExists: true, IndexName: "follow_index"}},
			gqlWant: `DROP EDGE INDEX IF EXISTS follow_index`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
