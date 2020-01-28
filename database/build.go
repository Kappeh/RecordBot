package database

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Edition gets the edition of the build
func (b Build) Edition() (Edition, bool, error) {
	db, err := Instance()
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	e, ok, err := db.Edition(b.EditionID)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to get edition")
	}
	return e, ok, nil
}

// BuildClass gets the build class of the build
func (b Build) BuildClass() (BuildClass, bool, error) {
	db, err := Instance()
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	bc, ok, err := db.BuildClass(b.BuildClassID)
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to get build class")
	}
	return bc, ok, nil
}

// UpdateRequestBuild get the build which is being requested to update
func (b Build) UpdateRequestBuild() (Build, bool, error) {
	db, err := Instance()
	if err != nil {
		return Build{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	b, ok, err := db.Build(b.UpdateRequestBuildID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to get build")
	}
	return b, ok, nil
}

// GuildBuildMessage gets the guild build message for a specified guild
func (b Build) GuildBuildMessage(guildID string) (GuildBuildMessage, bool, error) {
	db, err := Instance()
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	gbm, ok, err := db.GuildBuildMessage(guildID, b.ID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "couldn't get guild build message")
	}
	return gbm, ok, nil
}

// GuildBuildMessages get the guild build messages for all guilds
func (b Build) GuildBuildMessages() ([]GuildBuildMessage, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	idInt, err := strconv.Atoi(b.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query database
	rows, err := db.db.Query(`
		SELECT GuildID, ChannelID, MessageID, Timestamp, EditedTimestamp
		FROM GuildBuildMessages
		WHERE BuildID = ?
	`, idInt)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []GuildBuildMessage{}
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
		results = append(results, GuildBuildMessage{
			GuildID:         strconv.Itoa(guildIDint),
			BuildID:         b.ID,
			ChannelID:       strconv.Itoa(channelIDint),
			MessageID:       strconv.Itoa(messageIDint),
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// BuildVersion gets the build version for a specified version
func (b Build) BuildVersion(versionID string) (BuildVersion, bool, error) {
	db, err := Instance()
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	bv, ok, err := db.BuildVersion(b.ID, versionID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "couldn't get build version")
	}
	return bv, ok, nil
}

// BuildVersions gets the build versions for all versions
func (b Build) BuildVersions() ([]BuildVersion, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	idInt, err := strconv.Atoi(b.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query database
	rows, err := db.db.Query(`
		SELECT VersionID, StatusID, Notes, Timestamp, EditedTimestamp
		FROM BuildVersions
		WHERE BuildID = ?
	`, idInt)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []BuildVersion{}
	var (
		versionIDint          int
		statusIDint           int
		notes                 string
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&versionIDint, &statusIDint, &notes,
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
		results = append(results, BuildVersion{
			BuildID:         b.ID,
			VersionID:       strconv.Itoa(versionIDint),
			StatusID:        strconv.Itoa(statusIDint),
			Notes:           notes,
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// BuildRecord gets the build record for the build and a specified record
func (b Build) BuildRecord(recordID string) (BuildRecord, bool, error) {
	db, err := Instance()
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id and recordID to ints
	buildIDint, err := strconv.Atoi(b.ID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert id to integer")
	}
	recordIDint, err := strconv.Atoi(recordID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert record id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT ID, Verified, VerifierID, VerifiedTimestamp, Reported, ReporterID,
			ReportedTimestamp, JointBuildRecord, JointBuildRecordID, SubmitterID,
			Timestamp, EditedTimestamp
		FROM BuildRecords
		WHERE BuildID = ? AND RecordID = ?
	`, buildIDint, recordIDint)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if build record exists
	if !rows.Next() {
		return BuildRecord{}, false, nil
	}
	// Extract data
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
	if err = rows.Scan(
		&idInt, &verifiedInt, &verifierIDint, &verifiedTimestampString,
		&reportedInt, &reporterIDint, &reportedTimestampString,
		&jointBuildRecordInt, &jointBuildRecordIDint, &submitterIDint,
		&timestampString, &editedTimestampString,
	); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if verifiedTimestamp, err = time.Parse(timeLayout, verifiedTimestampString); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to parse verified timestamp")
	}
	if reportedTimestamp, err = time.Parse(timeLayout, reportedTimestampString); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to parse reported timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	return BuildRecord{
		ID:                 strconv.Itoa(idInt),
		BuildID:            b.ID,
		RecordID:           recordID,
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
	}, true, nil
}

// BuildRecords gets the build records for the build and all records
func (b Build) BuildRecords() ([]BuildRecord, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	buildIDint, err := strconv.Atoi(b.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT ID, RecordID, Verified, VerifierID, VerifiedTimestamp,
			Reported, ReporterID, ReportedTimestamp, JointBuildRecord,
			JointBuildRecordID, SubmitterID, Timestamp, EditedTimestamp
		FROM BuildRecords
		WHERE BuildID = ?
	`, buildIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []BuildRecord{}
	var (
		idInt                   int
		recordIDint             int
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
			&idInt, &recordIDint, &verifiedInt, &verifierIDint,
			&verifiedTimestampString, &reportedInt, &reporterIDint,
			&reportedTimestampString, &jointBuildRecordInt, &jointBuildRecordIDint,
			&submitterIDint, &timestampString, &editedTimestampString,
		); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Convert timestamps
		if verifiedTimestamp, err = time.Parse(timeLayout, verifiedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse verified timestmap")
		}
		if reportedTimestamp, err = time.Parse(timeLayout, reportedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse reported timestmap")
		}
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestmap")
		}
		if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse edited timestmap")
		}
		// Add to results
		results = append(results, BuildRecord{
			ID:                 strconv.Itoa(idInt),
			BuildID:            b.ID,
			RecordID:           strconv.Itoa(recordIDint),
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
