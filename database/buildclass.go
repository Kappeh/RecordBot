package database

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Builds gets all of the build of the build class
func (b BuildClass) Builds() ([]Build, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	buildClassIDint, err := strconv.Atoi(b.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT ID, Verified, VerifierID, VerifiedTimestamp, Reported,
			ReporterID, ReportedTimestamp, UpdateRequest, UpdateRequestBuildID,
			EditionID, Name, Description, Creators, CreationTimestamp, Width,
			Height, Depth, NormalCloseDuration, NormalOpenDuration,
			VisibleCloseDuration, VisibleOpenDuration, DelayCloseDuration,
			DelayOpenDuration, ResetCloseDuration, ResetOpenDuration,
			ExtensionDuration, RetractionDuration, ExtensionDelayDuration,
			RetractionDelayDuration, ImageURL, YoutubeURL, WorldDownloadURL,
			ServerIPAddress, ServerCoordinates, ServerCommand, SubmitterID,
			Timestamp, EditedTimestamp
		FROM Builds
		WHERE BuildClassID = ?
	`, buildClassIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []Build{}
	var (
		idInt                   int
		verifiedInt             int
		verifierIDint           int
		verifiedTimestampString string
		reportedInt             int
		reporterIDint           int
		reportedTimestampString string
		updateRequestInt        int
		updateRequestBuildIDint int
		editionIDint            int
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
			&idInt, &verifiedInt, &verifierIDint, &verifiedTimestampString,
			&reportedInt, &reporterIDint, &reportedTimestampString, &updateRequestInt,
			&updateRequestBuildIDint, &editionIDint, &name, &description, &creators,
			&creationTimestampString, &width, &height, &depth, &normalCloseDuration,
			&normalOpenDuration, &visibleCloseDuration, &visibleOpenDuration,
			&delayCloseDuration, &delayOpenDuration, &resetCloseDuration,
			&resetOpenDuration, &extensionDuration, &retractionDuration,
			&extensionDelayDuration, &retractionDelayDuration, &imageURL, &youtubeURL,
			&worldDownloadURL, &serverIPAddress, &serverCoordinates, &serverCommand,
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
			ID:                      strconv.Itoa(idInt),
			Verified:                verifiedInt != 0,
			VerifierID:              strconv.Itoa(verifierIDint),
			VerifiedTimestamp:       Timestamp(verifiedTimestamp),
			Reported:                reportedInt != 0,
			ReporterID:              strconv.Itoa(reporterIDint),
			ReportedTimestamp:       Timestamp(reportedTimestamp),
			UpdateRequest:           updateRequestInt != 0,
			UpdateRequestBuildID:    strconv.Itoa(updateRequestBuildIDint),
			EditionID:               strconv.Itoa(editionIDint),
			BuildClassID:            b.ID,
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

// Records get all of the records of the build class
func (b BuildClass) Records() ([]Record, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	buildClassIDint, err := strconv.Atoi(b.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT ID, Verified, VerifierID, VerifiedTimestamp, UpdateRequest,
			UpdateRequestRecordID, EditionID, RecordTypeID, Name,
			Description, SubmitterID, Timestamp, EditedTimestamp
		FROM Records
		WHERE BuildClassID = ?
	`, buildClassIDint)
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
			&idInt, &verifiedInt, &verifierIDint, &verifiedTimestampString,
			&updateRequestInt, &updateRequestRecordIDint, &editionIDint,
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
			ID:                    strconv.Itoa(idInt),
			Verified:              verifiedInt != 0,
			VerifierID:            strconv.Itoa(verifierIDint),
			VerifiedTimestamp:     Timestamp(verifiedTimestamp),
			UpdateRequest:         updateRequestInt != 0,
			UpdateRequestRecordID: strconv.Itoa(updateRequestRecordIDint),
			EditionID:             strconv.Itoa(editionIDint),
			BuildClassID:          b.ID,
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
