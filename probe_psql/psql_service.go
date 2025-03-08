package probe_psql

import (
	"log"
	"reflect"

	"github.com/jackc/pgx"
)

var conf = pgx.ConnConfig{
	Host:     "127.0.0.1",
	Port:     5432,
	Database: "service",
	User:     "service",
	Password: "service_ex",
	LogLevel: pgx.LogLevelDebug,
}

func structPointersSpread(struct_ interface{}) []interface{} {
	reflectValue := reflect.ValueOf(struct_).Elem()
	numFields := reflectValue.NumField()
	fieldPointers := make([]interface{}, numFields, numFields)

	for i := 0; i < numFields; i++ {
		field := reflectValue.Field(i)
		fieldPointers[i] = field.Addr().Interface()
	}

	return fieldPointers
}

func logQuery(query string) {
	switch conf.LogLevel {
	case pgx.LogLevelDebug:
		log.Printf("INFO: Query to %s\nQuery:\n%s\n", conf.Host, query)
	}
}

func probePsqlExtractDataDueReflect[T interface{}](query string) ([]T, error) {
	conn, err := pgx.Connect(conf)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
		return nil, err
	}
	defer conn.Close()

	logQuery(query)
	rows, err := conn.Query(query)
	if err != nil {
		log.Fatalf("Query failed: %v\n", err)
		return nil, err
	}

	result := make([]T, 0)
	for rows.Next() {
		rowData := new(T)
		if err := rows.Scan(structPointersSpread(rowData)...); err != nil {
			log.Fatalf("Query failed: %v\n", err)
			return nil, err
		}
		result = append(result, *rowData)
	}

	return result, nil
}
