package models

import "time"

const (
	InviteStatusPending = "pending"
	InviteStatusUsed    = "used"
)

type InviteRecord struct {
	NodeAddress string    `json:"node_address"`
	InviteID    string    `json:"invite_id"`
	ExpireIn    time.Time `json:"expire_in"`
	Status      string    `json:"status"`
}
