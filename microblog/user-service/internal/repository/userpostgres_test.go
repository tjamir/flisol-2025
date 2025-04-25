package repository

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/tjamir/flisol-2025/microblog/user-service/internal/model"
)

func TestPostgresUserRepository_CreateAndGetUser(t *testing.T) {
	db := NewPostgresDB()

	// inicia uma transação para rollback no final
	tx := db.Begin()
	defer tx.Rollback()

	repo := &PostgresUserRepository{db: tx}

	user := &model.User{
		ID:            uuid.New().String(),
		Username:     "testuser",
		Email:        "test@example.com",
		PasswordHash: "hashedpassword",
	}

	err := repo.CreateUser(user)
	require.NoError(t, err)

	// Testa GetUserByEmail
	gotUser, err := repo.GetUserByEmail(user.Email)
	require.NoError(t, err)
	assert.Equal(t, user.ID, gotUser.ID)
	assert.Equal(t, user.Username, gotUser.Username)

	// Testa GetUserByID
	gotUserByID, err := repo.GetUserByID(user.ID)
	require.NoError(t, err)
	assert.Equal(t, user.Email, gotUserByID.Email)
}
