package database

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Versions gets all versions for the edition
func (e Edition) Versions() ([]Version, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	idInt, err := strconv.Atoi(e.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT ID, MajorVersion, MinorVersion, Patch, Name, Description,
			VersionTimestamp, Timestamp, EditedTimestamp
		FROM Versions
		WHERE EditionID = ?
	`, idInt)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []Version{}
	var (
		versionIDint           int
		majorVersion           int
		minorVersion           int
		patch                  int
		name                   string
		description            string
		versionTimestampString string
		timestampString        string
		editedTimestampString  string
		versionTimestamp       time.Time
		timestamp              time.Time
		editedTimestamp        time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&versionIDint, &majorVersion, &minorVersion, &patch, &name,
			&description, &versionTimestampString, &timestampString,
			&editedTimestampString,
		); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Parse timestamps
		if versionTimestamp, err = time.Parse(timeLayout, versionTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse version timestamp")
		}
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}
		if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse edited timestamp")
		}
		// Add to results
		results = append(results, Version{
			ID:               strconv.Itoa(idInt),
			EditionID:        e.ID,
			MajorVersion:     majorVersion,
			MinorVersion:     minorVersion,
			Patch:            patch,
			Name:             name,
			Description:      description,
			VersionTimestamp: Timestamp(versionTimestamp),
			Timestamp:        Timestamp(timestamp),
			EditedTimestamp:  Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// Builds gets all builds in the edition
func (e Edition) Builds() ([]Build, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	idInt, err := strconv.Atoi(e.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT ID, Verified, VerifierID, VerifiedTimestamp, Reported,
			ReporterID, ReportedTimestamp, UpdateRequest, UpdateRequestBuildID,
			BuildClassID, Name, Description, Creators, CreationTimestamp,
			Width, Height, Depth, NormalCloseDuration, NormalOpenDuration,
			VisibleCloseDuration, VisibleOpenDuration, DelayCloseDuration,
			DelayOpenDuration, ResetCloseDuration, ResetOpenDuration,
			ExtensionDuration, RetractionDuration, ExtensionDelayDuration,
			RetractionDelayDuration, ImageURL, YoutubeURL, WorldDownloadURL,
			ServerIPAddress, ServerCoordinates, ServerCommand, SubmitterID,
			Timestamp, EditedTimestamp
		FROM Builds
		WHERE EditionID = ?
	`, idInt)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []Build{}
	var (
		buildIDint              int
		verifiedInt             int
		verifierIDint           int
		verifiedTimestampString string
		reportedInt             int
		reporterIDint           int
		reportedTimestampString string
		updateRequestInt        int
		updateRequestBuildIDint int
		buildClassIDint         int
		name                    string
		description             string
		creators                string
		creationTimestampString string
		width                   int
		height                  int
		depth                   int
		normalCloseDuration     int
		normalOpenDuration      int
		visibleCloseDuration    int
		visibleOpenDuration     int
		delayCloseDuration      int
		delayOpenDuration       int
		resetCloseDuration      int
		resetOpenDuration       int
		extensionDuration       int
		retractionDuration      int
		extensionDelayDuration  int
		retractionDelayDuration int
		imageURL                string
		youtubeURL              string
		worldDownloadURL        string
		serverIPAddress         string
		serverCoordinates       string
		serverCommand           string
		submitterIDint          int
		timestampString         string
		editedTimestampString   string
		verifiedTimestamp       time.Time
		reportedTimestamp       time.Time
		creationTimestamp       time.Time
		timestamp               time.Time
		editedTimestamp         time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&buildIDint, &verifiedInt, &verifierIDint, &verifiedTimestampString, &reportedInt,
			&reporterIDint, &reportedTimestampString, &updateRequestInt, &updateRequestBuildIDint,
			&buildClassIDint, &name, &description, &creators, &creationTimestampString, &width,
			&height, &depth, &normalCloseDuration, &normalOpenDuration, &visibleCloseDuration,
			&visibleOpenDuration, &delayCloseDuration, &delayOpenDuration, &resetCloseDuration,
			&resetOpenDuration, &extensionDuration, &retractionDuration, &extensionDelayDuration,
			&retractionDelayDuration, &imageURL, &youtubeURL, &worldDownloadURL, &serverIPAddress,
			&serverCoordinates, &serverCommand, &submitterIDint, &timestampString, &editedTimestampString,
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
		if creationTimestamp, err = time.Parse(timeLayout, creationTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse creation timestamp")
		}
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}
		if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse edited timestamp")
		}
		// Add to results
		results = append(results, Build{
			ID:                      strconv.Itoa(buildIDint),
			Verified:                verifiedInt != 0,
			VerifierID:              strconv.Itoa(verifierIDint),
			VerifiedTimestamp:       Timestamp(verifiedTimestamp),
			Reported:                reportedInt != 0,
			ReporterID:              strconv.Itoa(reporterIDint),
			ReportedTimestamp:       Timestamp(reportedTimestamp),
			UpdateRequest:           updateRequestInt != 0,
			UpdateRequestBuildID:    strconv.Itoa(updateRequestBuildIDint),
			EditionID:               e.ID,
			BuildClassID:            strconv.Itoa(buildClassIDint),
			Name:                    name,
			Description:             description,
			Creators:                creators,
			CreationTimestamp:       Timestamp(creationTimestamp),
			Width:                   width,
			Height:                  height,
			Depth:                   depth,
			NormalCloseDuration:     normalCloseDuration,
			NormalOpenDuration:      normalOpenDuration,
			VisibleCloseDuration:    visibleCloseDuration,
			VisibleOpenDuration:     visibleOpenDuration,
			DelayCloseDuration:      delayCloseDuration,
			DelayOpenDuration:       delayOpenDuration,
			ResetCloseDuration:      resetCloseDuration,
			ResetOpenDuration:       resetOpenDuration,
			ExtensionDuration:       extensionDuration,
			RetractionDuration:      retractionDuration,
			ExtensionDelayDuration:  extensionDelayDuration,
			RetractionDelayDuration: retractionDelayDuration,
			ImageURL:                imageURL,
			YoutubeURL:              youtubeURL,
			WorldDownloadURL:        worldDownloadURL,
			ServerIPAddress:         serverIPAddress,
			ServerCoordinates:       serverCoordinates,
			ServerCommand:           serverCommand,
			SubmitterID:             strconv.Itoa(submitterIDint),
			Timestamp:               Timestamp(timestamp),
			EditedTimestamp:         Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// Records gets all records in the edition
func (e Edition) Records() ([]Record, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to integer
	idInt, err := strconv.Atoi(e.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT ID, Verified, VerifierID, VerifiedTimestamp, UpdateRequest,
			UpdateRequestRecordID, BuildClassID, RecordTypeID, Name,
			Description, SubmitterID, Timestamp, EditedTimestamp
		FROM Records
		WHERE EditionID = ?
	`, idInt)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []Record{}
	var (
		recordIDint              int
		verifiedInt              int
		verifierIDint            int
		verifiedTimestampString  string
		updateRequestInt         int
		updateRequestRecordIDint int
		buildClassIDint          int
		recordTypeIDint          int
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
			&recordIDint, &verifiedInt, &verifierIDint, &verifiedTimestampString,
			&updateRequestInt, &updateRequestRecordIDint, &buildClassIDint,
			&recordTypeIDint, &name, &description, &submitterIDint,
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
			ID:                    strconv.Itoa(recordIDint),
			Verified:              verifiedInt != 0,
			VerifierID:            strconv.Itoa(verifierIDint),
			VerifiedTimestamp:     Timestamp(verifiedTimestamp),
			UpdateRequest:         updateRequestInt != 0,
			UpdateRequestRecordID: strconv.Itoa(updateRequestRecordIDint),
			EditionID:             e.ID,
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
