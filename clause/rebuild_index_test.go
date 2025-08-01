package clause_test

import (
	"fmt"
	"github.com/haysons/norm/clause"
	"github.com/haysons/norm/resolver"
	"testing"
)

func TestRebuildIndex(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{clause.RebuildIndex{IndexType: resolver.IndexTypeTag, IndexNames: []string{"single_person_index"}}},
			gqlWant: `REBUILD TAG INDEX single_person_index`,
		},
		{
			clauses: []clause.Interface{clause.RebuildIndex{IndexType: resolver.IndexTypeTag, IndexNames: []string{"idx1", "idx2"}}},
			gqlWant: `REBUILD TAG INDEX idx1, idx2`,
		},
		{
			clauses: []clause.Interface{clause.RebuildIndex{IndexType: resolver.IndexTypeEdge, IndexNames: []string{"idx1"}}},
			gqlWant: `REBUILD EDGE INDEX idx1`,
		},
		{
			clauses: []clause.Interface{clause.RebuildIndex{IndexType: resolver.IndexTypeEdge, IndexNames: []string{"idx1", "idx2"}}},
			gqlWant: `REBUILD EDGE INDEX idx1, idx2`,
		},
		{
			clauses: []clause.Interface{clause.RebuildIndex{IndexType: resolver.IndexTypeEdge}},
			errWant: clause.ErrInvalidClauseParams,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
