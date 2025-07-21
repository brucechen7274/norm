package statement

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestCreateVertexTags(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().CreateVertexTags(&vm1{}, true)
			},
			want: `CREATE TAG IF NOT EXISTS player(name string, age int);`,
		},
		{
			stmt: func() *Statement {
				return New().CreateVertexTags(&vm2{}, true)
			},
			want: `CREATE TAG IF NOT EXISTS no_property();`,
		},
		{
			stmt: func() *Statement {
				return New().CreateVertexTags(&vm3{}, true)
			},
			want: `CREATE TAG IF NOT EXISTS player_with_default(name string, age int DEFAULT 20);`,
		},
		{
			stmt: func() *Statement {
				return New().CreateVertexTags(vm4{}, true)
			},
			want: `CREATE TAG IF NOT EXISTS woman(name string, age int, married bool, salary double, create_time timestamp) TTL_DURATION = 100, TTL_COL = "create_time";`,
		},
		{
			stmt: func() *Statement {
				return New().CreateVertexTags(vm5{})
			},
			want: `CREATE TAG woman(name fixed_string DEFAULT "hayson", age int32 DEFAULT 20, create_time datetime DEFAULT datetime(1625469277));`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, ngql)
			}
		})
	}
}

func TestDropTag(t *testing.T) {
	tests := []struct {
		stmt    func() *Statement
		want    string
		wantErr bool
	}{
		{
			stmt: func() *Statement {
				return New().DropTag("test")
			},
			want: `DROP TAG test;`,
		},
		{
			stmt: func() *Statement {
				return New().DropTag("test", true)
			},
			want: `DROP TAG IF EXISTS test;`,
		},
	}
	for i, tt := range tests {
		t.Run(fmt.Sprintf("#_%d", i), func(t *testing.T) {
			s := tt.stmt()
			ngql, err := s.NGQL()
			if tt.wantErr {
				assert.Error(t, err)
				return
			}
			if assert.NoError(t, err) {
				assert.Equal(t, tt.want, ngql)
			}
		})
	}
}

type vm1 struct {
	VID  string `norm:"vertex_id"`
	Name string
	Age  int
}

func (t vm1) VertexID() string {
	return t.VID
}

func (t vm1) VertexTagName() string {
	return "player"
}

type vm2 struct {
	VID string `norm:"vertex_id"`
}

func (t vm2) VertexID() string {
	return t.VID
}

func (t vm2) VertexTagName() string {
	return "no_property"
}

type vm3 struct {
	VID  string `norm:"vertex_id"`
	Name string `norm:"prop:name;type:string"`
	Age  int    `norm:"prop:age;type:int;default:20"`
}

func (t vm3) VertexID() string {
	return t.VID
}

func (t vm3) VertexTagName() string {
	return "player_with_default"
}

type vm4 struct {
	VID        string `norm:"vertex_id"`
	Name       string
	Age        int
	Married    bool
	Salary     float64
	CreateTime time.Time `norm:"type:timestamp;ttl:100"`
}

func (t *vm4) VertexID() string {
	return t.VID
}

func (t *vm4) VertexTagName() string {
	return "woman"
}

type vm5 struct {
	VID        string    `norm:"vertex_id"`
	Name       string    `norm:"prop:name;type:fixed_string;default:hayson"`
	Age        int       `norm:"prop:age;type:int32;default:20"`
	CreateTime time.Time `norm:"default:datetime(1625469277)"`
}

func (t *vm5) VertexID() string {
	return t.VID
}

func (t *vm5) VertexTagName() string {
	return "woman"
}
