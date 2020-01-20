package database

// GuildRecordTypeChannel gets the guild record type channel for the record type for
// a specified guild
func (r RecordType) GuildRecordTypeChannel(guildID string) (GuildRecordTypeChannel, error) {
	return GuildRecordTypeChannel{}, nil
}

// GuildRecordTypeChannels gets the guild record type channels for the record type
// for all guilds
func (r RecordType) GuildRecordTypeChannels() ([]GuildRecordTypeChannel, error) {
	return nil, nil
}

// Records get the all the records that fall into the record type
func (r RecordType) Records() ([]Record, error) {
	return nil, nil
}
