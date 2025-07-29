package resolver

import (
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	TagSettingKey       = "norm"        // norm struct tag key
	TagSettingColName   = "col"         // name of the field in the record
	TagSettingVertexID  = "vertex_id"   // marks the field as a vertex ID
	TagSettingEdgeSrcID = "edge_src_id" // marks the field as an edge source ID
	TagSettingEdgeDstID = "edge_dst_id" // marks the field as an edge destination ID
	TagSettingEdgeRank  = "edge_rank"   // marks the field as an edge rank
	TagSettingPropName  = "prop"        // property name for a vertex or edge
	TagSettingDataType  = "type"        // specifies the data type (see: https://docs.nebula-graph.com.cn/3.6.0/3.ngql-guide/3.data-types/1.numeric/)
	TagSettingNotNull   = "not_null"    // declares the field as NOT NULL
	TagSettingDefault   = "default"     // declares a default value for the field
	TagSettingComment   = "comment"     // declares a comment/description for the field
	TagSettingTTL       = "ttl"         // marks the field as TTL (time-to-live) for expiration
	TagSettingIndex     = "index"       // defines index configuration on the field
	TagSettingIgnore    = "-"           // norm will ignore this field
)

func ParseTagSetting(s string) map[string]string {
	m := make(map[string]string)
	tags := strings.Split(s, ";")
	for _, tag := range tags {
		kv := strings.Split(tag, ":")
		k := strings.TrimSpace(strings.ToLower(kv[0]))
		if k == "" {
			continue
		}
		if len(kv) >= 2 {
			m[k] = strings.Join(kv[1:], ":")
		} else {
			m[k] = k
		}
	}
	return m
}

func GetPropName(field reflect.StructField) string {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	propName := setting[TagSettingPropName]
	if propName == "" {
		propName = camelCaseToUnderscore(field.Name)
	}
	return propName
}

func GetValueSdkType(field reflect.StructField) string {
	dataTypeRaw := GetFieldDataType(field)
	dataTypeLower := strings.ToLower(dataTypeRaw)
	switch dataTypeLower {
	case "int", "int64", "int32", "int16", "int8":
		return NebulaSdkTypeInt
	case "float", "double":
		return NebulaSdkTypeFloat
	case "bool":
		return NebulaSdkTypeBool
	case "string":
		return NebulaSdkTypeString
	case "date":
		return NebulaSdkTypeDate
	case "time":
		return NebulaSdkTypeTime
	case "datetime", "timestamp":
		return NebulaSdkTypeDatetime
	default:
		if strings.HasPrefix(dataTypeLower, "fixed_string") {
			return NebulaSdkTypeString
		}
		return dataTypeRaw
	}
}

func GetFieldDataType(field reflect.StructField) string {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	dataType := setting[TagSettingDataType]
	if dataType != "" {
		return dataType
	}
	fieldType := field.Type
	switch fieldType.Kind() {
	case reflect.Bool:
		return "bool"
	case reflect.Int, reflect.Uint:
		return "int"
	case reflect.Int64, reflect.Uint64:
		return "int64"
	case reflect.Int32, reflect.Uint32:
		return "int32"
	case reflect.Int16, reflect.Uint16:
		return "int16"
	case reflect.Int8, reflect.Uint8:
		return "int8"
	case reflect.Float32:
		return "float"
	case reflect.Float64:
		return "double"
	case reflect.String:
		return "string"
	case reflect.Struct:
		if fieldType.PkgPath() == "time" && fieldType.Name() == "Time" {
			return "datetime"
		}
	default:
		return ""
	}
	return ""
}

func IsFieldNotNull(field reflect.StructField) bool {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	notNull := setting[TagSettingNotNull]
	return notNull != ""
}

func GetFieldDefault(field reflect.StructField) string {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	return setting[TagSettingDefault]
}

func GetFieldComment(field reflect.StructField) string {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	return setting[TagSettingComment]
}

func GetFieldTTL(field reflect.StructField) string {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	return setting[TagSettingTTL]
}

type FieldIndex struct {
	Name     string
	Prop     string
	DataType string
	Length   int
	Priority int
}

// GetFieldIndex retrieves the index configuration from the given struct field tag.
func GetFieldIndex(field reflect.StructField, propName, dataType string) *FieldIndex {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	indexSetting, ok := setting[TagSettingIndex]
	if !ok {
		return nil
	}

	settingArr := strings.Split(indexSetting, ",")
	indexName := settingArr[0]
	if indexName == TagSettingIndex || indexName == "" {
		indexName = "idx_" + propName
	}
	fieldIndex := &FieldIndex{
		Name:     indexName,
		Prop:     propName,
		DataType: dataType,
		Length:   0,
		Priority: 10,
	}
	if len(settingArr) > 1 {
		for k, v := range ParseTagSetting(strings.Join(settingArr[1:], ";")) {
			if k == "priority" {
				priority, err := strconv.Atoi(v)
				if err == nil {
					fieldIndex.Priority = priority
				}
			}
			if k == "length" {
				length, err := strconv.Atoi(v)
				if err == nil {
					fieldIndex.Length = length
				}
			}
		}
	}
	return fieldIndex
}

func FieldIgnore(field reflect.StructField) bool {
	setting := ParseTagSetting(field.Tag.Get(TagSettingKey))
	return setting[TagSettingIgnore] != ""
}

func camelCaseToUnderscore(s string) string {
	var output []rune
	for i, r := range s {
		if i == 0 {
			output = append(output, unicode.ToLower(r))
			continue
		}
		if unicode.IsUpper(r) {
			output = append(output, '_')
		}
		output = append(output, unicode.ToLower(r))
	}
	return string(output)
}

var timezoneDefault = time.Local

func SetTimezone(loc *time.Location) {
	if loc != nil {
		timezoneDefault = loc
	}
}
