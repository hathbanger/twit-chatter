package server

import (
	"fmt"
	"log"
	"flag"
	// "fmt"
	// "log"
	"os"
	"os/signal"
	"syscall"

	"github.com/coreos/pkg/flagutil"
	"github.com/dghubble/go-twitter/twitter"
	"github.com/dghubble/oauth1"

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

func twitChatter() websocket.Handler {

	var clientMessage string
	return websocket.Handler(func(ws *websocket.Conn) {
	flags := flag.NewFlagSet("user-auth", flag.ExitOnError)
	consumerKey := flags.String("consumer-key", "letM9qyusLEQLqHn9uz4AKltJ", "Twitter Consumer Key")
	consumerSecret := flags.String("consumer-secret", "Lexfs85yTEo2IXJUpHd8CZgWgdvx3J2HFanlSXLRKsRjgGUFXU", "Twitter Consumer Secret")
	accessToken := flags.String("access-token", "28226407-y4Hvy6ftSJh15ebWImKNYd8oZyZQTweT81vn1eW6x", "Twitter Access Token")
	accessSecret := flags.String("access-secret", "yb4NV7CiY3UV2FtQ6p90N1TSjWObE4AIbuKtbA3SvAxcx", "Twitter Access Secret")
	flags.Parse(os.Args[1:])
	flagutil.SetFlagsFromEnv(flags, "TWITTER")

	if *consumerKey == "" || *consumerSecret == "" || *accessToken == "" || *accessSecret == "" {
		log.Fatal("Consumer key/secret and Access token/secret required")
	}

	config := oauth1.NewConfig(*consumerKey, *consumerSecret)
	token := oauth1.NewToken(*accessToken, *accessSecret)
	// OAuth1 http.Client will automatically authorize Requests


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


			httpClient := config.Client(oauth1.NoContext, token)

			// Twitter Client
			client2 := twitter.NewClient(httpClient)

			// Convenience Demux demultiplexed stream messages
			demux := twitter.NewSwitchDemux()
			demux.Tweet = func(tweet *twitter.Tweet) {
				fmt.Println(tweet.Text)
					websocket.Message.Send(ws, tweet.Text)
			}
			demux.DM = func(dm *twitter.DirectMessage) {
				fmt.Println(dm.SenderID)
			}
			demux.Event = func(event *twitter.Event) {
				fmt.Printf("%#v\n", event)
			}

			fmt.Println("Starting Stream...")

			// FILTER
			filterParams := &twitter.StreamFilterParams{
				Track:         []string{msg},
				StallWarnings: twitter.Bool(true),
			}
			stream, err := client2.Streams.Filter(filterParams)
			if err != nil {
				log.Fatal(err)
			}

			// Receive messages until stopped or stream quits
			go demux.HandleChan(stream.Messages)

			// Wait for SIGINT and SIGTERM (HIT CTRL-C)
			ch := make(chan os.Signal)
			signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM)
			log.Println(<-ch)

			fmt.Println("Stopping Stream...")
			stream.Stop()	

		}
	
	})


}
