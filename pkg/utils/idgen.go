package utils

var startId = 444

func GenerateNextId() int {
	startId++
	return startId
}
