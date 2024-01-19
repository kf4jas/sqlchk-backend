package db

import (
	"database/sql"
	"errors"
	"fmt"
	lite "github.com/mattn/go-sqlite3"
	"github.com/spf13/viper"
	"log"
	"regexp"
    "encoding/json"
    "bytes"
    "strings"
)

const file string = "sqlite.db"

// printQueryResultsqlite - a very ugly function that allows me to return various things
func PrintQueryResultSqlite(db *sql.DB, query string) ([]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		if liteErr, ok := err.(*lite.Error); ok {
			log.Println(liteErr.Code.Name())
			return nil, errors.New(liteErr.Code.Name())
		}
		return nil, errors.New("unknown") // fiber.StatusInternalServerError
	}
	defer rows.Close()
	rowsout, err := ProcessRows(rows)
	return rowsout, err
}

func CheckifTableExistsSQLite(table string) bool {
	queryValue := "SELECT name FROM sqlite_master WHERE type='table' AND name='{table_name}';"
	db := OpenConn("sqlite3")
	rowsout, err := PrintQueryResult(db, queryValue)
	if err != nil {
		log.Fatal("6", err)
		// return false
	}
	for _, v := range rowsout {
		m := v.(map[string]interface{})
		if m["table"] == table {
			return true
		}
	}
	return false
}

func CheckifColumnExistsSQLite(table, column string) bool {
	if !IsSQLName(table) || !IsSQLName(column) {
		fmt.Println("is not them")
		return false
	}
    // PRAGMA table_info(tablename)
	queryValue := "SELECT column_name FROM information_schema.columns WHERE table_name='" + table + "' and column_name='" + column + "';"
	fmt.Println(queryValue)
	db := OpenConn()
	rowsout, err := PrintQueryResult(db, queryValue)
	if err != nil {
		log.Println("7", err)
		return true
	}
	for _, v := range rowsout {
		fmt.Println(v)
		m := v.(map[string]interface{})
		if m["column_name"] == column {
			return true
		}
	}
	return false
}

func ProcessJSONSQLite(table_name string, body []byte) error {
		var output map[string]string
		var outq = []string{}
		if ! IsSQLName(table_name) {
			return errors.New("bad table name")
		}
		dec := json.NewDecoder(bytes.NewReader(body))
		err := dec.Decode(&output)
		if err != nil {
			return err
		}
		if ! CheckifTableExists(table_name) {
			outq = append(outq, "CREATE TABLE "+table_name+" (id INTEGER PRIMARY KEY AUTOINCREMENT);")
		}
		fieldsArr := []string{}
		valuesArr := []string{}
		for k, v := range output {
			if ! CheckifColumnExists(table_name, k) {
				outq = append(outq, "ALTER TABLE "+table_name+" ADD COLUMN "+k+" VARCHAR(126);")
			}
			fieldsArr = append(fieldsArr, k)
			valuesArr = append(valuesArr, v)
		}
		fields := strings.Join(fieldsArr, ",")
		values := "'" + strings.Join(valuesArr, "', '") + "'"
		outq = append(outq, "INSERT INTO "+table_name+" ("+fields+") VALUES ("+values+");")
		cdb := OpenConn()
		for _, s := range outq {
			result, err := cdb.Exec(s)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(result)
			fmt.Println(s)
		}
        return nil
}
