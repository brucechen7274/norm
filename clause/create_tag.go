package clause

import (
	"fmt"
	"github.com/haysons/norm/resolver"
	"reflect"
	"strconv"
)

type CreateTag struct {
	IfNotExist bool
	Vertex     reflect.Value
}

const CreateTagName = "CREATE_TAG"

func (ct CreateTag) Name() string {
	return CreateTagName
}

func (ct CreateTag) MergeIn(clause *Clause) {
	clause.Expression = ct
}

func (ct CreateTag) Build(nGQL Builder) error {
	ct.Vertex = reflect.Indirect(ct.Vertex)
	if ct.Vertex.Kind() != reflect.Struct {
		return fmt.Errorf("norm: %w, build create tag clause failed, dest must be struct or pointer", ErrInvalidClauseParams)
	}
	vertexType := ct.Vertex.Type()
	vertexSchema, err := resolver.ParseVertex(vertexType)
	if err != nil {
		return err
	}
	vertexTags := vertexSchema.GetTags()
	for i, tag := range vertexTags {
		if err = ct.buildCreateTag(nGQL, tag); err != nil {
			return err
		}
		if i != len(vertexTags)-1 {
			nGQL.WriteString("; ")
		}
	}
	return nil
}

func (ct CreateTag) buildCreateTag(nGQL Builder, tag *resolver.VertexTag) error {
	nGQL.WriteString("CREATE TAG ")
	if ct.IfNotExist {
		nGQL.WriteString("IF NOT EXISTS ")
	}
	nGQL.WriteString(tag.TagName)
	nGQL.WriteByte('(')
	propsLen := len(tag.GetProps())
	ttlCols := make([]string, 0, 1)
	ttlDuration := ""
	for i, prop := range tag.GetProps() {
		nGQL.WriteString(prop.Name)
		nGQL.WriteByte(' ')
		nGQL.WriteString(prop.NebulaType)
		if prop.NotNull {
			nGQL.WriteString(" NOT NULL")
		}
		if prop.Default != "" {
			nGQL.WriteString(" DEFAULT ")
			nGQL.WriteString(prop.Default)
		}
		if prop.Comment != "" {
			nGQL.WriteString(" COMMENT ")
			nGQL.WriteByte('\'')
			nGQL.WriteString(prop.Comment)
			nGQL.WriteByte('\'')
		}
		if prop.TTL != "" {
			ttlCols = append(ttlCols, prop.Name)
			ttlDuration = prop.TTL
		}
		if i < propsLen-1 {
			nGQL.WriteString(", ")
		}
	}
	nGQL.WriteByte(')')
	if len(ttlCols) > 1 {
		return fmt.Errorf("norm: %w, build create tag clause failed, must only one ttl col", ErrInvalidClauseParams)
	}
	if len(ttlCols) == 1 && ttlDuration != "" {
		nGQL.WriteString(" TTL_DURATION ")
		nGQL.WriteString(ttlDuration)
		nGQL.WriteString(", ")
		nGQL.WriteString(strconv.Quote(ttlCols[0]))
	}
	return nil
}
