package db

import (
	"database/sql"
	"fmt"
	"github.com/spf13/viper"
	// "log"
	"regexp"
)

type DBDriver interface {
	PrintQueryResult(query string) ([]interface{}, error)
	CheckifTableExists(table string) bool
	CheckifColumnExists(table, column string) bool
	ProcessJSON(table_name string, body []byte) error
	OpenConn()
}

type Backend struct{}

// ProcessRows - this turns the rows into key / value objects
func (b Backend) ProcessRows(rows *sql.Rows) ([]interface{}, error) {
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
		rout := b.orderedRows(cols, row...)
		out = append(out, rout)
	}
	return out, rows.Err()
}

// orderRows - lists of lists
func (b Backend) orderRows(row ...interface{}) []interface{} {
	out := make([]interface{}, 0)
	for _, r := range row {
		out = append(out, r)
	}
	return out
}

// orderedRows - Json like objects
func (b Backend) orderedRows(cols []string, row ...interface{}) map[string]interface{} {
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

// IsSQLName - this checks if the name is a valid table or column name
func (b Backend) IsSQLName(input string) bool {
	re := regexp.MustCompile("^[A-Za-z_][A-Za-z0-9_-]*$")
	if re.MatchString(input) {
		return true
	}
	return false
}

func GetDriver() DBDriver {
	connStr := viper.GetString("connStr")
	switch connStr[0:5] {
	//~ case "mysql":
	//~ return PostgresDriver{}
	case "postg":
		return PostgresDriver{ConnStr: connStr}
	default:
		//sqlite3
		return SQLiteDriver{ConnStr: connStr}
	}
}
