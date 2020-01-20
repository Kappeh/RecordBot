package database

// Edition gets the edition of the build
func (b Build) Edition() (Edition, error) {
	return Edition{}, nil
}

// BuildClass gets the build class of the build
func (b Build) BuildClass() (BuildClass, error) {
	return BuildClass{}, nil
}

// UpdateRequestBuild get the build which is being requested to update
func (b Build) UpdateRequestBuild() (Build, error) {
	return Build{}, nil
}

// GuildBuildMessage gets the guild build message for a specified guild
func (b Build) GuildBuildMessage(guildID string) (GuildBuildMessage, error) {
	return GuildBuildMessage{}, nil
}

// GuildBuildMessages get the guild build messages for all guilds
func (b Build) GuildBuildMessages() ([]GuildBuildMessage, error) {
	return nil, nil
}

// BuildVersion gets the build version for a specified version
func (b Build) BuildVersion(versionID string) (BuildVersion, error) {
	return BuildVersion{}, nil
}

// BuildVersions gets the build versions for all versions
func (b Build) BuildVersions() ([]BuildVersion, error) {
	return nil, nil
}

// BuildRecord gets the build record for the build and a specified record
func (b Build) BuildRecord(recordID string) (BuildRecord, error) {
	return BuildRecord{}, nil
}

// BuildRecords gets the build records for the build and all records
func (b Build) BuildRecords() ([]BuildRecord, error) {
	return nil, nil
}
