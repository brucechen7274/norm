package resolver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

type prop struct {
	name       string
	index      []int
	nebulaType string
}

func TestParseVertex(t *testing.T) {
	tests := []struct {
		dest                 any
		wantVIDType          VIDType
		wantVIDIndex         []int
		wantVIDMethodIndex   int
		wantVIDReceiverIsPtr bool
		wantTag              map[string][]prop
		wantErr              bool
	}{
		{dest: vertex1{}, wantVIDType: VIDTypeString, wantVIDIndex: []int{2}, wantVIDMethodIndex: 0, wantVIDReceiverIsPtr: true, wantTag: map[string][]prop{
			"vertex_tag1": {
				{"name", []int{0}, "string"},
				{"age", []int{1}, "int"},
			},
		}},
		{dest: vertex2{}, wantVIDType: VIDTypeString, wantVIDIndex: nil, wantVIDMethodIndex: 1, wantVIDReceiverIsPtr: true, wantTag: map[string][]prop{
			"vertex_tag2": {
				{"name", []int{0}, "string"},
				{"age", []int{1}, "int"},
				{"gender", []int{2}, "string"},
			},
		}},
		{dest: &vertex3{}, wantVIDType: VIDTypeInt64, wantVIDIndex: nil, wantVIDMethodIndex: 1, wantVIDReceiverIsPtr: false, wantTag: map[string][]prop{
			"vertex_tag3": {
				{"name", []int{0}, "string"},
				{"age", []int{1}, "int"},
			},
		}},
		{dest: &vertex4{}, wantVIDType: VIDTypeString, wantVIDIndex: []int{4}, wantVIDMethodIndex: 0, wantVIDReceiverIsPtr: true, wantTag: map[string][]prop{
			"vertex_tag1": {
				{"name", []int{2, 0}, "string"},
				{"age", []int{2, 1}, "int"},
			},
			"vertex_tag2": {
				{"name", []int{3, 0}, "string"},
				{"age", []int{3, 1}, "int"},
				{"gender", []int{3, 2}, "string"},
			},
		}},
		{dest: &vertex5{}, wantVIDType: VIDTypeString, wantVIDIndex: []int{0, 0}, wantVIDMethodIndex: 0, wantVIDReceiverIsPtr: true, wantTag: map[string][]prop{
			"vertex_tag5": {
				{"name", []int{1}, "string"},
				{"age", []int{2}, "int"},
			},
		}},
		{dest: &vertex6{}, wantVIDType: VIDTypeString, wantVIDIndex: []int{0, 0}, wantVIDMethodIndex: 1, wantVIDReceiverIsPtr: false, wantTag: map[string][]prop{
			"vertex_tag1": {
				{"name", []int{1, 0}, "string"},
				{"age", []int{1, 1}, "int"},
			},
			"vertex_tag2": {
				{"name", []int{2, 0}, "string"},
				{"age", []int{2, 1}, "int"},
				{"gender", []int{2, 2}, "string"},
			},
		}},
		{dest: &vertex7{}, wantVIDType: VIDTypeString, wantVIDIndex: []int{0, 2}, wantVIDMethodIndex: 0, wantVIDReceiverIsPtr: true, wantTag: map[string][]prop{
			"vertex_tag1": {
				{"name", []int{1}, "string"},
				{"age", []int{2}, "int"},
			},
		}},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			got, err := ParseVertex(reflect.TypeOf(tt.dest))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.wantVIDType, got.vidType)
			assert.Equal(t, tt.wantVIDIndex, got.vidFieldIndex)
			assert.Equal(t, tt.wantVIDMethodIndex, got.vidMethodIndex)
			assert.Equal(t, tt.wantVIDReceiverIsPtr, got.vidReceiverIsPtr)
			gotTag := make(map[string][]prop)
			for _, tag := range got.tags {
				props := make([]prop, 0)
				for _, p := range tag.props {
					props = append(props, prop{
						name:       p.Name,
						index:      p.StructField.Index,
						nebulaType: p.SdkType,
					})
				}
				gotTag[tag.TagName] = props
			}
			assert.Equal(t, tt.wantTag, gotTag)
		})
	}
}

func TestGetVertexInfo(t *testing.T) {
	v1 := vertex4{
		Tag3: vertex1{
			Name: "name11",
		},
		Tag4: &vertex2{
			Name: "name21",
		},
		VID: "v1",
	}
	v2 := &vertex4{
		Tag3: vertex1{
			Name: "name12",
		},
		Tag4: &vertex2{
			Name: "name22",
		},
		VID: "v2",
	}
	v3 := &vertex4{
		Tag3: vertex1{
			Name: "name13",
		},
		Tag4: &vertex2{
			Name: "name23",
		},
		VID: "v3",
	}
	v4 := &vertex5{
		vertexBase: vertexBase{VID: "v4"},
		Name:       "name14",
		Age:        18,
	}
	v5 := &vertex6{
		vertexBase: &vertexBase{VID: "v5"},
		Tag1: vertex1{
			Name: "name15",
		},
		Tag2: &vertex2{
			Name: "name25",
		},
	}
	v6 := &vertex7{
		vertex1: vertex1{
			Name: "name26",
			Age:  18,
			VID:  "v6",
		},
		Name: "name27",
		Age:  28,
	}
	tests := []struct {
		v            any
		wantVIDExpr  string
		wantPropExpr string
	}{
		{v: v1, wantVIDExpr: `"v1"`, wantPropExpr: `"name11" "name21"`},
		{v: v2, wantVIDExpr: `"v2"`, wantPropExpr: `"name12" "name22"`},
		{v: v3, wantVIDExpr: `"v3"`, wantPropExpr: `"name13" "name23"`},
		{v: v4, wantVIDExpr: `"v4"`, wantPropExpr: `"name14" 18`},
		{v: v5, wantVIDExpr: `"v5"`, wantPropExpr: `"name15" "name25"`},
		{v: v6, wantVIDExpr: `"v6"`, wantPropExpr: `"name27" 28`},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			vertexSchema, err := ParseVertex(reflect.TypeOf(tt.v))
			if !assert.NoError(t, err) {
				return
			}
			vertexValue := reflect.ValueOf(tt.v)
			vidExpr := vertexSchema.GetVIDExpr(vertexValue)
			assert.Equal(t, tt.wantVIDExpr, vidExpr)
			propStr := ""
			vertexValue = reflect.Indirect(vertexValue)
			for _, tag := range vertexSchema.GetTags() {
				for _, p := range tag.GetProps() {
					f := vertexValue.FieldByIndex(p.StructField.Index)
					if !f.IsZero() {
						s, _ := FormatSimpleValue("", f)
						propStr += " " + s
					}
				}
			}
			assert.Equal(t, tt.wantPropExpr, strings.TrimSpace(propStr))
		})
	}
}

type vertex1 struct {
	Name     string `norm:"prop:name"`
	Age      int    `norm:"prop:age"`
	VID      string `norm:"vertex_id"`
	gender   int
	Pleasure string `norm:"-"`
}

func (v *vertex1) VertexID() string {
	return v.VID
}

func (v *vertex1) VertexTagName() string {
	return "vertex_tag1"
}

type vertex2 struct {
	Name   string `norm:"prop:name"`
	Age    int
	Gender int `norm:"type:string"`
}

func (v *vertex2) A() string {
	return v.Name
}

func (v *vertex2) VertexTagName() string {
	return "vertex_tag2"
}

func (v *vertex2) VertexID() string {
	return v.Name
}

type vertex3 struct {
	Name string `norm:"prop:name"`
	Age  int64
}

func (v vertex3) A() string {
	return v.Name
}

func (v vertex3) VertexTagName() string {
	return "vertex_tag3"
}

func (v vertex3) VertexID() int64 {
	return v.Age
}

type vertex4 struct {
	tag1 *vertex3
	Tag2 *vertex3 `norm:"-"`
	Tag3 vertex1
	Tag4 *vertex2
	VID  string `norm:"vertex_id"`
}

func (v *vertex4) VertexID() string {
	return v.VID
}

type vertexBase struct {
	VID string `norm:"vertex_id"`
}

func (v *vertexBase) VertexID() string {
	return v.VID
}

type vertex5 struct {
	vertexBase
	Name string `norm:"prop:name"`
	Age  int64
}

func (v vertex5) VertexTagName() string {
	return "vertex_tag5"
}

type vertex6 struct {
	*vertexBase
	Tag1 vertex1
	Tag2 *vertex2
}

func (v vertex6) A() string {
	return v.VertexID()
}

type vertex7 struct {
	vertex1
	Name string `norm:"prop:name"`
	Age  int    `norm:"prop:age"`
}
