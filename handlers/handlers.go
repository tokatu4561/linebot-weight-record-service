package handlers

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"tokatu4561/line-bot/record-service/line"
	"tokatu4561/line-bot/record-service/models"

	_ "github.com/lib/pq"
	"github.com/line/line-bot-sdk-go/linebot"
)

func WeightRegist(w http.ResponseWriter, r *http.Request) {
	line, err := line.LineConnection()
	if err != nil {
		log.Fatalln(err)
	}

	lineEvents, err := line.Client.ParseRequest(r)
	if err != nil {
		log.Fatalln(err)
	}

	for _, event := range lineEvents {
		// イベントがメッセージの受信だった場合
		if event.Type == linebot.EventTypeMessage {
			switch message := event.Message.(type) {		
			case *linebot.TextMessage:
				log.Println(message)
				// replyMessage := message.Text
				err = recordWeight(line, event)
				if err != nil {
					log.Fatalln(err)
				}
			case *linebot.LocationMessage:
				break		
			default:
			}
		}
	}
}

func recordWeight(line *line.Line, event *linebot.Event) error {
	lineID := event.Source.UserID

	// DB接続
	var connectionString string = fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable", 
	os.Getenv("DB_USER"),
	os.Getenv("DB_DBNAME"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_HOST"),
	os.Getenv("DB_PORT"))
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return err
	}
	defer db.Close()
	
	DBModel := &models.DBModel{
		DB: db,
	}

	user, err := DBModel.GetOneUser(lineID)
	if err != nil {
		// ユーザーが存在しなかった場合新規作成
		err = DBModel.AddUser(lineID)
		if err != nil {
			return err
		}	
		user, _ = DBModel.GetOneUser(lineID)
	}

	// 以降体重を記録し、メッセージを返信
	weight, err := strconv.ParseFloat(event.Message.(*linebot.TextMessage).Text, 64)
	if err != nil {
		return err
	}

	lastWeight, _ := DBModel.GetLatestWeight(user.ID)
	err = DBModel.AddWeightRecord(user.ID, weight)
	if err != nil {
		return err
	}
	minWeight, err := DBModel.GetMinWeight(user.ID)
	if err != nil {
		return err
	}
	
	replyMessage := getReplyMessageOfWeight(weight, minWeight, lastWeight)
	
	_, err = line.Client.ReplyMessage(event.ReplyToken, linebot.NewTextMessage(replyMessage)).Do()

	return err
}

// 引数の体重からユーザーに返信するリプライメッセージを返す
func getReplyMessageOfWeight(currentWeight float64, minWeight float64, lastWeight float64) string {
	var replyMessage string

	// 前回との比較　
	// TODO:　異常値はバリデーションしたい(前日から過度に体重増えているなど)
	diffWeight := currentWeight - lastWeight;
	log.Println(diffWeight)

	if diffWeight == 0 {
        replyMessage = fmt.Sprintf("%.1fkg! 前回と同じ!", currentWeight)
	} else if diffWeight > 0 {
		replyMessage =  fmt.Sprintf("%.1fkg!\n%.1fkg太ったよ。頑張ろう!\nちなみに最も痩せていた体重は%.1fkg", 
									currentWeight, diffWeight, minWeight)	
	} else {
		diffWeight = 0 - diffWeight
        replyMessage =  fmt.Sprintf("%.1fkg!\n%.1fkg痩せたよ!お疲れ!\nちなみに最も痩せていた体重は%.1fkg",
									currentWeight, diffWeight, minWeight)	
	}

	return replyMessage
}