package dto

type PayoutDTO struct {
	ID      string
	Created string
	Gross   string
	Fee     string
	Net     string
}

func NewPayoutDTO(id, created, gross, fee, net string) *PayoutDTO {
	return &PayoutDTO{
		ID:      id,
		Created: created,
		Gross:   gross,
		Fee:     fee,
		Net:     net,
	}
}
