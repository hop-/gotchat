package core

import (
	"testing"
)

func TestBaseEntity_GetId(t *testing.T) {
	entity := BaseEntity{Id: 42}
	if entity.GetId() != 42 {
		t.Errorf("Expected Id to be 42, got %d", entity.GetId())
	}
}

func TestGetFieldNamesOfEntity(t *testing.T) {
	fieldNames := GetFieldNamesOfEntity[User]()
	expected := []string{"unique_id", "name", "password", "last_login"}
	if len(fieldNames) != len(expected) {
		t.Errorf("Expected %d field names, got %d", len(expected), len(fieldNames))
	}
	for _, fieldName := range fieldNames {
		found := false
		for _, name := range expected {
			if fieldName == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected field name '%s' not found in %v", fieldName, expected)
		}
	}
}
