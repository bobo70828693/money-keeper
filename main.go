package main

import (
	"fmt"
	"github.com/line/line-bot-sdk-go/v7/linebot"
	"log"
	"moneykeeper/services"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

//type CostLogData struct {
//	GroupId string `json:group_id`
//	UserId  string `json:user_id`
//	Name    string `json:name`
//	Comment string `json:comment`
//	Amount  int    `json:amount`
//}

var bot *linebot.Client

func main() {
	bot, err := linebot.New(os.Getenv("CHANNEL_SECRET_KEY"), os.Getenv("LINE_ACCESS_KEY"))
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("/callback", func(w http.ResponseWriter, req *http.Request) {
		events, err := bot.ParseRequest(req)
		if err != nil {
			if err == linebot.ErrInvalidSignature {
				w.WriteHeader(400)
			} else {
				w.WriteHeader(500)
			}
			return
		}

		for _, event := range events {
			if event.Type == linebot.EventTypeMessage {
				switch message := event.Message.(type) {
				case *linebot.TextMessage:
					if event.Source.Type == linebot.EventSourceTypeGroup {
						groupId := event.Source.GroupID
						userId := event.Source.UserID

						if validateMsg(message.Text) {
							costLogData := services.HandleMsg(groupId, userId, message.Text)
							costLogData.CreateCostLog()
						} else {
							if message.Text == "目前花費" {
								amountResult := services.GetCurrentMonthCostLogAmount(groupId)
								costResultStr := fmt.Sprintf("目前 %d 月份花費\n", time.Now().Month())

								for i, amountItem := range amountResult {
									costResultStr += fmt.Sprintf("[%s] 花費 %d", amountItem["name"], amountItem["amount"])
									if i != len(amountResult)-1 {
										costResultStr += "\n"
									}
								}

								if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(costResultStr)).Do(); err != nil {
									log.Print(err)
								}
							} else {
								if _, err = bot.ReplyMessage(event.ReplyToken, linebot.NewTextMessage("格式錯誤")).Do(); err != nil {
									log.Print(err)
								}
							}
						}
					}
					break
				}
			}
		}
	})

	http.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		w.Write([]byte("hello11"))
	})

	if err := http.ListenAndServe(":"+os.Getenv("PORT"), nil); err != nil {
		log.Fatal(err)
	}
}

func validateMsg(text string) bool {
	//-- {User} {FoodType} {Amount}
	strArray := strings.Split(text, " ")
	if len(strArray) != 3 {
		return false
	}

	_, err := strconv.Atoi(strings.Trim(strArray[2], " "))
	if err != nil {
		return false
	}

	return true
}
