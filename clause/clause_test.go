package clause_test

import (
	"github.com/haysons/nebulaorm/clause"
	"github.com/haysons/nebulaorm/statement"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func testBuildClauses(t *testing.T, clauses []clause.Interface, gqlWant string, errWant error) {
	buildNames := make([]string, len(clauses))
	buildNamesMap := make(map[string]bool, len(clauses))
	stmtPart := statement.NewPart()
	for _, c := range clauses {
		if _, ok := buildNamesMap[c.Name()]; !ok {
			buildNames = append(buildNames, c.Name())
			buildNamesMap[c.Name()] = true
		}
		stmtPart.AddClause(c)
	}
	stmtPart.SetClausesBuild(buildNames)
	gqlBuilder := new(strings.Builder)
	err := stmtPart.Build(gqlBuilder)
	gql := gqlBuilder.String()
	assert.ErrorIs(t, err, errWant)
	if err == nil {
		assert.Equal(t, gql, gqlWant)
	}
}
