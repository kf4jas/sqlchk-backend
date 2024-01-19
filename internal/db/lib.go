package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"log"
	"regexp"
    "encoding/json"
    "bytes"
    "strings"
)


// ProcessRows - this turns the rows into key / value objects
func ProcessRows(rows *sql.Rows) ([]interface{}, error) {
	out := make([]interface{}, 0)
	cols, _ := rows.Columns()
	row := make([]interface{}, len(cols))
	rowPtr := make([]interface{}, len(cols))
	for i := range row {
		rowPtr[i] = &row[i]
	}
	fmt.Println(cols)
	icols := make([]interface{}, len(cols))
	for i, v := range cols {
		icols[i] = v
	}

	for rows.Next() {
		err := rows.Scan(rowPtr...)
		if err != nil {
			fmt.Println("cannot scan row:", err)
		}
		fmt.Println(row...)
		rout := orderedRows(cols, row...)
		out = append(out, rout)
	}
	return out, rows.Err()
}

// orderRows - lists of lists
func orderRows(row ...interface{}) []interface{} {
	out := make([]interface{}, 0)
	for _, r := range row {
		out = append(out, r)
	}
	return out
}

// orderedRows - Json like objects
func orderedRows(cols []string, row ...interface{}) map[string]interface{} {
	out := make(map[string]interface{}, 0)
	for i, r := range cols {
		switch row[i].(type) {
		case []uint8:
			out[r] = string(row[i].([]uint8))
		default:
			out[r] = row[i]
		}
	}
	return out
}

func IsSQLName(input string) bool {
	re := regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_-]*$")
	if re.MatchString(input) {
		return true
	}
	return false
}

func OpenConn(sqltype string) *sql.DB {
	connStr := viper.GetString("connStr") // can be the file name or a postgres connect string
	db, err := sql.Open(sqltype, connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
