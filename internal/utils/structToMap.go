package utils

import (
	"encoding/json"
)

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

//func StructToMap(obj interface{}) map[string]interface{} {
//    // 1. Create an empty map named result to store the fields and their values.
//    result := make(map[string]interface{})
//
//    // 2. Get the reflect.Value and reflect.Type of the input object.
//    val := reflect.ValueOf(obj)
//
//    // 3. If the input object is a pointer, dereference it to get the underlying value.
//    if val.Kind() == reflect.Ptr {
//        val = val.Elem()
//    }
//
//    typ := val.Type()
//
//    // 4. Iterate through the fields of the struct using a for loop.
//    for i := 0; i < val.NumField(); i++ {
//        // 5. For each field, get its name and kind (e.g., string, int, struct).
//        fieldName := typ.Field(i).Name
//        fieldValueKind := val.Field(i).Kind()
//        var fieldValue interface{}
//
//        // 6. If the field is a struct, recursively call structToMap to get the map representation of the nested struct.
//        // Otherwise, get the field value directly.
//        if fieldValueKind == reflect.Struct {
//            fieldValue = StructToMap(val.Field(i).Interface())
//        } else {
//            fieldValue = val.Field(i).Interface()
//        }
//
//        // 7. Add the field name and value to the result map.
//        result[fieldName] = fieldValue
//    }
//
//    return result
//}
