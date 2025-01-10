package resolver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"strings"
	"testing"
)

func TestParseEdge(t *testing.T) {
	tests := []struct {
		dest     interface{}
		want     *EdgeSchema
		wantProp []prop
		wantErr  bool
	}{
		{dest: edge1{}, want: &EdgeSchema{srcVIDType: VIDTypeString, srcVIDFieldIndex: []int{2}, dstVIDType: VIDTypeString, dstVIDFieldIndex: []int{3}, rankFieldIndex: []int{4}, edgeTypeName: "edge1"}, wantProp: []prop{
			{name: "name", index: []int{0}}, {name: "age", index: []int{1}}, {name: "gender", index: []int{5}, nebulaType: "string"},
		}},
		{dest: &edge2{}, want: &EdgeSchema{srcVIDType: VIDTypeInt64, srcVIDFieldIndex: []int{0}, dstVIDType: VIDTypeString, dstVIDFieldIndex: []int{1}, rankFieldIndex: nil, edgeTypeName: "edge2"}, wantProp: []prop{
			{name: "name", index: []int{2}}, {name: "age", index: []int{3}},
		}},
		{dest: edge3{}, wantErr: true},
		{dest: record1{}, wantErr: true},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			got, err := ParseEdge(reflect.TypeOf(tt.dest))
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if !assert.NoError(t, err) {
				return
			}
			assert.Equal(t, tt.want.srcVIDType, got.srcVIDType)
			assert.Equal(t, tt.want.srcVIDFieldIndex, got.srcVIDFieldIndex)
			assert.Equal(t, tt.want.dstVIDType, got.dstVIDType)
			assert.Equal(t, tt.want.dstVIDFieldIndex, got.dstVIDFieldIndex)
			assert.Equal(t, tt.want.edgeTypeName, got.edgeTypeName)
			assert.Equal(t, tt.want.rankFieldIndex, got.rankFieldIndex)
			propsGot := make([]prop, 0)
			for _, p := range got.props {
				propsGot = append(propsGot, prop{
					name:       p.Name,
					index:      p.StructField.Index,
					nebulaType: p.NebulaType,
				})
			}
			assert.Equal(t, tt.wantProp, propsGot)
		})
	}
}

func TestGetEdgeInfo(t *testing.T) {
	e1 := &edge1{
		Name:   "name1",
		Age:    18,
		SrcID:  "101",
		DstID:  "102",
		Rank:   1,
		Gender: 1,
	}
	e2 := edge1{
		Name:   "name2",
		Age:    19,
		SrcID:  "201",
		DstID:  "202",
		Rank:   2,
		Gender: 1,
	}
	e3 := &edge1{
		Name:   "name3",
		Age:    20,
		SrcID:  "301",
		DstID:  "302",
		Rank:   3,
		Gender: 1,
	}
	edgeSchema, err := ParseEdge(reflect.TypeOf(e1))
	if err != nil {
		t.Errorf("ParseEdge() error = %v", err)
		return
	}
	tests := []struct {
		e             interface{}
		wantSrcIDExpr string
		wantDstIDExpr string
		wantRank      int64
		wantPropExpr  string
	}{
		{e: e1, wantSrcIDExpr: `"101"`, wantDstIDExpr: `"102"`, wantRank: 1, wantPropExpr: `"name1" 18 1`},
		{e: e2, wantSrcIDExpr: `"201"`, wantDstIDExpr: `"202"`, wantRank: 2, wantPropExpr: `"name2" 19 1`},
		{e: e3, wantSrcIDExpr: `"301"`, wantDstIDExpr: `"302"`, wantRank: 3, wantPropExpr: `"name3" 20 1`},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			edgeValue := reflect.ValueOf(tt.e)
			assert.Equal(t, tt.wantSrcIDExpr, edgeSchema.GetSrcVIDExpr(edgeValue))
			assert.Equal(t, tt.wantDstIDExpr, edgeSchema.GetDstVIDExpr(edgeValue))
			assert.Equal(t, tt.wantRank, edgeSchema.GetRank(edgeValue))
			propStr := ""
			edgeValue = reflect.Indirect(edgeValue)
			for _, p := range edgeSchema.GetProps() {
				f := edgeValue.FieldByIndex(p.StructField.Index)
				if !f.IsZero() {
					s, _ := FormatSimpleValue("", f)
					propStr += " " + s
				}
			}
			assert.Equal(t, tt.wantPropExpr, strings.TrimSpace(propStr))
		})
	}
}

type edge1 struct {
	Name     string `norm:"prop:name"`
	Age      int    `norm:"prop:age"`
	SrcID    string `norm:"edge_src_id"`
	DstID    string `norm:"edge_dst_id"`
	Rank     int    `norm:"edge_rank"`
	Gender   int    `norm:"datatype:string"`
	Pleasure string `norm:"-"`
}

func (e *edge1) EdgeTypeName() string {
	return "edge1"
}

type edge2 struct {
	SrcID int64  `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Name  string
	Age   int
}

func (e edge2) EdgeTypeName() string {
	return "edge2"
}

type edge3 struct {
	SrcID int64 `norm:"edge_src_id"`
	Name  string
	Age   int
}

func (e edge3) EdgeTypeName() string {
	return "edge3"
}
