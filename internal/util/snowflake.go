package util

import (
	"github.com/sony/sonyflake"
)

func getMachineID() (uint16, error) {
	return 11, nil
}

var sn = sonyflake.NewSonyflake(sonyflake.Settings{
	MachineID: getMachineID,
})

func NextID() (uint64, error) {
	return sn.NextID()
}
