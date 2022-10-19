package util

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

func ReadFile(data interface{}, file string) error {
	// Open our jsonFile
	jsonFile, err := os.Open(file)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}
	// defer the closing of our jsonFile so that we can parse it later on
	defer jsonFile.Close()

	// read our opened xmlFile as a byte array.
	byteValue, err := ioutil.ReadAll(jsonFile)
	// if we os.Open returns an error then handle it
	if err != nil {
		return err
	}

	// we unmarshal our byteArray which contains our
	// jsonFile's content into 'users' which we defined above

	return json.Unmarshal(byteValue, data)
}

func SaveFile(fileName string, data []byte) error {

	file, err := os.Create(fileName)
	if err != nil {
		return err
	}

	file.Write(data)
	file.Close()

	return nil
}
