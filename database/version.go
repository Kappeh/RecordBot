package database

import (
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// Edition gets the edition of the version
func (v Version) Edition() (Edition, bool, error) {
	db, err := Instance()
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	e, ok, err := db.Edition(v.EditionID)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to determine if edition exists")
	}
	return e, ok, nil
}

// BuildVersion gets the build version of the version for a specified build
func (v Version) BuildVersion(buildID string) (BuildVersion, bool, error) {
	db, err := Instance()
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	bv, ok, err := db.BuildVersion(buildID, v.ID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to determine if build version exists")
	}
	return bv, ok, nil
}

// BuildVersions gets the build versions of the version for all builds
func (v Version) BuildVersions() ([]BuildVersion, error) {
	db, err := Instance()
	if err != nil {
		return nil, errors.Wrap(err, "couldn't get database instance")
	}
	// Convert id to int
	versionIDint, err := strconv.Atoi(v.ID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert version id to integer")
	}
	// Query the database
	rows, err := db.db.Query(`
		SELECT BuildID, StatusID, Notes, Timestamp, EditedTimestamp
		FROM BuildVersions
		WHERE VersionID = ?
	`, versionIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []BuildVersion{}
	var (
		buildIDint            int
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
			&buildIDint, &statusIDint, &notes,
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
			BuildID:         strconv.Itoa(buildIDint),
			VersionID:       v.ID,
			StatusID:        strconv.Itoa(statusIDint),
			Notes:           notes,
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}
