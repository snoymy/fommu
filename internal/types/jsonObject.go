package types

import (
	"database/sql/driver"
	"encoding/json"
	"errors"
)

type JsonObject map[string]interface{}

func (j JsonObject) Scan(value any) error {
	if value == nil {
        j = nil
		return nil
	}

    var data []byte
    switch v := value.(type) {
        case string:
            data = []byte(v)
        case []byte:
            data = v
        default:
	        return errors.New("Unable to scan value to type JsonObject")
    }

    err := json.Unmarshal(data, &j)
    if err != nil {
        return err
    }

    return nil
}

func (j JsonObject) Value() (driver.Value, error) {
    bytes, err := json.Marshal(j)
    if err != nil {
        return nil, err
    }

    return string(bytes), nil
}
