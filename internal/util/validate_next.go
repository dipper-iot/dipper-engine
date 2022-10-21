package util

func ValidateNext(listNext []string) (listNew []string) {
	listNew = make([]string, 0)
	for _, next := range listNext {
		if next != "" {
			listNew = append(listNew, next)
		}
	}
	return
}
