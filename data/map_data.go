package data

import "encoding/json"

func MapToStruct(source map[string]interface{}, dest interface{}) error {
	dataByte, err := json.Marshal(source)
	if err != nil {
		return err
	}

	return json.Unmarshal(dataByte, dest)
}
