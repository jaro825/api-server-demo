package types

import (
	"encoding/json"
	"github.com/google/uuid"
)

type User struct {
	ID   uuid.UUID `json:"id"`
	Name string    `json:"name"`
}

func (u User) MarshalBinary() ([]byte, error) {
	return json.Marshal(u)
}
