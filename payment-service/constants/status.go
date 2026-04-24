package constants

type PaymentStatus int
type PaymentStatusString string

const (
	Initial    PaymentStatus = 0
	Pending    PaymentStatus = 100
	Settlement PaymentStatus = 200
	Expired    PaymentStatus = 300

	InitialString    PaymentStatusString = "initial"
	PendingString    PaymentStatusString = "pending"
	SettlementString PaymentStatusString = "settlement"
	ExpiredString    PaymentStatusString = "expire"
)

var mapStatusStringToInt = map[PaymentStatusString]PaymentStatus{
	InitialString:    Initial,
	PendingString:    Pending,
	SettlementString: Settlement,
	ExpiredString:    Expired,
}

var mapStatusIntToString = map[PaymentStatus]PaymentStatusString{
	Initial:    InitialString,
	Pending:    PendingString,
	Settlement: SettlementString,
	Expired:    ExpiredString,
}

func (p *PaymentStatus) GetStatusString() PaymentStatusString {
	return mapStatusIntToString[*p]
}
func (p *PaymentStatusString) GetStatusInt() PaymentStatus {
	return mapStatusStringToInt[*p]
}
