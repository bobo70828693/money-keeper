package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

func Initial() (bool, *sql.DB) {

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", os.Getenv("MYSQL_USER"), os.Getenv("MYSQL_PASSWORD"), os.Getenv("MYSQL_URL"), os.Getenv("MYSQL_PORT"), os.Getenv("MYSQL_DATABASE")))
	if err != nil {
		log.Fatal(err)
		return false, nil
	}

	err = db.Ping()
	if err != nil {
		log.Fatal(err)
		return false, nil
	}

	return true, db
}

func Save(db *sql.DB, table string, saveData map[string]interface{}) {
	var query string
	columns := make([]string, 0, len(saveData))
	values := make([]interface{}, 0, len(saveData))

	for column, value := range saveData {
		columns = append(columns, column)
		values = append(values, value)
	}

	placeholders := make([]string, len(columns))
	for i := range placeholders {
		placeholders[i] = "?"
	}

	query = fmt.Sprintf("INSERT INTO `%s` (%s) VALUES (%s)", table, strings.Join(columns, ","), strings.Join(placeholders, ","))

	stmt, err := db.Prepare(query)
	if err != nil {
		log.Fatal(err)
		return
	}
	defer stmt.Close()

	result, err := stmt.Exec(values...)
	if err != nil {
		log.Fatal(err)
		return
	}

	fmt.Println(result)
}

func GetAll(db *sql.DB, table string, whereData map[string]interface{}) *sql.Rows {
	var query string

	columns := make([]string, 0)
	values := make([]interface{}, 0)

	for column, val := range whereData {
		columns = append(columns, column)
		values = append(values, val)
	}

	query = fmt.Sprintf("SELECT * FROM  `%s` WHERE ", table)

	for i, column := range columns {
		if i > 0 {
			query += " AND "
		}
		query += fmt.Sprintf("%s=?", column)
	}

	fmt.Println("Query : ", query)

	rows, err := db.Query(query, values...)
	if err != nil {
		log.Fatal(err)
	}
	return rows
}

func Close(db *sql.DB) {
	db.Close()
}
