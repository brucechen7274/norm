package resolver

import (
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"
	"testing"
	"time"

	nebula "github.com/vesoft-inc/nebula-go/v3"
)

// BenchmarkResolver comprehensively tests performance metrics of resolver package
// Including memory allocation, memory usage, and operation time

// Test data structure
type BenchmarkPlayer struct {
	ID       string    `norm:"vertex_id"`
	Name     string    `norm:"prop:name"`
	Age      int       `norm:"prop:age"`
	Score    float64   `norm:"prop:score"`
	Active   bool      `norm:"prop:active"`
	Birthday time.Time `norm:"prop:birthday"`
}

func (p BenchmarkPlayer) VertexID() string      { return p.ID }
func (p BenchmarkPlayer) VertexTagName() string { return "benchmark_player" }

type BenchmarkFriendship struct {
	SrcID   string `norm:"edge_src_id"`
	DstID   string `norm:"edge_dst_id"`
	Rank    int    `norm:"edge_rank"`
	Level   string `norm:"prop:level"`
	Since   int    `norm:"prop:since"`
	Comment string `norm:"prop:comment"`
}

func (f BenchmarkFriendship) EdgeSrcID() string    { return f.SrcID }
func (f BenchmarkFriendship) EdgeDstID() string    { return f.DstID }
func (f BenchmarkFriendship) EdgeRank() int        { return f.Rank }
func (f BenchmarkFriendship) EdgeTypeName() string { return "benchmark_friendship" }

// Create mock Nebula ValueWrapper for testing
func createMockNebulaValue() *nebula.ValueWrapper {
	// This is a simplified mock, actual testing may require more complex implementation
	// For benchmark testing, we mainly focus on resolver performance
	return &nebula.ValueWrapper{}
}

// BenchmarkFormatSimpleValue_SimpleTypes tests formatting performance of simple types
func BenchmarkFormatSimpleValue_SimpleTypes(b *testing.B) {
	tests := []struct {
		name    string
		sdkType string
		value   any
		allocs  int64 // Expected allocation count
	}{
		{"Int", "", 42, 0},
		{"String", "", "hello world", 1},
		{"Float", "", 3.14159, 0},
		{"Bool", "", true, 0},
		{"DateTime", "", time.Now(), 0},
		{"NilPointer", "", (*string)(nil), 0},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			var result string
			var err error

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result, err = FormatSimpleValue(tt.sdkType, reflect.ValueOf(tt.value))
				if err != nil {
					b.Fatal(err)
				}
			}
			_ = result // Avoid compiler optimization
		})
	}
}

// BenchmarkFormatSimpleValue_ComplexTypes tests formatting performance of complex types
func BenchmarkFormatSimpleValue_ComplexTypes(b *testing.B) {
	// Pre-allocate test data
	intSlice := make([]int, 100)
	for i := range intSlice {
		intSlice[i] = i
	}

	stringSlice := make([]string, 50)
	for i := range stringSlice {
		stringSlice[i] = fmt.Sprintf("item_%d", i)
	}

	mapData := make(map[string]any)
	for i := 0; i < 30; i++ {
		mapData[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
	}

	tests := []struct {
		name    string
		sdkType string
		value   any
	}{
		{"LargeList", NebulaSdkTypeList, intSlice},
		{"StringList", NebulaSdkTypeList, stringSlice},
		{"LargeMap", NebulaSdkTypeMap, mapData},
		{"ListAsSet", NebulaSdkTypeSet, intSlice},
		{"MapAsSet", NebulaSdkTypeSet, mapData},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			var result string
			var err error

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result, err = FormatSimpleValue(tt.sdkType, reflect.ValueOf(tt.value))
				if err != nil {
					b.Fatal(err)
				}
			}
			_ = result
		})
	}
}

// BenchmarkFormatSimpleValue_NestedStructures tests formatting performance of nested structures
func BenchmarkFormatSimpleValue_NestedStructures(b *testing.B) {
	// Create nested structures
	nestedList := make([][]string, 20)
	for i := range nestedList {
		nestedList[i] = make([]string, 10)
		for j := range nestedList[i] {
			nestedList[i][j] = fmt.Sprintf("nested_%d_%d", i, j)
		}
	}

	nestedMap := make(map[string][]any)
	for i := 0; i < 15; i++ {
		nestedMap[fmt.Sprintf("key_%d", i)] = []any{
			i, fmt.Sprintf("str_%d", i), float64(i) / 3.0,
		}
	}

	tests := []struct {
		name    string
		sdkType string
		value   any
	}{
		{"NestedList", NebulaSdkTypeList, nestedList},
		{"NestedMap", NebulaSdkTypeMap, nestedMap},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			var result string
			var err error

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				result, err = FormatSimpleValue(tt.sdkType, reflect.ValueOf(tt.value))
				if err != nil {
					b.Fatal(err)
				}
			}
			_ = result
		})
	}
}

// BenchmarkScanValue_SimpleTypes tests scanning performance of simple types
func BenchmarkScanValue_SimpleTypes(b *testing.B) {
	r := NewResolver()

	// Mock target values
	var destInt int
	var destString string
	var destFloat float64
	var destBool bool
	var destTime time.Time

	tests := []struct {
		name     string
		mockType string
		dest     any
	}{
		{"ToInt", "int", &destInt},
		{"ToString", "string", &destString},
		{"ToFloat", "float", &destFloat},
		{"ToBool", "bool", &destBool},
		{"ToTime", "datetime", &destTime},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			mockValue := createMockNebulaValue()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := r.ScanValue(mockValue, reflect.ValueOf(tt.dest).Elem())
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkScanValue_ComplexTypes tests scanning performance of complex types
func BenchmarkScanValue_ComplexTypes(b *testing.B) {
	r := NewResolver()

	// Pre-allocate target structures
	var destSlice []int
	var destMap map[string]any
	var destSet []string
	var destPlayer BenchmarkPlayer

	tests := []struct {
		name     string
		mockType string
		dest     any
	}{
		{"ToSlice", "list", &destSlice},
		{"ToMap", "map", &destMap},
		{"ToSet", "set", &destSet},
		{"ToVertex", "vertex", &destPlayer},
	}

	for _, tt := range tests {
		b.Run(tt.name, func(b *testing.B) {
			b.ReportAllocs()
			mockValue := createMockNebulaValue()

			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				err := r.ScanValue(mockValue, reflect.ValueOf(tt.dest).Elem())
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	}
}

// BenchmarkSchemaOperations tests performance of schema operations
func BenchmarkSchemaOperations(b *testing.B) {
	r := NewResolver()

	playerType := reflect.TypeOf(BenchmarkPlayer{})
	friendshipType := reflect.TypeOf(BenchmarkFriendship{})

	recordType := reflect.TypeOf(struct {
		ID   string `norm:"col:id"`
		Name string `norm:"col:name"`
		Age  int    `norm:"col:age"`
	}{})

	b.Run("GetVertexSchema", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := r.getVertexSchema(playerType)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("GetEdgeSchema", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := r.getEdgeSchema(friendshipType)
			if err != nil {
				b.Fatal(err)
			}
		}
	})

	b.Run("GetRecordSchema", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := r.getRecordSchema(recordType)
			if err != nil {
				b.Fatal(err)
			}
		}
	})
}

// BenchmarkMemoryUsage tests memory usage
func BenchmarkMemoryUsage(b *testing.B) {
	// Test memory usage when formatting large amounts of data
	b.Run("MemoryUsage_LargeList", func(b *testing.B) {
		largeList := make([]int, 10000)
		for i := range largeList {
			largeList[i] = i
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := FormatSimpleValue(NebulaSdkTypeList, reflect.ValueOf(largeList))
			if err != nil {
				b.Fatal(err)
			}
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "bytes/op")
	})

	b.Run("MemoryUsage_LargeMap", func(b *testing.B) {
		largeMap := make(map[string]any)
		for i := 0; i < 1000; i++ {
			largeMap[fmt.Sprintf("key_%d", i)] = fmt.Sprintf("value_%d", i)
		}

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			_, err := FormatSimpleValue(NebulaSdkTypeMap, reflect.ValueOf(largeMap))
			if err != nil {
				b.Fatal(err)
			}
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "bytes/op")
	})
}

// BenchmarkConcurrentOperations tests performance of concurrent operations
func BenchmarkConcurrentOperations(b *testing.B) {
	r := NewResolver()

	b.Run("ConcurrentFormat", func(b *testing.B) {
		b.ReportAllocs()
		b.RunParallel(func(pb *testing.PB) {
			data := []int{1, 2, 3, 4, 5}
			for pb.Next() {
				_, err := FormatSimpleValue(NebulaSdkTypeList, reflect.ValueOf(data))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	b.Run("ConcurrentSchemaAccess", func(b *testing.B) {
		b.ReportAllocs()
		playerType := reflect.TypeOf(BenchmarkPlayer{})
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := r.getVertexSchema(playerType)
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// BenchmarkStringBuilderUsage specifically tests StringBuilder usage patterns
func BenchmarkStringBuilderUsage(b *testing.B) {
	b.Run("CurrentImplementation", func(b *testing.B) {
		b.ReportAllocs()

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Simulate current implementation - create new Builder each time
			var builder strings.Builder
			builder.WriteString("[")
			for j := 0; j < 100; j++ {
				builder.WriteString(strconv.Itoa(j))
				if j < 99 {
					builder.WriteString(", ")
				}
			}
			builder.WriteString("]")
			_ = builder.String()
		}
	})
}

// BenchmarkStringConcatenation compares different string concatenation methods
func BenchmarkStringConcatenation(b *testing.B) {
	items := make([]string, 100)
	for i := range items {
		items[i] = fmt.Sprintf("item_%d", i)
	}

	b.Run("StringConcat", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			result := ""
			for _, item := range items {
				result += item + ", "
			}
			_ = result
		}
	})

	b.Run("StringBuilder", func(b *testing.B) {
		b.ReportAllocs()
		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			var builder strings.Builder
			for _, item := range items {
				builder.WriteString(item)
				builder.WriteString(", ")
			}
			_ = builder.String()
		}
	})
}

// BenchmarkOverall comprehensive performance test
func BenchmarkOverall(b *testing.B) {
	// Create complex test data
	complexData := struct {
		IDs     []string
		Players []BenchmarkPlayer
		Mapping map[string]any
	}{
		IDs:     make([]string, 1000),
		Players: make([]BenchmarkPlayer, 100),
		Mapping: make(map[string]any),
	}

	for i := range complexData.IDs {
		complexData.IDs[i] = fmt.Sprintf("id_%d", i)
	}

	for i := range complexData.Players {
		complexData.Players[i] = BenchmarkPlayer{
			ID:       fmt.Sprintf("player_%d", i),
			Name:     fmt.Sprintf("Player %d", i),
			Age:      20 + (i % 30),
			Score:    float64(i) * 0.1,
			Active:   i%2 == 0,
			Birthday: time.Now().AddDate(-20-(i%30), 0, 0),
		}
	}

	for i := 0; i < 500; i++ {
		complexData.Mapping[fmt.Sprintf("key_%d", i)] = []any{
			i, fmt.Sprintf("value_%d", i), float64(i) / 10.0, i%2 == 0,
		}
	}

	b.Run("ComplexDataFormat", func(b *testing.B) {
		b.ReportAllocs()

		var m1, m2 runtime.MemStats
		runtime.GC()
		runtime.ReadMemStats(&m1)

		b.ResetTimer()
		for i := 0; i < b.N; i++ {
			// Format various complex types
			_, err1 := FormatSimpleValue(NebulaSdkTypeList, reflect.ValueOf(complexData.IDs))
			_, err2 := FormatSimpleValue(NebulaSdkTypeMap, reflect.ValueOf(complexData.Mapping))
			if err1 != nil || err2 != nil {
				b.Fatal(err1, err2)
			}
		}

		runtime.GC()
		runtime.ReadMemStats(&m2)

		b.ReportMetric(float64(m2.TotalAlloc-m1.TotalAlloc)/float64(b.N), "bytes/op")
	})
}
