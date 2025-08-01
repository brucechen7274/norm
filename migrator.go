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
//
// In addition, if index definitions are declared in the vertex struct and the corresponding
// indexes do not exist in the current graph space, they will be created.
//
// Note: Index creation is asynchronous in NebulaGraph, so newly created indexes cannot
// be immediately rebuilt within this method. It is recommended to call RebuildVertexTagIndexes
// manually after migration to ensure indexes are properly built.
//
// For safety reasons, existing indexes will not be dropped.
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
		if err = m.autoCreateTagIndexes(tag); err != nil {
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
	propsExist := make(map[string]*PropDesc)
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

func (m *Migrator) autoCreateTagIndexes(tag *resolver.VertexTag) error {
	indexes := tag.GetIndexes()
	for _, index := range indexes {
		hasIndex, err := m.HasVertexTagIndex(index.Name)
		if err != nil {
			return err
		}
		if hasIndex {
			continue
		}
		if err = m.CreateVertexTagsIndex(index, true); err != nil {
			return err
		}
	}
	return nil
}

// AutoMigrateEdges automatically migrates edge schemas.
// If the specified edge does not exist in the current graph space, it will be created.
// If it already exists, each property will be compared to determine whether updates are needed.
// For safety, this method will not delete existing edges or edge properties.
//
// In addition, if index definitions are declared in the edge struct and the corresponding
// indexes do not exist in the current graph space, they will be created.
//
// Note: Index creation is asynchronous in NebulaGraph, so newly created indexes cannot
// be immediately rebuilt within this method. It is recommended to call RebuildEdgeIndexes
// manually after migration to ensure indexes are properly built.
//
// For safety reasons, existing edges or their properties and indexes will not be dropped.
func (m *Migrator) AutoMigrateEdges(edges ...any) error {
	for _, edge := range edges {
		edgeSchema, err := resolver.ParseEdge(reflect.TypeOf(edge))
		if err != nil {
			return err
		}
		hasEdge, err := m.HasEdge(edgeSchema.GetTypeName())
		if err != nil {
			return err
		}
		if !hasEdge {
			// Create the edge if it doesn't exist
			err = m.CreateEdge(edgeSchema, true)
		} else {
			// Otherwise, apply ALTER operations to update the edge
			err = m.autoAlterEdge(edgeSchema)
		}
		if err != nil {
			return err
		}
		if err = m.autoCreateEdgeIndexes(edgeSchema); err != nil {
			return err
		}
	}
	return nil
}

// autoAlterEdge automatically applies property changes to an existing edge.
// It compares the current edge schema with the existing one in the space,
// and constructs an ALTER operation for any new or changed properties.
func (m *Migrator) autoAlterEdge(edge *resolver.EdgeSchema) error {
	edgeProps, err := m.DescEdge(edge.GetTypeName())
	if err != nil {
		return err
	}
	propsExist := make(map[string]*PropDesc)
	for _, edgeProp := range edgeProps {
		propsExist[edgeProp.Field] = edgeProp
	}
	alterOp := clause.AlterOperate{}
	for _, propNew := range edge.GetProps() {
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
	return m.AlterEdge(edge, alterOp)
}

func (m *Migrator) autoCreateEdgeIndexes(edge *resolver.EdgeSchema) error {
	indexes := edge.GetIndexes()
	for _, index := range indexes {
		hasIndex, err := m.HasEdgeIndex(index.Name)
		if err != nil {
			return err
		}
		if hasIndex {
			continue
		}
		if err = m.CreateEdgeIndex(index, true); err != nil {
			return err
		}
	}
	return nil
}

// isPropChanged determines whether a property definition has changed
// by comparing type, nullability, and default value.
func (m *Migrator) isPropChanged(propExist *PropDesc, propNew *resolver.Prop) bool {
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

type PropDesc struct {
	Field   string `norm:"col:Field"`
	Type    string `norm:"col:Type"`
	Null    string `norm:"col:Null"`
	Default string `norm:"col:Default"`
	Comment string `norm:"col:Comment"`
}

// DescVertexTag retrieves detailed information about a vertex tag,
// including field names, data types, and other metadata.
func (m *Migrator) DescVertexTag(tagName string) ([]*PropDesc, error) {
	tagProps := make([]*PropDesc, 0)
	err := m.db.Raw("DESCRIBE TAG " + tagName).
		Find(&tagProps)
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

// HasEdge checks whether the specified edge exists in the current graph space.
func (m *Migrator) HasEdge(edgeTypeName string) (bool, error) {
	edges := make([]string, 0)
	err := m.db.Raw("SHOW EDGES").
		FindCol("Name", &edges)
	if err != nil {
		return false, err
	}
	for _, edge := range edges {
		if edge == edgeTypeName {
			return true, nil
		}
	}
	return false, nil
}

// DescEdge returns the detailed property description of the specified edge.
func (m *Migrator) DescEdge(edgeTypeName string) ([]*PropDesc, error) {
	edgeProps := make([]*PropDesc, 0)
	err := m.db.Raw("DESCRIBE EDGE " + edgeTypeName).
		Find(&edgeProps)
	if err != nil {
		return nil, err
	}
	return edgeProps, nil
}

// CreateEdge creates an edge schema in the space.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) CreateEdge(edge any, ifNotExists ...bool) error {
	tx := m.db.getInstance()
	tx.Statement.CreateEdge(edge, ifNotExists...)
	return tx.Exec()
}

// DropEdge drops an edge schema by its type name.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) DropEdge(edgeTypeName string, ifExists ...bool) error {
	tx := m.db.getInstance()
	tx.Statement.DropEdge(edgeTypeName, ifExists...)
	return tx.Exec()
}

// AlterEdge alters the definition of an existing edge type.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) AlterEdge(edge any, op clause.AlterOperate) error {
	tx := m.db.getInstance()
	tx.Statement.AlterEdge(edge, op)
	return tx.Exec()
}

// HasVertexTagIndex checks whether a tag index with the given name exists.
func (m *Migrator) HasVertexTagIndex(indexName string) (bool, error) {
	indexNames := make([]string, 0)
	err := m.db.Raw("SHOW TAG INDEXES").
		FindCol("Index Name", &indexNames)
	if err != nil {
		return false, err
	}
	for _, name := range indexNames {
		if name == indexName {
			return true, nil
		}
	}
	return false, nil
}

// CreateVertexTagsIndex creates indexes for all tags of the given vertex struct.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) CreateVertexTagsIndex(vertex any, ifNotExists ...bool) error {
	tx := m.db.getInstance()
	tx.Statement.CreateVertexTagsIndex(vertex, ifNotExists...)
	return tx.Exec()
}

// RebuildVertexTagIndexes rebuilds the specified vertex tag indexes.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) RebuildVertexTagIndexes(indexNames ...string) error {
	tx := m.db.getInstance()
	tx.Statement.RebuildVertexTagIndexes(indexNames...)
	return tx.Exec()
}

// DropVertexTagIndex drops a vertex tag index by its name.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) DropVertexTagIndex(indexName string, ifExists ...bool) error {
	tx := m.db.getInstance()
	tx.Statement.DropVertexTagIndex(indexName, ifExists...)
	return tx.Exec()
}

// HasEdgeIndex checks whether an edge index with the given name exists.
func (m *Migrator) HasEdgeIndex(indexName string) (bool, error) {
	indexNames := make([]string, 0)
	err := m.db.Raw("SHOW EDGE INDEXES").
		FindCol("Index Name", &indexNames)
	if err != nil {
		return false, err
	}
	for _, name := range indexNames {
		if name == indexName {
			return true, nil
		}
	}
	return false, nil
}

// CreateEdgeIndex creates an index on the specified edge's properties.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) CreateEdgeIndex(edge any, ifNotExists ...bool) error {
	tx := m.db.getInstance()
	tx.Statement.CreateEdgeIndex(edge, ifNotExists...)
	return tx.Exec()
}

// RebuildEdgeIndexes rebuilds one or more edge indexes by their names.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) RebuildEdgeIndexes(indexNames ...string) error {
	tx := m.db.getInstance()
	tx.Statement.RebuildEdgeIndexes(indexNames...)
	return tx.Exec()
}

// DropEdgeIndex drops an edge index by name.
// see more information on the method of the same name in statement.Statement
func (m *Migrator) DropEdgeIndex(indexName string, ifExists ...bool) error {
	tx := m.db.getInstance()
	tx.Statement.DropEdgeIndex(indexName, ifExists...)
	return tx.Exec()
}
