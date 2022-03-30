package opcodes

import "discord.go/presence"

type Heartbeat struct {
	Trace []string `json:"trace"`
}

type Identify struct {
	Token      string                  `json:"token"`
	Intents    int                     `json:"intents"`
	Properties IdentifyConnection      `json:"properties"`
	Presence   presence.UpdatePresence `json:"presence"`
}

type IdentifyConnection struct {
	OperatingSystem string `json:"$os"`
	Browser         string `json:"$browser"`
	Device          string `json:"$device"`
}

type OpcodeZero struct {
	Type string `json:"t"`
	Op   int    `json:"op"`
	D    any    `json:"d"`
}

type OpcodeOne struct {
	Op int `json:"op"`
	D  any `json:"d"`
}

type OpcodeTwo struct {
	Op       int      `json:"op"`
	Identity Identify `json:"d"`
}

type OpcodeThree struct {
	Op       int                     `json:"op"`
	Presence presence.UpdatePresence `json:"d"`
}

type OpcodeTen struct {
	Operation int `json:"op"`
	Heartbeat struct {
		HeartbeatInterval int `json:"heartbeat_interval"`
	} `json:"d"`
}
