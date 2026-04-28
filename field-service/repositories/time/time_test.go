package repositories

import (
	"context"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func setupMockDB(t *testing.T) (*gorm.DB, sqlmock.Sqlmock) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("failed to open sqlmock: %s", err)
	}

	dialector := postgres.New(postgres.Config{
		Conn: db,
	})
	gormDB, err := gorm.Open(dialector, &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open gorm: %s", err)
	}

	return gormDB, mock
}

func TestTimeRepository_FindAll(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTimeRepository(db)
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "start_time", "end_time"}).
			AddRow(1, "08:00", "09:00").
			AddRow(2, "09:00", "10:00")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "times"`)).
			WillReturnRows(rows)

		result, err := repo.FindAll(ctx)

		assert.NoError(t, err)
		assert.Len(t, result, 2)
	})
}

func TestTimeRepository_FindByUUID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewTimeRepository(db)
	ctx := context.Background()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "uuid", "start_time"}).
			AddRow(1, id, "08:00")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "times" WHERE uuid = $1 ORDER BY "times"."id" LIMIT $2`)).
			WithArgs(id.String(), 1).
			WillReturnRows(rows)

		result, err := repo.FindByUUID(ctx, id.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "08:00", result.StartTime)
	})
}
