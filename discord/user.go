package discord

type User struct {
	Verified      bool   `json:"verified"`
	Username      string `json:"username"`
	MfaEnabled    bool   `json:"mfa_enabled"`
	Id            string `json:"id"`
	Discriminator string `json:"discriminator"`
}
