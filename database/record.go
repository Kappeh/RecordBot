package database

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Edition gets the edition of the record
func (r Record) Edition() (Edition, bool, error) {
	db, err := Instance()
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	e, ok, err := db.Edition(r.EditionID)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to determine if edition exists")
	}
	return e, ok, nil
}

// BuildClass gets the build class of the record
func (r Record) BuildClass() (BuildClass, bool, error) {
	db, err := Instance()
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	bc, ok, err := db.BuildClass(r.BuildClassID)
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to determine if build class exists")
	}
	return bc, ok, nil
}

// UpdateRequestRecord gets the record which record is requesting to update
func (r Record) UpdateRequestRecord() (Record, bool, error) {
	if !r.UpdateRequest {
		return Record{}, false, nil
	}
	db, err := Instance()
	if err != nil {
		return Record{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	record, ok, err := db.Record(r.UpdateRequestRecordID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to determine if record exists")
	}
	return record, ok, nil
}

// BuildRecords gets the build records of the record for a specified build
func (r Record) BuildRecords(buildID string) ([]BuildRecord, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert ids to ints
	recordIDint, err := strconv.Atoi(r.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert build id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT ID, Verified, VerifierID, VerifiedTimestamp, Reported,
			ReporterID, ReportedTimestamp, JointBuildRecord, JointBuildRecordID,
			SubmitterID, Timestamp, EditedTimestamp
		FROM BuildRecords
		WHERE BuildID = ? AND RecordID = ?
	`, buildIDint, recordIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []BuildRecord{}
	var (
		idInt                   int
		verifiedInt             int
		verifierIDint           int
		verifiedTimestampString string
		reportedInt             int
		reporterIDint           int
		reportedTimestampString string
		jointBuildRecordInt     int
		jointBuildRecordIDint   int
		submitterIDint          int
		timestampString         string
		editedTimestampString   string
		verifiedTimestamp       time.Time
		reportedTimestamp       time.Time
		timestamp               time.Time
		editedTimestamp         time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&idInt, &verifiedInt, &verifierIDint, &verifiedTimestampString,
			&reportedInt, &reporterIDint, &reportedTimestampString,
			&jointBuildRecordInt, &jointBuildRecordIDint, &submitterIDint,
			&timestampString, &editedTimestampString,
		); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Parse timestamps
		if verifiedTimestamp, err = time.Parse(timeLayout, verifiedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse verified timestamp")
		}
		if reportedTimestamp, err = time.Parse(timeLayout, reportedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse reported timestamp")
		}
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}
		if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse edited timestamp")
		}
		// Add to results
		results = append(results, BuildRecord{
			ID:                 strconv.Itoa(idInt),
			BuildID:            buildID,
			RecordID:           r.ID,
			Verified:           verifiedInt != 0,
			VerifierID:         strconv.Itoa(verifierIDint),
			VerifiedTimestamp:  Timestamp(verifiedTimestamp),
			Reported:           reportedInt != 0,
			ReporterID:         strconv.Itoa(reporterIDint),
			ReportedTimestamp:  Timestamp(reportedTimestamp),
			JointBuildRecord:   jointBuildRecordInt != 0,
			JointBuildRecordID: strconv.Itoa(jointBuildRecordIDint),
			SubmitterID:        strconv.Itoa(submitterIDint),
			Timestamp:          Timestamp(timestamp),
			EditedTimestamp:    Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// BuildRecordsAll gets the build records of the record for all builds
func (r Record) BuildRecordsAll() ([]BuildRecord, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	recordIDint, err := strconv.Atoi(r.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT ID, BuildID, Verified, VerifierID, VerifiedTimestamp, Reported,
			ReporterID, ReportedTimestamp, JointBuildRecord, JointBuildRecordID,
			SubmitterID, Timestamp, EditedTimestamp
		FROM BuildRecords
		WHERE RecordID = ?
	`, recordIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []BuildRecord{}
	var (
		idInt                   int
		buildIDint              int
		verifiedInt             int
		verifierIDint           int
		verifiedTimestampString string
		reportedInt             int
		reporterIDint           int
		reportedTimestampString string
		jointBuildRecordInt     int
		jointBuildRecordIDint   int
		submitterIDint          int
		timestampString         string
		editedTimestampString   string
		verifiedTimestamp       time.Time
		reportedTimestamp       time.Time
		timestamp               time.Time
		editedTimestamp         time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&idInt, &buildIDint, &verifiedInt, &verifierIDint, &verifiedTimestampString,
			&reportedInt, &reporterIDint, &reportedTimestampString,
			&jointBuildRecordInt, &jointBuildRecordIDint, &submitterIDint,
			&timestampString, &editedTimestampString,
		); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Parse timestamps
		if verifiedTimestamp, err = time.Parse(timeLayout, verifiedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse verified timestamp")
		}
		if reportedTimestamp, err = time.Parse(timeLayout, reportedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse reported timestamp")
		}
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}
		if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse edited timestamp")
		}
		// Add to results
		results = append(results, BuildRecord{
			ID:                 strconv.Itoa(idInt),
			BuildID:            strconv.Itoa(buildIDint),
			RecordID:           r.ID,
			Verified:           verifiedInt != 0,
			VerifierID:         strconv.Itoa(verifierIDint),
			VerifiedTimestamp:  Timestamp(verifiedTimestamp),
			Reported:           reportedInt != 0,
			ReporterID:         strconv.Itoa(reporterIDint),
			ReportedTimestamp:  Timestamp(reportedTimestamp),
			JointBuildRecord:   jointBuildRecordInt != 0,
			JointBuildRecordID: strconv.Itoa(jointBuildRecordIDint),
			SubmitterID:        strconv.Itoa(submitterIDint),
			Timestamp:          Timestamp(timestamp),
			EditedTimestamp:    Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// GuildRecordMessage gets the guild record message for the record for a specified guild
func (r Record) GuildRecordMessage(guildID string) (GuildRecordMessage, bool, error) {
	db, err := Instance()
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	grm, ok, err := db.GuildRecordMessage(guildID, r.ID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to determine if guild record message exists")
	}
	return grm, ok, nil
}

// GuildRecordMessages gets the guild record message for the record for all guilds
func (r Record) GuildRecordMessages() ([]GuildRecordMessage, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	idInt, err := strconv.Atoi(r.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT GuildID, ChannelID, MessageID, Timestamp, EditedTimestamp
		FROM GuildRecordMessages
		WHERE RecordID = ?
	`, idInt)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []GuildRecordMessage{}
	var (
		guildIDint            int
		channelIDint          int
		messageIDint          int
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&guildIDint, &channelIDint, &messageIDint,
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
		results = append(results, GuildRecordMessage{
			GuildID:         strconv.Itoa(guildIDint),
			RecordID:        r.ID,
			ChannelID:       strconv.Itoa(channelIDint),
			MessageID:       strconv.Itoa(messageIDint),
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}
