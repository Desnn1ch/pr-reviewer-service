package mapper

import (
	"github.com/google/uuid"

	"github.com/Desnn1ch/pr-reviewer-service/internal/domain/entity"
	req "github.com/Desnn1ch/pr-reviewer-service/internal/interface/http-server/dto/request"
	resp "github.com/Desnn1ch/pr-reviewer-service/internal/interface/http-server/dto/response"
)

func SetIsActiveRequestToArgs(r req.SetIsActive) (uuid.UUID, bool, error) {
	id, err := uuid.Parse(r.UserID)
	if err != nil {
		return uuid.Nil, false, err
	}

	return id, r.IsActive, nil
}

func UserToResponse(u entity.User, teamName string) resp.User {
	return resp.User{
		UserID:   u.ID.String(),
		Username: u.Name,
		TeamName: teamName,
		IsActive: u.IsActive,
	}
}
