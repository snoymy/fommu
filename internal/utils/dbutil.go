package utils

import (
	"fmt"
	"reflect"

	"github.com/jackc/pgx/v5"
)

func GetTable[T any](rows pgx.Rows) (out []T) {
    var table []T
    for rows.Next() {
        var data T
        s := reflect.ValueOf(&data).Elem()
        numCols := s.NumField()
        columns := make([]interface{}, numCols)

        for i := 0; i < numCols; i++ {
            field := s.Field(i)
            columns[i] = field.Addr().Interface()
        }

        if err := rows.Scan(columns...); err != nil {
            fmt.Println("Case Read Error ", err)
        }

        table = append(table, data)
    }
    return table
}
