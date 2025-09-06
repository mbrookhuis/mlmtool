package uuid

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGenerateUniqueID(t *testing.T) {
	// Unit Test Case 1: Check if the ID is non-empty
	id := GenerateUniqueID()
	assert.NotEmpty(t, id, "ID should not be empty")

	// Unit Test Case 2: Check if the ID is a valid UUID
	_, err := uuid.Parse(id)
	require.NoError(t, err, "ID should be a valid UUID")

	// Unit Test Case 3: Check uniqueness
	id2 := GenerateUniqueID()
	assert.NotEqual(t, id, id2, "Each generated ID should be unique")

	// Unit Test Case 4: Length check (UUIDs are 36 chars long)
	assert.Len(t, id, 36, "UUID length should be 36 characters")

	// Unit Test Case 5: Check if it starts with hex digits (UUID v4)
	assert.Regexp(t, "^[a-fA-F0-9]{8}-", id, "UUID should start with hex digits")
}
