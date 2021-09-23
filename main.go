package main

import (
	"encoding/json"
	//"encoding/json"
	"fmt"
	stomp "github.com/drawdy/stomp-ws-go"
	"github.com/gorilla/websocket"
	"log"
	"net/url"
	"strconv"
)

var exampid = "gosend: "


func say(conn stomp.STOMPConnector) {
	sh := stomp.Headers{
		stomp.HK_DESTINATION, "/user/subscribe/room/live_1",
		"liveId", "live_1",
		stomp.HK_ID, stomp.Uuid()}
	cs, e := conn.Subscribe(sh)
	if e != nil {
		log.Fatal(e)
	}

	md := <-cs
	if md.Error != nil {
		log.Fatalf("receive greeting message caught error: %v", md.Error)
	}
	fmt.Printf("----> receive new message: %v\n", md.Message.BodyString())
}

func send(conn stomp.STOMPConnector, loop int){
	mapD := map[string]string{
		"liveId": "live_1",
		"userId": "user_"  + strconv.Itoa(loop),
		"deviceId": "device_" + strconv.Itoa(loop),
		"type": "JOIN"}
	mapB, _ := json.Marshal(mapD)
	message := string(mapB)
	// Suppress content length here, so JMS will treat this as a 'text' message.
	s := stomp.Headers{
		stomp.HK_DESTINATION, "/publish/room/join/live_1",
		"liveId", "live_1",
		stomp.HK_HOST, "server",
		stomp.HK_ID, stomp.Uuid(),
		stomp.HK_CONTENT_LENGTH, strconv.Itoa(len(message)),
		stomp.HK_CONTENT_TYPE, "application/json",
	}
	e := conn.Send(s, message)

	if e != nil {
		log.Fatalln(e) // Handle this ...
	}
	fmt.Println(exampid, "send complete:", string(mapB))
}

// Connect to a STOMP 1.1 broker, send some messages and disconnect.
func main() {
	fmt.Println(exampid + "starts ...")

	loop := 1
	for loop < 50000 {
		url_c := url.URL{
			Scheme: "ws",
			Host:   "server:8060",
			Path:   "/ws",
		}

		//duration := time.Duration(50)*time.Millisecond
		//time.Sleep(duration)

		// Open a net connection
		n, resp, e := websocket.DefaultDialer.Dial(url_c.String(), nil)
		if e != nil {
			log.Fatalln(e) // Handle this ......
		}
		fmt.Println("response status" + resp.Status)
		fmt.Println(exampid + "dial complete ...")

		// Connect to broker
		eh := stomp.Headers{
			stomp.HK_ACCEPT_VERSION, "1.2,1.1,1.0",
			stomp.HK_HOST, "server",
			stomp.HK_HEART_BEAT, "15000,15000",
			"liveId", "live_1",
			"deviceId", "device_" + strconv.Itoa(loop),
			"userId", "user_"  + strconv.Itoa(loop)}
		conn, e := stomp.ConnectOverWS(n, eh)
		if e != nil {
			log.Fatalln(e) // Handle this ......
		}
		fmt.Println(exampid + "stomp connect complete ...")

		go say(conn)
		go send(conn, loop)

		loop += 1
	}


	//uh := stomp.Headers{"destination", "subscribe/room/room_1", "liveId", "live_1", "deviceId", "device_1", "userId", "user_1"} // Unsubscribe headers
	//e = conn.Unsubscribe(uh)
	//if e != nil {
	//	log.Fatal(e)
	//}
	//
	//
	//// Disconnect from the Stomp server
	//eh = stomp.Headers{}
	//e = conn.Disconnect(eh)
	//if e != nil {
	//	log.Fatalln(e) // Handle this ......
	//}
	//fmt.Println(exampid + "stomp disconnect complete ...")
	//// Close the network connection
	//e = n.Close()
	//if e != nil {
	//	log.Fatalln(e) // Handle this ......
	//}
	//fmt.Println(exampid + "network close complete ...")
	//
	//fmt.Println(exampid + "ends ...")
}