package team

type TeamMember struct {
	UserId   string `json:"user_id"`
	UserName string `json:"username"`
	IsActive bool   `json:"is_active"`
}

func NewTeamMember(userId string, userName string, isActive bool) *TeamMember {
	return &TeamMember{
		UserId:   userId,
		UserName: userName,
		IsActive: isActive,
	}
}
