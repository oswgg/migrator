package utils

func Contains(slice *[]string, item any) bool {
	for _, element := range *slice {
		if element == item {
			return true
		}
	}

	return false
}
