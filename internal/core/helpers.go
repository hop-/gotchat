package core

import (
	"fmt"

	"github.com/google/uuid"
)

func generateUuid() string {
	id, err := uuid.NewRandom()
	if err != nil {
		fmt.Println("Error generating UUID:", err)
		return ""
	}
	fmt.Println("Generated UUID:", id.String())
	return id.String()
}
