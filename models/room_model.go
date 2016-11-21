package models

import (
	"time"
	"fmt"

	"labix.org/v2/mgo/bson"
	"github.com/go-chatter/store"
)


type Room struct {
	//BaseModel
	Id 		bson.ObjectId          `json:"id",bson:"_id,omitempty"`
	Timestamp 	time.Time	       `json:"time",bson:"time,omitempty"`
	Name	string           `json:"name",bson:"name,omitempty"`
	Users []User
}

func NewRoom(name string) *Room {
	r := new(Room)
	r.Id = bson.NewObjectId()
	r.Name = name

	return r
}


func (r *Room) SaveRoom() error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}


	collection := session.DB("test").C("rooms")


	err = collection.Insert(&Room{
		Id: r.Id,
		Timestamp: r.Timestamp,
		Name: r.Name,
		Users: r.Users})
	if err != nil {
		return err
	}

	return nil
}

// func (r *Room) UserSendMessage(u User) error {
// 	r.Users = append(r.Users, u)

// 	fmt.Println("room:", r.Users)
// 	fmt.Println("User:", u)	
// }

func (r *Room) UserJoinRoom(u User) error {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}


	r.Users = append(r.Users, u)
	fmt.Println("User:", r.Users)	

	// Collection Stack
	collection := session.DB("test").C("rooms")

	// // Query
	query := bson.M{"name": "first"}
	update := bson.M{"$push": bson.M{"Users": u}}


	// Update
	err = collection.Update(query, update)
	if err != nil {
	    panic(err)
	}

	return nil
}

func FindRoom(name string) (Room, error) {
	session, err := store.ConnectToDb()
	defer session.Close()
	if err != nil {
		panic(err)
	}

	collection, err := store.ConnectToCollection(session, "rooms")
	if err != nil {
		panic(err)
	}

	room := Room{}
	err = collection.Find(bson.M{"name": name}).One(&room)
	if err != nil {
		return room, err
	}

	return room, err
} 



