package resolver

import (
	"errors"
	"github.com/haysons/nebulaorm/internal/utils"
	"reflect"
	"sort"
)

// RecordSchema parses the record structure provided by the business layer for subsequent assignment of the Record \
// object returned by the nebula graph to the business layer
type RecordSchema struct {
	Name          string
	colFieldIndex map[string][]int
}

func ParseRecord(destType reflect.Type) (*RecordSchema, error) {
	if destType.Kind() == reflect.Ptr {
		destType = destType.Elem()
	}
	if destType.Kind() != reflect.Struct {
		return nil, errors.New("nebulaorm: parse record schema failed, dest should be a struct or struct pointer")
	}
	record := &RecordSchema{
		Name:          destType.Name(),
		colFieldIndex: make(map[string][]int),
	}
	for _, structField := range getDestFields(destType) {
		colName := getColName(structField)
		if _, ok := record.colFieldIndex[colName]; !ok {
			record.colFieldIndex[colName] = structField.Index
		}
	}
	return record, nil
}

// GetFieldIndexByColName get the index position of a field
func (r *RecordSchema) GetFieldIndexByColName(colName string) []int {
	return r.colFieldIndex[colName]
}

func getColName(field reflect.StructField) string {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	colName := setting[TagSettingColName]
	if colName == "" {
		colName = camelCaseToUnderscore(field.Name)
	}
	return colName
}

func getDestFields(destType reflect.Type) []reflect.StructField {
	fields := make([]reflect.StructField, 0)
	for _, field := range utils.StructFields(destType) {
		if !field.IsExported() || FieldIgnore(field) {
			continue
		}
		fields = append(fields, field)
	}
	sort.Slice(fields, func(i, j int) bool {
		if len(fields[i].Index) != len(fields[j].Index) {
			return len(fields[i].Index) < len(fields[j].Index)
		}
		for k := 0; k < len(fields[i].Index); k++ {
			if fields[i].Index[k] < fields[j].Index[k] {
				return true
			} else if fields[i].Index[k] > fields[j].Index[k] {
				return false
			}
		}
		return true
	})
	return fields
}
