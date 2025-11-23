package pullrequest

import (
	"time"
)

type PullRequest struct {
	PullRequestShort *PullRequestShort `json:"pull_request_short"`
	AssignedReviwers []string          `json:"assigned_reviewers"`
	CreatedAt        time.Time         `json:"created_at"`
	MergedAt         time.Time         `json:"merged_at"`
}

func NewPullRequest(pullRequestShort *PullRequestShort) *PullRequest {
	reviwers := make([]string, 0, 2)
	return &PullRequest{
		PullRequestShort: pullRequestShort,
		AssignedReviwers: reviwers,
	}
}
