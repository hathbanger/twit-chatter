package server

import (
	"net/http"

	"time"
	"fmt"
	"encoding/json"
	"io/ioutil"
	"github.com/go-chatter/models"
	"github.com/labstack/echo"

	"labix.org/v2/mgo/bson"
	log "github.com/cihub/seelog"	
)

func CreateRoom(c echo.Context) error {
	method := c.Request().Method()
	uri := c.Request().URI()
	log.Debugf("%s %s", method, uri)

	json_body, err := ioutil.ReadAll(c.Request().Body())
	room := models.Room{}
	err = json.Unmarshal(json_body, &room)
	if err != nil {
		fmt.Println(err)
	}

	room.Timestamp = time.Now()
	room.Id = bson.NewObjectId()

	err = room.SaveRoom()
	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, room)
}

func GetRoom(c echo.Context) error {
	name := c.FormValue("name")

	room, err := models.FindRoom(name)
	if err != nil {
		panic(err)
	}

	if room.Id != "" /*&& room.roomname != "" */ {
		return c.JSON(http.StatusOK, room)
	} else {
		return c.JSON(http.StatusNotFound, "not found")
	}
}

func JoinRoom(c echo.Context) error {
	name := c.FormValue("name")
	username := c.FormValue("username")

	room, err := models.FindRoom(name)
	if err != nil {
		panic(err)
	}
	user, err := models.FindUser(username)
	if err != nil {
		panic(err)
	}
	
	room.UserJoinRoom(user)
	
	fmt.Println("room:", room.Users)

	if room.Id != "" /*&& room.roomname != "" */ {
		return c.JSON(http.StatusOK, room)
	} else {
		return c.JSON(http.StatusNotFound, "not found")
	}
}
