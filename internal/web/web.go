package web

import (
	"database/sql"
	"errors"
	"embed"
    "fmt"
	"net/http"
    "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/recover"
    "github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/lib/pq"
	"log"
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
        Root: http.FS(EmbedDirStatic),
        PathPrefix: "frontend/public",
        Browse: false, // security 
    }))

	app.Get("/query", func(c *fiber.Ctx) error {
		queryValue := c.Query("q")
		// connValue := c.Query("c")
		//~ return c.SendString(queryValue+" "+connValue)

		connStr := "postgresql://joee:password@localhost/joee?sslmode=require"
		// connStr := connValue
		// Connect to database
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}

		// Safety Code v1.1
		// safeQuery := SafetyChecks(payload.Query)
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

	app.Post("/query", func(c *fiber.Ctx) error {

		payload := Payload{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}

		// connStr := "postgresql://joee:password@localhost/joee?sslmode=require"
		connStr := payload.ConnStr
		// Connect to database
		db, err := sql.Open("postgres", connStr)
		if err != nil {
			log.Fatal(err)
		}

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
	out := make([]interface{}, 0)
	rows, err := db.Query(query)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			log.Println(pqErr.Code.Name())
			return out, errors.New(pqErr.Code.Name())
		}
		return out, errors.New("unknown") // fiber.StatusInternalServerError
	}
	defer rows.Close()
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
		err = rows.Scan(rowPtr...)
		if err != nil {
			fmt.Println("cannot scan row:", err)
		}
		fmt.Println(row...)
		rout := orderedRows(cols, row...)
		out = append(out, rout)
	}
	return out, rows.Err()
}

func orderRows(row ...interface{}) []interface{} {
	out := make([]interface{}, 0)
	for _, r := range row {
		out = append(out, r)
	}
	return out
}

func orderedRows(cols []string, row ...interface{}) map[string]interface{} {
	out := make(map[string]interface{}, 0)
	for i, r := range cols {
		out[r] = row[i]
	}
	return out
}
