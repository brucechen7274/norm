package statement

import (
	"fmt"
	"github.com/haysons/nebulaorm/clause"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatement(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				stmt := New()
				stmt.SetClausesBuild([]string{clause.FromName, clause.OrderName})
				stmt.AddClause(&clause.From{VID: "team1"})
				stmt.AddClause(&clause.Order{Expr: "$-.id"})
				return stmt
			},
			want: `FROM "team1" ORDER BY $-.id;`,
		},
		{
			stmt: func() *Statement {
				stmt := New()
				stmt.SetClausesBuild([]string{clause.GoName, clause.OrderName})
				stmt.SetPartType(PartTypeFetch)
				stmt.AddClause(&clause.Go{StepStart: 1, StepEnd: 2})
				stmt.AddClause(&clause.Order{Expr: "$-.id"})
				return stmt
			},
			want: `GO 1 TO 2 STEPS ORDER BY $-.id;`,
		},
		{
			stmt: func() *Statement {
				stmt := New()
				stmt.Go(1).Pipe().Pipe().GroupBy("$-.id")
				return stmt
			},
			want: `GO 1 STEPS | GROUP BY $-.id;`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, ngql)
			}
		})
	}
}
