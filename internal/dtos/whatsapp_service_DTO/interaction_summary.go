package whatsappservicedto

type UserContactInfo struct {
	Telefono string
	Email    string
}

type InteractionSummary struct {
	NumberPhoneID int64
	NumberPhone   int64
	Contacts      []UserContactInfo
}
