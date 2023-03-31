package goexer

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

// One item inside container.
type Item struct {
	Value interface{}
	Type  string // Just for better string representation.
	Name  string // Just for better string representation.
}

// Marshaled representation.
type Marshaled map[string]interface{}

// Return string representation for Item.
func (i *Item) String() string {
	return fmt.Sprintf("Item{ Name: '%s' Type: '%s' Value: %v}", i.Name, i.Type, i.Get())
}

// Return raw, unconverted Value as interface{}.
func (i *Item) GetRaw() interface{} {
	return i.Value
}

// Return converted to original type, if possible, value of item.

//nolint:gocyclo
func (i *Item) Get() any {
	switch i.Type {
	case "bool":
		return cast.ToBool(i.Value)
	case "time.Time":
		return cast.ToTime(i.Value)
	case "time.Duration":
		return cast.ToDuration(i.Value)
	case "float64":
		return cast.ToFloat64(i.Value)
	case "float32":
		return cast.ToFloat32(i.Value)
	case "int64":
		return cast.ToInt64(i.Value)
	case "int32":
		return cast.ToInt32(i.Value)
	case "int16":
		return cast.ToInt16(i.Value)
	case "int8":
		return cast.ToInt8(i.Value)
	case "int":
		return cast.ToInt(i.Value)
	case "uint64":
		return cast.ToUint64(i.Value)
	case "uint32":
		return cast.ToUint32(i.Value)
	case "uint16":
		return cast.ToUint16(i.Value)
	case "uint8":
		return cast.ToUint8(i.Value)
	case "uint":
		return cast.ToUint(i.Value)

	case "string":
		return cast.ToString(i.Value)
	case "map[string]string":
		return cast.ToStringMapString(i.Value)
	case "map[string][]string":
		return cast.ToStringMapStringSlice(i.Value)
	case "map[string]bool":
		return cast.ToStringMapBool(i.Value)
	case "map[string]int":
		return cast.ToStringMapInt(i.Value)
	case "map[string]int64":
		return cast.ToStringMapInt64(i.Value)
	case "[]bool":
		return cast.ToBoolSlice(i.Value)
	case "[]string":
		return cast.ToStringSlice(i.Value)
	case "[]int":
		return cast.ToIntSlice(i.Value)
	case "[]time.Duration":
		return cast.ToDurationSlice(i.Value)
	default:
		return i.Value
	}
}

// Container for extended data.
type Container struct {
	items map[string]Item
}

// Return value of Item.
func (c *Container) Get(key string) any {
	item, ok := c.items[key]
	if !ok {
		return nil
	}

	return item.Get()
}

// Return raw value of Item.
func (c *Container) GetRaw(key string) any {
	item, ok := c.items[key]
	if !ok {
		return nil
	}

	return item.GetRaw()
}

// Return value of Item with ok status.
func (c *Container) GetE(key string) (any, bool) {
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	return item.Get(), true
}

// Return raw value of Item with ok status.
func (c *Container) GetRawE(key string) (any, bool) {
	item, ok := c.items[key]
	if !ok {
		return nil, false
	}

	return item.GetRaw(), true
}

// Set value for item.
func (c *Container) Set(key string, value interface{}) *Container {
	item := Item{
		Name:  key,
		Value: value,
	}

	item.Type = reflect.TypeOf(value).String()

	c.items[key] = item

	return c
}

// Create new variable of type Container.
func NewContainer() *Container {
	cc := Container{}
	cc.items = make(map[string]Item)

	return &cc
}

// Return count of fields (items) inside Container.
func (c *Container) Size() int {
	return len(c.items)
}

// Return all keys of fields inside Container.
func (c *Container) Keys() []string {
	return lo.Keys(c.items)
}

// Return marshaled representation.
func (c *Container) Marshaled() Marshaled {
	m := make(map[string]interface{}, len(c.items)+1)
	for k, v := range c.items {
		m[k] = v.Get()
	}

	return m
}

// return json representation.
func (c *Container) JSON() ([]byte, error) {
	j, err := json.Marshal(c.Marshaled())
	if err != nil {
		//nolint:wrapcheck
		return nil, err
	}

	return j, nil
}
