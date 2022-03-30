package presence

type UpdatePresence struct {
	Since      int64      `json:"since"`
	Activities []Activity `json:"activities"`
	Status     Presence   `json:"status"`
	Afk        bool       `json:"afk"`
}

type Activity struct {
	Name      string       `json:"name"`
	Type      ActivityType `json:"type"`
	CreatedAt int64        `json:"created_at"`
}

type Presence string
type ActivityType int64

const (
	Online       Presence = "online"
	DoNotDisturb          = "dnd"
	Idle                  = "idle"
	Invisible             = "invisible"
	Offline               = "offline"
)

const (
	Game      ActivityType = 0
	Streaming              = 1
	Listening              = 2
	Watching               = 3
	Custom                 = 4
	Competing              = 5
)
