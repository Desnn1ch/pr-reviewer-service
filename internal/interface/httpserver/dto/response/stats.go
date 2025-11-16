package response

type ReviewerStat struct {
	UserID          string `json:"user_id"`
	Username        string `json:"username"`
	TeamName        string `json:"team_name"`
	AssignedOpenPRs int    `json:"assigned_open_prs"`
}

type ReviewerStats struct {
	Items []ReviewerStat `json:"items"`
}
