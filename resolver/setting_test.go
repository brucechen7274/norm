package resolver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

func TestParseTagSetting(t *testing.T) {
	tests := []struct {
		tag  string
		want map[string]string
	}{
		{tag: "a:1", want: map[string]string{"a": "1"}},
		{tag: "a", want: map[string]string{"a": "a"}},
		{tag: "", want: map[string]string{}},
		{tag: "a:1;b:2", want: map[string]string{"a": "1", "b": "2"}},
		{tag: "a:1;c;b:2", want: map[string]string{"a": "1", "b": "2", "c": "c"}},
		{tag: "a:1;;b:2;;", want: map[string]string{"a": "1", "b": "2"}},
		{tag: "a:1:2:3;b", want: map[string]string{"a": "1:2:3", "b": "b"}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			assert.Equal(t, tt.want, ParseTagSetting(tt.tag))
		})
	}
}

func Test_camelCaseToUnderscore(t *testing.T) {
	tests := []struct {
		s    string
		want string
	}{
		{s: "", want: ""},
		{s: "aBc", want: "a_bc"},
		{s: "ABC", want: "a_b_c"},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			assert.Equal(t, tt.want, camelCaseToUnderscore(tt.s))
		})
	}
}

func TestGetFieldIndex(t *testing.T) {
	tests := []struct {
		field reflect.StructField
		want  *FieldIndex
	}{
		{field: reflect.StructField{Name: "Name"}, want: nil},
		{field: reflect.StructField{Name: "Name", Tag: `norm:"prop:name"`}, want: nil},
		{field: reflect.StructField{Name: "Name", Tag: `norm:"index"`}, want: &FieldIndex{Name: "idx_name", Prop: "name", Priority: 10}},
		{field: reflect.StructField{Name: "Name", Tag: `norm:"index:idx_name"`}, want: &FieldIndex{Name: "idx_name", Prop: "name", Priority: 10}},
		{field: reflect.StructField{Name: "Name", Tag: `norm:"index:,priority:5"`}, want: &FieldIndex{Name: "idx_name", Prop: "name", Priority: 5}},
		{field: reflect.StructField{Name: "Name", Tag: `norm:"index:,length:5"`}, want: &FieldIndex{Name: "idx_name", Prop: "name", Length: 5, Priority: 10}},
		{field: reflect.StructField{Name: "Name", Tag: `norm:"index:idx_name,priority:5,length:10"`}, want: &FieldIndex{Name: "idx_name", Prop: "name", Priority: 5, Length: 10}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			propName := GetPropName(tt.field)
			assert.Equal(t, tt.want, GetFieldIndex(tt.field, propName, ""))
		})
	}
}
