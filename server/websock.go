package server

import (
	"fmt"
	"log"

	"golang.org/x/net/websocket"
)


var (
	Msg       = websocket.Message
	ActiveClients = make(map[ClientConn]int) // map containing clients
)

type ClientConn struct {
	websocket *websocket.Conn
	clientIP  string
}

func hello() websocket.Handler {

	var clientMessage string
	return websocket.Handler(func(ws *websocket.Conn) {
		for {
			// Read
			msg := ""
			err := websocket.Message.Receive(ws, &msg)
			// websocket.Message.Send(ws, &msg)
			
			if err != nil {
				log.Fatal(err)
			}
			fmt.Printf("%s\n", msg)	
				
			client := ws.Request().RemoteAddr

			sockCli := ClientConn{ws, client}
			ActiveClients[sockCli] = 0
			log.Println("Number of clients connected ...", len(ActiveClients))			
			clientMessage = sockCli.clientIP + " Said: " + msg
			for cs, _ := range ActiveClients {
				if err = Msg.Send(cs.websocket, clientMessage); err != nil {
					// we could not send the message to a peer
					log.Println("Could not send message to ", cs.clientIP, err.Error())
				}
			}		
		}
		// for {
		// 	// Write
		// 	err := websocket.Message.Send(ws, "Hello, Client!")
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}

		// 	// Read
		// 	msg := ""
		// 	err = websocket.Message.Receive(ws, &msg)
		// 	if err != nil {
		// 		log.Fatal(err)
		// 	}
		// 	fmt.Printf("%s\n", msg)
		// }		
	})
}
