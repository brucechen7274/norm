package resolver

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func TestFormatSimpleValue(t *testing.T) {
	a := "hello"
	tests := []struct {
		nebulaType string
		value      []any
		want       string
		wantErr    bool
	}{
		{
			value: []any{1, int8(1), int16(1), int32(1), int64(1), uint(1), uint8(1), uint16(1), uint32(1), uint64(1)},
			want:  "1",
		},
		{
			value: []any{-1, int8(-1), int16(-1), int32(-1), int64(-1)},
			want:  "-1",
		},
		{
			value: []any{1000.234, float32(1000.234)},
			want:  "1000.234",
		},
		{
			value: []any{-1000.234, float32(-1000.234)},
			want:  "-1000.234",
		},
		{
			nebulaType: NebulaSdkTypeFloat,
			value:      []any{100},
			want:       "100",
		},
		{
			nebulaType: NebulaSdkTypeInt,
			value:      []any{100.1234},
			want:       "100",
		},
		{
			value: []any{true},
			want:  "true",
		},
		{
			value: []any{false},
			want:  "false",
		},
		{
			value: []any{"hello world 你好 世界"},
			want:  `"hello world 你好 世界"`,
		},
		{
			value: []any{"Hello \\ world 你好 \" \t t 世界"},
			want:  `"Hello \\ world 你好 \" \t t 世界"`,
		},
		{
			value: []any{`Hello \ world 你好 " t 世界`},
			want:  `"Hello \\ world 你好 \" t 世界"`,
		},
		{
			nebulaType: NebulaSdkTypeDatetime,
			value:      []any{`2024-08-01T00:00:00`},
			want:       `datetime("2024-08-01T00:00:00")`,
		},
		{
			nebulaType: NebulaSdkTypeDate,
			value:      []any{`2023-12-12`},
			want:       `date("2023-12-12")`,
		},
		{
			nebulaType: NebulaSdkTypeTime,
			value:      []any{`11:00:51.457000`},
			want:       `time("11:00:51.457000")`,
		},
		{
			value: []any{time.Date(2024, 8, 20, 11, 16, 30, 10000, time.Local)},
			want:  `datetime("2024-08-20T11:16:30.000010")`,
		},
		{
			nebulaType: NebulaSdkTypeDate,
			value:      []any{time.Date(2024, 8, 20, 11, 16, 30, 10000, time.Local)},
			want:       `date("2024-08-20")`,
		},
		{
			nebulaType: NebulaSdkTypeTime,
			value:      []any{time.Date(2024, 8, 20, 11, 16, 30, 10000, time.Local)},
			want:       `time("11:16:30.000010")`,
		},
		{
			value: []any{[]int{1, -1, 2, -2, 0}},
			want:  "[1, -1, 2, -2, 0]",
		},
		{
			value: []any{make([]int, 0)},
			want:  "[]",
		},
		{
			value: []any{[]string{"h", "e", "l", "l", "o"}},
			want:  `["h", "e", "l", "l", "o"]`,
		},
		{
			value: []any{[]time.Time{time.Date(2024, 8, 20, 11, 16, 30, 10000, time.Local), time.Date(2024, 8, 20, 11, 16, 16, 0, time.Local)}},
			want:  `[datetime("2024-08-20T11:16:30.000010"), datetime("2024-08-20T11:16:16")]`,
		},
		{
			value: []any{[5]string{"h", "e", "l", "l", "o"}},
			want:  `["h", "e", "l", "l", "o"]`,
		},
		{
			nebulaType: NebulaSdkTypeSet,
			value:      []any{[]int{1, -1, 2, -2, 0}, []int64{1, -1, 2, -2, 0}},
			want:       "set{1, -1, 2, -2, 0}",
		},
		{
			nebulaType: NebulaSdkTypeSet,
			value:      []any{[]int(nil)},
			want:       "set{}",
		},
		{
			value: []any{map[string]int{"c": 3}},
			want:  `map{c: 3}`,
		},
		{
			value: []any{map[string]any{"d": map[string]int{"age": 18}}},
			want:  `map{d: map{age: 18}}`,
		},
		{
			nebulaType: NebulaSdkTypeSet,
			value:      []any{map[int]struct{}{1: {}}},
			want:       `set{1}`,
		},
		{
			value: []any{(*int)(nil)},
			want:  `NULL`,
		},
		{
			value: []any{&a},
			want:  `"hello"`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("case #%d", i), func(t *testing.T) {
			for _, v := range tt.value {
				got, err := FormatSimpleValue(tt.nebulaType, reflect.ValueOf(v))
				if tt.wantErr {
					assert.Error(t, err)
					return
				}
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}
