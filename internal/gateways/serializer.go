package gateways

import "encoding/json"

type Serializer[T any] interface {
	Serialize(value []T) (string, error)
	Deserialize(value string) ([]T, error)
}

type JSONSerializer[T any] struct{}

func NewJSONSerializer[T any]() Serializer[T] {
	return JSONSerializer[T]{}
}

func (s JSONSerializer[T]) Serialize(value []T) (string, error) {
	bytes, err := json.Marshal(value)
	return string(bytes), err
}

func (s JSONSerializer[T]) Deserialize(value string) ([]T, error) {
	var data []T
	err := json.Unmarshal([]byte(value), &data)
	return data, err
}
