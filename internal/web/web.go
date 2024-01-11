package web

import (
	"bytes"
	"database/sql"
	"embed"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/lib/pq"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"regexp"
	"strings"
)

type Payload struct {
	Query   string `json:"query"`
	ConnStr string `json:"constr"`
}

type Response struct {
	Query  string        `json:"query"`
	Result []interface{} `json:"result"`
}

// Embed a directory from Main
var EmbedDirStatic embed.FS

func Start() {
	app := fiber.New()
	CustomConfig := recover.Config{
		// Next:              nil,
		EnableStackTrace: true,
		// StackTraceHandler: defaultStackTraceHandler,
	}
	app.Use(recover.New(CustomConfig))

	//app.Use(recover.New())

	app.Get("/", func(c *fiber.Ctx) error {
		return c.SendString("Hello, World!")
	})

	app.Use("/static", filesystem.New(filesystem.Config{
		Root:       http.FS(EmbedDirStatic),
		PathPrefix: "frontend/public",
		Browse:     false, // security
	}))

	app.Get("/query", func(c *fiber.Ctx) error {
		queryValue := c.Query("q")
		// Safety Code v1.1
		// safeQuery := SafetyChecks(payload.Query)
		// Connect to database
		db := OpenConn()

		// Raw SQL
		fmt.Printf("query: %v\n", queryValue)
		rowsout, err := printQueryResult(db, queryValue)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		fmt.Println(rowsout)
		response := Response{
			Query:  queryValue,
			Result: rowsout,
		}
		return c.JSON(response)
	})

	app.Post("/data", func(c *fiber.Ctx) error {
		var output map[string]string
		var outq = []string{}
		table_name := c.Query("table")
		if !IsSQLName(table_name) {
			return fiber.ErrBadRequest
		}
		dec := json.NewDecoder(bytes.NewReader(c.Body()))
		err := dec.Decode(&output)
		if err != nil {
			return err
		}
		if !CheckifTableExists(table_name) {
			outq = append(outq, "CREATE TABLE "+table_name+" (id SERIAL PRIMARY KEY);")
		}
		fieldsArr := []string{}
		valuesArr := []string{}
		for k, v := range output {
			if !CheckifColumnExists(table_name, k) {
				outq = append(outq, "ALTER TABLE "+table_name+" ADD COLUMN "+k+" VARCHAR(126);")
			}
			fieldsArr = append(fieldsArr, k)
			valuesArr = append(valuesArr, v)
		}
		fields := strings.Join(fieldsArr, ",")
		values := "'" + strings.Join(valuesArr, "', '") + "'"
		outq = append(outq, "INSERT INTO "+table_name+" ("+fields+") VALUES ("+values+");")
		db := OpenConn()
		for _, s := range outq {
			result, err := db.Exec(s)
			if err != nil {
				fmt.Println(err)
			}
			fmt.Println(result)
			fmt.Println(s)
		}
		return c.Send([]byte("Added"))
	})

	app.Post("/query", func(c *fiber.Ctx) error {

		payload := Payload{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		// connStr := "postgresql://joee:password@localhost/joee?sslmode=require"
		// connStr := payload.ConnStr
		// Connect to database
		db := OpenConn()
		// Safety Code v1.1
		// safeQuery := SafetyChecks(payload.Query)
		// Raw SQL
		fmt.Printf("Wow: %v\n", payload.Query)
		rowsout, err := printQueryResult(db, payload.Query)
		if err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, err.Error())
		}
		fmt.Println(rowsout)
		response := Response{
			Query:  payload.Query,
			Result: rowsout,
		}
		return c.JSON(response)
	})

	app.Static("/static/", "./public")
	app.Listen(":3030")
}

// printQueryResult - a very ugly function that allows me to return various things
func printQueryResult(db *sql.DB, query string) ([]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			return nil, errors.New(pqErr.Code.Name())
		}
		return nil, errors.New("unknown") // fiber.StatusInternalServerError
	}
	defer rows.Close()
	rowsout, err := ProcessRows(rows)
	return rowsout, err
}

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

func CheckifTableExists(table string) bool {
	queryValue := "select tablename as table from pg_tables where schemaname = 'public'"
	db := OpenConn()
	rowsout, err := printQueryResult(db, queryValue)
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

func CheckifColumnExists(table, column string) bool {
	if !IsSQLName(table) || !IsSQLName(column) {
		fmt.Println("is not them")
		return false
	}
	queryValue := "SELECT column_name FROM information_schema.columns WHERE table_name='" + table + "' and column_name='" + column + "';"
	fmt.Println(queryValue)
	db := OpenConn()
	rowsout, err := printQueryResult(db, queryValue)
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

func IsSQLName(input string) bool {
	re := regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_-]*$")
	if re.MatchString(input) {
		return true
	}
	return false
}

func OpenConn() *sql.DB {
	connStr := viper.GetString("connStr")
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}
	return db
}
