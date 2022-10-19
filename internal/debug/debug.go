package debug

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

func PrintJson(data interface{}, format string, arg ...interface{}) {
	fmt.Println()
	fmt.Println(fmt.Sprintf(format, arg...))
	fmt.Println("------------------------------------------------------------")
	dataStr, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		log.Error(err)
		return
	}
	fmt.Println(string(dataStr))
	fmt.Println("------------------------------------------------------------")
	fmt.Println()
}
