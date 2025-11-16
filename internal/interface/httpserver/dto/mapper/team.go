package mapper

import (
	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
	req "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/request"
	resp "github.com/Desnn1ch/pr-reviewer-service/internal/interface/httpserver/dto/response"
)

func TeamAddRequestToArgs(r req.TeamAdd) (string, []entity.User, error) {
	members := make([]entity.User, 0, len(r.Members))

	for _, m := range r.Members {
		var id uuid.UUID

		if m.UserID != "" {
			parsed, err := uuid.Parse(m.UserID)
			if err != nil {
				return "", nil, err
			}
			id = parsed
		} else {
			id = uuid.New()
		}

		members = append(members, entity.User{
			ID:       id,
			Name:     m.Username,
			IsActive: m.IsActive,
		})
	}

	return r.TeamName, members, nil
}

func TeamToResponse(team entity.Team, members []entity.User) resp.Team {
	respMembers := make([]resp.TeamMember, 0, len(members))

	for _, u := range members {
		respMembers = append(respMembers, resp.TeamMember{
			UserID:   u.ID.String(),
			Username: u.Name,
			IsActive: u.IsActive,
		})
	}

	return resp.Team{
		TeamName: team.Name,
		Members:  respMembers,
	}
}
