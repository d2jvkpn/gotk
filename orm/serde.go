package orm

import (
	"bytes"
	"encoding/json"
	"fmt"

	"database/sql/driver"
)

// result
type PageResult[T any] struct {
	Total int64 `json:"total"`
	Items []T   `json:"items"`
}

func NewPageResult[T any]() *PageResult[T] {
	return &PageResult[T]{Items: make([]T, 0)}
}

func (result PageResult[T]) AsMap() map[string]any {
	return map[string]any{
		"total": result.Total, "items": result.Items,
	}
}

// vector
type Vector[T any] []T // T should marshalable

func (vec *Vector[T]) Scan(value any) (err error) {
	bts, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	bts = bytes.TrimSpace(bts)

	result := make(Vector[T], 0, 5)
	if len(bts) > 0 { // empty []uint8 cause error
		err = json.Unmarshal(bts, &result)
	}

	*vec = result
	return err
}

func (vec Vector[T]) Value() (driver.Value, error) {
	return json.Marshal(vec)
}

// dict
type Hashmap[T any] map[string]T

func (dict *Hashmap[T]) Scan(value any) (err error) {
	bts, ok := value.([]byte)
	if !ok {
		return fmt.Errorf("failed to unmarshal JSONB value: %v", value)
	}
	bts = bytes.TrimSpace(bts)

	result := make(Hashmap[T], 5)
	if len(bts) > 0 { // empty []uint8 causes error
		err = json.Unmarshal(bts, &result)
	}

	*dict = result
	return err
}

func (dict Hashmap[T]) Value() (driver.Value, error) {
	return json.Marshal(dict)
}
