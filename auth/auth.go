package auth

import (
	"github.com/google/uuid"
)

func GenUuid() string {
	uuid := uuid.New()
	return uuid.String()
}
