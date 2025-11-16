package mapper

import (
	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
	req "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/request"
	resp "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
)

func CreatePRRequestToArgs(r req.CreatePR) (uuid.UUID, string, uuid.UUID, error) {
	prID, err := uuid.Parse(r.PullRequestID)
	if err != nil {
		return uuid.Nil, "", uuid.Nil, err
	}

	authorID, err := uuid.Parse(r.AuthorID)
	if err != nil {
		return uuid.Nil, "", uuid.Nil, err
	}

	return prID, r.PullRequestName, authorID, nil
}

func MergePRRequestToID(r req.MergePR) (uuid.UUID, error) {
	prID, err := uuid.Parse(r.PullRequestID)
	if err != nil {
		return uuid.Nil, err
	}
	return prID, nil
}

func ReassignRequestToArgs(r req.ReassignReviewer) (uuid.UUID, uuid.UUID, error) {
	prID, err := uuid.Parse(r.PullRequestID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	oldID, err := uuid.Parse(r.OldUserID)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return prID, oldID, nil
}

func PRToResponse(pr entity.PR) resp.PullRequest {
	reviewers := make([]string, 0, len(pr.Reviewers))
	for _, id := range pr.Reviewers {
		reviewers = append(reviewers, id.String())
	}

	return resp.PullRequest{
		PullRequestID:     pr.ID.String(),
		PullRequestName:   pr.Title,
		AuthorID:          pr.AuthorID.String(),
		Status:            string(pr.Status),
		AssignedReviewers: reviewers,
		CreatedAt:         pr.CreatedAt,
		MergedAt:          pr.MergedAt,
	}
}

func PRsToShortResponse(prs []entity.PR) []resp.PullRequestShort {
	res := make([]resp.PullRequestShort, 0, len(prs))

	for _, pr := range prs {
		res = append(res, resp.PullRequestShort{
			PullRequestID:   pr.ID.String(),
			PullRequestName: pr.Title,
			AuthorID:        pr.AuthorID.String(),
			Status:          string(pr.Status),
		})
	}

	return res
}
