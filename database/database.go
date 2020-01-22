package database

import (
	"github.com/pkg/errors"
)

// Instance gets the packages database instance
func Instance() (*Database, error) {
	if Error != nil {
		return nil, errors.Wrap(Error, "database connection failed during init")
	}
	return &databaseInstance, nil
}

// Close closes the database connection
func (d *Database) Close() error {
	return d.db.Close()
}

// UserStrikeCount gets the number of strikes that have been given
// to a user
func (d *Database) UserStrikeCount(userID string) (UserStrikeCount, error) {
	return UserStrikeCount{}, nil
}

// UserStrikeCounts gets the number of strikes given to each user
// that has been given at least one strike
func (d *Database) UserStrikeCounts() ([]UserStrikeCount, error) {
	return nil, nil
}

// UserStrike gets the information of a strike given to a user
func (d *Database) UserStrike(userID, strikeID string) (UserStrike, bool, error) {
	return UserStrike{}, false, nil
}

// UserStrikes gets the information of all strikes given to a user
func (d *Database) UserStrikes(userID string) ([]UserStrike, error) {
	return nil, nil
}

// UserStrikeCreate creates a strike
func (d *Database) UserStrikeCreate(userID, reason, givenByID string) (UserStrike, error) {
	return UserStrike{}, nil
}

// UserStrikeDelete a strike given to a user
func (d *Database) UserStrikeDelete(userID, strikeID string) (UserStrike, bool, error) {
	return UserStrike{}, false, nil
}

// UserStrikeEdit edits a strike given to a user
func (d *Database) UserStrikeEdit(userID, strikeID, reason string) (UserStrike, bool, error) {
	return UserStrike{}, false, nil
}

// GuildSetting gets the setting information for a guild
func (d *Database) GuildSetting(guildID string) (GuildSetting, bool, error) {
	return GuildSetting{}, false, nil
}

// GuildSettings gets the setting information for all guilds
func (d *Database) GuildSettings() ([]GuildSetting, error) {
	return nil, nil
}

// GuildSettingCreate creates setting information for a guild
func (d *Database) GuildSettingCreate(guildID, buildChannelID, ticketCategoryID string) (GuildSetting, bool, error) {
	return GuildSetting{}, false, nil
}

// GuildSettingDelete deletes the setting information for a guild
func (d *Database) GuildSettingDelete(guildID string) (GuildSetting, bool, error) {
	return GuildSetting{}, false, nil
}

// GuildSettingEdit edits the setting information for a guild
func (d *Database) GuildSettingEdit(guildID, buildChannelID, ticketCategoryID string) (GuildSetting, bool, error) {
	return GuildSetting{}, false, nil
}

// Edition gets the edition information with the specified id
func (d *Database) Edition(editionID string) (Edition, bool, error) {
	return Edition{}, false, nil
}

// Editions gets the edition information for all editions in the database
func (d *Database) Editions() ([]Edition, error) {
	return nil, nil
}

// EditionCreate creates an edition in the database
func (d *Database) EditionCreate(name, description string) (Edition, error) {
	return Edition{}, nil
}

// EditionDelete removes an edition from the database
func (d *Database) EditionDelete(editionID string) (Edition, bool, error) {
	return Edition{}, false, nil
}

// EditionEdit edits the edition information for a specified edition
func (d *Database) EditionEdit(editionID, name, description string) (Edition, bool, error) {
	return Edition{}, false, nil
}

// BuildClass gets the information for a build class in the database
func (d *Database) BuildClass(buildClassID string) (BuildClass, bool, error) {
	return BuildClass{}, false, nil
}

// BuildClasses gets the information for all build classes in the database
func (d *Database) BuildClasses() ([]BuildClass, error) {
	return nil, nil
}

// BuildClassCreate creates a new build class
func (d *Database) BuildClassCreate(name, description, embedColour string) (BuildClass, error) {
	return BuildClass{}, nil
}

// BuildClassDelete removes an existing build class
func (d *Database) BuildClassDelete(buildClassID string) (BuildClass, bool, error) {
	return BuildClass{}, false, nil
}

// BuildClassEdit edits an existing build class
func (d *Database) BuildClassEdit(buildClassID, name, description, embedColour string) (BuildClass, bool, error) {
	return BuildClass{}, false, nil
}

// RecordType gets the information for the specified record type
func (d *Database) RecordType(recordTypeID string) (RecordType, bool, error) {
	return RecordType{}, false, nil
}

// RecordTypes get the information for all record types
func (d *Database) RecordTypes() ([]RecordType, error) {
	return nil, nil
}

// RecordTypeCreate creates a new record type
func (d *Database) RecordTypeCreate(name, description string) (RecordType, error) {
	return RecordType{}, nil
}

// RecordTypeDelete removes an existing record type
func (d *Database) RecordTypeDelete(recordTypeID string) (RecordType, bool, error) {
	return RecordType{}, false, nil
}

// RecordTypeEdit edits an existing record type
func (d *Database) RecordTypeEdit(recordTypeID, name, description string) (RecordType, bool, error) {
	return RecordType{}, false, nil
}

// GuildRecordTypeChannel gets information for a specified guild and record type
func (d *Database) GuildRecordTypeChannel(guildID, recordTypeID string) (GuildRecordTypeChannel, bool, error) {
	return GuildRecordTypeChannel{}, false, nil
}

// GuildRecordTypeChannels gets information for all guilds and record types
func (d *Database) GuildRecordTypeChannels(guildID string) ([]GuildRecordTypeChannel, error) {
	return nil, nil
}

// GuildRecordTypeChannelCreate creates guild record type channel information for
// a specified guild and record type
func (d *Database) GuildRecordTypeChannelCreate(guildID, recordTypeID, channelID string) (GuildRecordTypeChannel, bool, error) {
	return GuildRecordTypeChannel{}, false, nil
}

// GuildRecordTypeChannelDelete removes guild record type channel information for
// a specified guild and record type
func (d *Database) GuildRecordTypeChannelDelete(guildID, recordTypeID string) (GuildRecordTypeChannel, bool, error) {
	return GuildRecordTypeChannel{}, false, nil
}

// GuildRecordTypeChannelEdit edits guild record type channel information for
// a specified guild and record type
func (d *Database) GuildRecordTypeChannelEdit(guildID, recordTypeID, channelID string) (GuildRecordTypeChannel, bool, error) {
	return GuildRecordTypeChannel{}, false, nil
}

// Build gets the information for a specified build
func (d *Database) Build(buildID string) (Build, bool, error) {
	return Build{}, false, nil
}

// Builds gets the information for all builds in the database
func (d *Database) Builds() ([]Build, error) {
	return nil, nil
}

// BuildCreate creates a new build
func (d *Database) BuildCreate(build Build) (Build, error) {
	return Build{}, nil
}

// BuildDelete removes build information from the database
func (d *Database) BuildDelete(buildID string) (Build, bool, error) {
	return Build{}, false, nil
}

// BuildEdit edits the information for a build in the database
func (d *Database) BuildEdit(buidID string, build Build) (Build, bool, error) {
	return Build{}, false, nil
}

// Version gets information for the specified version
func (d *Database) Version(versionID string) (Version, bool, error) {
	return Version{}, false, nil
}

// Versions gets information for all versions
func (d *Database) Versions() ([]Version, error) {
	return nil, nil
}

// VersionCreate creates a new version in the database
func (d *Database) VersionCreate(version Version) (Version, error) {
	return Version{}, nil
}

// VersionDelete removes a version from the database
func (d *Database) VersionDelete(versionID string) (Version, bool, error) {
	return Version{}, false, nil
}

// VersionEdit edits the version information for a specified version
func (d *Database) VersionEdit(versionID string, version Version) (Version, bool, error) {
	return Version{}, false, nil
}

// Record gets the information for a specified record
func (d *Database) Record(recordID string) (Record, bool, error) {
	return Record{}, false, nil
}

// Records gets information for all records in the database
func (d *Database) Records() ([]Record, error) {
	return nil, nil
}

// RecordCreate creates a new record
func (d *Database) RecordCreate(record Record) (Record, error) {
	return Record{}, nil
}

// RecordDelete removes a specified record from the database
func (d *Database) RecordDelete(recordID string) (Record, bool, error) {
	return Record{}, false, nil
}

// RecordEdit edits the information for a record in the database
func (d *Database) RecordEdit(recordID string, record Record) (Record, bool, error) {
	return Record{}, false, nil
}

// GuildBuildMessage gets the guild build message information for a specified
// guild and build
func (d *Database) GuildBuildMessage(guildID, buildID string) (GuildBuildMessage, bool, error) {
	return GuildBuildMessage{}, false, nil
}

// GuildBuildMessages get the guild build message information for a
// specified guild
func (d *Database) GuildBuildMessages(guildID string) ([]GuildBuildMessage, error) {
	return nil, nil
}

// GuildBuildMessageCreate creates guild build message information in the database
func (d *Database) GuildBuildMessageCreate(guildID, buildID, channelID, messageID string) (GuildBuildMessage, bool, error) {
	return GuildBuildMessage{}, false, nil
}

// GuildBuildMessageDelete removes guild build message information from the database
func (d *Database) GuildBuildMessageDelete(guildID, buildID string) (GuildBuildMessage, bool, error) {
	return GuildBuildMessage{}, false, nil
}

// GuildBuildMessageEdit edits the build build message information for a specified
// guild and build
func (d *Database) GuildBuildMessageEdit(guildID, buildID, channelID, messageID string) (GuildBuildMessage, bool, error) {
	return GuildBuildMessage{}, false, nil
}

// BuildVersion gets specified build version information for
// a build and a version
func (d *Database) BuildVersion(buildID, versionID string) (BuildVersion, bool, error) {
	return BuildVersion{}, false, nil
}

// BuildVersionCreate creates information in the database for a specified
// build and version
func (d *Database) BuildVersionCreate(buildID, versionID, statusID, notes string) (BuildVersion, bool, error) {
	return BuildVersion{}, false, nil
}

// BuildVersionDelete removes build version information from the database
// for a specified build and version
func (d *Database) BuildVersionDelete(buildID, versionID string) (BuildVersion, bool, error) {
	return BuildVersion{}, false, nil
}

// BuildVersionEdit edits build version information from the database
// for a specified build and version
func (d *Database) BuildVersionEdit(buildID, versionID, statusID, notes string) (BuildVersion, bool, error) {
	return BuildVersion{}, false, nil
}

// Status gets a specified status's information
func (d *Database) Status(statusID string) (Status, bool, error) {
	return Status{}, false, nil
}

// Statuses gets all statuses and their information
func (d *Database) Statuses() ([]Status, error) {
	return nil, nil
}

// StatusCreate creates a new status
func (d *Database) StatusCreate(name, description string) (Status, error) {
	return Status{}, nil
}

// StatusDelete removes a status
func (d *Database) StatusDelete(statusID string) (Status, bool, error) {
	return Status{}, false, nil
}

// StatusEdit edits a status
func (d *Database) StatusEdit(statusID, name, description string) (Status, bool, error) {
	return Status{}, false, nil
}

// BuildRecord gets build record information for a build record
func (d *Database) BuildRecord(buildRecordID string) (BuildRecord, bool, error) {
	return BuildRecord{}, false, nil
}

// BuildRecordCreate creates new build record information
func (d *Database) BuildRecordCreate(buildRecord BuildRecord) (BuildRecord, error) {
	return BuildRecord{}, nil
}

// BuildRecordDelete removes build record information from the database
func (d *Database) BuildRecordDelete(buildRecordID string) (BuildRecord, bool, error) {
	return BuildRecord{}, false, nil
}

// BuildRecordEdit edits build record information within the database
func (d *Database) BuildRecordEdit(buildRecordID string, buildRecord BuildRecord) (BuildRecord, bool, error) {
	return BuildRecord{}, false, nil
}

// GuildRecordMessage gets the guild record message information for a specified
// guild and record
func (d *Database) GuildRecordMessage(guildID, recordID string) (GuildBuildMessage, bool, error) {
	return GuildBuildMessage{}, false, nil
}

// GuildRecordMessages gets the guild record message information for a specified guild
func (d *Database) GuildRecordMessages(guildID string) ([]GuildRecordMessage, error) {
	return nil, nil
}

// GuildRecordMessageCreate creates guild record message information for a specified
// guild and record
func (d *Database) GuildRecordMessageCreate(guildID, recordID, channelID, messageID string) (GuildRecordMessage, bool, error) {
	return GuildRecordMessage{}, false, nil
}

// GuildRecordMessageDelete removes guild record message information for a specified
// guild and record
func (d *Database) GuildRecordMessageDelete(guildID, recordID string) (GuildRecordMessage, bool, error) {
	return GuildRecordMessage{}, false, nil
}

// GuildRecordMessageEdit edits guild record message information for a specified
// guild and record
func (d *Database) GuildRecordMessageEdit(guildID, recordID, channelID, messageID string) (GuildRecordMessage, bool, error) {
	return GuildRecordMessage{}, false, nil
}

// GuildTicketChannel gets information for a specified ticket within a guild
func (d *Database) GuildTicketChannel(guildID, channelID string) (GuildTicketChannel, bool, error) {
	return GuildTicketChannel{}, false, nil
}

// GuildTicketChannels gets information for all tickets within a guild
func (d *Database) GuildTicketChannels(guildID string) ([]GuildTicketChannel, error) {
	return nil, nil
}

// GuildTicketChannelCreate creates a new ticket channel within the database
func (d *Database) GuildTicketChannelCreate(guildID, channelID string, ticketType TicketType) (GuildTicketChannel, bool, error) {
	return GuildTicketChannel{}, false, nil
}

// GuildTicketChannelDelete removes an existing ticket channel from the database
func (d *Database) GuildTicketChannelDelete(guildID, channelID string) (GuildTicketChannel, bool, error) {
	return GuildTicketChannel{}, false, nil
}

// GuildTicketChannelEdit edits the ticket channel information for a specified ticket
func (d *Database) GuildTicketChannelEdit(guildID, channelID string, ticketType TicketType) (GuildTicketChannel, bool, error) {
	return GuildTicketChannel{}, false, nil
}
