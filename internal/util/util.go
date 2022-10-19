package util

import (
	"encoding/json"
	log "github.com/sirupsen/logrus"
)

func ConvertToByte(input interface{}) (dataByte []byte, err error) {
	dataByte = make([]byte, 0)
	switch data := input.(type) {
	case []byte:
		dataByte = data
	case string:
		dataByte = []byte(data)
	default:
		dataByte, err = json.Marshal(input)
		if err != nil {
			log.Error(err)
			return
		}
	}
	return
}
