package util

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestGenerateID(t *testing.T) {
	id, err := GenerateID()
	t.Logf("Generated ID: %d", id)
	id_2, err := GenerateID()
	t.Logf("Generated ID: %d", id_2)
	require.NoError(t, err)
	require.NotEmpty(t, id)
	require.NotEmpty(t, id_2)
	require.NotEqual(t, id, id_2, "IDs should be unique")
}
