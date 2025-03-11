package createinvitation

type Invitation struct {
	Email   string `json:"email"`
	Address string `json:"address"`
	Phone   string `json:"phone"`
}

func CreateInvitation(invitation Invitation) (Invitation, error) {
	return invitation, nil
}
