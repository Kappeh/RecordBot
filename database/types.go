package database

import (
	"database/sql"
	"time"
)

// Database is an instance of a database connection
type Database struct{ db *sql.DB }

// Timestamp is a time
type Timestamp time.Time

// TicketType indicates what type of ticket a guild ticket channel is
type TicketType int

const (
	// TicketGeneral is a ticket created for talking in private
	// to guild administrators and moderators (usually for help)
	TicketGeneral TicketType = iota
	// TicketSubmitBuild is a ticket which guides a user through
	// the process of submitting a build
	TicketSubmitBuild
	// TicketSubmitRecord is a ticket which guides a user through
	// the process of submitting a record
	TicketSubmitRecord
	// TicketSubmitBuildUpdate is a ticket which guides a user
	// through the process of submitting an update to a build
	TicketSubmitBuildUpdate
	// TicketSubmitRecordUpdate is a ticket which guides a user
	// through the process of submitting an update to a record
	TicketSubmitRecordUpdate
)

// Table Elements

// UserStrike is a strike which has been given to a user
type UserStrike struct {
	// UserID is the id of the user that the strike belongs to
	UserID string
	// StrikeID is the unique id of the strike for the user
	// This is not the number of strikes the user has had
	// and do not treat it as such
	StrikeID string

	// Reason is the reason given by the author for the strike
	// being given to the user
	Reason string

	// AuthorID is the id of the user that gave the strike
	AuthorID string

	// Timestamp is the time which the strike was initially given
	Timestamp Timestamp
	// EditedTimestamp is the time which the strike was last edited
	EditedTimestamp Timestamp
}

// GuildSetting contains guild specific settings
type GuildSetting struct {
	// GuildID is the id of the guild the settings apply to
	GuildID string

	// BuildChannelID is the id of the channel within the guild
	// in which build messages are to be sent to and reside in
	BuildChannelID string
	// TicketChannelCategoryID is the id of the category in which
	// ticket channels are to be made within
	TicketChannelCategoryID string

	// Timestamp is the time in which the guild setting was first
	// created. This coinsides with when the application was
	// first initializes (joined) for a specific guild
	Timestamp Timestamp
	// EditedTimestamp is the time which the guild settings
	// were last edited
	EditedTimestamp Timestamp
}

// Edition is a Minecraft edition
type Edition struct {
	// ID is the id of the edition within the database
	ID string

	// Name is the name of the edition
	// e.g. Door, Logic, Farm, etc
	Name string
	// Description is a description of the edition
	Description string

	// Timestamp is the time in which the edition was created
	Timestamp Timestamp
	// EditedTimestamp is the time in which the edition was last edited
	EditedTimestamp Timestamp
}

// BuildClass is a classification category of a build
type BuildClass struct {
	// ID is the id of the build class
	ID string

	// Name is the name of the build class
	Name string
	// Description is a description of the build class
	Description string
	// EmbedColour is the colour to be used for discord
	// embeds for messages related to the build class
	EmbedColour string

	// Timestamp is the time which the build class was created
	Timestamp Timestamp
	// EditedTimestamp is the time which the build class was last edited
	EditedTimestamp Timestamp
}

// RecordType is a major classification category for records
// e.g. Smallest, Fastest, Smallest Overserless, etc
type RecordType struct {
	// ID is the id of the record type in the database
	ID string

	// Name is the name of the record type
	Name string
	// Description is a description of the record type
	Description string

	// Timestamp is the time which the record type was created
	Timestamp Timestamp
	// EditedTimestamp is the time which the record type was last edited
	EditedTimestamp Timestamp
}

// GuildRecordTypeChannel contains the discord channel in which
// records of a specific type should be posted to within a
// specific guild
type GuildRecordTypeChannel struct {
	// GuildID is the id of the guild
	GuildID string
	// RecordTypeID is the id of the record type
	RecordTypeID string

	// ChannelID is the id of the discord channel which records
	// should be posted to
	ChannelID string

	// Timestamp is the time which the guild record type channel
	// was created
	Timestamp Timestamp
	// EditedTimestamp is the time which the guild record type channel
	// was last updated
	EditedTimestamp Timestamp
}

// Build contains relevant information about a minecraft redston build
type Build struct {
	// ID is the id of the build in the database
	ID string

	// Verified indicates whether the build has been verified
	// by a system administrator/moderator
	Verified bool
	// VerifierID is the id of the user that verified the build
	// if the build has been verified
	VerifierID string
	// VerifiedTimestamp is the time which the build was verified
	// if the build has been verified
	VerifiedTimestamp Timestamp

	// Reported indicates if the build has been reported as fake
	Reported bool
	// ReporterID is the id of the user that reported the build as
	// fake if the build has been reported as fake
	ReporterID string
	// ReportedTimestamp is the time which the build was reported
	// as fake if the build has been reported as fake
	ReportedTimestamp Timestamp

	// UpdateRequest indicates whether this build is an updated
	// version of an existing build in the database
	UpdateRequest bool
	// UpdateRequestBuildID is the id of the build in the database
	// which this build is an update of if this build is an updated
	// version of an existing build in the database
	UpdateRequestBuildID string

	// EditionID is the id of the edition which this build was
	// built in
	EditionID string
	// BuildClassID is the id of the build class this build
	// is classified into
	BuildClassID string
	// Name is the name of this build
	Name string
	// Description is a description of the build
	Description string
	// Creators is a string representation of a list
	// of peoples' in game names which created this build
	Creators string
	// CreationTimestamp is the time which the build was created
	CreationTimestamp Timestamp

	// Width is the width of the build
	Width int
	// Height is the height of the build
	Height int
	// Depth is the depth of the build
	Depth int

	// NormalCloseDuration is the normal closing time (in gameticks)
	// in accordance with 2.1.2.3 - b.i
	NormalCloseDuration int
	// NormalOpenDuration is the normal opening time (in gameticks)
	// in accordance with 2.1.2.3 - b.ii
	NormalOpenDuration int
	// VisibleCloseDuration is the visible closing time (in gameticks)
	// in accordance with 2.1.2.3 - a.i
	VisibleCloseDuration int
	// VisibleOpenDuration is the visible opening time (in gameticks)
	// in accordance with 2.1.2.3 - a.ii
	VisibleOpenDuration int
	// DelayCloseDuration is the closing input delay (in gameticks)
	// in accordance with 2.1.2.3 - c.i
	DelayCloseDuration int
	// DelayOpenDuration is the opening input delay (in gameticks)
	// in accordance with 2.1.2.3 - c.ii
	DelayOpenDuration int
	// ResetCloseDuration is the closing reset time (in gameticks)
	// in accordance with 2.1.2.3 - d.i
	ResetCloseDuration int
	// ResetOpenDuration is the opening reset time (in gameticks)
	// in accordance with 2.1.2.3 - d.ii
	ResetOpenDuration int

	// ExtensionDuration is the normal extension time (in gameticks)
	// in accordance with 2.3.2.2 - a.i
	ExtensionDuration int
	// RetractionDuration is the normal retraction time (in gameticks)
	// in accordance with 2.3.2.2 - a.ii
	RetractionDuration int
	// ExtensionDelayDuration is the extension input delay (in gameticks)
	// in accordance with 2.3.2.2 - b.i
	ExtensionDelayDuration int
	// RetractionDelayDuration is the retraction input delay (in gameticks)
	// in accordance with 2.3.2.2 - b.ii
	RetractionDelayDuration int

	// ImageURL is the URL to an image of the build
	ImageURL string
	// YoutubeURL is the URL to a youtube video of the build
	YoutubeURL string
	// WorldDownloadURL is the URL to a world download for the build
	WorldDownloadURL string
	// ServerIPAddress is the ip address of a server which contains the build
	ServerIPAddress string
	// ServerCoordinates is the coordinates of the build on the
	// server specified by ServerIPAdress
	ServerCoordinates string
	// ServerCommand is a command that can be used to get to the build
	// on the server specified by ServerIPAdress
	ServerCommand string

	// SubmitterID is the id of the user which submitted this build
	SubmitterID string

	// Timestamp is the time which the build was first created
	Timestamp Timestamp
	// EditedTimestamp is the time which the build was last edted
	EditedTimestamp Timestamp
}

// Version is a Minecraft version
type Version struct {
	// ID is the id of the version in the database
	ID string

	// EditionID is the id of the edition of the version
	EditionID string
	// MajorVersion is the major version of the version
	MajorVersion int
	// MinorVersion is the minor version of the version
	MinorVersion int
	// Patch if the patch version of the version
	Patch int
	// Name is the name of the version
	Name string
	// Description is a description of the version
	Description string

	// VersionTimestamp is the time the version was released
	VersionTimestamp Timestamp

	// Timestamp is the time the version was added to the database
	Timestamp Timestamp
	// EditedTimestamp is the time the version was last edited
	// in the database
	EditedTimestamp Timestamp
}

// Record is a record category that ranks builds in an order
// relative to some criteria
type Record struct {
	// ID is the id of the record in the database
	ID string

	// Verified indicates whether the record has been verified
	// by a system administrator/moderator
	Verified bool
	// VerifierID is the id of the user that verified the record
	// if the record has been verified
	VerifierID string
	// VerifiedTimestamp is the time the record was verified
	// if the record has been verified
	VerifiedTimestamp Timestamp

	// UpdateRequest indicates whether this record is an updated
	// version of an existing record in the database
	UpdateRequest bool
	// UpdateRequestRecordID is the id of the record in the database
	// that this record is an updated version of if this record
	// is an updated version of an existing record in the database
	UpdateRequestRecordID string

	// EditionID is the id of the edition this record falls within
	EditionID string
	// BuildClassID is the id of the build class this record falls within
	BuildClassID string
	// RecordTypeID is the id of the record type this record falls within
	RecordTypeID string
	// Name is the name of the record
	Name string
	// Description is a description of the record
	Description string

	// SubmitterID is the id of the user that submitted the record
	SubmitterID string

	// Timestamp is the time this record was created in the database
	Timestamp Timestamp
	// EditedTimestamp is the time this record was last edited
	// within the database
	EditedTimestamp Timestamp
}

// GuildBuildMessage contains the information about the discord message
// within a guild which displays this build
type GuildBuildMessage struct {
	// GuildID is the id of the guild the message is in
	GuildID string
	// BuildID is the id of the build within the message
	BuildID string

	// ChannelID is the id of the channel containing the message
	// that contains the build
	ChanneID string
	// MessageID is the id of the message contaning the build
	MessageID string

	// Timestamp is the time the guild build message was created
	Timestamp Timestamp
	// EditedTimestamp is the time the guild build message was last edited
	EditedTimestamp Timestamp
}

// BuildVersion is an indicator of the status of a build within a Minecraft
// version. Notes can be included to describe the status is more detail
type BuildVersion struct {
	// BuildID is the id of the build within the database
	BuildID string
	// VersionID is the id of the version within the database
	VersionID string

	// StatusID is the id of the status within the database of the
	// build within the specified version
	StatusID string
	// Notes is a plaintext set of notes describing the status of the
	// build in more detail
	Notes string

	// Timestamp is the time the build version was created
	Timestamp Timestamp
	// EditedTimestamp is the time the build version was last edited
	EditedTimestamp Timestamp
}

// Status indicates the status of a build
type Status struct {
	// ID is the id of the status within the database
	ID string

	// Name is the name of the status
	// e.g. Working, Broken, etc
	Name string
	// Description is a description of the status
	Description string

	// Timestamp is the time the status was created
	Timestamp Timestamp
	// EditedTimestamp is the time the status was last edited
	EditedTimestamp Timestamp
}

// BuildRecord indicates a build which holds/held a record
type BuildRecord struct {
	// ID is the id of the build record in the database
	ID string

	// BuildID is the id of the build in the database
	// that holds/held the record
	BuildID string
	// RecordID is the id of the record in the database
	// that the build holds/held
	RecordID string

	// Verified indicates whether the build record has been
	// verified by a system administrator/moderator
	Verified bool
	// VerifierID is the id of the user that verified the build record
	// if the build record has been verified
	VerifierID string
	// VerifiedTimestamp is the time the build record was verified
	// if the build record was verified
	VerifiedTimestamp Timestamp

	// Reported indicates whether the build record has
	// been reported as fake
	Reported bool
	// ReporterID is the id of the user that reported the build record
	// as fake if the build record has been reported as fake
	ReporterID string
	// ReportedTimestamp is the time the build record was reported
	// as fake is the build record has been reported as fake
	ReportedTimestamp Timestamp

	// JointBuildRecord indicated whether this build record ties with another
	// build record in the database (The first build record submitted
	// that holds this record will not be indicated as a join record
	// through JoinRecord)
	JointBuildRecord bool
	// JoinBuildRecordID is the id of the first build record submitted
	// that ties with this record if this build record ties with another
	// build record in the database
	JointBuildRecordID string

	// SubmitterID is the if of the user that submitted this build record
	SubmitterID string

	// Timestamp is the time the build record was created
	Timestamp Timestamp
	// EditedTimestamp is the time the build record was last edited
	EditedTimestamp Timestamp
}

// GuildRecordMessage contains information about the discord message
// that contains that displays a specific record's information
type GuildRecordMessage struct {
	// GuildID is the id of the guild the message is within
	GuildID string
	// RecordID is the id of the record which has its information
	// displayed in the message
	RecordID string

	// ChannelID is the id of the channel which contains the message
	// which contains the record information
	ChannelID string
	// MessageID is the id of the message which contains the record
	// information
	MessageID string

	// Timestamp is the time the guild record message was created
	Timestamp Timestamp
	// EditedTimestamp is the time the guild record message was last edited
	EditedTimestamp Timestamp
}

// GuildTicketChannel is a ticket within a discord guild
type GuildTicketChannel struct {
	// GuildID is the id of the discord guild that contains the ticket
	GuildID string
	// ChannelID is the id of the channel marked as a ticket
	ChannelID string

	// TicketID is the id of the ticket within the guild
	TicketID string
	// TicketType is the type of the ticket
	TicketType TicketType

	// CreatorID is the id of the user that created the ticket
	CreatorID string

	// Timestamp is the time the ticket was opened
	Timestamp Timestamp
}

// Other Elements

// UserStrikeCount indicates how many strikes a user has
type UserStrikeCount struct {
	// UserID is the id of the user
	UserID string

	// Count is the amount of strikes the user has
	Count int
}
