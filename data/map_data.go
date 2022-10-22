package data

import (
	"encoding/json"
	"time"
)

func MapToStruct(source map[string]interface{}, dest interface{}) error {
	dataByte, err := json.Marshal(source)
	if err != nil {
		return err
	}

	return json.Unmarshal(dataByte, dest)
}

func CreateOutput(input *InputEngine, id string) (output *OutputEngine) {

	timeData := time.Now()

	output = new(OutputEngine)
	output.BranchMain = input.BranchMain
	output.ChanId = input.ChanId
	output.IdNode = input.IdNode
	output.FromEngine = id
	output.SessionId = input.SessionId
	output.Time = &timeData
	output.Data = input.Data
	output.Type = TypeOutputEngineError

	return
}
