package models

import (
	"log"

	"github.com/containerops/wrench/db"
)

func GetRandomOneFromSet(setName string) (string, error) {
	return db.Client.SRandMember(setName).Result()
}

func GetMsgFromList(listName string, start, end int64) ([]string, error) {
	return db.Client.LRange(listName, start, end).Result()
}

func SaveMsgToSet(setName, msg string) error {
	_, err := db.Client.SAdd(setName, msg).Result()
	return err
}

func PushMsgToList(listName, msg string) error {
	_, err := db.Client.RPush(listName, msg).Result()
	return err
}

func SubscribeChannel(channelName string) chan string {
	msgChan := make(chan string, 30)
	go receiveMsgFromChannel(channelName, msgChan)
	return msgChan
}

func receiveMsgFromChannel(channelName string, msgChan chan string) {
	psc, err := db.Client.Subscribe(channelName)
	if err != nil {
		log.Fatal(err)
	}

	var isEnd = false
	for {
		msg, err := psc.ReceiveMessage()
		if err != nil {
			log.Fatal(err)
		}

		msgChan <- msg.Payload
		if msg.String() == "bye" {
			isEnd = true
		}
		if isEnd {
			break
		}
	}
}

func PublishMsg(channelName, msg string) error {
	_, err := db.Client.Publish(channelName, msg).Result()
	return err
}

func WrapListMsg(listName string) error {
	_, err := db.Client.LTrim(listName, int64(0), int64(0)).Result()
	if err != nil {
		return err
	}
	_, err = db.Client.LPop(listName).Result()
	return err
}

func MoveFromListByValue(listName, value string, count int64) error {
	_, err := db.Client.LRem(listName, count, value).Result()
	return err
}

func GetListLength(listName string) (int64, error) {
	return db.Client.LLen(listName).Result()
}
