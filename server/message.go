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

func CreateMessage(c echo.Context) error {
	method := c.Request().Method()
	uri := c.Request().URI()
	log.Debugf("%s %s", method, uri)

	json_body, err := ioutil.ReadAll(c.Request().Body())
	message := models.Message{}
	err = json.Unmarshal(json_body, &message)
	if err != nil {
		fmt.Println(err)
	}

	message.Timestamp = time.Now()
	message.Id = bson.NewObjectId()

	err = message.SaveMessage()
	if err != nil {
		fmt.Println(err)
	}

	return c.JSON(http.StatusOK, message)
}

