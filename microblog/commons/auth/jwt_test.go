package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateAndValidateToken(t *testing.T) {
	userID := "test-user-id"

	// Gera o token
	token, err := GenerateToken(userID)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// Valida o token
	validatedUserID, err := ValidateToken(token)
	assert.NoError(t, err)
	assert.Equal(t, userID, validatedUserID)
}

func TestValidateInvalidToken(t *testing.T) {
	// Token inválido
	_, err := ValidateToken("invalid.token.here")
	assert.Error(t, err)
}
