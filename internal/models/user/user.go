package user

import "github.com/lohpesch09/Avito-backend-trainee-assignment-autumn-2025/internal/models/team"

type User struct {
	TeamMember *team.TeamMember `json:"team_member"`
	TeamName   string           `json:"team_name"`
}

func NewUser(teamMember *team.TeamMember, teamName string) *User {
	return &User{
		TeamMember: teamMember,
		TeamName:   teamName,
	}
}
