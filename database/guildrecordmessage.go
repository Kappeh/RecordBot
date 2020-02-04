package database

import (
	"github.com/pkg/errors"
)

// Record gets the record of the guild record message
func (g GuildRecordMessage) Record() (Record, bool, error) {
	db, err := Instance()
	if err != nil {
		return Record{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	r, ok, err := db.Record(g.RecordID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to determine if record exists")
	}
	return r, ok, nil
}
