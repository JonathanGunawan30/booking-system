package repositories

import (
	"context"
	"field-service/domain/dto"
	"regexp"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// --- Helper for Mock DB ---

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

// --- Scenarios ---

func TestFieldRepositories_FindByUUID(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFieldRepositories(db)
	ctx := context.Background()
	id := uuid.New()

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "uuid", "name"}).
			AddRow(1, id, "Field 1")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "fields" WHERE uuid = $1 AND "fields"."deleted_at" IS NULL ORDER BY "fields"."id" LIMIT $2`)).
			WithArgs(id.String(), 1).
			WillReturnRows(rows)

		result, err := repo.FindByUUID(ctx, id.String())

		assert.NoError(t, err)
		assert.NotNil(t, result)
		if result != nil {
			assert.Equal(t, "Field 1", result.Name)
		}
	})

	t.Run("not found", func(t *testing.T) {
		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "fields" WHERE uuid = $1 AND "fields"."deleted_at" IS NULL ORDER BY "fields"."id" LIMIT $2`)).
			WithArgs(id.String(), 1).
			WillReturnError(gorm.ErrRecordNotFound)

		result, err := repo.FindByUUID(ctx, id.String())

		assert.Error(t, err)
		assert.Nil(t, result)
	})
}

func TestFieldRepositories_FindAllWithPagination(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFieldRepositories(db)
	ctx := context.Background()
	sortCol := "id"
	sortOrder := "asc"
	param := &dto.FieldRequestParam{
		Page:       1,
		Limit:      10,
		SortColumn: &sortCol,
		SortOrder:  &sortOrder,
	}

	t.Run("success", func(t *testing.T) {
		rows := sqlmock.NewRows([]string{"id", "name"}).
			AddRow(1, "Field 1").
			AddRow(2, "Field 2")

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT * FROM "fields" WHERE "fields"."deleted_at" IS NULL ORDER BY id asc LIMIT $1`)).
			WithArgs(10).
			WillReturnRows(rows)

		mock.ExpectQuery(regexp.QuoteMeta(`SELECT count(*) FROM "fields" WHERE "fields"."deleted_at" IS NULL`)).
			WillReturnRows(sqlmock.NewRows([]string{"count"}).AddRow(2))

		fields, total, err := repo.FindAllWithPagination(ctx, param)

		assert.NoError(t, err)
		assert.Len(t, fields, 2)
		assert.Equal(t, int64(2), total)
	})
}


func TestFieldRepositories_Delete(t *testing.T) {
	db, mock := setupMockDB(t)
	repo := NewFieldRepositories(db)
	ctx := context.Background()
	id := uuid.New().String()

	t.Run("success", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "fields" SET "deleted_at"=$1 WHERE uuid = $2 AND "fields"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), id).
			WillReturnResult(sqlmock.NewResult(0, 1))
		mock.ExpectCommit()

		err := repo.Delete(ctx, id)

		assert.NoError(t, err)
	})

	t.Run("error", func(t *testing.T) {
		mock.ExpectBegin()
		mock.ExpectExec(regexp.QuoteMeta(`UPDATE "fields" SET "deleted_at"=$1 WHERE uuid = $2 AND "fields"."deleted_at" IS NULL`)).
			WithArgs(sqlmock.AnyArg(), id).
			WillReturnError(assert.AnError)
		mock.ExpectRollback()

		err := repo.Delete(ctx, id)

		assert.Error(t, err)
	})
}

