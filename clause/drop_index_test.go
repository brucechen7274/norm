package clause_test

import (
	"fmt"
	"github.com/haysons/norm/clause"
	"github.com/haysons/norm/resolver"
	"testing"
)

func TestDropIndex(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.DropIndex{IndexType: resolver.IndexTypeTag, IndexName: "player_index_0"}},
			gqlWant: `DROP TAG INDEX player_index_0`,
		},
		{
			clauses: []clause.Interface{clause.DropIndex{IndexType: resolver.IndexTypeTag, IfExists: true, IndexName: "player_index_0"}},
			gqlWant: `DROP TAG INDEX IF EXISTS player_index_0`,
		},
		{
			clauses: []clause.Interface{clause.DropIndex{IndexType: resolver.IndexTypeEdge, IndexName: "follow_index"}},
			gqlWant: `DROP EDGE INDEX follow_index`,
		},
		{
			clauses: []clause.Interface{clause.DropIndex{IndexType: resolver.IndexTypeEdge, IfExists: true, IndexName: "follow_index"}},
			gqlWant: `DROP EDGE INDEX IF EXISTS follow_index`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
