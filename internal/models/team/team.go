package team

import validation "github.com/go-ozzo/ozzo-validation"

type Team struct {
	TeamName string       `json:"team_name"`
	Members  []TeamMember `json:"members"`
}

func NewTeam(teamName string) *Team {
	return &Team{
		TeamName: teamName,
	}
}

func (t *Team) Validation() error {
	return validation.ValidateStruct(
		t,
		validation.Field(&t.TeamName, validation.Required, validation.Length(2, 50)),
		validation.Field(&t.Members, validation.Required),
	)
}
