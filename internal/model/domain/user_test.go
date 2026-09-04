package domain_test

import (
	"AlfianChabib/go-job-board-api/internal/model/domain"
	"sync"
	"testing"

	"gorm.io/gorm/schema"
)

func TestUserSchema(t *testing.T) {
	s, err := schema.Parse(&domain.User{}, &sync.Map{}, schema.NamingStrategy{})
	if err != nil {
		t.Fatalf("Failed to parse User schema: %v", err)
	}

	rel, exists := s.Relationships.Relations["Auth"]
	if !exists || rel == nil {
		t.Errorf("Expected Auth relationship on User struct")
	}
}
