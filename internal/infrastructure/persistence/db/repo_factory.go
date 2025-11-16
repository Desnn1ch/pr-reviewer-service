package db

import (
	"github.com/Desnn1ch/pr-reviewer-service/internal/app"
)

type Repositories struct {
	Teams app.TeamRepo
	Users app.UserRepo
	PRs   app.PRRepo
	Tx    app.TxManager
}

func NewRepositories(db *DB) Repositories {
	return Repositories{
		Teams: NewTeamRepo(db),
		Users: NewUserRepo(db),
		PRs:   NewPRRepo(db),
		Tx:    db,
	}
}
