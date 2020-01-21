package database

import (
	"database/sql"
	"os"
	"path/filepath"

	// Sqlite3 database driver
	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
)

var (
	// Error is the error that occured during init
	// if an error occured during init
	// Otherwise, it's nil
	Error error
	// databaseInstance is the instance of the database
	databaseInstance Database
)

// tables is a list of sql queries where each query
// creates a table within the database
var tables = [...]string{
	`	CREATE TABLE IF NOT EXISTS UserStrikes (
			UserID 			INTEGER NOT NULL,
			StrikeID 		INTEGER NOT NULL,
			Reason 			TEXT	NOT NULL,
			AuthorID 		INTEGER NOT NULL,
			Timestamp 		TEXT	NOT NULL,
			EditedTimestamp TEXT	NOT NULL,

			PRIMARY KEY (UserID, StrikeID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS GuildSettings (
			GuildID 				INTEGER NOT NULL,
			BuildChannelID 			INTEGER NOT NULL,
			TicketChannelCategoryID	INTEGER NOT NULL,
			Timestamp				TEXT	NOT NULL,
			EditedTimestamp			TEXT	NOT NULL,

			PRIMARY KEY (GuildID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS Editions (
			ID				INTEGER NOT NULL,
			Name 			TEXT	NOT NULL,
			Description 	TEXT	NOT NULL,
			Timestamp 		TEXT	NOT NULL,
			EditedTimestamp TEXT	NOT NULL,

			PRIMARY KEY (ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS BuildClasses (
			ID 				INTEGER NOT NULL,
			Name 			TEXT 	NOT NULL,
			Description 	TEXT 	NOT NULL,
			EmbedColour 	TEXT 	NOT NULL,
			Timestamp 		TEXT	NOT NULL,
			EditedTimestamp TEXT	NOT NULL,

			PRIMARY KEY (ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS RecordTypes (
			ID 				INTEGER NOT NULL,
			Name 			TEXT 	NOT NULL,
			Description 	TEXT 	NOT NULL,
			Timestamp 		TEXT	NOT NULL,
			EditedTimestamp TEXT	NOT NULL,

			PRIMARY KEY (ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS GuildRecordTypeChannels (
			GuildID 		INTEGER NOT NULL,
			RecordTypeID 	INTEGER NOT NULL,
			ChannelID 		INTEGER NOT NULL,
			Timestamp 		TEXT	NOT NULL,
			EditedTimestamp TEXT	NOT NULL,

			PRIMARY KEY (GuildID, RecordTypeID),
			FOREIGN KEY (RecordTypeID) REFERENCES RecordTypes(ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS Builds (
			ID 						INTEGER NOT NULL,
			Verified 				INTEGER NOT NULL,
			VerifierID 				INTEGER NOT NULL,
			VerifiedTimestamp 		INTEGER NOT NULL,
			Reported 				INTEGER NOT NULL,
			ReporterID 				INTEGER NOT NULL,
			ReportedTimestamp 		INTEGER NOT NULL,
			UpdateRequest 			INTEGER NOT NULL,
			UpdateRequestBuildID 	INTEGER NOT NULL,
			EditionID 				INTEGER NOT NULL,
			BuildClassID 			INTEGER NOT NULL,
			Name 					TEXT 	NOT NULL,
			Description 			TEXT 	NOT NULL,
			Creators 				TEXT 	NOT NULL,
			CreationTimestamp 		INTEGER NOT NULL,
			Width 					INTEGER NOT NULL,
			Height 					INTEGER NOT NULL,
			Depth 					INTEGER NOT NULL,
			NormalCloseDuration 	INTEGER NOT NULL,
			NormalOpenDuration 		INTEGER NOT NULL,
			VisibleCloseDuration 	INTEGER NOT NULL,
			VisibleOpenDuration 	INTEGER NOT NULL,
			DelayCloseDuration 		INTEGER NOT NULL,
			DelayOpenDuration 		INTEGER NOT NULL,
			ResetCloseDuration 		INTEGER NOT NULL,
			ResetOpenDuration 		INTEGER NOT NULL,
			ExtensionDuration 		INTEGER NOT NULL,
			RetractionDuration 		INTEGER NOT NULL,
			ExtensionDelayDuration 	INTEGER NOT NULL,
			RetractionDelayDuration INTEGER NOT NULL,
			ImageURL 				TEXT 	NOT NULL,
			YoutubeURL 				TEXT 	NOT NULL,
			WorldDownloadURL 		TEXT 	NOT NULL,
			ServerIPAddress 		TEXT 	NOT NULL,
			ServerCoordinates 		TEXT 	NOT NULL,
			ServerCommand 			TEXT 	NOT NULL,
			SubmitterID 			INTEGER NOT NULL,
			Timestamp 				TEXT	NOT NULL,
			EditedTimestamp 		TEXT	NOT NULL,

			PRIMARY KEY (ID),
			FOREIGN KEY (UpdateRequestBuildID)	REFERENCES Builds(ID),
			FOREIGN KEY (EditionID) 			REFERENCES Editions(ID),
			FOREIGN KEY (BuildClassID)			REFERENCES BuildClasses(ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS Versions (
			ID 					INTEGER NOT NULL,
			EditionID 			INTEGER NOT NULL,
			MajorVersion 		INTEGER NOT NULL,
			MinorVersion 		INTEGER NOT NULL,
			Path 				INTEGER NOT NULL,
			Name 				TEXT 	NOT NULL,
			Description 		TEXT 	NOT NULL,
			VersionTimestamp 	INTEGER NOT NULL,
			Timestamp 			TEXT	NOT NULL,
			EditedTimestamp 	TEXT	NOT NULL,

			PRIMARY KEY (ID),
			FOREIGN KEY (EditionID) REFERENCES Editions(ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS Records (
			ID 						INTEGER NOT NULL,
			Verified 				INTEGER NOT NULL,
			VerifierID 				INTEGER NOT NULL,
			VerifiedTimestamp 		INTEGER NOT NULL,
			UpdateRequest 			INTEGER NOT NULL,
			UpdateRequestRecordID 	INTEGER NOT NULL,
			EditionID 				INTEGER NOT NULL,
			BuildClassID 			INTEGER NOT NULL,
			RecordTypeID 			INTEGER NOT NULL,
			Name 					TEXT 	NOT NULL,
			Description 			TEXT 	NOT NULL,
			SubmitterID 			INTEGER NOT NULL,
			Timestamp 				TEXT	NOT NULL,
			EditedTimestamp 		TEXT	NOT NULL,

			PRIMARY KEY (ID),
			FOREIGN KEY (UpdateRequestRecordID) REFERENCES Records(ID),
			FOREIGN KEY (EditionID) 			REFERENCES Editions(ID),
			FOREIGN KEY (BuildClassID) 			REFERENCES BuildClasses(ID),
			FOREIGN KEY (RecordTypeID) 			REFERENCES RecordTypes(ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS GuildBuildMessages (
			GuildID 		INTEGER NOT NULL,
			BuildID 		INTEGER NOT NULL,
			ChannelID 		INTEGER NOT NULL,
			MessageID 		INTEGER NOT NULL,
			Timestamp 		TEXT	NOT NULL,
			EditedTimestamp TEXT	NOT NULL,

			PRIMARY KEY (GuildID, BuildID),
			FOREIGN KEY (BuildID) REFERENCES Builds(ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS BuildVersions (
			BuildID 		INTEGER NOT NULL,
			VersionID 		INTEGER NOT NULL,
			StatusID 		INTEGER NOT NULL,
			Notes 			TEXT 	NOT NULL,
			Timestamp 		TEXT	NOT NULL,
			EditedTimestamp TEXT	NOT NULL,

			PRIMARY KEY (BuildID, VersionID),
			FOREIGN KEY (BuildID) 	REFERENCES Builds(ID),
			FOREIGN KEY (VersionID) REFERENCES Versions(ID),
			FOREIGN KEY (StatusID) 	REFERENCES Statuses(ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS Statuses (
			ID 				INTEGER NOT NULL,
			Name 			TEXT 	NOT NULL,
			Description 	TEXT 	NOT NULL,
			Timestamp 		TEXT	NOT NULL,
			EditedTimestamp TEXT	NOT NULL,

			PRIMARY KEY (ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS BuildRecords (
			ID 					INTEGER NOT NULL,
			BuildID 			INTEGER NOT NULL,
			RecordID 			INTEGER NOT NULL,
			Verified 			INTEGER NOT NULL,
			VerifierID 			INTEGER NOT NULL,
			VerifiedTimestamp 	INTEGER NOT NULL,
			Reported 			INTEGER NOT NULL,
			ReporterID 			INTEGER NOT NULL,
			ReportedTimestamp 	INTEGER NOT NULL,
			JointBuildRecord 	INTEGER NOT NULL,
			JointBuildRecordID 	INTEGER NOT NULL,
			SubmitterID 		INTEGER NOT NULL,
			Timestamp 			TEXT	NOT NULL,
			EditedTimestamp 	TEXT	NOT NULL,

			PRIMARY KEY (ID),
			FOREIGN KEY (BuildID) 				REFERENCES Builds(ID),
			FOREIGN KEY (RecordID) 				REFERENCES Records(ID),
			FOREIGN KEY (JointBuildRecordID) 	REFERENCES BuildRecords(ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS GuildRecordMessages (
			GuildID 		INTEGER NOT NULL,
			RecordID 		INTEGER NOT NULL,
			ChannelID 		INTEGER NOT NULL,
			MessageID 		INTEGER NOT NULL,
			Timestamp 		TEXT	NOT NULL,
			EditedTimestamp TEXT	NOT NULL,

			PRIMARY KEY (GuildID, RecordID),
			FOREIGN KEY (RecordID) REFERENCES Records(ID)
		)
	`,
	`	CREATE TABLE IF NOT EXISTS GuildTicketChannels (
			GuildID 	INTEGER NOT NULL,
			TicketID 	INTEGER NOT NULL,
			ChannelID 	INTEGER NOT NULL,
			TicketType 	INTEGER NOT NULL,
			CreatorID 	INTEGER NOT NULL,
			Timestamp 	TEXT	NOT NULL,

			PRIMARY KEY (GuildID, TicketID)
		)
	`,
}

// databaseFileExists tests to see if there is already
// a file at the database path
func databaseFileExists() (bool, error) {
	// Get absolute path of databasePath
	path, err := filepath.Abs(databasePath)
	if err != nil {
		return false, err
	}
	// Getting status of the file path
	if _, err := os.Stat(path); err == nil {
		// The file exists
		return true, nil
	} else if !os.IsNotExist(err) {
		// An error occured
		return false, err
	}
	// The file doesn't exist
	return false, nil
}

// executeQuery performs an sql query on an sql database
// and returns the result
func executeQuery(db *sql.DB, query string) (sql.Result, error) {
	// Prepare the query into a statement
	s, err := db.Prepare(query)
	if err != nil {
		return nil, errors.Wrap(err, "failed to prepare query")
	}
	// Execute the statement
	result, err := s.Exec()
	if err != nil {
		return nil, errors.Wrap(err, "failed to execute statement")
	}
	return result, nil
}

// createTables executes all of the queries in the tables array
func createTables(db *sql.DB) error {
	// For each query
	for _, table := range tables {
		// Execute the query
		_, err := executeQuery(db, table)
		if err != nil {
			return errors.Wrap(err, "failed to execute query")
		}
	}
	return nil
}

// init creates the database connection
func init() {
	// Test if the database already exists
	exists, err := databaseFileExists()
	if err != nil {
		Error = errors.Wrap(err, "failed to determine status of database file")
	}

	// Create database connection
	db, err := sql.Open("sqlite3", databasePath)
	if err != nil {
		Error = errors.Wrap(err, "failed to open database connection")
		return
	}

	// If database was just created, create tables
	if !exists {
		err := createTables(db)
		if err != nil {
			Error = errors.Wrap(err, "failed to create table")
			return
		}
	}

	// Create the database instance
	databaseInstance = Database{db: db}
}
