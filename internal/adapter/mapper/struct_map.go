package mapper

import "encoding/json"

func StructToMap(obj interface{}) (map[string]interface{}, error) {
    var result map[string]interface{}

    jsonBytes, err := json.Marshal(obj)
    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(jsonBytes, &result)
    if err != nil {
        return nil, err
    }

    return result, nil
}

func MapToStruct[T any](obj map[string]interface{}) (*T, error) {
    var result *T

    jsonBytes, err := json.Marshal(obj)
    if err != nil {
        return nil, err
    }

    err = json.Unmarshal(jsonBytes, &result)
    if err != nil {
        return nil, err
    }

    return result, nil
}
