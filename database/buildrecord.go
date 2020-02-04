package database

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Build gets the build of the build record
func (b BuildRecord) Build() (Build, bool, error) {
	db, err := Instance()
	if err != nil {
		return Build{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	build, ok, err := db.Build(b.BuildID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed getting build")
	}
	return build, ok, nil
}

// Record gets the record of the build record
func (b BuildRecord) Record() (Record, bool, error) {
	db, err := Instance()
	if err != nil {
		return Record{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	record, ok, err := db.Record(b.RecordID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "couldn't get record")
	}
	return record, ok, nil
}

// FirstJointBuildRecord gets the first joint build record
// It get's the root node of a dependency tree of build records
func (b BuildRecord) FirstJointBuildRecord() (BuildRecord, bool, error) {
	db, err := Instance()
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	buildRecordIDint, err := strconv.Atoi(b.ID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		WITH CTE (RootID, BuildID, RecordID, Verified, VerifierID, VerifiedTimestamp, 
				Reported, ReporterID, ReportedTimestamp, JointBuildRecord,
				JointBuildRecordID, SubmitterID, Timestamp, EditedTimestamp, LeafID)		
		AS (
			SELECT ID, BuildID, RecordID, Verified, VerifierID, VerifiedTimestamp, 
				Reported, ReporterID, ReportedTimestamp, JointBuildRecord,
				JointBuildRecordID, SubmitterID, Timestamp, EditedTimestamp, ID
			FROM BuildRecords
			WHERE JointBuildRecord = 0
			UNION ALL
			SELECT CTE.RootID, CTE.BuildID, CTE.RecordID,
				CTE.Verified, CTE.VerifierID, CTE.VerifiedTimestamp,
				CTE.Reported, CTE.ReporterID, CTE.ReportedTimestamp,
				CTE.JointBuildRecord, CTE.JointBuildRecordID,
				CTE.SubmitterID, CTE.Timestamp, CTE.EditedTimestamp,
				BuildRecords.ID
			FROM BuildRecords INNER JOIN CTE
			ON BuildRecords.JointBuildRecordID = CTE.LeafID AND BuildRecords.JointBuildRecord = 1
		)
		SELECT RootID, BuildID, RecordID, Verified, VerifierID, VerifiedTimestamp, 
			Reported, ReporterID, ReportedTimestamp, JointBuildRecord,
			JointBuildRecordID, SubmitterID, Timestamp, EditedTimestamp
		FROM CTE
		WHERE LeafID = ?
	`, buildRecordIDint)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if row exists
	// Row should always exist
	if !rows.Next() {
		return BuildRecord{}, false, nil
	}
	// Extract data
	var (
		idInt                   int
		buildIDint              int
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
	if err = rows.Scan(
		&idInt, &buildIDint, &recordIDint, &verifiedInt, &verifierIDint,
		&verifiedTimestampString, &reportedInt, &reporterIDint,
		&reportedTimestampString, &jointBuildRecordInt, &jointBuildRecordIDint,
		&submitterIDint, &timestampString, &editedTimestampString,
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
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return BuildRecord{
		ID:                 strconv.Itoa(idInt),
		BuildID:            strconv.Itoa(buildIDint),
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
	}, true, nil
}

// JointBuildRecords gets all joint build records
func (b BuildRecord) JointBuildRecords() ([]BuildRecord, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to ints
	buildRecordIDint, err := strconv.Atoi(b.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	// TODO: Check if the nested 'SELECT ... FROM CTE'
	// causes a performance issue
	rows, err := db.db.Query(`
		WITH CTE (ID, BuildID, RecordID, Verified, VerifierID, VerifiedTimestamp, 
			Reported, ReporterID, ReportedTimestamp, JointBuildRecord,
			JointBuildRecordID, SubmitterID, Timestamp, EditedTimestamp, rootID)
		AS (
			SELECT ID, BuildID, RecordID, Verified, VerifierID, VerifiedTimestamp, 
				Reported, ReporterID, ReportedTimestamp, JointBuildRecord,
				JointBuildRecordID, SubmitterID, Timestamp, EditedTimestamp, ID
			FROM BuildRecords
			WHERE JointBuildRecord = 0
			UNION ALL
			SELECT BuildRecords.ID, BuildRecords.BuildID, BuildRecords.RecordID,
				BuildRecords.Verified, BuildRecords.VerifierID, BuildRecords.VerifiedTimestamp,
				BuildRecords.Reported, BuildRecords.ReporterID, BuildRecords.ReportedTimestamp,
				BuildRecords.JointBuildRecord, BuildRecords.JointBuildRecordID,
				BuildRecords.SubmitterID, BuildRecords.Timestamp, BuildRecords.EditedTimestamp,
				CTE.RootID
			FROM BuildRecords INNER JOIN CTE
			ON BuildRecords.JointBuildRecordID = CTE.ID AND BuildRecords.JointBuildRecord = 1
		)
		SELECT ID, BuildID, RecordID, Verified, VerifierID, VerifiedTimestamp, 
			Reported, ReporterID, ReportedTimestamp, JointBuildRecord,
			JointBuildRecordID, SubmitterID, Timestamp, EditedTimestamp
		FROM CTE
		WHERE RootID = (
			SELECT RootID
			FROM CTE
			WHERE ID = ?
		)
	`, buildRecordIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []BuildRecord{}
	var (
		idInt                   int
		buildIDint              int
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
		// Extract the data
		if err = rows.Scan(
			&idInt, &buildIDint, &recordIDint, &verifiedInt, &verifierIDint,
			&verifiedTimestampString, &reportedInt, &reporterIDint,
			&reportedTimestampString, &jointBuildRecordInt, &jointBuildRecordIDint,
			&submitterIDint, &timestampString, &editedTimestampString,
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
