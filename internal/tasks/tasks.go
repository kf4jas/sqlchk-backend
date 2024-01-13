package tasks

import (
	"github.com/RichardKnop/machinery/v1/log"
    "sqlchk/internal/db"
)

func AddToLog(data string) error {
    log.INFO.Println("Adding to log this:",data)
    return nil
}


func ProcessJSONTask(table_name string,content []byte) error {
    err := db.ProcessJSON(table_name, content)
    if err != nil {
           log.ERROR.Println("Error:",err) 
           return err
    }
    log.INFO.Println("Added Content")
    return nil
}
