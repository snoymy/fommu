package types

import (
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"time"
)

type Nullable[T any] struct {
    value T
    valid bool
}

func NewNullable[T any](value T) Nullable[T] {
    return Nullable[T]{value: value, valid: true}
}

func Null[T any]() Nullable[T] {
    return Nullable[T]{valid: false}
}

func (n Nullable[T]) ValueOrError() (T, error) {
    if !n.valid {
        var zero T
        return zero, errors.New("Unable to get null value from Nullable type")
    }

    return n.value, nil
}

func (n Nullable[T]) ValueOrFail() T {
    if !n.valid {
        panic("Unable to get null value from Nullable type")
    }

    return n.value
}

func (n Nullable[T]) ValueOrZero() T {
    if !n.valid {
        var zero T
        return zero
    }

    return n.value
}

func (n Nullable[T]) Value() (driver.Value, error) {
	if !n.valid {
		return nil, nil
	}

	if valuer, ok := interface{}(n.value).(driver.Valuer); ok {
		return valuer.Value()
	}

	return convertToDriverValue(n.value)
}

func (n *Nullable[T]) Set(value T) {
    n.value = value
    n.valid = true
}

func (n Nullable[T]) IsNull() bool {
    return !n.valid
}

func (n *Nullable[T]) SetNull() {
    var zero T
    n.value = zero
    n.valid = false
}

func (n *Nullable[T]) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		n.valid = false
		return nil
	}

	var value T
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	n.value = value
	n.valid = true
	return nil
}

func (n Nullable[T]) MarshalJSON() ([]byte, error) {
	if !n.valid {
		return []byte("null"), nil
	}

	return json.Marshal(n.value)
}

func (n *Nullable[T]) Scan(value any) error {
	if value == nil {
        var zero T
		n.value = zero
		n.valid = false
		return nil
	}

	if scanner, ok := interface{}(&n.value).(sql.Scanner); ok {
		if err := scanner.Scan(value); err != nil {
			return err
		}
		n.valid = true
		return nil
	}

	var err error
	n.value, err = convertToType[T](value)
	n.valid = err == nil
	return err
}

func (n Nullable[T]) Format(f fmt.State, c rune) {
    if n.valid {
        fmt.Fprint(f, n.value)
    } else {
        fmt.Fprint(f, "null")
    }
}

func convertToType[T any](value any) (T, error) {
	var zero T
	if value == nil {
		return zero, nil
	}

	valueType := reflect.TypeOf(value)
	targetType := reflect.TypeOf(zero)
	if valueType == targetType {
		return value.(T), nil
	}

	isNumeric := func(kind reflect.Kind) bool {
		return kind >= reflect.Int && kind <= reflect.Float64
	}

	// Check if the value is a numeric type and if T is also a numeric type.
	if isNumeric(valueType.Kind()) && isNumeric(targetType.Kind()) {
		convertedValue := reflect.ValueOf(value).Convert(targetType)
		return convertedValue.Interface().(T), nil
	}

	return zero, errors.New("unsupported type conversion")
}

func convertToDriverValue(v any) (driver.Value, error) {
	if valuer, ok := v.(driver.Valuer); ok {
		return valuer.Value()
	}

	rv := reflect.ValueOf(v)
	switch rv.Kind() {
	case reflect.Pointer:
		if rv.IsNil() {
			return nil, nil
		}
		return convertToDriverValue(rv.Elem().Interface())

	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return rv.Int(), nil

	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32:
		return int64(rv.Uint()), nil

	case reflect.Uint64:
		u64 := rv.Uint()
		if u64 >= 1<<63 {
			return nil, fmt.Errorf("uint64 values with high bit set are not supported")
		}
		return int64(u64), nil

	case reflect.Float32, reflect.Float64:
		return rv.Float(), nil

	case reflect.Bool:
		return rv.Bool(), nil

	case reflect.Slice:
		if rv.Type().Elem().Kind() == reflect.Uint8 {
			return rv.Bytes(), nil
		}
		return nil, fmt.Errorf("unsupported slice type: %s", rv.Type().Elem().Kind())

	case reflect.String:
		return rv.String(), nil

	case reflect.Struct:
		if t, ok := v.(time.Time); ok {
			return t, nil
		}
		return nil, fmt.Errorf("unsupported struct type: %s", rv.Type())

	default:
		return nil, fmt.Errorf("unsupported type: %T", v)
	}
}
