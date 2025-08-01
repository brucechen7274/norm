package resolver

import (
	"errors"
	"fmt"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"reflect"
	"strconv"
)

// EdgeTypeNamer specifies the name of the edge type. a structure that implements this interface will be treated as an edge.
type EdgeTypeNamer interface {
	EdgeTypeName() string
}

type EdgeSchema struct {
	srcVIDType       VIDType
	srcVIDFieldIndex []int
	dstVIDType       VIDType
	dstVIDFieldIndex []int
	edgeTypeName     string
	rankFieldIndex   []int
	props            []*Prop
	propByName       map[string]*Prop
	indexNames       []string
	indexFields      map[string][]*IndexField
}

// ParseEdge parse edge struct
func ParseEdge(destType reflect.Type) (*EdgeSchema, error) {
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	if destType.Kind() != reflect.Struct {
		return nil, errors.New("norm: parse edge failed, dest should be a struct or a struct pointer")
	}
	edge := &EdgeSchema{
		srcVIDFieldIndex: nil,
		dstVIDFieldIndex: nil,
		rankFieldIndex:   nil,
		props:            make([]*Prop, 0),
		propByName:       make(map[string]*Prop),
	}
	// whether it implements the EdgeTypeNamer interface
	destValue := reflect.New(destType).Interface()
	edgeTypeNamer, ok := destValue.(EdgeTypeNamer)
	if !ok {
		return nil, errors.New("norm: parse edge failed, need to implement interface resolver.EdgeTypeNamer")
	}
	edge.edgeTypeName = edgeTypeNamer.EdgeTypeName()
	for _, field := range getDestFields(destType) {
		setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
		if _, isSrcID := setting[TagSettingEdgeSrcID]; isSrcID {
			if edge.srcVIDFieldIndex == nil {
				switch field.Type.Kind() {
				case reflect.String:
					edge.srcVIDType = VIDTypeString
				case reflect.Int64:
					edge.srcVIDType = VIDTypeInt64
				default:
					return nil, errors.New("norm: parse edge failed, src_id field should be a string or int64")
				}
				edge.srcVIDFieldIndex = field.Index
			}
			continue
		}
		if _, isDstID := setting[TagSettingEdgeDstID]; isDstID {
			if edge.dstVIDFieldIndex == nil {
				switch field.Type.Kind() {
				case reflect.String:
					edge.dstVIDType = VIDTypeString
				case reflect.Int64:
					edge.dstVIDType = VIDTypeInt64
				default:
					return nil, errors.New("norm: parse edge failed, dst_id field should be a string or int64")
				}
				edge.dstVIDFieldIndex = field.Index
			}
			continue
		}
		if _, isRank := setting[TagSettingEdgeRank]; isRank {
			if edge.rankFieldIndex == nil {
				if !(field.Type.Kind() == reflect.Int64 || field.Type.Kind() == reflect.Int || field.Type.Kind() == reflect.Int32 || field.Type.Kind() == reflect.Int8 || field.Type.Kind() == reflect.Int16) {
					return nil, errors.New("norm: parse edge failed, rank field should be int")
				}
				edge.rankFieldIndex = field.Index
			}
			continue
		}
		// parsing Edge Properties
		propName := GetPropName(field)
		sdkType := GetValueSdkType(field)
		dataType := GetFieldDataType(field)
		notNull := IsFieldNotNull(field)
		propDefault := GetFieldDefault(field)
		comment := GetFieldComment(field)
		ttl := GetFieldTTL(field)
		index := GetFieldIndex(field, edge.edgeTypeName, propName, dataType)
		prop := &Prop{
			Name:        propName,
			StructField: field,
			Type:        field.Type,
			SdkType:     sdkType,
			DataType:    dataType,
			NotNull:     notNull,
			Default:     propDefault,
			Comment:     comment,
			TTL:         ttl,
		}
		if _, ok = edge.propByName[propName]; ok {
			continue
		}
		edge.SetProps(prop)
		edge.SetIndexFields(index)
	}
	if edge.srcVIDFieldIndex == nil || edge.dstVIDFieldIndex == nil {
		return nil, errors.New("norm: parse edge failed, edge must contains src_id field and dst_id field")
	}
	return edge, nil
}

// GetTypeName get edge type name
func (e *EdgeSchema) GetTypeName() string {
	return e.edgeTypeName
}

// SetTypeName set edge type name
func (e *EdgeSchema) SetTypeName(edgeTypeName string) {
	e.edgeTypeName = edgeTypeName
}

// GetSrcVID get the src_id of the edge
func (e *EdgeSchema) GetSrcVID(edgeValue reflect.Value) any {
	if e.srcVIDFieldIndex != nil {
		edgeValue = reflect.Indirect(edgeValue)
		return edgeValue.FieldByIndex(e.srcVIDFieldIndex).Interface()
	}
	return nil
}

// GetSrcVIDExpr get the src_id expr of the edge
func (e *EdgeSchema) GetSrcVIDExpr(edgeValue reflect.Value) string {
	srcID := e.GetSrcVID(edgeValue)
	if srcID == nil {
		return ""
	}
	switch e.srcVIDType {
	case VIDTypeString:
		return strconv.Quote(srcID.(string))
	case VIDTypeInt64:
		return strconv.FormatInt(srcID.(int64), 10)
	}
	return ""
}

// GetDstVID get the dst_id of the edge
func (e *EdgeSchema) GetDstVID(edgeValue reflect.Value) any {
	if e.dstVIDFieldIndex != nil {
		edgeValue = reflect.Indirect(edgeValue)
		return edgeValue.FieldByIndex(e.dstVIDFieldIndex).Interface()
	}
	return nil
}

// GetDstVIDExpr get the dst_id expr of the edge
func (e *EdgeSchema) GetDstVIDExpr(edgeValue reflect.Value) string {
	dstID := e.GetDstVID(edgeValue)
	if dstID == nil {
		return ""
	}
	switch e.dstVIDType {
	case VIDTypeString:
		return strconv.Quote(dstID.(string))
	case VIDTypeInt64:
		return strconv.FormatInt(dstID.(int64), 10)
	}
	return ""
}

// GetRank get the rank value of the edge
func (e *EdgeSchema) GetRank(edgeValue reflect.Value) int64 {
	if e.rankFieldIndex != nil {
		edgeValue = reflect.Indirect(edgeValue)
		return edgeValue.FieldByIndex(e.rankFieldIndex).Int()
	}
	return 0
}

// GetProps get a list of attributes for the current edge
func (e *EdgeSchema) GetProps() []*Prop {
	return e.props
}

// SetProps set attributes of the edge
func (e *EdgeSchema) SetProps(props ...*Prop) {
	if e.propByName == nil {
		e.propByName = make(map[string]*Prop)
	}
	for _, prop := range props {
		e.props = append(e.props, prop)
		e.propByName[prop.Name] = prop
	}
}

func (e *EdgeSchema) GetIndexes() []*Index {
	indexes := make([]*Index, 0, len(e.indexNames))
	for _, indexName := range e.indexNames {
		indexes = append(indexes, &Index{
			Name:   indexName,
			Type:   IndexTypeEdge,
			Target: e.GetTypeName(),
			Fields: e.indexFields[indexName],
		})
	}
	return indexes
}

func (e *EdgeSchema) SetIndexFields(fields ...*IndexField) {
	if e.indexFields == nil {
		e.indexFields = make(map[string][]*IndexField)
	}
	for _, field := range fields {
		if field == nil {
			continue
		}
		_, ok := e.indexFields[field.Name]
		if !ok {
			e.indexNames = append(e.indexNames, field.Name)
		}
		e.indexFields[field.Name] = append(e.indexFields[field.Name], field)
	}
}

// Scan assign a value to a target struct
func (e *EdgeSchema) Scan(rl *nebula.Relationship, destValue reflect.Value) error {
	destValue = reflect.Indirect(destValue)
	if !destValue.CanSet() {
		return fmt.Errorf("norm: edge schema scan dest value failed, %w", ErrValueCannotSet)
	}
	if e.srcVIDFieldIndex != nil {
		srcID := rl.GetSrcVertexID()
		if err := ScanSimpleValue(&srcID, destValue.FieldByIndex(e.srcVIDFieldIndex)); err != nil {
			return err
		}
	}
	if e.dstVIDFieldIndex != nil {
		dstID := rl.GetDstVertexID()
		if err := ScanSimpleValue(&dstID, destValue.FieldByIndex(e.dstVIDFieldIndex)); err != nil {
			return err
		}
	}
	if e.rankFieldIndex != nil {
		rank := rl.GetRanking()
		destValue.FieldByIndex(e.rankFieldIndex).SetInt(rank)
	}
	for propName, propValue := range rl.Properties() {
		eProp, ok := e.propByName[propName]
		if !ok {
			continue
		}
		if err := ScanSimpleValue(propValue, destValue.FieldByIndex(eProp.StructField.Index)); err != nil {
			return err
		}
	}
	return nil
}
