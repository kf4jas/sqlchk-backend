package db

import (
	"database/sql"
	"errors"
	"fmt"
	"github.com/lib/pq"
	"log"
	// "regexp"
	"bytes"
	"encoding/json"
	"strings"
)

type PostgresDriver struct {
	Bknd    Backend
	Db      *sql.DB
	ConnStr string
}

// printQueryResult - a very ugly function that allows me to return various things
func (p PostgresDriver) PrintQueryResult(query string) ([]interface{}, error) {
	rows, err := p.Db.Query(query)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			return nil, errors.New(pqErr.Code.Name())
		}
		return nil, errors.New("unknown") // fiber.StatusInternalServerError
	}
	defer rows.Close()
	rowsout, err := p.Bknd.ProcessRows(rows)
	return rowsout, err
}

func (p PostgresDriver) CheckifTableExists(table string) bool {
	queryValue := "select tablename as table from pg_tables where schemaname = 'public'"
	p.OpenConn()
	rowsout, err := p.PrintQueryResult(queryValue)
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

func (p PostgresDriver) CheckifColumnExists(table, column string) bool {
	if !p.Bknd.IsSQLName(table) || !p.Bknd.IsSQLName(column) {
		fmt.Println("is not them")
		return false
	}
	queryValue := "SELECT column_name FROM information_schema.columns WHERE table_name='" + table + "' and column_name='" + column + "';"
	fmt.Println(queryValue)
	p.OpenConn()
	rowsout, err := p.PrintQueryResult(queryValue)
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

func (p PostgresDriver) ProcessJSON(table_name string, body []byte) error {
	var output map[string]string
	var outq = []string{}
	if !p.Bknd.IsSQLName(table_name) {
		return errors.New("bad table name")
	}
	dec := json.NewDecoder(bytes.NewReader(body))
	err := dec.Decode(&output)
	if err != nil {
		return err
	}
	if !p.CheckifTableExists(table_name) {
		outq = append(outq, "CREATE TABLE "+table_name+" (id SERIAL PRIMARY KEY);")
	}
	fieldsArr := []string{}
	valuesArr := []string{}
	for k, v := range output {
		if !p.CheckifColumnExists(table_name, k) {
			outq = append(outq, "ALTER TABLE "+table_name+" ADD COLUMN "+k+" VARCHAR(126);")
		}
		fieldsArr = append(fieldsArr, k)
		valuesArr = append(valuesArr, v)
	}
	fields := strings.Join(fieldsArr, ",")
	values := "'" + strings.Join(valuesArr, "', '") + "'"
	outq = append(outq, "INSERT INTO "+table_name+" ("+fields+") VALUES ("+values+");")
	p.OpenConn()
	for _, s := range outq {
		result, err := p.Db.Exec(s)
		if err != nil {
			fmt.Println(err)
		}
		fmt.Println(result)
		fmt.Println(s)
	}
	return nil
}

func (p PostgresDriver) OpenConn() {
	var err error
	p.Db, err = sql.Open("postgres", p.ConnStr)
	if err != nil {
		log.Fatal(err)
	}
	return
}
