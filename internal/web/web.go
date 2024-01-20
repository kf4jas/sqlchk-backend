package web

import (
	"embed"
	"fmt"
    "github.com/spf13/viper"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"net/http"
	"sqlchk/internal/db"
	"sqlchk/internal/utils"
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

	app.Get("/log", func(c *fiber.Ctx) error {
		utils.SendToLog("this log")
		return nil
	})

	app.Get("/query", func(c *fiber.Ctx) error {
		queryValue := c.Query("q")
		// Safety Code v1.1
		// safeQuery := SafetyChecks(payload.Query)
		// Connect to database
		cdb := db.GetDriver()

		// Raw SQL
		fmt.Printf("query: %v\n", queryValue)
		db := cdb.OpenConn()
        fmt.Println("Opened a connection")
        rowsout, err := cdb.PrintQueryResult(db, queryValue)
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
        mqMode := viper.GetBool("mq_mode")
	    if mqMode {
            utils.SendProcessJSONTask(table_name, c.Body())
            return c.Send([]byte("Added m(q)"))
        }
        cdb := db.GetDriver()
        err := cdb.ProcessJSON(table_name,c.Body())
        if err != nil {
            return fiber.NewError(fiber.StatusInternalServerError, err.Error())
        }
        return c.Send([]byte("Added"))
	})

	app.Post("/query", func(c *fiber.Ctx) error {

		payload := Payload{}
		if err := c.BodyParser(&payload); err != nil {
			return err
		}
        cdb := db.GetDriver()
		fmt.Printf("Wow: %v\n", payload.Query)
        db := cdb.OpenConn()
        fmt.Println("Opened a connection")
		rowsout, err := cdb.PrintQueryResult(db, payload.Query)
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
