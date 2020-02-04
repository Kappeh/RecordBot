package database

import (
	"github.com/pkg/errors"
)

// Build gets the build of the build version
func (b BuildVersion) Build() (Build, bool, error) {
	db, err := Instance()
	if err != nil {
		return Build{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	build, ok, err := db.Build(b.BuildID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to determine if build exists")
	}
	return build, ok, nil
}

// Version gets the version of the build version
func (b BuildVersion) Version() (Version, bool, error) {
	db, err := Instance()
	if err != nil {
		return Version{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	version, ok, err := db.Version(b.VersionID)
	if err != nil {
		return Version{}, false, errors.Wrap(err, "failed to determine if version exists")
	}
	return version, ok, nil
}

// Status gets the status of the build version
func (b BuildVersion) Status() (Status, bool, error) {
	db, err := Instance()
	if err != nil {
		return Status{}, false, errors.Wrap(err, "couldn't get database instance")
	}
	status, ok, err := db.Status(b.StatusID)
	if err != nil {
		return Status{}, false, errors.Wrap(err, "failed to determine if status exists")
	}
	return status, ok, nil
}
