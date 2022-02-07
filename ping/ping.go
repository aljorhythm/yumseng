package ping

import "time"

type PingResponse struct {
	Tag             string    `json:"tag"`
	ServerStartTime time.Time `json:"server_start_time"`
}
