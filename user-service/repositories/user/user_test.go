package repositories

import (
	"context"
	"regexp"
	"testing"
	"user-service/domain/dto"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func TestUserRepository(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}
	defer db.Close()

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db}), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %s", err)
	}

	repo := NewUserRepository(gormDB)
	ctx := context.Background()

	t.Run("FindByUsername Success", func(t *testing.T) {
		username := "testuser"
		rows := sqlmock.NewRows([]string{"id", "username"}).AddRow(1, username)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE username = $1`)).
			WithArgs(username, 1).WillReturnRows(rows)

		user, err := repo.FindByUsername(ctx, username)
		assert.NoError(t, err)
		assert.Equal(t, username, user.Username)
	})

	t.Run("FindByEmail Success", func(t *testing.T) {
		email := "test@mail.com"
		rows := sqlmock.NewRows([]string{"id", "email"}).AddRow(1, email)
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "users" WHERE email = $1`)).
			WithArgs(email, 1).WillReturnRows(rows)

		user, err := repo.FindByEmail(ctx, email)
		assert.NoError(t, err)
		assert.Equal(t, email, user.Email)
	})

	t.Run("Register Success", func(t *testing.T) {
		req := &dto.RegisterRequest{
			Name: "Test", Username: "test", Email: "test@mail.com", Password: "hash", RoleID: 1,
		}
		
		mock.ExpectBegin()
		mock.ExpectQuery(regexp.QuoteMeta(`INSERT INTO "users"`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
		mock.ExpectCommit()

		user, err := repo.Register(ctx, req)
		assert.NoError(t, err)
		assert.NotNil(t, user)
	})
}
