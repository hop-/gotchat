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
func TestIsFieldExist(t *testing.T) {
	tests := []struct {
		name     string
		field    string
		entity   string
		expected bool
	}{
		{"User existing field", "Name", "User", true},
		{"User existing field with tag", "UniqueId", "User", true},
		{"User inherited field", "Id", "User", true},
		{"User non-existing field", "NonExistentField", "User", false},
		{"Account existing field", "Password", "Account", true},
		{"Account inherited field", "Id", "Account", true},
		{"Account non-existing field", "InvalidField", "Account", false},
		{"Message existing field", "Content", "Message", true},
		{"Message non-existing field", "InvalidContent", "Message", false},
		{"Channel existing field", "Name", "Channel", true},
		{"ConnectionDetails existing field", "EncryptionKey", "ConnectionDetails", true},
		{"Attendance existing field", "JoinedAt", "Attendance", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var result bool
			switch tt.entity {
			case "User":
				result = IsFieldExist[User](tt.field)
			case "Account":
				result = IsFieldExist[Account](tt.field)
			case "Message":
				result = IsFieldExist[Message](tt.field)
			case "Channel":
				result = IsFieldExist[Channel](tt.field)
			case "ConnectionDetails":
				result = IsFieldExist[ConnectionDetails](tt.field)
			case "Attendance":
				result = IsFieldExist[Attendance](tt.field)
			}

			if result != tt.expected {
				t.Errorf("IsFieldExist[%s](%q) = %v, want %v", tt.entity, tt.field, result, tt.expected)
			}
		})
	}
}
