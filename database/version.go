package database

// Edition gets the edition of the version
func (v Version) Edition() (Edition, error) {
	return Edition{}, nil
}

// BuildVersion gets the build version of the version for a specified build
func (v Version) BuildVersion(buildID string) (BuildVersion, error) {
	return BuildVersion{}, nil
}

// BuildVersions gets the build versions of the version for all builds
func (v Version) BuildVersions() ([]BuildVersion, error) {
	return nil, nil
}
