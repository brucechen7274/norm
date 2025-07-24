package resolver

import (
	"errors"
	"fmt"
	"github.com/haysons/norm/internal/utils"
	nebula "github.com/vesoft-inc/nebula-go/v3"
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	NebulaSdkTypeNull     = "null"
	NebulaSdkTypeBool     = "bool"
	NebulaSdkTypeInt      = "int"
	NebulaSdkTypeFloat    = "float"
	NebulaSdkTypeString   = "string"
	NebulaSdkTypeDate     = "date"
	NebulaSdkTypeTime     = "time"
	NebulaSdkTypeDatetime = "datetime"
	NebulaSdkTypeVertex   = "vertex"
	NebulaSdkTypeEdge     = "edge"
	NebulaSdkTypeList     = "list"
	NebulaSdkTypeMap      = "map"
	NebulaSdkTypeSet      = "set"
	NebulaSdkTypeEmpty    = "empty"
)

var (
	ErrValueCannotSet = errors.New("reflect value can not be set")
)

// Resolver responsible for parsing and converting data types in nebula graph and defined data types in golang
type Resolver struct {
	vertexSchema map[string]*VertexSchema
	edgeSchema   map[string]*EdgeSchema
	recordSchema map[string]*RecordSchema
}

func NewResolver() *Resolver {
	return &Resolver{
		vertexSchema: make(map[string]*VertexSchema),
		edgeSchema:   make(map[string]*EdgeSchema),
		recordSchema: make(map[string]*RecordSchema),
	}
}

// ScanValue scan nebula graph value into dest value.
func (r *Resolver) ScanValue(nebulaValue *nebula.ValueWrapper, destValue reflect.Value) error {
	if !destValue.CanSet() && destValue.Kind() != reflect.Map {
		return fmt.Errorf("norm: scan dest value failed, %w", ErrValueCannotSet)
	}
	switch nebulaValue.GetType() {
	case NebulaSdkTypeVertex:
		vNode, _ := nebulaValue.AsNode()
		destValue = utils.PtrValue(destValue)
		switch destValue.Kind() {
		case reflect.Struct:
			destType := destValue.Type()
			vertexSchema, err := r.getVertexSchema(destType)
			if err != nil {
				return err
			}
			return vertexSchema.Scan(vNode, destValue)
		default:
		}
	case NebulaSdkTypeEdge:
		vRelationShip, _ := nebulaValue.AsRelationship()
		destValue = utils.PtrValue(destValue)
		switch destValue.Kind() {
		case reflect.Struct:
			destType := destValue.Type()
			edgeSchema, err := r.getEdgeSchema(destType)
			if err != nil {
				return err
			}
			return edgeSchema.Scan(vRelationShip, destValue)
		default:
		}
	case NebulaSdkTypeList:
		vList, _ := nebulaValue.AsList()
		destValue = utils.PtrValue(destValue)
		switch destValue.Kind() {
		case reflect.Slice, reflect.Array:
			return utils.SliceSetElem(destValue, len(vList), func(i int, elem reflect.Value) (bool, error) {
				if i >= len(vList) {
					return false, nil
				}
				if err := r.ScanValue(&vList[i], elem); err != nil {
					return false, err
				}
				return true, nil
			})
		default:
		}
	case NebulaSdkTypeMap:
		vMap, _ := nebulaValue.AsMap()
		destValue = utils.PtrValue(destValue)
		switch destValue.Kind() {
		case reflect.Map:
			destType := destValue.Type()
			if destType.Key().Kind() != reflect.String {
				return errors.New("norm: scan dest value failed, map key must be string")
			}
			elemType := destType.Elem()
			if destValue.IsNil() && destValue.CanSet() {
				destValue.Set(reflect.MakeMap(destType))
			}
			for key, value := range vMap {
				elemValue := reflect.New(elemType).Elem()
				if err := r.ScanValue(&value, elemValue); err != nil {
					return err
				}
				destValue.SetMapIndex(reflect.ValueOf(key), elemValue)
			}
			return nil
		default:
		}
	case NebulaSdkTypeSet:
		vSet, _ := nebulaValue.AsDedupList()
		destValue = utils.PtrValue(destValue)
		switch destValue.Kind() {
		case reflect.Slice, reflect.Array:
			return utils.SliceSetElem(destValue, len(vSet), func(i int, elem reflect.Value) (bool, error) {
				if i >= len(vSet) {
					return false, nil
				}
				if err := r.ScanValue(&vSet[i], elem); err != nil {
					return false, err
				}
				return true, nil
			})
		case reflect.Map:
			keyType := destValue.Type().Key()
			valueType := destValue.Type().Elem()
			for _, value := range vSet {
				keyValue := reflect.New(keyType).Elem()
				if err := r.ScanValue(&value, keyValue); err != nil {
					return err
				}
				destValue.SetMapIndex(keyValue, reflect.Zero(valueType))
			}
			return nil
		default:
		}
	default:
		return ScanSimpleValue(nebulaValue, destValue)
	}
	return fmt.Errorf("norm: can not scan nebula type %s into golang type %v", nebulaValue.GetType(), destValue.Type())
}

// ScanRecord scan the record value into a struct.
func (r *Resolver) ScanRecord(record *nebula.Record, colNames []string, destValue reflect.Value) error {
	destValue = reflect.Indirect(destValue)
	if destValue.Kind() != reflect.Struct {
		return errors.New("norm: scan record failed, dest should be a struct or a struct pointer")
	}
	destType := destValue.Type()
	recordSchema, err := r.getRecordSchema(destType)
	if err != nil {
		return err
	}
	for _, colName := range colNames {
		colValue, err := record.GetValueByColName(colName)
		if err != nil {
			return err
		}
		fieldIndex := recordSchema.GetFieldIndexByColName(colName)
		if len(fieldIndex) == 0 {
			continue
		}
		fieldValue := destValue.FieldByIndex(fieldIndex)
		if err = r.ScanValue(colValue, fieldValue); err != nil {
			return err
		}
	}
	return nil
}

func (r *Resolver) getRecordSchema(destType reflect.Type) (*RecordSchema, error) {
	key := r.getSchemaKey(destType)
	if s, ok := r.recordSchema[key]; ok {
		return s, nil
	}
	recordSchema, err := ParseRecord(destType)
	if err != nil {
		return nil, err
	}
	r.recordSchema[key] = recordSchema
	return recordSchema, nil
}

func (r *Resolver) getVertexSchema(destType reflect.Type) (*VertexSchema, error) {
	key := r.getSchemaKey(destType)
	if s, ok := r.vertexSchema[key]; ok {
		return s, nil
	}
	vertexSchema, err := ParseVertex(destType)
	if err != nil {
		return nil, err
	}
	r.vertexSchema[key] = vertexSchema
	return vertexSchema, nil
}

func (r *Resolver) getEdgeSchema(destType reflect.Type) (*EdgeSchema, error) {
	key := r.getSchemaKey(destType)
	if e, ok := r.edgeSchema[key]; ok {
		return e, nil
	}
	edgeSchema, err := ParseEdge(destType)
	if err != nil {
		return nil, err
	}
	r.edgeSchema[key] = edgeSchema
	return edgeSchema, nil
}

func (r *Resolver) getSchemaKey(destType reflect.Type) string {
	return destType.PkgPath() + "." + destType.Name()
}

// ScanSimpleValue assign values to simple data types
func ScanSimpleValue(nebulaValue *nebula.ValueWrapper, destValue reflect.Value) error {
	if !destValue.CanSet() {
		return fmt.Errorf("norm: scan dest value failed, %w", ErrValueCannotSet)
	}
	if nebulaValue.GetType() == NebulaSdkTypeNull {
		destValue.SetZero()
		return nil
	}
	destValue = utils.PtrValue(destValue)
	if destValue.Kind() == reflect.Interface && destValue.NumMethod() == 0 {
		valueIface, err := GetValueIface(nebulaValue)
		if err != nil {
			return err
		}
		destValue.Set(reflect.ValueOf(valueIface))
		return nil
	}
	switch nebulaValue.GetType() {
	case NebulaSdkTypeEmpty:
		switch destValue.Kind() {
		case reflect.String:
			destValue.SetString("_EMPTY_")
			return nil
		default:
		}
	case NebulaSdkTypeBool:
		switch destValue.Kind() {
		case reflect.Bool:
			vBool, _ := nebulaValue.AsBool()
			destValue.SetBool(vBool)
			return nil
		case reflect.String:
			vBool, _ := nebulaValue.AsBool()
			destValue.SetString(strconv.FormatBool(vBool))
			return nil
		default:
		}
	case NebulaSdkTypeInt:
		switch destValue.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vInt, _ := nebulaValue.AsInt()
			destValue.SetInt(vInt)
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vInt, _ := nebulaValue.AsInt()
			destValue.SetUint(uint64(vInt))
			return nil
		case reflect.Float32, reflect.Float64:
			vInt, _ := nebulaValue.AsInt()
			destValue.SetFloat(float64(vInt))
			return nil
		case reflect.String:
			vInt, _ := nebulaValue.AsInt()
			destValue.SetString(strconv.Itoa(int(vInt)))
			return nil
		default:
		}
	case NebulaSdkTypeFloat:
		switch destValue.Kind() {
		case reflect.Float32, reflect.Float64:
			vFloat, _ := nebulaValue.AsFloat()
			destValue.SetFloat(vFloat)
			return nil
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			vFloat, _ := nebulaValue.AsFloat()
			destValue.SetInt(int64(vFloat))
			return nil
		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			vFloat, _ := nebulaValue.AsFloat()
			destValue.SetUint(uint64(vFloat))
			return nil
		case reflect.String:
			vFloat, _ := nebulaValue.AsFloat()
			destValue.SetString(strconv.FormatFloat(vFloat, 'g', -1, 64))
			return nil
		default:
		}
	case NebulaSdkTypeString:
		switch destValue.Kind() {
		case reflect.String:
			vString, _ := nebulaValue.AsString()
			destValue.SetString(vString)
			return nil
		default:
		}
	case NebulaSdkTypeDate:
		switch destValue.Kind() {
		case reflect.String:
			vDate, _ := nebulaValue.AsDate()
			dateUTC := time.Date(int(vDate.GetYear()), time.Month(vDate.GetMonth()), int(vDate.GetDay()), 0, 0, 0, 0, time.UTC)
			dateObj := dateUTC.In(timezoneDefault)
			destValue.SetString(dateObj.Format("2006-01-02"))
			return nil
		case reflect.Struct:
			destType := destValue.Type()
			if destType.PkgPath() == "time" && destType.Name() == "Time" {
				vDate, _ := nebulaValue.AsDate()
				dateUTC := time.Date(int(vDate.GetYear()), time.Month(vDate.GetMonth()), int(vDate.GetDay()), 0, 0, 0, 0, time.UTC)
				dateObj := dateUTC.In(timezoneDefault)
				destValue.Set(reflect.ValueOf(dateObj))
				return nil
			}
		default:
		}
	case NebulaSdkTypeTime:
		vTimeW, _ := nebulaValue.AsTime()
		vTime, _ := vTimeW.GetLocalTimeWithTimezoneName(timezoneDefault.String())
		switch destValue.Kind() {
		case reflect.String:
			dateObj := time.Date(2020, 1, 1, int(vTime.GetHour()), int(vTime.GetMinute()), int(vTime.GetSec()), int(vTime.GetMicrosec()*1000), timezoneDefault)
			destValue.SetString(dateObj.Format("15:04:05.000000"))
			return nil
		default:
		}
	case NebulaSdkTypeDatetime:
		vDateTimeW, _ := nebulaValue.AsDateTime()
		vDateTime, _ := vDateTimeW.GetLocalDateTimeWithTimezoneName(timezoneDefault.String())
		switch destValue.Kind() {
		case reflect.String:
			dateObj := time.Date(int(vDateTime.GetYear()), time.Month(vDateTime.GetMonth()), int(vDateTime.GetDay()), int(vDateTime.GetHour()), int(vDateTime.GetMinute()), int(vDateTime.GetSec()), int(vDateTime.GetMicrosec()*1000), timezoneDefault)
			destValue.SetString(dateObj.Format("2006-01-02T15:04:05.000000"))
			return nil
		case reflect.Struct:
			destType := destValue.Type()
			if destType.PkgPath() == "time" && destType.Name() == "Time" {
				dateObj := time.Date(int(vDateTime.GetYear()), time.Month(vDateTime.GetMonth()), int(vDateTime.GetDay()), int(vDateTime.GetHour()), int(vDateTime.GetMinute()), int(vDateTime.GetSec()), int(vDateTime.GetMicrosec()*1000), timezoneDefault)
				destValue.Set(reflect.ValueOf(dateObj))
				return nil
			}
		default:
		}
	}
	return fmt.Errorf("norm: can not set value, nebula type %s into golang type %v", nebulaValue.GetType(), destValue.Type())
}

// FormatSimpleValue format variable values to nebula graph data format
func FormatSimpleValue(sdkType string, value reflect.Value) (string, error) {
	switch value.Kind() {
	case reflect.Bool:
		switch sdkType {
		case NebulaSdkTypeBool, "":
			if value.Bool() {
				return "true", nil
			} else {
				return "false", nil
			}
		}
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		switch sdkType {
		case NebulaSdkTypeInt, NebulaSdkTypeFloat, "":
			return strconv.FormatInt(value.Int(), 10), nil
		}
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		switch sdkType {
		case NebulaSdkTypeInt, NebulaSdkTypeFloat, "":
			return strconv.FormatUint(value.Uint(), 10), nil
		}
	case reflect.Float32:
		switch sdkType {
		case NebulaSdkTypeFloat, "":
			return strconv.FormatFloat(value.Float(), 'g', -1, 32), nil
		case NebulaSdkTypeInt:
			valueF := value.Float()
			return strconv.Itoa(int(valueF)), nil
		}
	case reflect.Float64:
		switch sdkType {
		case NebulaSdkTypeFloat, "":
			return strconv.FormatFloat(value.Float(), 'g', -1, 64), nil
		case NebulaSdkTypeInt:
			valueF := value.Float()
			return strconv.Itoa(int(valueF)), nil
		}
	case reflect.String:
		switch sdkType {
		case NebulaSdkTypeString, "":
			str := value.String()
			str = strconv.Quote(str)
			return str, nil
		case NebulaSdkTypeDatetime:
			datetimeStr := `datetime("` + value.String() + `")`
			return datetimeStr, nil
		case NebulaSdkTypeDate:
			dateStr := `date("` + value.String() + `")`
			return dateStr, nil
		case NebulaSdkTypeTime:
			timeStr := `time("` + value.String() + `")`
			return timeStr, nil
		}
	case reflect.Struct:
		switch sdkType {
		case NebulaSdkTypeDatetime, "":
			t, ok := value.Interface().(time.Time)
			if ok {
				if t.Nanosecond() == 0 {
					return t.Format(`datetime("2006-01-02T15:04:05")`), nil
				} else {
					return t.Format(`datetime("2006-01-02T15:04:05.000000")`), nil
				}
			}
		case NebulaSdkTypeDate:
			t, ok := value.Interface().(time.Time)
			if ok {
				return t.Format(`date("2006-01-02")`), nil
			}
		case NebulaSdkTypeTime:
			t, ok := value.Interface().(time.Time)
			if ok {
				return t.Format(`time("15:04:05.000000")`), nil
			}
		}
	case reflect.Slice, reflect.Array:
		switch sdkType {
		case NebulaSdkTypeList, "":
			listStr := strings.Builder{}
			listStr.WriteString("[")
			for i := 0; i < value.Len(); i++ {
				elemStr, err := FormatSimpleValue("", value.Index(i))
				if err != nil {
					return "", err
				}
				listStr.WriteString(elemStr)
				if i < value.Len()-1 {
					listStr.WriteString(", ")
				}
			}
			listStr.WriteString("]")
			return listStr.String(), nil
		case NebulaSdkTypeSet:
			setStr := strings.Builder{}
			setStr.WriteString("set{")
			for i := 0; i < value.Len(); i++ {
				elemStr, err := FormatSimpleValue("", value.Index(i))
				if err != nil {
					return "", err
				}
				setStr.WriteString(elemStr)
				if i < value.Len()-1 {
					setStr.WriteString(", ")
				}
			}
			setStr.WriteString("}")
			return setStr.String(), nil
		}
	case reflect.Map:
		switch sdkType {
		case NebulaSdkTypeMap, "":
			mapStr := strings.Builder{}
			mapStr.WriteString("map{")
			mapLen := value.Len()
			mapIter := value.MapRange()
			var i int
			for mapIter.Next() {
				i++
				k := mapIter.Key()
				if k.Kind() != reflect.String {
					return "", fmt.Errorf("norm: format value failed, can not convert map key to string")
				}
				v := mapIter.Value()
				mapStr.WriteString(k.String())
				mapStr.WriteString(": ")
				vStr, err := FormatSimpleValue("", v)
				if err != nil {
					return "", err
				}
				mapStr.WriteString(vStr)
				if i < mapLen {
					mapStr.WriteString(", ")
				}
			}
			mapStr.WriteString("}")
			return mapStr.String(), nil
		case NebulaSdkTypeSet:
			setStr := strings.Builder{}
			setStr.WriteString("set{")
			mapLen := value.Len()
			mapIter := value.MapRange()
			var i int
			for mapIter.Next() {
				i++
				kStr, err := FormatSimpleValue("", mapIter.Key())
				if err != nil {
					return "", err
				}
				setStr.WriteString(kStr)
				if i < mapLen {
					setStr.WriteString(", ")
				}
			}
			setStr.WriteString("}")
			return setStr.String(), nil
		}
	case reflect.Ptr, reflect.Interface:
		if !value.IsNil() {
			return FormatSimpleValue("", value.Elem())
		} else {
			switch sdkType {
			case NebulaSdkTypeNull, "":
				return "NULL", nil
			case NebulaSdkTypeEmpty:
				return "_EMPTY_", nil
			}
		}
	case reflect.Invalid:
		return "", errors.New("norm: format value failed, invalid type, eg: the undefined type nil")
	default:
	}
	return "", fmt.Errorf("norm: format value failed, golang type: %s, nebula type: %s", value.Type(), sdkType)
}

// GetValueIface get the nebula graph return value
func GetValueIface(nebulaValue *nebula.ValueWrapper) (any, error) {
	switch nebulaValue.GetType() {
	case NebulaSdkTypeNull, NebulaSdkTypeEmpty:
		return nil, nil
	case NebulaSdkTypeBool:
		return nebulaValue.AsBool()
	case NebulaSdkTypeInt:
		return nebulaValue.AsInt()
	case NebulaSdkTypeFloat:
		return nebulaValue.AsFloat()
	case NebulaSdkTypeString:
		return nebulaValue.AsString()
	case NebulaSdkTypeDate:
		nDate, _ := nebulaValue.AsDate()
		dateUTC := time.Date(int(nDate.GetYear()), time.Month(nDate.GetMonth()), int(nDate.GetDay()), 0, 0, 0, 0, time.UTC)
		date := dateUTC.In(timezoneDefault)
		return date, nil
	case NebulaSdkTypeTime:
		nTimeW, _ := nebulaValue.AsTime()
		nTime, _ := nTimeW.GetLocalTimeWithTimezoneName(timezoneDefault.String())
		dateObj := time.Date(2020, 1, 1, int(nTime.GetHour()), int(nTime.GetMinute()), int(nTime.GetSec()), int(nTime.GetMicrosec()*1000), timezoneDefault)
		return dateObj.Format("15:04:05.000000"), nil
	case NebulaSdkTypeDatetime:
		nDatetimeW, _ := nebulaValue.AsDateTime()
		nDatetime, _ := nDatetimeW.GetLocalDateTimeWithTimezoneName(timezoneDefault.String())
		return time.Date(int(nDatetime.GetYear()), time.Month(nDatetime.GetMonth()), int(nDatetime.GetDay()), int(nDatetime.GetHour()), int(nDatetime.GetMinute()), int(nDatetime.GetSec()), int(nDatetime.GetMicrosec()*1000), timezoneDefault), nil
	case NebulaSdkTypeVertex:
		return nebulaValue.AsNode()
	case NebulaSdkTypeEdge:
		return nebulaValue.AsRelationship()
	case NebulaSdkTypeList:
		nList, _ := nebulaValue.AsList()
		res := make([]any, 0, len(nList))
		for _, v := range nList {
			vIface, err := GetValueIface(&v)
			if err != nil {
				return nil, err
			}
			res = append(res, vIface)
		}
		return res, nil
	case NebulaSdkTypeMap:
		nMap, _ := nebulaValue.AsMap()
		res := make(map[string]any, len(nMap))
		for k, v := range nMap {
			vIface, err := GetValueIface(&v)
			if err != nil {
				return nil, err
			}
			res[k] = vIface
		}
		return res, nil
	case NebulaSdkTypeSet:
		nList, _ := nebulaValue.AsDedupList()
		res := make([]any, 0, len(nList))
		for _, v := range nList {
			vIface, err := GetValueIface(&v)
			if err != nil {
				return nil, err
			}
			res = append(res, vIface)
		}
		return res, nil
	case "path":
		return nebulaValue.AsPath()
	case "geography":
		return nebulaValue.AsGeography()
	}
	return nil, fmt.Errorf("norm: can not get nebula type %s interface value", nebulaValue.GetType())
}
