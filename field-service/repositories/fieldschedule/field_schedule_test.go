package repositories

import (
	"context"
	"field-service/constants"
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

func TestFieldScheduleRepositories_FindByUUID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFieldScheduleRepositories(db)
	ctx := context.Background()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "uuid", "field_id"}).
			AddRow(1, id, 1)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "field_schedules" WHERE uuid = $1 AND "field_schedules"."deleted_at" IS NULL ORDER BY "field_schedules"."id" LIMIT $2`)).
			WithArgs(id.String(), 1).
			WillReturnRows(rows)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "fields" WHERE "fields"."id" = $1 AND "fields"."deleted_at" IS NULL`)).
			WithArgs(1).
			WillReturnRows(sqlmock.NewRows([]string{"id", "name"}).AddRow(1, "Field 1"))

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "times" WHERE "times"."id" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"id"}))

		result, err := repo.FindByUUID(ctx, id.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
	})
}

func TestFieldScheduleRepositories_UpdateStatus(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFieldScheduleRepositories(db)
	ctx := context.Background()
	id := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "field_schedules" SET "status"=$1,"updated_at"=$2 WHERE uuid = $3 AND "field_schedules"."deleted_at" IS NULL`)).
			WithArgs(constants.Booked, sqlmock.AnyArg(), id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.UpdateStatus(ctx, id, constants.Booked)

		assert.NoError(t, err)
	})
}
