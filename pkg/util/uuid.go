package util

import "github.com/google/uuid"

func ParseUUID(id string) uuid.UUID {
	uid, err := uuid.Parse(id)
	if err != nil {
		return uuid.Nil
	}

	return uid
}
