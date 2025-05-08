package utils

import (
	"time"

	"github.com/google/uuid"
)

func NewUUID() string {
	return uuid.New().String()
}

func GenerateTimestamp() time.Time {
	return time.Now()
}