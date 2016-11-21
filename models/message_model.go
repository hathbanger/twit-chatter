package models

import (
	"time"
	// "fmt"

	"labix.org/v2/mgo/bson"
	"github.com/go-chatter/store"
)

type Message struct {
	Id 		bson.ObjectId          `json:"id",bson:"_id,omitempty"`
	Timestamp 	time.Time	       `json:"time",bson:"time,omitempty"`
	Username	string           `json:"username",bson:"username,omitempty"`
	Body	string           `json:"body",bson:"body,omitempty"`
}



func NewMessage(username string, body string) *Message {
	m := new(Message)
	m.Id = bson.NewObjectId()
	m.Username = username
	m.Body = body

	return m
}


func (m *Message) SaveMessage() error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}


	collection := session.DB("test").C("messages")

	message := Message {
		Id:		m.Id,
		Timestamp:	m.Timestamp,
		Username:	m.Username,
		Body: 	m.Body,
	}


	err = collection.Insert(message)
	if err != nil {
		return err
	}

	return nil
}

