package goexer_test

import (
	"reflect"
	"testing"
	"time"

	"github.com/Tolyar/goexer"
)

func TestContainerSetGet(t *testing.T) {
	container := goexer.NewContainer()

	t.Parallel()

	test := []struct {
		Value any
		Type  string
	}{
		{goexer.Item{}, "interface{}"},
		{[]time.Duration{time.Duration(10)}, "[]time.Duration"},
		{int(10), "int"},
		{[]string{"test"}, "[]string"},
		{[]bool{true, false}, "[]bool"},
		{[]int{10, 20}, "[]int"},
		{map[string]int64{"Test": int64(10)}, "map[string]int64"},
		{map[string]int{"Test": int(10)}, "map[string]int"},
		{map[string]bool{"Test": true}, "map[string]bool"},
		{map[string][]string{"Test": {"test"}}, "map[string][]string"},
		{map[string]string{"Test": "test"}, "map[string]string"},
		{"test", "string"},

		{uint(10), "uint"},
		{uint8(10), "uint8"},
		{uint16(10), "uint16"},
		{uint32(10), "uint32"},
		{uint64(10), "uint64"},

		{true, "bool"},
		{time.Now(), "time.Time"},
		{time.Duration(10), "time.Duration"},
		{float64(10.1), "float64"},
		{float32(10.1), "float32"},

		{int64(-1), "int64"},
		{int32(-1), "int32"},
		{int16(-1), "int16"},
		{int8(-1), "int8"},
		{int(-1), "int"},
	}

	for _, tt := range test {
		container.Set(tt.Type, tt.Value)
		vv := container.Get(tt.Type)
		if !reflect.DeepEqual(vv, tt.Value) {
			t.Errorf("Set(%s, %v) != Get(). Get returns '%v'", tt.Type, tt.Value, vv)
		}
	}

	for _, tt := range test {
		container.Set(tt.Type, tt.Value)
		vv, ok := container.GetE(tt.Type)
		if !reflect.DeepEqual(vv, tt.Value) {
			t.Errorf("Set(%s, %v) != Get(). Get returns '%v'", tt.Type, tt.Value, vv)
		}

		if !ok {
			t.Errorf("Want ok==true, got ok==false for %v", tt.Type)
		}
	}

	// Test for absent key.
	if container.Get("not-exists") != nil {
		t.Errorf("Get: want 'nil', got '%v'", container.Get("not-exists"))
	}

	_, ok := container.GetE("not-exists")
	// Test for absent key.
	if ok == true {
		t.Errorf("GetE: want 'false', got 'true'")
	}
}

func TestItemString(t *testing.T) {
	t.Parallel()
	item := goexer.Item{Name: "test", Value: int(10), Type: "int"}
	if item.String() != "Item{ Name: 'test' Type: 'int' Value: 10}" {
		t.Errorf("Want '%s' , got '%s'", "zzzzz", item.String())
	}
}

func TestContainerGetRaw(t *testing.T) {
	t.Parallel()

	cc := goexer.NewContainer()
	cc.Set("test123", "test123")

	rr := cc.GetRaw("test123")
	vv, ok := rr.(string)
	if !ok || vv != "test123" {
		t.Errorf("Want 'test123' but got '%v'", rr)
	}

	rr = cc.GetRaw("not-exists")
	if rr != nil {
		t.Errorf("GetRaw not exists want 'nil', got '%v' with type '%s'", rr, reflect.TypeOf(rr).String())
	}
}

func TestContainerGetRawE(t *testing.T) {
	t.Parallel()

	cc := goexer.NewContainer()
	cc.Set("test", "test")

	rr, ok := cc.GetRawE("test")
	if !ok {
		t.Errorf("Can't find key '%s'", "test")
	}

	vv, ok := rr.(string)
	if !ok || vv != "test" {
		t.Errorf("Want 'test' but got '%v'", rr)
	}

	if vv, ok := cc.GetRawE("not-exists"); ok || vv != nil {
		t.Errorf("GetRawE not exists want 'nil', got '%v'", cc.GetRaw("not-exists"))
	}
}
