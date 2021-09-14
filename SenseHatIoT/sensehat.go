package SenseHatIoT

/*
 * sensehat.go
 * A structure for Sense Hat temperature measurements!
 *
 * Principal author(s) : Abhishek Tiwary
 */

import (
	"time"
)

type Measurement struct {
	TimeObj     time.Time `json:"-"`
	Timestamp   string    `json:"timestamp"`
	Temperature float64   `json:"temperature"`
}
