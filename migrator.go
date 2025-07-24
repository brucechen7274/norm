package norm

import (
	"github.com/haysons/norm/clause"
	"github.com/haysons/norm/resolver"
	"reflect"
	"strings"
)

type Migrator struct {
	db *DB
}

// Migrator creates a new Migrator instance based on the current DB object
func (db *DB) Migrator() *Migrator {
	return &Migrator{db: db}
}

// NewMigrator creates a new Migrator instance based on the specified DB object
func NewMigrator(db *DB) *Migrator {
	return &Migrator{db: db}
}

// AutoMigrateVertexes automatically migrates all tags associated with the given vertices
// in the current graph space. If a tag does not exist, it will be created.
// If the tag exists, each property will be checked for changes.
// New properties or changed types, null constraints, or default values will trigger ALTER operations.
// For safety, this method does not delete any existing tags or their properties.
func (m *Migrator) AutoMigrateVertexes(vertexes ...any) error {
	for _, vertex := range vertexes {
		vertexSchema, err := resolver.ParseVertex(reflect.TypeOf(vertex))
		if err != nil {
			return err
		}
		if err = m.autoMigrateVertex(vertexSchema); err != nil {
			return err
		}
	}
	return nil
}

// autoMigrateVertex handles the tag migration process for a single vertex.
// It creates new tags if they do not exist, or alters existing ones.
func (m *Migrator) autoMigrateVertex(vertex *resolver.VertexSchema) error {
	for _, tag := range vertex.GetTags() {
		hasTag, err := m.HasVertexTag(tag.TagName)
		if err != nil {
			return err
		}
		if !hasTag {
			// Create the tag if it doesn't exist
			err = m.CreateVertexTags(tag, true)
		} else {
			// Otherwise, apply ALTER operations to update the tag
			err = m.autoAlterVertexTags(tag)
		}
		if err != nil {
			return err
		}
	}
	return nil
}

// autoAlterVertexTags compares the current tag schema with the one in the database,
// and generates ALTER operations for any property additions or changes.
func (m *Migrator) autoAlterVertexTags(tag *resolver.VertexTag) error {
	tagProps, err := m.DescVertexTag(tag.TagName)
	if err != nil {
		return err
	}
	propsExist := make(map[string]*VertexTagPropDesc)
	for _, tagProp := range tagProps {
		propsExist[tagProp.Field] = tagProp
	}
	alterOp := clause.AlterOperate{}
	for _, propNew := range tag.GetProps() {
		// Add the property if it does not exist
		propExist, ok := propsExist[propNew.Name]
		if !ok {
			alterOp.AddProps = append(alterOp.AddProps, propNew.Name)
			continue
		}
		// Apply change if the property differs from existing definition
		if m.isPropChanged(propExist, propNew) {
			alterOp.ChangeProps = append(alterOp.ChangeProps, propNew.Name)
		}
	}
	if len(alterOp.AddProps) == 0 && len(alterOp.ChangeProps) == 0 {
		return nil
	}
	return m.AlterVertexTag(tag, alterOp)
}

// isPropChanged determines whether a property definition has changed
// by comparing type, nullability, and default value.
func (m *Migrator) isPropChanged(propExist *VertexTagPropDesc, propNew *resolver.Prop) bool {
	propType := func(t string) string {
		t = strings.ToLower(t)
		// "int" is treated as an alias for "int64"
		if t == "int" {
			t = "int64"
		}
		return t
	}

	notNull := func(s string) bool {
		if strings.ToLower(s) == "yes" {
			return false
		}
		return true
	}

	defaultValue := func(s string) string {
		if s == "_EMPTY_" {
			return ""
		}
		if s == "" {
			return "''"
		}
		return s
	}

	if propType(propNew.DataType) != propType(propExist.Type) ||
		propNew.NotNull != notNull(propExist.Null) ||
		propNew.Default != defaultValue(propExist.Default) {
		return true
	}

	return false
}

// HasVertexTag checks whether a given tag exists in the current graph space.
// Returns true if the tag exists, false otherwise.
func (m *Migrator) HasVertexTag(tagName string) (bool, error) {
	tags := make([]string, 0)
	err := m.db.Raw("SHOW TAGS").
		FindCol("Name", &tags)
	if err != nil {
		return false, err
	}
	for _, tag := range tags {
		if tag == tagName {
			return true, nil
		}
	}
	return false, nil
}

type VertexTagPropDesc struct {
	Field   string `norm:"col:Field"`
	Type    string `norm:"col:Type"`
	Null    string `norm:"col:Null"`
	Default string `norm:"col:Default"`
	Comment string `norm:"col:Comment"`
}

// DescVertexTag retrieves detailed information about a vertex tag,
// including field names, data types, and other metadata.
func (m *Migrator) DescVertexTag(tagName string) ([]*VertexTagPropDesc, error) {
	tagProps := make([]*VertexTagPropDesc, 0)
	err := m.db.Raw("DESCRIBE TAG " + tagName).
		Take(&tagProps)
	if err != nil {
		return nil, err
	}
	return tagProps, nil
}

// CreateVertexTags creates all tags associated with a vertex.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) CreateVertexTags(vertex any, ifNotExists ...bool) error {
	tx := m.db.getInstance()
	tx.Statement.CreateVertexTags(vertex, ifNotExists...)
	return tx.Exec()
}

// DropVertexTag drops a vertex tag by its name.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) DropVertexTag(tagName string, ifExists ...bool) error {
	tx := m.db.getInstance()
	tx.Statement.DropVertexTag(tagName, ifExists...)
	return tx.Exec()
}

// AlterVertexTag modifies a tag of the given vertex using the specified operation.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) AlterVertexTag(vertex any, op clause.AlterOperate, opts ...clause.Option) error {
	tx := m.db.getInstance()
	tx.Statement.AlterVertexTag(vertex, op, opts...)
	return tx.Exec()
}
