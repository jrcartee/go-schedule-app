package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"strings"
)

type Config struct {
	Driver   string
	Name     string
	User     string
	Password string
}

var DefaultConfig = Config{
	Driver:   "postgres",
	Name:     "go_test",
	User:     "go_test_user",
	Password: "go_test_secret",
}

func Setup(c Config) (*sql.DB, func() error) {
	connStr := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", c.User, c.Password, c.Name)
	db, err := sql.Open(c.Driver, connStr)
	if err != nil {
		log.Fatalf("Error while connecting to DB:\n %s", err)
	}
	err = db.Ping()
	if err != nil {
		log.Fatalf("Error while connecting to DB:\n %s", err)
	}
	return db, db.Close
}

type Queryer interface {
	Query(string, ...interface{}) (*sql.Rows, error)
	QueryRow(string, ...interface{}) *sql.Row
	Prepare(string) (*sql.Stmt, error)
	Exec(string, ...interface{}) (sql.Result, error)
}

type Scannable interface {
	Scan(...interface{}) error
}

func GenerateBindvars(num int) string {
	var args_id_slice []string
	for i := 1; i <= num; i++ {
		arg := fmt.Sprintf("$%d", i)
		args_id_slice = append(args_id_slice, arg)
	}
	return strings.Join(args_id_slice, ", ")
}

func FormatForInsert(fields []string) string {
	fStr := strings.Join(fields, ", ")
	v := GenerateBindvars(len(fields))
	return fmt.Sprintf("(%s) VALUES(%s)", fStr, v)
}

func FormatForUpdate(fields []string) string {
	var args_id_slice []string
	for i := 1; i <= len(fields); i++ {
		arg := fmt.Sprintf("%s=$%d", fields[i-1], i)
		args_id_slice = append(args_id_slice, arg)
	}
	return strings.Join(args_id_slice, ", ")

}
