package constants

type FieldStatusString string

const (
	AvailableStatus FieldStatusString = "available"
	BookedStatus    FieldStatusString = "booked"
)

func (f *FieldStatusString) GetStatus() string {
	return string(*f)
}
