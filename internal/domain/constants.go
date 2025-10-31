package domain

const (
	// user roles
	RoleCustomer = "CUSTOMER"
	RoleAdmin    = "ADMIN"

	BookingStatusPending   = "PENDING"
	BookingStatusConfirmed = "CONFIRMED"
	BookingStatusCancelled = "CANCELLED"
	BookingStatusCompleted = "COMPLETED"

	PaymentStatusPending = "PENDING"
	PaymentStatusSuccess = "SUCCESS"
	PaymentStatusFailed  = "FAILED"

	PaymentMethodVA           = "VIRTUAL_ACCOUNT"
	PaymentMethodCreditCard   = "CREDIT_CARD"
	PaymentMethodEWallet      = "E_WALLET"
	PaymentMethodBankTransfer = "BANK_TRANSFER"
)
