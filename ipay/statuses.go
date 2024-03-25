package ipay

type PaymentStatus int8

const (
	PaymentStatusUnknown             PaymentStatus = 0
	PaymentStatusRegistered          PaymentStatus = 1
	PaymentStatusPreAuthorized       PaymentStatus = 3
	PaymentStatusFailed              PaymentStatus = 4
	PaymentStatusSuccess             PaymentStatus = 5
	PaymentStatusCanceled            PaymentStatus = 9
	PaymentStatusManualProcessing    PaymentStatus = 11
	PaymentStatusSuccessWithoutClaim PaymentStatus = 13
	PaymentStatusSecurityRefusal     PaymentStatus = 14
)

func (s PaymentStatus) String() string {
	switch s {
	case PaymentStatusRegistered:
		return "Registered"
	case PaymentStatusPreAuthorized:
		return "PreAuthorized"
	case PaymentStatusFailed:
		return "Failed"
	case PaymentStatusSuccess:
		return "Success"
	case PaymentStatusCanceled:
		return "Canceled"
	case PaymentStatusManualProcessing:
		return "ManualProcessing"
	case PaymentStatusSuccessWithoutClaim:
		return "SuccessWithoutClaim"
	case PaymentStatusSecurityRefusal:
		return "SecurityRefusal"
	default:
		return "Unknown"
	}
}
