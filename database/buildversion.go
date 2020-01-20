package database

// Build gets the build of the build version
func (b BuildVersion) Build() (Build, error) {
	return Build{}, nil
}

// Version gets the version of the build version
func (b BuildVersion) Version() (Version, error) {
	return Version{}, nil
}

// Status gets the status of the build version
func (b BuildVersion) Status() (Status, error) {
	return Status{}, nil
}
