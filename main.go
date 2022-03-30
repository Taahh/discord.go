package main

import (
	"discord.go/bot"
	"discord.go/discord"
	"log"
	"unsafe"
)

func main() {
	client := bot.Client{}
	client.Login("token", func(c bot.Client) {
		log.Println("Logged in as", c.GetFullIdentifier())
	})
	for {
		client.On("MESSAGE_CREATE", func(value any) {
			message := **(**discord.ReceivedMessage)(unsafe.Pointer(uintptr(unsafe.Pointer(&value)) + unsafe.Sizeof(uint(0))))
			if message.Author.Bot {
				return
			}
			log.Println("Message:", message.Content)
			message.Reply("Hi")
		})
	}
}
