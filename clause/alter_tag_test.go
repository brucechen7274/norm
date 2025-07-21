package clause_test

import (
	"fmt"
	"github.com/haysons/norm/clause"
	"github.com/haysons/norm/resolver"
	"testing"
)

func TestAlterTag(t *testing.T) {
	tests := []struct {
		clauses []clause.Interface
		gqlWant string
		errWant error
	}{
		{
			clauses: []clause.Interface{
				func() clause.AlterTag {
					tag := &resolver.VertexTag{TagName: "t1"}
					tag.SetProps(
						&resolver.Prop{Name: "p3", DataType: "int32"},
						&resolver.Prop{Name: "p4", DataType: "fixed_string(10)"},
					)
					return clause.AlterTag{
						Tag: tag,
						AlterTagOperate: clause.AlterTagOperate{
							AddProps: []string{"p3", "p4"},
						},
					}
				}(),
			},
			gqlWant: `ALTER TAG t1 ADD (p3 int32, p4 fixed_string(10))`,
		},
		{
			clauses: []clause.Interface{
				func() clause.AlterTag {
					tag := &resolver.VertexTag{TagName: "t1"}
					tag.SetProps(
						&resolver.Prop{Name: "p2", DataType: "int32", TTL: "2"},
					)
					return clause.AlterTag{
						Tag: tag,
						AlterTagOperate: clause.AlterTagOperate{
							UpdateTTL: true,
						},
					}
				}(),
			},
			gqlWant: `ALTER TAG t1 TTL_DURATION = 2, TTL_COL = "p2"`,
		},
		{
			clauses: []clause.Interface{
				func() clause.AlterTag {
					tag := &resolver.VertexTag{TagName: "t1"}
					tag.SetProps(
						&resolver.Prop{Name: "p5", DataType: "double", NotNull: true, Default: "0.4", Comment: "p5"},
					)
					return clause.AlterTag{
						Tag: tag,
						AlterTagOperate: clause.AlterTagOperate{
							AddProps: []string{"p5"},
						},
					}
				}(),
			},
			gqlWant: `ALTER TAG t1 ADD (p5 double NOT NULL DEFAULT 0.4 COMMENT "p5")`,
		},
		{
			clauses: []clause.Interface{
				func() clause.AlterTag {
					tag := &resolver.VertexTag{TagName: "t1"}
					tag.SetProps(
						&resolver.Prop{Name: "p3", DataType: "int64"},
						&resolver.Prop{Name: "p4", DataType: "string"},
					)
					return clause.AlterTag{
						Tag: tag,
						AlterTagOperate: clause.AlterTagOperate{
							ChangeProps: []string{"p3", "p4"},
						},
					}
				}(),
			},
			gqlWant: `ALTER TAG t1 CHANGE (p3 int64, p4 string)`,
		},
		{
			clauses: []clause.Interface{
				func() clause.AlterTag {
					tag := &resolver.VertexTag{TagName: "t1"}
					return clause.AlterTag{
						Tag: tag,
						AlterTagOperate: clause.AlterTagOperate{
							DropProps: []string{"p3", "p4"},
						},
					}
				}(),
			},
			gqlWant: `ALTER TAG t1 DROP (p3, p4)`,
		},
		{
			clauses: []clause.Interface{
				func() clause.AlterTag {
					tag := &resolver.VertexTag{TagName: "t1"}
					tag.SetProps(
						&resolver.Prop{Name: "p3", DataType: "int32"},
						&resolver.Prop{Name: "p4", DataType: "fixed_string(10)"},
					)
					return clause.AlterTag{
						Tag: tag,
						AlterTagOperate: clause.AlterTagOperate{
							AddProps:  []string{"p3", "p4"},
							DropProps: []string{"p1", "p2"},
						},
					}
				}(),
			},
			gqlWant: `ALTER TAG t1 ADD (p3 int32, p4 fixed_string(10)), DROP (p1, p2)`,
		},
		{
			clauses: []clause.Interface{
				func() clause.AlterTag {
					tag := &resolver.VertexTag{TagName: "t1"}
					tag.SetProps(
						&resolver.Prop{Name: "p5", DataType: "string"},
						&resolver.Prop{Name: "p6", DataType: "int"},
						&resolver.Prop{Name: "p4", DataType: "fixed_string(12)"},
					)
					return clause.AlterTag{
						Tag: tag,
						AlterTagOperate: clause.AlterTagOperate{
							AddProps:    []string{"p5", "p6"},
							DropProps:   []string{"p3"},
							ChangeProps: []string{"p4"},
						},
					}
				}(),
			},
			gqlWant: `ALTER TAG t1 ADD (p5 string, p6 int), DROP (p3), CHANGE (p4 fixed_string(12))`,
		},
		{
			clauses: []clause.Interface{
				func() clause.AlterTag {
					tag := &resolver.VertexTag{TagName: "t1"}
					tag.SetProps(
						&resolver.Prop{Name: "p9", DataType: "string"},
						&resolver.Prop{Name: "p10", DataType: "int", TTL: "20"},
						&resolver.Prop{Name: "p4", DataType: "string"},
					)
					return clause.AlterTag{
						Tag: tag,
						AlterTagOperate: clause.AlterTagOperate{
							AddProps:    []string{"p9", "p10"},
							ChangeProps: []string{"p4"},
							UpdateTTL:   true,
						},
					}
				}(),
			},
			gqlWant: `ALTER TAG t1 ADD (p9 string, p10 int), CHANGE (p4 string) TTL_DURATION = 20, TTL_COL = "p10"`,
		},
		{
			clauses: []clause.Interface{
				func() clause.AlterTag {
					tag := &resolver.VertexTag{TagName: "t1"}
					tag.SetProps(
						&resolver.Prop{Name: "p9", DataType: "string", NotNull: true, Default: "''"},
					)
					return clause.AlterTag{
						Tag: tag,
						AlterTagOperate: clause.AlterTagOperate{
							DropProps:   []string{"p7", "p8"},
							ChangeProps: []string{"p9"},
						},
					}
				}(),
			},
			gqlWant: `ALTER TAG t1 DROP (p7, p8), CHANGE (p9 string NOT NULL DEFAULT "")`,
		},
	}

	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			testBuildClauses(t, tt.clauses, tt.gqlWant, tt.errWant)
		})
	}
}
