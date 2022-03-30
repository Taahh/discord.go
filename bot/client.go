package bot

import (
	"discord.go/discord"
	"discord.go/opcodes"
	"discord.go/presence"
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"math/rand"
	"net/url"
	"reflect"
	"time"
)

var DiscordGateway = url.URL{Scheme: "wss", Host: "gateway.discord.gg", Path: "?v=9&encoding=json"}

type Client struct {
	T         string `json:"t"`
	S         int64  `json:"s"`
	Operation int    `json:"op"`
	Payload   struct {
		V            int          `json:"v"`
		UserSettings any          `json:"user_settings"`
		User         discord.User `json:"user"`
	} `json:"d"`
	authenticated bool
	listeners     map[string][]Event /*func(value any)*/

}

func (c Client) GetUser() discord.User {
	return c.Payload.User
}

func (c Client) GetFullIdentifier() string {
	return c.Payload.User.Username + "#" + c.Payload.User.Discriminator
}

func (c *Client) Login(token string, callback func(client Client)) {
	var dialer = websocket.Dialer{
		ReadBufferSize:  4096,
		WriteBufferSize: 4096,
	}
	conn, _, err := dialer.Dial(DiscordGateway.String(), nil)
	if err != nil {
		log.Fatal(err)
		return
	}
	go func() {
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				fmt.Println(err)
				continue
			}

			log.Printf("Message Type %v Data %v\n", messageType, string(p))

			switch messageType {
			case 1:
				data := opcodes.OpcodeTen{}
				err := json.Unmarshal(p, &data)
				if err != nil {
					log.Fatal(err.Error())
				} else {
					fmt.Println("Data:", data)
					random := rand.Float64()
					var delay = float64(data.Heartbeat.HeartbeatInterval) * random
					time.AfterFunc(time.Duration(delay), func() {
						switch data.Operation {
						case 10:
							opcode := opcodes.OpcodeOne{Op: 1, D: nil}
							send(conn, opcode, func() {
								go func() {
									for _ = range time.Tick(time.Duration(data.Heartbeat.HeartbeatInterval) * time.Millisecond) {
										heartbeat(conn)
									}
								}()
							})
							break
						case 11:
							if c.authenticated {
								break
							}
							opcode := opcodes.OpcodeTwo{
								Op: 2,
								Identity: opcodes.Identify{
									Token:      token,
									Intents:    (1 << 0) + (1 << 1) + (1 << 2) + (1 << 9),
									Properties: opcodes.IdentifyConnection{OperatingSystem: "windows", Device: "discord.go", Browser: "discord.go"},
									Presence:   presence.UpdatePresence{},
								},
							}
							send(conn, opcode, func() {
								c.authenticated = true
								log.Println("Authenticated client")
								send(conn, opcodes.OpcodeThree{
									Op: 3,
									Presence: presence.UpdatePresence{Since: time.Now().Unix(), Status: presence.Idle, Afk: false, Activities: []presence.Activity{
										{Name: "The Vampire Diaries on Netflix", Type: presence.Game, CreatedAt: time.Now().Unix()},
									}},
								}, func() {
									heartbeat(conn)
								})
							})
							break
						case 0:
							opcode := opcodes.OpcodeZero{}
							err := json.Unmarshal(p, &opcode)
							if err != nil {
								log.Fatal(err)
							} else {
								switch opcode.Type {
								case "READY":
									client := Client{}
									err := json.Unmarshal(p, &client)
									if err != nil {
										log.Fatal(err)
									} else {
										c.Payload = client.Payload
										callback(client)
									}
									break
								case "MESSAGE_CREATE":
									message := discord.MessageStructure{}
									err := json.Unmarshal(p, &message)
									if err != nil {
										log.Fatal(err)
									} else {
										log.Printf("INCOMING MESSAGE FROM %v SAYING %v", message.Payload.Author.Username, message.Payload.Content)
										payload := &message.Payload
										payload.BotToken = token
										log.Println("set token", payload.BotToken)
										c.Emit(opcode.Type, &message.Payload)
									}
								}

							}
						default:
							log.Println("Not implemented")
						}
					})
				}
			}
		}
	}()
}

type Event func(value any)

func (c *Client) RetrieveUser(userId string) {

}

func (c *Client) On(event string, callback Event) {
	if c.listeners == nil {
		c.listeners = make(map[string][]Event /*func(value any)*/)
	}
	func1 := reflect.ValueOf(callback)
	if _, ok := c.listeners[event]; ok {
		for _, cb := range c.listeners[event] {
			func2 := reflect.ValueOf(cb)
			if func1.Pointer() == func2.Pointer() {
				return
			}
		}
		c.listeners[event] = append(c.listeners[event], callback)
	} else {
		c.listeners[event] = []Event /*func(value any)*/ {callback}
	}
}

func (c Client) Emit(event string, p any) {
	if _, ok := c.listeners[event]; ok {
		for _, callback := range c.listeners[event] {
			callback(p)
		}
	}
}

func heartbeat(c *websocket.Conn) {
	opcode := opcodes.OpcodeOne{Op: 1, D: nil}
	log.Println("Writing", opcode)
	err := c.WriteJSON(opcode)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Heartbeat")
	}
}

func send(c *websocket.Conn, v any, callback func()) {
	log.Println("Writing", v)
	err := c.WriteJSON(v)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Wrote JSON")
		callback()
	}
}
