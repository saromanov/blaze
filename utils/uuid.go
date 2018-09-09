package utils

import "github.com/satori/go.uuid"

func GetUUID() (string, error) {
	return uuid.NewV4()
}