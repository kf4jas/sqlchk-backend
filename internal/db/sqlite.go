package db

import (
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	// "github.com/spf13/viper"
	"log"
	// "regexp"
	"bytes"
	"encoding/json"
	"strings"
)

type SQLiteDriver struct {
	Bknd    Backend
	ConnStr string
}

const file string = "sqlite.db"

// printQueryResult - a very ugly function that allows me to return various things
func (s SQLiteDriver) PrintQueryResult(db *sql.DB, query string) ([]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		log.Fatal(err) // fiber.StatusInternalServerError
	}
	defer rows.Close()
	rowsout, err := s.Bknd.ProcessRows(rows)
	return rowsout, err
}

func (s SQLiteDriver) CheckifTableExists(table string) bool {
	queryValue := "SELECT name FROM sqlite_master WHERE type='table';"
	db := s.OpenConn()
	rowsout, err := s.PrintQueryResult(db, queryValue)
	if err != nil {
		log.Fatal("6", err)
		// return false
	}
    fmt.Println("rowsout: ",rowsout)
	for _, v := range rowsout {
        fmt.Println(v)
		m := v.(map[string]interface{})
		if m["name"] == table {
			return true
		}
	}
	return false
}

func (s SQLiteDriver) CheckifColumnExists(table, column string) bool {
	if !s.Bknd.IsSQLName(table) || !s.Bknd.IsSQLName(column) {
		fmt.Println("is not them")
		return false
	}
	// PRAGMA table_info(tablename)
	// queryValue := "SELECT column_name FROM information_schema.columns WHERE table_name='" + table + "' and column_name='" + column + "';"
	queryValue := "SELECT name FROM PRAGMA_TABLE_INFO('"+table+"');"
	fmt.Println(queryValue)
	db := s.OpenConn()
	rowsout, err := s.PrintQueryResult(db, queryValue)
	if err != nil {
		log.Println("7", err)
		return true
	}
	for _, v := range rowsout {
		fmt.Println("cols",v)
		m := v.(map[string]interface{})
		if m["name"] == column {
			return true
		}
	}
	return false
}

func (s SQLiteDriver) ProcessJSON(table_name string, body []byte) error {
	var output map[string]string
	var outq = []string{}
	if !s.Bknd.IsSQLName(table_name) {
		return errors.New("bad table name")
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	err := dec.Decode(&output)
	if err != nil {
		return err
	}
	if !s.CheckifTableExists(table_name) {
		outq = append(outq, "CREATE TABLE "+table_name+" (id INTEGER PRIMARY KEY AUTOINCREMENT);")
	}
	fieldsArr := []string{}
	valuesArr := []string{}
	for k, v := range output {
		if !s.CheckifColumnExists(table_name, k) {
			outq = append(outq, "ALTER TABLE "+table_name+" ADD COLUMN "+k+" VARCHAR(126);")
		}
		fieldsArr = append(fieldsArr, k)
		valuesArr = append(valuesArr, v)
	}
	fields := strings.Join(fieldsArr, ",")
	values := "'" + strings.Join(valuesArr, "', '") + "'"
	outq = append(outq, "INSERT INTO "+table_name+" ("+fields+") VALUES ("+values+");")
	db := s.OpenConn()
	for _, st := range outq {
		result, err := db.Exec(st)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
		fmt.Println(st)
	}
	return nil
}

func (s SQLiteDriver) OpenConn() *sql.DB {
	var err error
	db, err := sql.Open("sqlite3", s.ConnStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
