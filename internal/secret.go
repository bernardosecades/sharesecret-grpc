package sharesecret

import "time"

type Secret struct {
	ID        string
	Content   string
	CustomPwd bool
	CreatedAt time.Time
	ExpiredAt time.Time
}
