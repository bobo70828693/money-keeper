package costLog

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/now"
	"log"
	"moneykeeper/db"
	"moneykeeper/models"
)

func Create(createData map[string]interface{}) {
	_, dbq := db.Initial()
	defer db.Close(dbq)
	db.Save(dbq, "cost_log", createData)
}

func GetCurrentMonthCostCategory(whereData map[string]interface{}) []models.CostLogData {
	_, dbq := db.Initial()
	defer db.Close(dbq)

	dateStart := now.BeginningOfMonth().Format("2006-01-02") + " 00:00:00"
	dateEnd := now.EndOfMonth().Format("2006-01-02") + " 23:59:59"

	fmt.Println(dateStart, dateEnd)

	columns := make([]string, 0)
	values := make([]interface{}, 0)

	values = append(values, dateStart, dateEnd)
	for column, val := range whereData {
		columns = append(columns, column)
		values = append(values, val)
	}

	query := "SELECT `name`, SUM(`amount`) FROM `cost_log` WHERE created_at >= ? AND created_at <= ?"

	for _, column := range columns {
		query += fmt.Sprintf(" AND %s=?", column)
	}

	query += " GROUP BY `name`"

	rows, err := dbq.Query(query, values...)
	if err != nil {
		log.Fatal(err)
	}

	results := make([]models.CostLogData, 0)

	for rows.Next() {
		var costLog models.CostLogData

		err := rows.Scan(&costLog.Name, &costLog.Amount)

		if err != nil {
			if err == sql.ErrNoRows {
				return nil
			}
			log.Fatal(err)
		}

		results = append(results, costLog)
	}

	return results
}
