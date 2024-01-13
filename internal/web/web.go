package web

import (
	"embed"
	"fmt"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
    "sqlchk/internal/db"
    "sqlchk/internal/utils"
	"net/http"
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
    app.Get("/log",func(c *fiber.Ctx) error {
        utils.SendToLog("this log")
        return nil
    })
	app.Get("/query", func(c *fiber.Ctx) error {
		queryValue := c.Query("q")
		// Safety Code v1.1
		// safeQuery := SafetyChecks(payload.Query)
		// Connect to database
		cdb := db.OpenConn()

		// Raw SQL
		fmt.Printf("query: %v\n", queryValue)
		rowsout, err := db.PrintQueryResult(cdb, queryValue)
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
		table_name := c.Query("table")
        utils.SendProcessJSONTask(table_name, c.Body())
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
		cdb := db.OpenConn()
		// Safety Code v1.1
		// safeQuery := SafetyChecks(payload.Query)
		// Raw SQL
		fmt.Printf("Wow: %v\n", payload.Query)
		rowsout, err := db.PrintQueryResult(cdb, payload.Query)
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