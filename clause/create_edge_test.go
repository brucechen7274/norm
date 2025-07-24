package clause_test

import (
	"fmt"
	"github.com/haysons/norm/clause"
	"github.com/haysons/norm/resolver"
	"testing"
)

func TestCreateEdge(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{
				func() clause.CreateEdge {
					edge := &resolver.EdgeSchema{}
					edge.SetTypeName("follow")
					edge.SetProps(
						&resolver.Prop{Name: "degree", DataType: "int"},
					)
					return clause.CreateEdge{
						IfNotExists: true,
						Edge:        edge,
					}
				}(),
			},
			gqlWant: `CREATE EDGE IF NOT EXISTS follow(degree int)`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateEdge {
					edge := &resolver.EdgeSchema{}
					edge.SetTypeName("no_property")
					return clause.CreateEdge{
						Edge: edge,
					}
				}(),
			},
			gqlWant: `CREATE EDGE no_property()`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateEdge {
					edge := &resolver.EdgeSchema{}
					edge.SetTypeName("follow_with_default")
					edge.SetProps(
						&resolver.Prop{Name: "degree", DataType: "int", Default: "20"},
					)
					return clause.CreateEdge{
						IfNotExists: true,
						Edge:        edge,
					}
				}(),
			},
			gqlWant: `CREATE EDGE IF NOT EXISTS follow_with_default(degree int DEFAULT 20)`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateEdge {
					edge := &resolver.EdgeSchema{}
					edge.SetTypeName("e1")
					edge.SetProps(
						&resolver.Prop{Name: "p1", DataType: "string"},
						&resolver.Prop{Name: "p2", DataType: "int", TTL: "100"},
						&resolver.Prop{Name: "p3", DataType: "timestamp"},
					)
					return clause.CreateEdge{
						IfNotExists: true,
						Edge:        edge,
					}
				}(),
			},
			gqlWant: `CREATE EDGE IF NOT EXISTS e1(p1 string, p2 int, p3 timestamp) TTL_DURATION = 100, TTL_COL = "p2"`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateEdge {
					edge := &resolver.EdgeSchema{}
					edge.SetTypeName("e1")
					edge.SetProps(
						&resolver.Prop{Name: "name", DataType: "string", TTL: "100"},
						&resolver.Prop{Name: "age", DataType: "int", TTL: "100"},
					)
					return clause.CreateEdge{
						IfNotExists: true,
						Edge:        edge,
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
