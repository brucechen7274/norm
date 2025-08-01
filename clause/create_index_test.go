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
						Index: &resolver.Index{
							Name:   "player_index",
							Type:   resolver.IndexTypeTag,
							Target: "player",
						},
					}
				}(),
			},
			gqlWant: `CREATE TAG INDEX player_index ON player()`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateIndex {
					return clause.CreateIndex{
						Index: &resolver.Index{
							Name:   "follow_index",
							Type:   resolver.IndexTypeEdge,
							Target: "follow",
						},
					}
				}(),
			},
			gqlWant: `CREATE EDGE INDEX follow_index ON follow()`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateIndex {
					return clause.CreateIndex{
						IfNotExists: true,
						Index: &resolver.Index{
							Name:   "var",
							Type:   resolver.IndexTypeTag,
							Target: "var_string",
							Fields: []*resolver.IndexField{
								{
									Prop:     "p1",
									DataType: "string",
									Length:   10,
								},
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
						IfNotExists: true,
						Index: &resolver.Index{
							Name:   "follow_index_0",
							Type:   resolver.IndexTypeEdge,
							Target: "follow",
							Fields: []*resolver.IndexField{
								{
									Prop:     "degree",
									DataType: "int",
								},
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
						IfNotExists: true,
						Index: &resolver.Index{
							Name:   "player_index_1",
							Type:   resolver.IndexTypeTag,
							Target: "player",
							Fields: []*resolver.IndexField{
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
						IfNotExists: true,
						Index: &resolver.Index{
							Name:   "player_index_1",
							Type:   resolver.IndexTypeTag,
							Target: "player",
							Fields: []*resolver.IndexField{
								{
									Prop:     "name",
									DataType: "string",
								},
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
