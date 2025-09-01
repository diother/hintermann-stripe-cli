package dto

type DonationDTO struct {
	ID          string
	Created     string
	ClientName  string
	ClientEmail string
	PayoutID    string
	Gross       string
	Fee         string
	Net         string
}

func NewDonationDTO(id, created, clientName, clientEmail, payoutID, gross, fee, net string) *DonationDTO {
	return &DonationDTO{
		ID:          id,
		Created:     created,
		ClientName:  clientName,
		ClientEmail: clientEmail,
		PayoutID:    payoutID,
		Gross:       gross,
		Fee:         fee,
		Net:         net,
	}
}
