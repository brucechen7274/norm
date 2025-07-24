package clause

import (
	"fmt"
	"github.com/haysons/norm/resolver"
	"slices"
	"strconv"
)

type AlterTag struct {
	Tag *resolver.VertexTag
	AlterOperate
}

type AlterOperate struct {
	AddProps    []string
	DropProps   []string
	ChangeProps []string
	UpdateTTL   bool
}

const AlterTagName = "ALTER_TAG"

func (at AlterTag) Name() string {
	return AlterTagName
}

func (at AlterTag) MergeIn(clause *Clause) {
	clause.Expression = at
}

func (at AlterTag) Build(nGQL Builder) error {
	nGQL.WriteString("ALTER TAG ")
	nGQL.WriteString(at.Tag.TagName)
	addProps := make([]*resolver.Prop, 0, len(at.AddProps))
	changeProps := make([]*resolver.Prop, 0, len(at.DropProps))
	ttlCols := make([]string, 0, 1)
	ttlDuration := ""
	for _, prop := range at.Tag.GetProps() {
		if slices.Contains(at.AddProps, prop.Name) {
			addProps = append(addProps, prop)
		}
		if slices.Contains(at.ChangeProps, prop.Name) {
			changeProps = append(changeProps, prop)
		}
		if prop.TTL != "" {
			ttlCols = append(ttlCols, prop.Name)
			ttlDuration = prop.TTL
		}
	}
	if len(ttlCols) > 1 {
		return fmt.Errorf("norm: %w, build alter tag clause failed, must only one ttl col", ErrInvalidClauseParams)
	}
	if len(addProps) == 0 && len(at.DropProps) == 0 && len(changeProps) == 0 && !at.UpdateTTL {
		return fmt.Errorf("norm: %w, build alter tag clause failed, must has operate", ErrInvalidClauseParams)
	}
	return buildAlterProps(addProps, changeProps, at.DropProps, at.UpdateTTL, ttlCols, ttlDuration, nGQL)
}

func buildAlterProps(addProps, changeProps []*resolver.Prop, dropProps []string, updateTTL bool, ttlCols []string, ttlDuration string, nGQL Builder) error {
	if len(addProps) > 0 {
		nGQL.WriteString(" ADD ")
		_, _, err := buildProps(addProps, nGQL)
		if err != nil {
			return err
		}
	}

	if len(dropProps) > 0 {
		if len(addProps) > 0 {
			nGQL.WriteByte(',')
		}
		nGQL.WriteString(" DROP (")
		for i, prop := range dropProps {
			nGQL.WriteString(prop)
			if i != len(dropProps)-1 {
				nGQL.WriteString(", ")
			}
		}
		nGQL.WriteByte(')')
	}

	if len(changeProps) > 0 {
		if len(addProps) > 0 || len(dropProps) > 0 {
			nGQL.WriteByte(',')
		}
		nGQL.WriteString(" CHANGE ")
		_, _, err := buildProps(changeProps, nGQL)
		if err != nil {
			return err
		}
	}

	if updateTTL && len(ttlCols) == 1 && ttlDuration != "" {
		nGQL.WriteString(" TTL_DURATION = ")
		nGQL.WriteString(ttlDuration)
		nGQL.WriteString(", TTL_COL = ")
		nGQL.WriteString(strconv.Quote(ttlCols[0]))
	}
	return nil
}
