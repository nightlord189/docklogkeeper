package usecase

func (u *Usecase) convertToShortNames(arr []string) {
	for i := range arr {
		arr[i] = u.Log.GetShortContainerName(arr[i])
	}
}

func arrToMap[T comparable](arr []T) map[T]bool {
	result := make(map[T]bool, len(arr))
	for _, key := range arr {
		result[key] = true
	}
	return result
}
