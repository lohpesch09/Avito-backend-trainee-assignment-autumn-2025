package errors

type Code string

const (
	TEAM_EXISTS              Code = "TEAM_EXISTS"
	PR_EXISTS                Code = "PR_EXISTS"
	PR_MERGED                Code = "PR_MERGED"
	NOT_ASSIGNED             Code = "NOT_ASSIGNED"
	NO_CANDIDATE             Code = "NO_CANDIDATE"
	NOT_FOUND                Code = "NOT_FOUND"
	USER_EXISTS_ANOTHER_TEAM Code = "USER_EXISTS_ANOTHER_TEAM"
)

var (
	TeamExists *Error = &Error{
		Code:    TEAM_EXISTS,
		Message: "team_name already exists",
	}
	PrExists *Error = &Error{
		Code:    PR_EXISTS,
		Message: "PR id already exists",
	}
	PrMerged *Error = &Error{
		Code:    PR_MERGED,
		Message: "cannot reassign on merged PR",
	}
	NotAssigned *Error = &Error{
		Code:    NOT_ASSIGNED,
		Message: "reviewer is not assigned to this PR",
	}
	NoCandidate *Error = &Error{
		Code:    NO_CANDIDATE,
		Message: "no active replacement candidate in team",
	}
	NotFound *Error = &Error{
		Code:    NOT_FOUND,
		Message: "resource not found",
	}
	UserExistsAnotherTeam *Error = &Error{
		Code:    USER_EXISTS_ANOTHER_TEAM,
		Message: "user already exists in another team",
	}
)

type Error struct {
	Code    Code   `json:"code"`
	Message string `json:"message"`
}

func (e *Error) Error() string {
	return e.Message
}

func (e *Error) Is(err error) bool {
	if castErr, ok := err.(*Error); ok {
		return e.Code == castErr.Code
	}
	return false
}
