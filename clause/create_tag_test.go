package clause_test

import (
	"fmt"
	"github.com/haysons/norm/clause"
	"github.com/haysons/norm/resolver"
	"testing"
)

func TestCreateTag(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{
				func() clause.CreateTag {
					tag := &resolver.VertexTag{TagName: "player"}
					tag.SetProps(
						&resolver.Prop{Name: "name", DataType: "string", NotNull: true},
						&resolver.Prop{Name: "age", DataType: "int"},
					)
					return clause.CreateTag{
						IfNotExists: true,
						Tag:         tag,
					}
				}(),
			},
			gqlWant: `CREATE TAG IF NOT EXISTS player(name string NOT NULL, age int)`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateTag {
					tag := &resolver.VertexTag{TagName: "no_property"}
					return clause.CreateTag{
						IfNotExists: false,
						Tag:         tag,
					}
				}(),
			},
			gqlWant: `CREATE TAG no_property()`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateTag {
					tag := &resolver.VertexTag{TagName: "player_with_default"}
					tag.SetProps(
						&resolver.Prop{Name: "name", DataType: "string", Default: "default name"},
						&resolver.Prop{Name: "age", DataType: "int", Default: "20"},
					)
					return clause.CreateTag{
						IfNotExists: true,
						Tag:         tag,
					}
				}(),
			},
			gqlWant: `CREATE TAG IF NOT EXISTS player_with_default(name string DEFAULT "default name", age int DEFAULT 20)`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateTag {
					tag := &resolver.VertexTag{TagName: "woman"}
					tag.SetProps(
						&resolver.Prop{Name: "name", DataType: "string"},
						&resolver.Prop{Name: "age", DataType: "int"},
						&resolver.Prop{Name: "married", DataType: "bool"},
						&resolver.Prop{Name: "salary", DataType: "double"},
						&resolver.Prop{Name: "create_time", DataType: "timestamp", TTL: "100"},
					)
					return clause.CreateTag{
						IfNotExists: true,
						Tag:         tag,
					}
				}(),
			},
			gqlWant: `CREATE TAG IF NOT EXISTS woman(name string, age int, married bool, salary double, create_time timestamp) TTL_DURATION = 100, TTL_COL = "create_time"`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateTag {
					tag := &resolver.VertexTag{TagName: "date1"}
					tag.SetProps(
						&resolver.Prop{Name: "p1", DataType: "date"},
						&resolver.Prop{Name: "p2", DataType: "time"},
						&resolver.Prop{Name: "p3", DataType: "datetime"},
					)
					return clause.CreateTag{
						IfNotExists: true,
						Tag:         tag,
					}
				}(),
			},
			gqlWant: `CREATE TAG IF NOT EXISTS date1(p1 date, p2 time, p3 datetime)`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateTag {
					tag := &resolver.VertexTag{TagName: "date1"}
					tag.SetProps(
						&resolver.Prop{Name: "p1", DataType: "date", NotNull: true, Default: `date("2021-03-17")`},
						&resolver.Prop{Name: "p2", DataType: "time", Default: `time("17:53:59")`},
						&resolver.Prop{Name: "p3", DataType: "datetime", Default: `datetime("2017-03-04T22:30:40.003000[Asia/Shanghai]")`},
					)
					return clause.CreateTag{
						IfNotExists: true,
						Tag:         tag,
					}
				}(),
			},
			gqlWant: `CREATE TAG IF NOT EXISTS date1(p1 date NOT NULL DEFAULT date("2021-03-17"), p2 time DEFAULT time("17:53:59"), p3 datetime DEFAULT datetime("2017-03-04T22:30:40.003000[Asia/Shanghai]"))`,
		},
		{
			clauses: []clause.Interface{
				func() clause.CreateTag {
					tag := &resolver.VertexTag{TagName: "t1"}
					tag.SetProps(
						&resolver.Prop{Name: "name", DataType: "string", TTL: "100"},
						&resolver.Prop{Name: "age", DataType: "int", TTL: "100"},
					)
					return clause.CreateTag{
						IfNotExists: true,
						Tag:         tag,
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
