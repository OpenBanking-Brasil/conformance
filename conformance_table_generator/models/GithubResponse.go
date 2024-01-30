package models

type GithubResponse struct {
	Sha  string     `json:"sha"`
	URL  string     `json:"url"`
	Tree GithubTree `json:"tree"`
	Truncated bool  `json:"truncated"`
}
