package database

import (
	"github.com/pkg/errors"
)

// Build gets the build of the guild build message
func (g GuildBuildMessage) Build() (Build, bool, error) {
	db, err := Instance()
	if err != nil {
		return Build{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	b, ok, err := db.Build(g.BuildID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to determine if build exists")
	}
	return b, ok, nil
}
