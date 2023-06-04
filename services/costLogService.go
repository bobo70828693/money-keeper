package services

import (
	"moneykeeper/models"
	"moneykeeper/repositories/costLog"
	"strconv"
	"strings"
)

type CostLogData models.CostLogData
type amountItem map[string]interface{}

func (costLogData CostLogData) CreateCostLog() {
	data := map[string]interface{}{
		"group_id": costLogData.GroupId,
		"user_id":  costLogData.UserId,
		"name":     costLogData.Name,
		"comment":  costLogData.Comment,
		"amount":   costLogData.Amount,
	}

	costLog.Create(data)
}

func GetCurrentMonthCostLogAmount(groupId string) []amountItem {
	results := costLog.GetCurrentMonthCostCategory(map[string]interface{}{
		"group_id": groupId,
	})

	var amountList []amountItem
	for _, data := range results {
		amountList = append(amountList, amountItem{
			"name":   data.Name,
			"amount": data.Amount,
		})
	}

	return amountList
}

func HandleMsg(groupId string, userId string, text string) (data CostLogData) {
	strArray := strings.Split(text, " ")

	amount, _ := strconv.Atoi(strings.Trim(strArray[2], " "))

	return CostLogData{
		GroupId: groupId,
		UserId:  userId,
		Name:    strings.Trim(strArray[0], " "),
		Comment: strings.Trim(strArray[1], " "),
		Amount:  amount,
	}
}
