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
