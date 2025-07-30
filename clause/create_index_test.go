package clause_test

import (
	"fmt"
	"github.com/haysons/norm/clause"
	"github.com/haysons/norm/resolver"
	"testing"
)

func TestCreateIndex(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{
				func() clause.CreateIndex {
					return clause.CreateIndex{
						TargetType: clause.IndexTargetTag,
						IndexName:  "player_index",
						TargetName: "player",
					}
				}(),
			},
			gqlWant: `CREATE TAG INDEX player_index ON player()`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateIndex {
					return clause.CreateIndex{
						TargetType: clause.IndexTargetEdge,
						IndexName:  "follow_index",
						TargetName: "follow",
					}
				}(),
			},
			gqlWant: `CREATE EDGE INDEX follow_index ON follow()`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateIndex {
					return clause.CreateIndex{
						TargetType:  clause.IndexTargetTag,
						IfNotExists: true,
						IndexName:   "var",
						TargetName:  "var_string",
						Props: []*resolver.FieldIndex{
							{
								Prop:     "p1",
								DataType: "string",
								Length:   10,
							},
						},
					}
				}(),
			},
			gqlWant: `CREATE TAG INDEX IF NOT EXISTS var ON var_string(p1(10))`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateIndex {
					return clause.CreateIndex{
						TargetType:  clause.IndexTargetEdge,
						IfNotExists: true,
						IndexName:   "follow_index_0",
						TargetName:  "follow",
						Props: []*resolver.FieldIndex{
							{
								Prop:     "degree",
								DataType: "int",
							},
						},
					}
				}(),
			},
			gqlWant: `CREATE EDGE INDEX IF NOT EXISTS follow_index_0 ON follow(degree)`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateIndex {
					return clause.CreateIndex{
						TargetType:  clause.IndexTargetTag,
						IfNotExists: true,
						IndexName:   "player_index_1",
						TargetName:  "player",
						Props: []*resolver.FieldIndex{
							{
								Prop:     "name",
								DataType: "string",
								Length:   10,
							},
							{
								Prop:     "age",
								DataType: "int",
							},
						},
					}
				}(),
			},
			gqlWant: `CREATE TAG INDEX IF NOT EXISTS player_index_1 ON player(name(10), age)`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateIndex {
					return clause.CreateIndex{
						TargetType:  clause.IndexTargetTag,
						IfNotExists: true,
						IndexName:   "player_index_1",
						TargetName:  "player",
						Props: []*resolver.FieldIndex{
							{
								Prop:     "name",
								DataType: "string",
							},
						},
					}
				}(),
			},
			errWant: clause.ErrInvalidClauseParams,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
