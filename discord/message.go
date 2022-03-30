package discord

import (
	"bytes"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type Message struct {
	Content          string           `json:"content"`
	Embeds           []any            `json:"embeds"`
	MessageReference MessageReference `json:"message_reference"`
}

type MessageReference struct {
	MessageId string `json:"message_id"`
	ChannelId string `json:"channel_id"`
	GuildId   string `json:"guild_id"`
}

type ReceivedMessage struct {
	MessageType       int       `json:"type"`
	TextToSpeech      bool      `json:"tts"`
	Timestamp         time.Time `json:"timestamp"`
	ReferencedMessage any       `json:"referenced_message"`
	Pinned            bool      `json:"pinned"`
	Nonce             string    `json:"nonce"`
	Mentions          []any     `json:"mentions"`
	MentionRoles      []any     `json:"mention_roles"`
	MentionEveryone   bool      `json:"mention_everyone"`
	Member            struct {
		Roles       []string  `json:"roles"`
		Nick        string    `json:"nick"`
		Mute        bool      `json:"mute"`
		JoinedAt    time.Time `json:"joined_at"`
		HoistedRole string    `json:"hoisted_role"`
		Deaf        bool      `json:"deaf"`
	} `json:"member"`
	ID              string `json:"id"`
	Flags           int    `json:"flags"`
	Embeds          []any  `json:"embeds"`
	EditedTimestamp any    `json:"edited_timestamp"`
	Content         string `json:"content"`
	Components      []any  `json:"components"`
	ChannelID       string `json:"channel_id"`
	Author          struct {
		Username      string `json:"username"`
		PublicFlags   int    `json:"public_flags"`
		ID            string `json:"id"`
		Discriminator string `json:"discriminator"`
		Avatar        string `json:"avatar"`
		Bot           bool   `json:"bot"`
	} `json:"author"`
	Attachments []any  `json:"attachments"`
	GuildID     string `json:"guild_id"`
	BotToken    string
}

type MessageStructure struct {
	Type      string          `json:"t"`
	S         int             `json:"s"`
	Operation int             `json:"op"`
	Payload   ReceivedMessage `json:"d"`
}

func (m *ReceivedMessage) Reply(message string) {
	msg := Message{
		Content: message,
		Embeds:  []any{},
		MessageReference: MessageReference{
			MessageId: m.ID,
			ChannelId: m.ChannelID,
			GuildId:   m.GuildID,
		},
	}
	p, err := json.Marshal(msg)
	if err != nil {
		log.Fatal(err.Error())
	} else {
		post, err := http.NewRequest("POST", "https://discord.com/api/channels/"+m.ChannelID+"/messages", bytes.NewBuffer(p))
		if err != nil {
			log.Fatal(err.Error())
		} else {
			log.Println("Token", m.BotToken)
			post.Header.Add("Authorization", "Bot "+m.BotToken)
			post.Header.Add("Content-Type", "application/json")
			post.Header.Add("Accept", "application/json")
			resp, err := http.DefaultClient.Do(post)
			if err != nil {
				log.Fatal(err.Error())
			} else {
				defer resp.Body.Close()
				body, _ := ioutil.ReadAll(resp.Body)
				log.Println("Response:", string(body))
				log.Println("Response:", resp.StatusCode)
				log.Println("Response:", resp.StatusCode)
			}
		}
	}

}
