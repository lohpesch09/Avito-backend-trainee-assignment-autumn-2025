package pullrequest

import validation "github.com/go-ozzo/ozzo-validation"

type Status string

const (
	OPEN   Status = "OPEN"
	MERGED Status = "MERGED"
)

type PullRequestShort struct {
	PullRequestId   string `json:"pull_request_id"`
	PullRequestName string `json:"pull_request_name"`
	AuthorId        string `json:"author_id"`
	Status          Status `json:"status"`
}

func NewPullRequestShort(pullRequestId string, pullRequestName string, authorId string) *PullRequestShort {
	return &PullRequestShort{
		PullRequestId:   pullRequestId,
		PullRequestName: pullRequestName,
		AuthorId:        authorId,
		Status:          OPEN,
	}
}

func (pr *PullRequestShort) Validation() error {
	return validation.ValidateStruct(
		pr,
		validation.Field(&pr.PullRequestId, validation.Required, validation.Length(4, 50)),
		validation.Field(&pr.PullRequestName, validation.Required),
		validation.Field(&pr.AuthorId, validation.Required, validation.Length(2, 50)),
	)
}
