package payment

import (
	"context"
	"database/sql"
	"payment-service/domain/dto"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type PaymentRepositoryTestSuite struct {
	suite.Suite
	mock       sqlmock.Sqlmock
	db         *gorm.DB
	repository PaymentRepositoryInterface
}

func stringPtr(s string) *string {
	return &s
}

func (s *PaymentRepositoryTestSuite) SetupTest() {
	var (
		db  *sql.DB
		err error
	)

	db, s.mock, err = sqlmock.New()
	assert.NoError(s.T(), err)

	dialector := postgres.New(postgres.Config{
		Conn:       db,
		DriverName: "postgres",
	})

	s.db, err = gorm.Open(dialector, &gorm.Config{})
	assert.NoError(s.T(), err)

	s.repository = NewPaymentRepository(s.db)
}

func (s *PaymentRepositoryTestSuite) TearDownTest() {
	db, _ := s.db.DB()
	db.Close()
}

func TestPaymentRepository(t *testing.T) {
	suite.Run(t, new(PaymentRepositoryTestSuite))
}

func (s *PaymentRepositoryTestSuite) TestCreate_Success() {
	orderID := uuid.New().String()
	req := &dto.PaymentRequest{
		OrderID:     orderID,
		Amount:      100000,
		PaymentLink: "http://midtrans.com/pay",
		ExpiredAt:   time.Now().Add(time.Hour),
		Description: stringPtr("Test"),
	}

	s.mock.ExpectBegin()
	s.mock.ExpectQuery(`INSERT INTO "payments"`).
		WithArgs(sqlmock.AnyArg(), sqlmock.AnyArg(), req.Amount, sqlmock.AnyArg(), req.PaymentLink, sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), sqlmock.AnyArg(), req.Description).
		WillReturnRows(sqlmock.NewRows([]string{"id"}).AddRow(1))
	s.mock.ExpectCommit()

	result, err := s.repository.Create(context.Background(), s.db, req)

	assert.NoError(s.T(), err)
	assert.NotNil(s.T(), result)
	assert.Equal(s.T(), req.Amount, result.Amount)
}

func (s *PaymentRepositoryTestSuite) TestFindByUUID_Success() {
	uid := uuid.New()
	// GORM First adds LIMIT 1
	s.mock.ExpectQuery(`SELECT \* FROM "payments" WHERE uuid = \$1`).
		WithArgs(uid.String(), 1).
		WillReturnRows(sqlmock.NewRows([]string{"uuid", "amount"}).AddRow(uid, 100000))

	result, err := s.repository.FindByUUID(context.Background(), uid.String())

	assert.NoError(s.T(), err)
	assert.Equal(s.T(), uid, result.UUID)
}

func (s *PaymentRepositoryTestSuite) TestFindByUUID_NotFound() {
	uid := uuid.New()
	s.mock.ExpectQuery(`SELECT \* FROM "payments" WHERE uuid = \$1`).
		WithArgs(uid.String(), 1).
		WillReturnError(gorm.ErrRecordNotFound)

	result, err := s.repository.FindByUUID(context.Background(), uid.String())

	assert.Error(s.T(), err)
	assert.Nil(s.T(), result)
}
