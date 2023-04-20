package util

func DataToValue(data map[string]interface{}, mainBranch string) map[string]interface{} {
	if mainBranch == "" {
		mainBranch = "default"
	}
	result := map[string]interface{}{
		"meta_data": data,
	}

	dataBranch, ok := data[mainBranch]
	if ok {

		mapData, ok := dataBranch.(map[string]interface{})
		if ok {
			for key, val := range mapData {
				result[key] = val
			}
		}
	}

	return result
}

func ValueToData(data map[string]interface{}, mainBranch string) map[string]interface{} {
	if mainBranch == "" {
		mainBranch = "default"
	}
	result := map[string]interface{}{}

	mainData := map[string]interface{}{}
	for key, value := range data {
		if key == "meta_data" {
			mapData, ok := value.(map[string]interface{})
			if ok {
				for key, val := range mapData {
					result[key] = val
				}
			}
			continue
		}
		mainData[key] = value
	}

	result[mainBranch] = mainData
	return result
}
