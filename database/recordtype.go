package database

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// GuildRecordTypeChannel gets the guild record type channel for the record type for
// a specified guild
func (rt RecordType) GuildRecordTypeChannel(guildID string) (GuildRecordTypeChannel, bool, error) {
	db, err := Instance()
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	grtc, ok, err := db.GuildRecordTypeChannel(guildID, rt.ID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to determine if guild record type channel exists")
	}
	return grtc, ok, nil
}

// GuildRecordTypeChannels gets the guild record type channels for the record type
// for all guilds
func (rt RecordType) GuildRecordTypeChannels() ([]GuildRecordTypeChannel, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	recordTypeIDint, err := strconv.Atoi(rt.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT GuildID, ChannelID, Timestamp, EditedTimestamp
		FROM GuildRecordTypeChannels
		WHERE RecordTypeID = ?
	`, recordTypeIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []GuildRecordTypeChannel{}
	var (
		guildIDint            int
		channelIDint          int
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&guildIDint, &channelIDint,
			&timestampString, &editedTimestampString,
		); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Parse timestamps
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}
		if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse edited timestamp")
		}
		// Add to results
		results = append(results, GuildRecordTypeChannel{
			GuildID:         strconv.Itoa(guildIDint),
			RecordTypeID:    rt.ID,
			ChannelID:       strconv.Itoa(channelIDint),
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// Records get the all the records that fall into the record type
func (rt RecordType) Records() ([]Record, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	recordTypeIDint, err := strconv.Atoi(rt.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT ID, Verified, VerifierID, VerifiedTimestamp, UpdateRequest,
			UpdateRequestRecordID, EditionID, BuildClassID, Name, Description,
			SubmitterID, Timestamp, EditedTimestamp
		FROM Records
		WHERE RecordTypeID = ?
	`, recordTypeIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []Record{}
	var (
		idInt                    int
		verifiedInt              int
		verifierIDint            int
		verifiedTimestampString  string
		updateRequestInt         int
		updateRequestRecordIDint int
		editionIDint             int
		buildClassIDint          int
		name                     string
		description              string
		submitterIDint           int
		timestampString          string
		editedTimestampString    string
		verifiedTimestamp        time.Time
		timestamp                time.Time
		editedTimestamp          time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&idInt, &verifiedInt, &verifierIDint, &verifiedTimestampString,
			&updateRequestInt, &updateRequestRecordIDint, &editionIDint,
			&buildClassIDint, &name, &description, &submitterIDint,
			&timestampString, &editedTimestampString,
		); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Parse timestamps
		if verifiedTimestamp, err = time.Parse(timeLayout, verifiedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse verified timestamp")
		}
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}
		if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse edited timestamp")
		}
		// Add to results
		results = append(results, Record{
			ID:                    strconv.Itoa(idInt),
			Verified:              verifiedInt != 0,
			VerifierID:            strconv.Itoa(verifierIDint),
			VerifiedTimestamp:     Timestamp(verifiedTimestamp),
			UpdateRequest:         updateRequestInt != 0,
			UpdateRequestRecordID: strconv.Itoa(updateRequestRecordIDint),
			EditionID:             strconv.Itoa(editionIDint),
			BuildClassID:          strconv.Itoa(buildClassIDint),
			RecordTypeID:          strconv.Itoa(recordTypeIDint),
			Name:                  name,
			Description:           description,
			SubmitterID:           strconv.Itoa(submitterIDint),
			Timestamp:             Timestamp(timestamp),
			EditedTimestamp:       Timestamp(editedTimestamp),
		})
	}
	return results, nil
}
