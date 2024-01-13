package utils

import (
	"fmt"
    "time"
    "context"
	"github.com/google/uuid"
    "github.com/RichardKnop/machinery/v1"
	"github.com/RichardKnop/machinery/v1/config"
    "github.com/RichardKnop/machinery/v1/tasks"
    "github.com/RichardKnop/machinery/v1/log"

    tracers "github.com/RichardKnop/machinery/example/tracers"
	task "sqlchk/internal/tasks"
    opentracing "github.com/opentracing/opentracing-go"
    opentracing_log "github.com/opentracing/opentracing-go/log"
	"go.uber.org/zap"
)

var (
	Logger *zap.SugaredLogger
)

func init() {
	logger, _ := zap.NewProduction()
	Logger = logger.Sugar()
}

func StartServer() (*machinery.Server, error) {
	cnf := &config.Config{
		DefaultQueue:    "machinery_tasks",
		ResultsExpireIn: 3600,
		Broker:          "redis://172.16.25.6:6379",
		ResultBackend:   "redis://172.16.25.6:6379",
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}

	server, err := machinery.NewServer(cnf)
	if err != nil {
		return nil, err
	}

	// Register tasks
	rtasks := map[string]interface{}{
		"addtolog": task.AddToLog,
        "processjsontask": task.ProcessJSONTask,
	}

	return server, server.RegisterTasks(rtasks)
}


func GetMachineryServer() *machinery.Server {
	Logger.Info("initing task server")

	taskserver, err := machinery.NewServer(&config.Config{
		Broker:        "redis://redismq:6379",
		ResultBackend: "redis://redismq:6379",
	})
	if err != nil {
		Logger.Fatal(err.Error())
	}

	taskserver.RegisterTasks(map[string]interface{}{
		"addtolog": task.AddToLog,
        "processjsontask": task.ProcessJSONTask,
	})

	return taskserver
}


func SendProcessJSONTask(table_name string,content []byte) error {
    cleanup, err := tracers.SetupTracer("sender")
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := StartServer()
	if err != nil {
		return err
	}

	var (
		ProcessJSONTask                     tasks.Signature
	)

	var initTasks = func() {
		ProcessJSONTask = tasks.Signature{
			Name: "processjsontask",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: table_name,
				},
                {
					Type:  "[]byte",
					Value: content,                    
                },
			},
		}
    }
	
	/*
	 * Lets start a span representing this run of the `send` command and
	 * set a batch id as baggage so it can travel all the way into
	 * the worker functions.
	 */
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracing_log.String("batch.id", batchID))

	log.INFO.Println("Starting batch:", batchID)
	/*
	 * First, let's try sending a single task
	 */
	initTasks()

	log.INFO.Println("Single task:")

	asyncResult, err := server.SendTaskWithContext(ctx, &ProcessJSONTask)
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		return fmt.Errorf("Getting task result failed with error: %s", err.Error())
	}
	log.INFO.Printf("1 + 1 = %v\n", tasks.HumanReadableResults(results))
	return nil
}


func SendToLog(data string) error {
	cleanup, err := tracers.SetupTracer("sender")
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := StartServer()
	if err != nil {
		return err
	}

	var (
		AddToLog                     tasks.Signature
	)

	var initTasks = func() {
		AddToLog = tasks.Signature{
			Name: "addtolog",
			Args: []tasks.Arg{
				{
					Type:  "string",
					Value: data,
				},
			},
		}
    }
	
	/*
	 * Lets start a span representing this run of the `send` command and
	 * set a batch id as baggage so it can travel all the way into
	 * the worker functions.
	 */
	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracing_log.String("batch.id", batchID))

	log.INFO.Println("Starting batch:", batchID)
	/*
	 * First, let's try sending a single task
	 */
	initTasks()

	log.INFO.Println("Single task:")

	asyncResult, err := server.SendTaskWithContext(ctx, &AddToLog)
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err := asyncResult.Get(time.Duration(time.Millisecond * 5))
	if err != nil {
		return fmt.Errorf("Getting task result failed with error: %s", err.Error())
	}
	log.INFO.Printf("1 + 1 = %v\n", tasks.HumanReadableResults(results))
	return nil
}
