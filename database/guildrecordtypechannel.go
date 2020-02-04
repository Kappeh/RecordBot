package database

import (
	"github.com/pkg/errors"
)

// RecordType gets the record type of the guild record type channel
func (g GuildRecordTypeChannel) RecordType() (RecordType, bool, error) {
	db, err := Instance()
	if err != nil {
		return RecordType{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	rt, ok, err := db.RecordType(g.RecordTypeID)
	if err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to determine if record type exists")
	}
	return rt, ok, nil
}
