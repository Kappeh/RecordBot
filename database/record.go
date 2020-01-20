package database

// Edition gets the edition of the record
func (r Record) Edition() (Edition, error) {
	return Edition{}, nil
}

// BuildClass gets the build class of the record
func (r Record) BuildClass() (BuildClass, error) {
	return BuildClass{}, nil
}

// UpdateRequestRecord gets the record which record is requesting to update
func (r Record) UpdateRequestRecord() (Record, error) {
	return Record{}, nil
}

// BuildRecord gets the build record of the record for a specified build
func (r Record) BuildRecord(buildID string) (BuildRecord, error) {
	return BuildRecord{}, nil
}

// BuildRecords gets the build records of the record for all builds
func (r Record) BuildRecords() ([]BuildRecord, error) {
	return nil, nil
}

// GuildRecordMessage gets the guild record message for the record for a specified guild
func (r Record) GuildRecordMessage(guildID string) (GuildRecordMessage, error) {
	return GuildRecordMessage{}, nil
}

// GuildRecordMessages gets the guild record message for the record for all guilds
func (r Record) GuildRecordMessages() ([]GuildRecordMessage, error) {
	return nil, nil
}
