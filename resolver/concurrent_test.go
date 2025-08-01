package resolver

import (
	"reflect"
	"sync"
	"testing"
	"time"
)

// TestVertex test vertex structure
type TestVertex struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name"`
	Age  int    `norm:"prop:age"`
}

func (t TestVertex) VertexID() string      { return t.VID }
func (t TestVertex) VertexTagName() string { return "test_vertex" }

// TestEdge test edge structure
type TestEdge struct {
	SrcID string `norm:"edge_src_id"`
	DstID string `norm:"edge_dst_id"`
	Rank  int    `norm:"edge_rank"`
	Label string `norm:"prop:label"`
}

func (t TestEdge) EdgeSrcID() string    { return t.SrcID }
func (t TestEdge) EdgeDstID() string    { return t.DstID }
func (t TestEdge) EdgeRank() int        { return t.Rank }
func (t TestEdge) EdgeTypeName() string { return "test_edge" }

// TestConcurrentVertexSchemaAccess tests race condition for concurrent vertex schema access
func TestConcurrentVertexSchemaAccess(t *testing.T) {
	r := NewResolver()
	var wg sync.WaitGroup
	const goroutines = 100

	// Concurrent access to the same type of schema
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := r.getVertexSchema(reflect.TypeOf(TestVertex{}))
			if err != nil {
				t.Errorf("getVertexSchema failed: %v", err)
			}
		}()
	}
	wg.Wait()
}

// TestConcurrentEdgeSchemaAccess tests race condition for concurrent edge schema access
func TestConcurrentEdgeSchemaAccess(t *testing.T) {
	r := NewResolver()
	var wg sync.WaitGroup
	const goroutines = 100

	// Concurrent access to the same type of schema
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := r.getEdgeSchema(reflect.TypeOf(TestEdge{}))
			if err != nil {
				t.Errorf("getEdgeSchema failed: %v", err)
			}
		}()
	}
	wg.Wait()
}

// TestMixedSchemaAccess tests mixed concurrent access to different schema types
func TestMixedSchemaAccess(t *testing.T) {
	r := NewResolver()
	var wg sync.WaitGroup
	const goroutines = 50

	// Mixed access to different types of schemas
	for i := 0; i < goroutines; i++ {
		wg.Add(2)
		go func() {
			defer wg.Done()
			_, err := r.getVertexSchema(reflect.TypeOf(TestVertex{}))
			if err != nil {
				t.Errorf("getVertexSchema failed: %v", err)
			}
		}()
		go func() {
			defer wg.Done()
			_, err := r.getEdgeSchema(reflect.TypeOf(TestEdge{}))
			if err != nil {
				t.Errorf("getEdgeSchema failed: %v", err)
			}
		}()
	}
	wg.Wait()
}

// TestConcurrentRecordSchemaAccess tests concurrent record schema access
func TestConcurrentRecordSchemaAccess(t *testing.T) {
	r := NewResolver()
	var wg sync.WaitGroup
	const goroutines = 100

	// Define test structure
	type TestRecord struct {
		ID   string `norm:"col:id"`
		Name string `norm:"col:name"`
		Age  int    `norm:"col:age"`
	}

	// Concurrent access to record schema
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			_, err := r.getRecordSchema(reflect.TypeOf(TestRecord{}))
			if err != nil {
				t.Errorf("getRecordSchema failed: %v", err)
			}
		}()
	}
	wg.Wait()
}

// TestCacheConsistency tests cache consistency
func TestCacheConsistency(t *testing.T) {
	r := NewResolver()

	// First access, should write to cache
	schema1, err := r.getVertexSchema(reflect.TypeOf(TestVertex{}))
	if err != nil {
		t.Fatalf("First getVertexSchema failed: %v", err)
	}

	// Second access, should read from cache
	schema2, err := r.getVertexSchema(reflect.TypeOf(TestVertex{}))
	if err != nil {
		t.Fatalf("Second getVertexSchema failed: %v", err)
	}

	// Verify that the two schema pointers are the same, indicating cache usage
	if schema1 != schema2 {
		t.Error("Schema cache is inconsistent - different instances returned for same type")
	}
}

// TestCacheConsistencyMixedTypes tests multi-type cache consistency
func TestCacheConsistencyMixedTypes(t *testing.T) {
	r := NewResolver()

	// Test vertex cache consistency
	vertexSchema1, err := r.getVertexSchema(reflect.TypeOf(TestVertex{}))
	if err != nil {
		t.Fatalf("First getVertexSchema failed: %v", err)
	}
	vertexSchema2, err := r.getVertexSchema(reflect.TypeOf(TestVertex{}))
	if err != nil {
		t.Fatalf("Second getVertexSchema failed: %v", err)
	}
	if vertexSchema1 != vertexSchema2 {
		t.Error("Vertex schema cache is inconsistent")
	}

	// Test edge cache consistency
	edgeSchema1, err := r.getEdgeSchema(reflect.TypeOf(TestEdge{}))
	if err != nil {
		t.Fatalf("First getEdgeSchema failed: %v", err)
	}
	edgeSchema2, err := r.getEdgeSchema(reflect.TypeOf(TestEdge{}))
	if err != nil {
		t.Fatalf("Second getEdgeSchema failed: %v", err)
	}
	if edgeSchema1 != edgeSchema2 {
		t.Error("Edge schema cache is inconsistent")
	}
}

// BenchmarkConcurrentSchemaAccess benchmarks concurrent schema access
func BenchmarkConcurrentSchemaAccess(b *testing.B) {
	r := NewResolver()

	b.Run("VertexSchema", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := r.getVertexSchema(reflect.TypeOf(TestVertex{}))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	b.Run("EdgeSchema", func(b *testing.B) {
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := r.getEdgeSchema(reflect.TypeOf(TestEdge{}))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})

	b.Run("RecordSchema", func(b *testing.B) {
		type TestRecord struct {
			ID   string `norm:"col:id"`
			Name string `norm:"col:name"`
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				_, err := r.getRecordSchema(reflect.TypeOf(TestRecord{}))
				if err != nil {
					b.Fatal(err)
				}
			}
		})
	})
}

// Create a complex structure to increase parsing time
type ComplexVertex struct {
	VID       string    `norm:"vertex_id"`
	Name      string    `norm:"prop:name"`
	Age       int       `norm:"prop:age"`
	Score     float64   `norm:"prop:score"`
	BirthDate time.Time `norm:"prop:birth_date"`
	Active    bool      `norm:"prop:active"`
}

func (c ComplexVertex) VertexID() string      { return c.VID }
func (c ComplexVertex) VertexTagName() string { return "complex_vertex" }

// TestSlowConcurrentSchemaAccess slow concurrent test, more likely to trigger race conditions
func TestSlowConcurrentSchemaAccess(t *testing.T) {
	r := NewResolver()
	var wg sync.WaitGroup
	const goroutines = 20

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < 10; j++ {
				_, err := r.getVertexSchema(reflect.TypeOf(ComplexVertex{}))
				if err != nil {
					t.Errorf("getVertexSchema failed: %v", err)
				}
				// Small delay to increase probability of race conditions
				time.Sleep(time.Microsecond * 10)
			}
		}()
	}
	wg.Wait()
}
