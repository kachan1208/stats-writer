package model

import "time"

type Message struct {
	AccountID string
	Subject   string
	Body      []byte
}
type Counter struct {
	AccountID string
	ActionID  uint
	Count     uint
	CreatedAt time.Time
}
