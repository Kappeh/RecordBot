package database

import (
	"strconv"
	"time"

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
	// Convert userID to int
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return UserStrikeCount{}, errors.Wrap(err, "failed to convert user id to interger")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT COUNT(1)
		FROM UserStrikes
		WHERE UserID = ?
	`, userIDint)
	if err != nil {
		return UserStrikeCount{}, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// The query should always return a value
	// Therefore the if block shouldn't be ran
	if !rows.Next() {
		return UserStrikeCount{}, errors.New("query didn't return a value")
	}
	// Extract data
	var count int
	if err = rows.Scan(&count); err != nil {
		return UserStrikeCount{}, errors.Wrap(err, "failed to extract data")
	}
	return UserStrikeCount{
		UserID: userID,
		Count:  count,
	}, nil
}

// UserStrikeCounts gets the number of strikes given to each user
// that has been given at least one strike
func (d *Database) UserStrikeCounts() ([]UserStrikeCount, error) {
	// Query the database
	rows, err := d.db.Query(`
		SELECT UserID, COUNT(1)
		FROM UserStrikes
		GROUP BY UserID
	`)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Space to store results
	results := []UserStrikeCount{}
	var (
		userID int
		count  int
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(&userID, &count); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Add to results
		results = append(results, UserStrikeCount{
			UserID: strconv.Itoa(userID),
			Count:  count,
		})
	}
	return results, nil
}

// UserStrike gets the information of a strike given to a user
func (d *Database) UserStrike(userID, strikeID string) (UserStrike, bool, error) {
	// Convert userID and strikeID to ints
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to convert user id to interger")
	}
	strikeIDint, err := strconv.Atoi(strikeID)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to convert strike id to interger")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT Reason, AuthorID, Timestamp, EditedTimestamp
		FROM UserStrikes
		WHERE UserID = ? AND StrikeID = ?
	`, userIDint, strikeIDint)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if the user strike exists
	if !rows.Next() {
		return UserStrike{}, false, nil
	}
	// Extract data
	var (
		reason                string
		authorID              int
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	if err = rows.Scan(&reason, &authorID, &timestampString, &editedTimestampString); err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return UserStrike{
		UserID:          userID,
		StrikeID:        strikeID,
		Reason:          reason,
		AuthorID:        strconv.Itoa(authorID),
		Timestamp:       Timestamp(timestamp),
		EditedTimestamp: Timestamp(editedTimestamp),
	}, true, nil
}

// UserStrikes gets the information of all strikes given to a user
func (d *Database) UserStrikes(userID string) ([]UserStrike, error) {
	// Convert userID to int
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert user id to interger")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT StrikeID, Reason, AuthorID, Timestamp, EditedTimestamp
		FROM UserStrikes
		WHERE UserID = ?
	`, userIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []UserStrike{}
	var (
		strikeID              int
		reason                string
		authorID              int
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&strikeID, &reason, &authorID,
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
		results = append(results, UserStrike{
			UserID:          userID,
			StrikeID:        strconv.Itoa(strikeID),
			Reason:          reason,
			AuthorID:        strconv.Itoa(authorID),
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// UserStrikeCreate creates a strike
func (d *Database) UserStrikeCreate(userID, reason, authorID string) (UserStrike, error) {
	// Convert userID and authorID to ints
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return UserStrike{}, errors.Wrap(err, "failed to convert user id to interger")
	}
	authorIDint, err := strconv.Atoi(authorID)
	if err != nil {
		return UserStrike{}, errors.Wrap(err, "failed to convert author id to interger")
	}
	// Get the next strike id for the user
	strikeIDint, err := d.nextStrikeID(userID)
	if err != nil {
		return UserStrike{}, errors.Wrap(err, "failed to get next strike id")
	}
	// Create the user strike
	us := UserStrike{
		UserID:          userID,
		StrikeID:        strconv.Itoa(strikeIDint),
		Reason:          reason,
		AuthorID:        authorID,
		Timestamp:       Timestamp(time.Now()),
		EditedTimestamp: Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO UserStrikes 
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return UserStrike{}, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		userIDint, strikeIDint, us.Reason, authorIDint,
		time.Time(us.Timestamp).Format(timeLayout),
		time.Time(us.EditedTimestamp).Format(timeLayout),
	); err != nil {
		return UserStrike{}, errors.Wrap(err, "database query failed")
	}
	return us, nil
}

// UserStrikeDelete a strike given to a user
func (d *Database) UserStrikeDelete(userID, strikeID string) (UserStrike, bool, error) {
	// Convert userID and strikeID to ints
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to convert user id to interger")
	}
	strikeIDint, err := strconv.Atoi(strikeID)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to convert strike id to interger")
	}
	// Get the user strike to return after deletion and
	// to check if it exists
	us, ok, err := d.UserStrike(userID, strikeID)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to get row from database")
	} else if !ok {
		// Record doesn't exist
		return UserStrike{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM UserStrikes
		WHERE UserID = ? AND StrikeID = ?
	`)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(userIDint, strikeIDint); err != nil {
		return UserStrike{}, false, errors.Wrap(err, "database query failed")
	}
	return us, true, nil
}

// UserStrikeEdit edits a strike given to a user
func (d *Database) UserStrikeEdit(userID, strikeID, reason string) (UserStrike, bool, error) {
	// Convert userID and strikeID into ints
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to convert user id to integer")
	}
	strikeIDint, err := strconv.Atoi(strikeID)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to convert strike id to integer")
	}
	// Get the user strike that's to be updated
	us, ok, err := d.UserStrike(userID, strikeID)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to get row from database")
	} else if !ok {
		// Record doesn't exist
		return UserStrike{}, false, nil
	}
	// Update information
	us.Reason = reason
	us.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE UserStrikes
		SET Reason = ?, EditedTimestamp = ?
		WHERE UserID = ? AND StrikeID = ?
	`)
	if err != nil {
		return UserStrike{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	if _, err = s.Exec(
		reason, time.Time(us.EditedTimestamp).Format(timeLayout),
		userIDint, strikeIDint,
	); err != nil {
		return UserStrike{}, false, errors.Wrap(err, "database query failed")
	}
	return us, true, nil
}

// GuildSetting gets the setting information for a guild
func (d *Database) GuildSetting(guildID string) (GuildSetting, bool, error) {
	// Convert guildID to int
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT BuildChannelID, TicketChannelCategoryID, Timestamp, EditedTimestamp
		FROM GuildSettings
		WHERE GuildID = ?
	`, guildIDint)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if the guild setting exists
	if !rows.Next() {
		return GuildSetting{}, false, nil
	}
	// Extract data
	var (
		buildChannelIDint          int
		ticketChannelCategoryIDint int
		timestampString            string
		editedTimestampString      string
		timestamp                  time.Time
		editedTimestamp            time.Time
	)
	if err = rows.Scan(
		&buildChannelIDint, &ticketChannelCategoryIDint,
		&timestampString, &editedTimestampString,
	); err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "database query failed")
	}
	// Parse timestamps
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return GuildSetting{
		GuildID:                 guildID,
		BuildChannelID:          strconv.Itoa(buildChannelIDint),
		TicketChannelCategoryID: strconv.Itoa(ticketChannelCategoryIDint),
		Timestamp:               Timestamp(timestamp),
		EditedTimestamp:         Timestamp(editedTimestamp),
	}, true, nil
}

// GuildSettings gets the setting information for all guilds
func (d *Database) GuildSettings() ([]GuildSetting, error) {
	// Query the database
	rows, err := d.db.Query(`
		SELECT *
		FROM GuildSettings
	`)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []GuildSetting{}
	var (
		guildIDint                 int
		buildChannelIDint          int
		ticketChannelCategoryIDint int
		timestampString            string
		editedTimestampString      string
		timestamp                  time.Time
		editedTimestamp            time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&guildIDint, &buildChannelIDint, &ticketChannelCategoryIDint,
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
		results = append(results, GuildSetting{
			GuildID:                 strconv.Itoa(guildIDint),
			BuildChannelID:          strconv.Itoa(buildChannelIDint),
			TicketChannelCategoryID: strconv.Itoa(ticketChannelCategoryIDint),
			Timestamp:               Timestamp(timestamp),
			EditedTimestamp:         Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// GuildSettingCreate creates setting information for a guild
func (d *Database) GuildSettingCreate(guildID, buildChannelID, ticketCategoryID string) (GuildSetting, bool, error) {
	// Convert guildID, buildChannelID and ticketCategoryID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	buildChannelIDint, err := strconv.Atoi(buildChannelID)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to convert build channel id to integer")
	}
	ticketCategoryIDint, err := strconv.Atoi(ticketCategoryID)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to convert ticket category id to integer")
	}
	// Check if guild setting already exists
	if _, ok, err := d.GuildSetting(guildID); err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to determine if guild setting exists")
	} else if ok {
		// Row already exist
		return GuildSetting{}, false, nil
	}
	// Create guild setting
	gs := GuildSetting{
		GuildID:                 guildID,
		BuildChannelID:          buildChannelID,
		TicketChannelCategoryID: ticketCategoryID,
		Timestamp:               Timestamp(time.Now()),
		EditedTimestamp:         Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO GuildSettings 
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		guildIDint, buildChannelIDint, ticketCategoryIDint,
		time.Time(gs.Timestamp).Format(timeLayout),
		time.Time(gs.EditedTimestamp).Format(timeLayout),
	); err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "database query failed")
	}
	return gs, false, nil
}

// GuildSettingDelete deletes the setting information for a guild
func (d *Database) GuildSettingDelete(guildID string) (GuildSetting, bool, error) {
	// Convert guildID to int
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	// Get the guild setting to return after deletion and
	// to check if it exists
	gs, ok, err := d.GuildSetting(guildID)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to determine if guild setting exists")
	} else if !ok {
		// Row doesn't exist
		return GuildSetting{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM GuildSettings
		WHERE GuildID = ?
	`)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(guildIDint); err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "database query failed")
	}
	return gs, true, nil
}

// GuildSettingEdit edits the setting information for a guild
func (d *Database) GuildSettingEdit(guildID, buildChannelID, ticketChannelCategoryID string) (GuildSetting, bool, error) {
	// Convert guildID, buildChannelID and ticketCategory to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	buildChannelIDint, err := strconv.Atoi(buildChannelID)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to convert build channel id to integer")
	}
	ticketChannelCategoryIDint, err := strconv.Atoi(ticketChannelCategoryID)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to convert ticket category id to integer")
	}
	// Get the guild setting to return after deletion and
	// to check if it exists
	gs, ok, err := d.GuildSetting(guildID)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to determine if guild setting exists")
	} else if !ok {
		// Row doesn't exist
		return GuildSetting{}, false, nil
	}
	// Update information
	gs.BuildChannelID = buildChannelID
	gs.TicketChannelCategoryID = ticketChannelCategoryID
	gs.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE GuildSettings
		SET BuildChannelID = ?, TicketChannelCategoryID = ?, EditedTimestamp = ?
		WHERE GuildID = ?
	`)
	if err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		buildChannelIDint, ticketChannelCategoryIDint,
		time.Time(gs.EditedTimestamp).Format(timeLayout),
		guildIDint,
	); err != nil {
		return GuildSetting{}, false, errors.Wrap(err, "database query failed")
	}
	return gs, true, nil
}

// Edition gets the edition information with the specified id
func (d *Database) Edition(editionID string) (Edition, bool, error) {
	// Convert editionID to int
	editionIDint, err := strconv.Atoi(editionID)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to convert edition if to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT Name, Description, Timestamp, EditedTimestamp
		FROM Editions
		WHERE ID = ?
	`, editionIDint)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if edition exists in database
	if !rows.Next() {
		return Edition{}, false, nil
	}
	// Extract data
	var (
		name                  string
		description           string
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	if err = rows.Scan(
		&name, &description,
		&timestampString,
		&editedTimestampString,
	); err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return Edition{
		ID:              editionID,
		Name:            name,
		Description:     description,
		Timestamp:       Timestamp(timestamp),
		EditedTimestamp: Timestamp(editedTimestamp),
	}, true, nil
}

// Editions gets the edition information for all editions in the database
func (d *Database) Editions() ([]Edition, error) {
	// Query the database
	rows, err := d.db.Query(`
		SELECT *
		FROM Editions
	`)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []Edition{}
	var (
		idInt                 int
		name                  string
		description           string
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&idInt, &name, &description,
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
		results = append(results, Edition{
			ID:              strconv.Itoa(idInt),
			Name:            name,
			Description:     description,
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// EditionCreate creates an edition in the database
func (d *Database) EditionCreate(name, description string) (Edition, error) {
	// Create edition
	e := Edition{
		Name:            name,
		Description:     description,
		Timestamp:       Timestamp(time.Now()),
		EditedTimestamp: Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO Editions (Name, Description, Timestamp, EditedTimestamp)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return Edition{}, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	res, err := s.Exec(
		name, description,
		time.Time(e.Timestamp).Format(timeLayout),
		time.Time(e.EditedTimestamp).Format(timeLayout),
	)
	if err != nil {
		return Edition{}, errors.Wrap(err, "database query failed")
	}
	// Update edition id
	idInt, err := res.LastInsertId()
	if err != nil {
		return Edition{}, errors.Wrap(err, "couldn't update edition id")
	}
	// strconv.Itoa(int(idInt)) would probably suffice
	// but doesn't give an error in case precision is lost
	// from int(idInt) int64 -> int
	e.ID = strconv.FormatInt(idInt, 10)
	return e, nil
}

// EditionDelete removes an edition from the database
func (d *Database) EditionDelete(editionID string) (Edition, bool, error) {
	// Convert editionID to int
	editionIDint, err := strconv.Atoi(editionID)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to convert edition id to integer")
	}
	// Get the edition to return after deletion and
	// to check if it exists
	e, ok, err := d.Edition(editionID)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to determine if edition exists")
	} else if !ok {
		// Row doesn't exist
		return Edition{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM Editions
		WHERE ID = ?
	`)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(editionIDint); err != nil {
		return Edition{}, false, errors.Wrap(err, "database query failed")
	}
	return e, true, nil
}

// EditionEdit edits the edition information for a specified edition
func (d *Database) EditionEdit(editionID, name, description string) (Edition, bool, error) {
	// Convert editionID to int
	editionIDint, err := strconv.Atoi(editionID)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to convert edition id to integer")
	}
	// Get the edition that is to be updated
	e, ok, err := d.Edition(editionID)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to determine if edition exists")
	} else if !ok {
		// Row doesn't exist
		return Edition{}, false, nil
	}
	// Update information
	e.Name = name
	e.Description = description
	e.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE Editions
		SET Name = ?, Description = ?, EditedTimestamp = ?
		WHERE ID = ?
	`)
	if err != nil {
		return Edition{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		name, description,
		time.Time(e.EditedTimestamp).Format(timeLayout),
		editionIDint,
	); err != nil {
		return Edition{}, false, errors.Wrap(err, "database query failed")
	}
	return e, false, nil
}

// BuildClass gets the information for a build class in the database
func (d *Database) BuildClass(buildClassID string) (BuildClass, bool, error) {
	// Convert buildClassID to int
	buildClassIDint, err := strconv.Atoi(buildClassID)
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to convert build class id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT Name, Description, EmbedColour, Timestamp, EditedTimestamp
		FROM BuildClasses
		WHERE ID = ?
	`, buildClassIDint)
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if the build class exists
	if !rows.Next() {
		return BuildClass{}, false, nil
	}
	// Create space to store result
	var (
		name                  string
		description           string
		embedColour           string
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// Extract data
	if err = rows.Scan(
		&name, &description, &embedColour,
		&timestampString, &editedTimestampString,
	); err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to parse edition timestamp")
	}
	return BuildClass{
		ID:              buildClassID,
		Name:            name,
		Description:     description,
		EmbedColour:     embedColour,
		Timestamp:       Timestamp(timestamp),
		EditedTimestamp: Timestamp(editedTimestamp),
	}, true, nil
}

// BuildClasses gets the information for all build classes in the database
func (d *Database) BuildClasses() ([]BuildClass, error) {
	// Query the database
	rows, err := d.db.Query(`
		SELECT *
		FROM BuildClasses
	`)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []BuildClass{}
	var (
		idInt                 int
		name                  string
		description           string
		embedColour           string
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&idInt, &name, &description, &embedColour,
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
		results = append(results, BuildClass{
			ID:              strconv.Itoa(idInt),
			Name:            name,
			Description:     description,
			EmbedColour:     embedColour,
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// BuildClassCreate creates a new build class
func (d *Database) BuildClassCreate(name, description, embedColour string) (BuildClass, error) {
	// Create build class
	bc := BuildClass{
		Name:            name,
		Description:     description,
		EmbedColour:     embedColour,
		Timestamp:       Timestamp(time.Now()),
		EditedTimestamp: Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO BuildClasses (Name, Description, EmbedColour, Timestamp, EditedTimestamp)
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return BuildClass{}, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	res, err := s.Exec(
		name, description, embedColour,
		time.Time(bc.Timestamp).Format(timeLayout),
		time.Time(bc.EditedTimestamp).Format(timeLayout),
	)
	if err != nil {
		return BuildClass{}, errors.Wrap(err, "database query failed")
	}
	// Update build class id
	idInt, err := res.LastInsertId()
	if err != nil {
		return BuildClass{}, errors.Wrap(err, "couldn't update edition id")
	}
	// strconv.Itoa(int(idInt)) would probably suffice
	// but doesn't give an error in case precision is lost
	// from int(idInt) int64 -> int
	bc.ID = strconv.FormatInt(idInt, 10)
	return bc, nil
}

// BuildClassDelete removes an existing build class
func (d *Database) BuildClassDelete(buildClassID string) (BuildClass, bool, error) {
	// Convert buildClassID to int
	buildClassIDint, err := strconv.Atoi(buildClassID)
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to convert build class id to int")
	}
	// Get the build class to return after deletion and
	// to check if it exists
	bc, ok, err := d.BuildClass(buildClassID)
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to determine if build class exists")
	} else if !ok {
		// Row doesn't exist
		return BuildClass{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM BuildClasses
		WHERE ID = ?
	`)
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(buildClassIDint); err != nil {
		return BuildClass{}, false, errors.Wrap(err, "database query failed")
	}
	return bc, true, nil
}

// BuildClassEdit edits an existing build class
func (d *Database) BuildClassEdit(buildClassID, name, description, embedColour string) (BuildClass, bool, error) {
	// Convert build class id to int
	buildClassIDint, err := strconv.Atoi(buildClassID)
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to convert build class id to integer")
	}
	// Get the build class that is to be updated
	bc, ok, err := d.BuildClass(buildClassID)
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to determine if build class exists")
	} else if !ok {
		// Record doesn't exists
		return BuildClass{}, false, nil
	}
	// Update information
	bc.Name = name
	bc.Description = description
	bc.EmbedColour = embedColour
	bc.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE BuildClasses
		SET Name = ?, Description = ?, EmbedColour = ?, EditedTimestamp = ?
		WHERE ID = ?
	`)
	if err != nil {
		return BuildClass{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		name, description, embedColour,
		time.Time(bc.EditedTimestamp).Format(timeLayout),
		buildClassIDint,
	); err != nil {
		return BuildClass{}, false, errors.Wrap(err, "database query failed")
	}
	return bc, true, nil
}

// RecordType gets the information for the specified record type
func (d *Database) RecordType(recordTypeID string) (RecordType, bool, error) {
	// Convert recordTypeID to int
	recordTypeIDint, err := strconv.Atoi(recordTypeID)
	if err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to convert record type id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT Name, Description, Timestamp, EditedTimestamp
		FROM RecordTypes
		WHERE ID = ?
	`, recordTypeIDint)
	if err != nil {
		return RecordType{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if the record type exists
	if !rows.Next() {
		return RecordType{}, false, nil
	}
	// Create space to store result
	var (
		name                  string
		description           string
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// Extract data
	if err = rows.Scan(
		&name, &description,
		&timestampString,
		&editedTimestampString,
	); err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return RecordType{
		ID:              recordTypeID,
		Name:            name,
		Description:     description,
		Timestamp:       Timestamp(timestamp),
		EditedTimestamp: Timestamp(editedTimestamp),
	}, true, nil
}

// RecordTypes get the information for all record types
func (d *Database) RecordTypes() ([]RecordType, error) {
	// Query the database
	rows, err := d.db.Query(`
		SELECT *
		FROM RecordTypes
	`)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []RecordType{}
	var (
		idInt                 int
		name                  string
		description           string
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&idInt, &name, &description,
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
		results = append(results, RecordType{
			ID:              strconv.Itoa(idInt),
			Name:            name,
			Description:     description,
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// RecordTypeCreate creates a new record type
func (d *Database) RecordTypeCreate(name, description string) (RecordType, error) {
	// Create record type
	rt := RecordType{
		Name:            name,
		Description:     description,
		Timestamp:       Timestamp(time.Now()),
		EditedTimestamp: Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO RecordTypes (Name, Description, Timestamp, EditedTimestamp)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return RecordType{}, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	res, err := s.Exec(
		name, description,
		time.Time(rt.Timestamp).Format(timeLayout),
		time.Time(rt.EditedTimestamp).Format(timeLayout),
	)
	if err != nil {
		return RecordType{}, errors.Wrap(err, "database query failed")
	}
	// Update record type id
	idInt, err := res.LastInsertId()
	if err != nil {
		return RecordType{}, errors.Wrap(err, "couldn't update record type id")
	}
	// strconv.Itoa(int(idInt)) would probably suffice
	// but doesn't give an error in case precision is lost
	// from int(idInt) int64 -> int
	rt.ID = strconv.FormatInt(idInt, 10)
	return rt, nil
}

// RecordTypeDelete removes an existing record type
func (d *Database) RecordTypeDelete(recordTypeID string) (RecordType, bool, error) {
	// Convert recordTypeID to int
	recordTypeIDint, err := strconv.Atoi(recordTypeID)
	if err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to convert record type id to integer")
	}
	// Get the record type to return after deletion and
	// to check if it exists
	rt, ok, err := d.RecordType(recordTypeID)
	if err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to determine if record type exists")
	} else if !ok {
		return RecordType{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM RecordTypes
		WHERE ID = ?
	`)
	if err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(recordTypeIDint); err != nil {
		return RecordType{}, false, errors.Wrap(err, "database query failed")
	}
	return rt, true, nil
}

// RecordTypeEdit edits an existing record type
func (d *Database) RecordTypeEdit(recordTypeID, name, description string) (RecordType, bool, error) {
	// Convert recordTypeID to int
	recordTypeIDint, err := strconv.Atoi(recordTypeID)
	if err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to convert record type id to integer")
	}
	// Get the record type that is to be updated
	rt, ok, err := d.RecordType(recordTypeID)
	if err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to determine if record type exists")
	} else if !ok {
		// Record doesn't exist
		return RecordType{}, false, nil
	}
	// Update information
	rt.Name = name
	rt.Description = description
	rt.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE RecordTypes
		SET Name = ?, Description = ?, EditedTimestamp = ?
		WHERE ID = ?
	`)
	if err != nil {
		return RecordType{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		name, description, time.Time(rt.EditedTimestamp).Format(timeLayout), recordTypeIDint,
	); err != nil {
		return RecordType{}, false, errors.Wrap(err, "database query failed")
	}
	return rt, false, nil
}

// GuildRecordTypeChannel gets information for a specified guild and record type
func (d *Database) GuildRecordTypeChannel(guildID, recordTypeID string) (GuildRecordTypeChannel, bool, error) {
	// Convert guildID and recordTypeID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	recordTypeIDint, err := strconv.Atoi(recordTypeID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to convert record type id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT ChannelID, Timestamp, EditedTimestamp
		FROM GuildRecordTypeChannels
		WHERE GuildID = ? AND RecordTypeID = ?
	`, guildIDint, recordTypeIDint)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if the guild record type channel exists
	if !rows.Next() {
		return GuildRecordTypeChannel{}, false, nil
	}
	// Create space to store result
	var (
		channelIDint          int
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// Extract data
	if err = rows.Scan(&channelIDint, &timestampString, &editedTimestampString); err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return GuildRecordTypeChannel{
		GuildID:         guildID,
		RecordTypeID:    recordTypeID,
		ChannelID:       strconv.Itoa(channelIDint),
		Timestamp:       Timestamp(timestamp),
		EditedTimestamp: Timestamp(editedTimestamp),
	}, true, nil
}

// GuildRecordTypeChannels gets information for all guilds and record types
func (d *Database) GuildRecordTypeChannels(guildID string) ([]GuildRecordTypeChannel, error) {
	// Convert guildID to int
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert guild id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT RecordTypeID, ChannelID, Timestamp, EditedTimestamp
		FROM GuildRecordTypeChannels
		WHERE GuildID = ?
	`, guildIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []GuildRecordTypeChannel{}
	var (
		recordTypeIDint       int
		channelIDint          int
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&recordTypeIDint, &channelIDint,
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
		results = append(results, GuildRecordTypeChannel{
			GuildID:         guildID,
			RecordTypeID:    strconv.Itoa(recordTypeIDint),
			ChannelID:       strconv.Itoa(channelIDint),
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// GuildRecordTypeChannelCreate creates guild record type channel information for
// a specified guild and record type
func (d *Database) GuildRecordTypeChannelCreate(guildID, recordTypeID, channelID string) (GuildRecordTypeChannel, bool, error) {
	// Convert guildID, recordTypeID and channelID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	recordTypeIDint, err := strconv.Atoi(recordTypeID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to convert record type id to integer")
	}
	channelIDint, err := strconv.Atoi(channelID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to convert channel id to integer")
	}
	// Check if guild record type channel already exists
	if _, ok, err := d.GuildRecordTypeChannel(guildID, recordTypeID); err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to determine if guild record type channel exists")
	} else if ok {
		// Row already exists
		return GuildRecordTypeChannel{}, false, nil
	}
	// Create guild record type channel
	grtc := GuildRecordTypeChannel{
		GuildID:         guildID,
		RecordTypeID:    recordTypeID,
		ChannelID:       channelID,
		Timestamp:       Timestamp(time.Now()),
		EditedTimestamp: Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO GuildRecordTypeChannels
		VALUES (?, ?, ?, ?, ?)
	`)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		guildIDint, recordTypeIDint, channelIDint,
		time.Time(grtc.Timestamp).Format(timeLayout),
		time.Time(grtc.EditedTimestamp).Format(timeLayout),
	); err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "database query failed")
	}
	return grtc, true, nil
}

// GuildRecordTypeChannelDelete removes guild record type channel information for
// a specified guild and record type
func (d *Database) GuildRecordTypeChannelDelete(guildID, recordTypeID string) (GuildRecordTypeChannel, bool, error) {
	// Convert guildID and recordTypeID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	recordTypeIDint, err := strconv.Atoi(recordTypeID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to convert record type id to integer")
	}
	// Get guild record type channel to return after
	// deletion and to check if it exists
	grtc, ok, err := d.GuildRecordTypeChannel(guildID, recordTypeID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to determine if guild record type channel exists")
	} else if !ok {
		// Row doesn't exist
		return GuildRecordTypeChannel{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM GuildRecordTypeChannels
		WHERE GuildID = ? AND RecordTypeID = ?
	`)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(guildIDint, recordTypeIDint); err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "database query failed")
	}
	return grtc, true, nil
}

// GuildRecordTypeChannelEdit edits guild record type channel information for
// a specified guild and record type
func (d *Database) GuildRecordTypeChannelEdit(guildID, recordTypeID, channelID string) (GuildRecordTypeChannel, bool, error) {
	// Convert guildID, recordTypeID and channelID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	recordTypeIDint, err := strconv.Atoi(recordTypeID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to convert record type id to integer")
	}
	channelIDint, err := strconv.Atoi(channelID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to convert channel id to integer")
	}
	// Get the guild record type channel that's to be updated
	grtc, ok, err := d.GuildRecordTypeChannel(guildID, recordTypeID)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to determine if guild record type channel exists")
	} else if !ok {
		// Record doesn't exist
		return GuildRecordTypeChannel{}, false, nil
	}
	// Update information
	grtc.ChannelID = channelID
	grtc.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE GuildRecordTypeChannels
		SET ChannelID = ?, EditedTimestamp = ?
		WHERE GuildID = ? AND RecordTypeID = ?
	`)
	if err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		channelIDint, time.Time(grtc.EditedTimestamp).Format(timeLayout),
		guildIDint, recordTypeIDint,
	); err != nil {
		return GuildRecordTypeChannel{}, false, errors.Wrap(err, "database query failed")
	}
	return grtc, false, nil
}

// Build gets the information for a specified build
func (d *Database) Build(buildID string) (Build, bool, error) {
	// Convert buildID to int
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT Verified, VerifierID, VerifiedTimestamp, Reported, ReporterID, 
			ReportedTimestamp, UpdateRequest, UpdateRequestBuildID, EditionID, 
			BuildClassID, Name, Description, Creators, CreationTimestamp, Width, 
			Height, Depth, NormalCloseDuration, NormalOpenDuration, VisibleCloseDuration, 
			VisibleOpenDuration, DelayCloseDuration, DelayOpenDuration, ResetCloseDuration, 
			ResetOpenDuration, ExtensionDuration, RetractionDuration, ExtensionDelayDuration, 
			RetractionDelayDuration, ImageURL, YoutubeURL, WorldDownloadURL, ServerIPAddress, 
			ServerCoordinates, ServerCommand, SubmitterID, Timestamp, EditedTimestamp
		FROM Builds
		WHERE ID = ?
	`, buildIDint)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Chec if build exists in database
	if !rows.Next() {
		return Build{}, false, nil
	}
	// Create space to store result
	var (
		verifiedInt             int
		verifierIDint           int
		verifiedTimestampString string
		reportedInt             int
		reporterIDint           int
		reportedTimestampString string
		updateRequestInt        int
		updateRequestBuildIDint int
		editionIDint            int
		buildClassIDint         int
		name                    string
		description             string
		creators                string
		creationTimestampString string
		width                   int
		height                  int
		depth                   int
		normalCloseDuration     int
		normalOpenDuration      int
		visibleCloseDuration    int
		visibleOpenDuration     int
		delayCloseDuration      int
		delayOpenDuration       int
		resetCloseDuration      int
		resetOpenDuration       int
		extensionDuration       int
		retractionDuration      int
		extensionDelayDuration  int
		retractionDelayDuration int
		imageURL                string
		youtubeURL              string
		worldDownloadURL        string
		serverIPAddress         string
		serverCoordinates       string
		serverCommand           string
		submitterIDint          int
		timestampString         string
		editedTimestampString   string
		verifiedTimestamp       time.Time
		reportedTimestamp       time.Time
		creationTimestamp       time.Time
		timestamp               time.Time
		editedTimestamp         time.Time
	)
	// Extract data
	if err = rows.Scan(
		&verifiedInt, &verifierIDint, &verifiedTimestampString, &reportedInt, &reporterIDint,
		&reportedTimestampString, &updateRequestInt, &updateRequestBuildIDint, &editionIDint,
		&buildClassIDint, &name, &description, &creators, &creationTimestampString, &width,
		&height, &depth, &normalCloseDuration, &normalOpenDuration, &visibleCloseDuration,
		&visibleOpenDuration, &delayCloseDuration, &delayOpenDuration, &resetCloseDuration,
		&resetOpenDuration, &extensionDuration, &retractionDuration, &extensionDelayDuration,
		&retractionDelayDuration, &imageURL, &youtubeURL, &worldDownloadURL, &serverIPAddress,
		&serverCoordinates, &serverCommand, &submitterIDint, &timestampString, &editedTimestampString,
	); err != nil {
		return Build{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if verifiedTimestamp, err = time.Parse(timeLayout, verifiedTimestampString); err != nil {
		return Build{}, false, errors.Wrap(err, "failed to parse verified timestamp")
	}
	if reportedTimestamp, err = time.Parse(timeLayout, reportedTimestampString); err != nil {
		return Build{}, false, errors.Wrap(err, "failed to parse reported timestamp")
	}
	if creationTimestamp, err = time.Parse(timeLayout, creationTimestampString); err != nil {
		return Build{}, false, errors.Wrap(err, "failed to parse creation timestamp")
	}
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return Build{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return Build{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	// Convert to build struct
	return Build{
		ID:                      buildID,
		Verified:                verifiedInt != 0,
		VerifierID:              strconv.Itoa(verifierIDint),
		VerifiedTimestamp:       Timestamp(verifiedTimestamp),
		Reported:                reportedInt != 0,
		ReporterID:              strconv.Itoa(reporterIDint),
		ReportedTimestamp:       Timestamp(reportedTimestamp),
		UpdateRequest:           updateRequestInt != 0,
		UpdateRequestBuildID:    strconv.Itoa(updateRequestBuildIDint),
		EditionID:               strconv.Itoa(editionIDint),
		BuildClassID:            strconv.Itoa(buildClassIDint),
		Name:                    name,
		Description:             description,
		Creators:                creators,
		CreationTimestamp:       Timestamp(creationTimestamp),
		Width:                   width,
		Height:                  height,
		Depth:                   depth,
		NormalCloseDuration:     normalCloseDuration,
		NormalOpenDuration:      normalOpenDuration,
		VisibleCloseDuration:    visibleCloseDuration,
		VisibleOpenDuration:     visibleOpenDuration,
		DelayCloseDuration:      delayCloseDuration,
		DelayOpenDuration:       delayOpenDuration,
		ResetCloseDuration:      resetCloseDuration,
		ResetOpenDuration:       resetOpenDuration,
		ExtensionDuration:       extensionDuration,
		RetractionDuration:      retractionDuration,
		ExtensionDelayDuration:  extensionDelayDuration,
		RetractionDelayDuration: retractionDelayDuration,
		ImageURL:                imageURL,
		YoutubeURL:              youtubeURL,
		WorldDownloadURL:        worldDownloadURL,
		ServerIPAddress:         serverIPAddress,
		ServerCoordinates:       serverCoordinates,
		ServerCommand:           serverCommand,
		SubmitterID:             strconv.Itoa(submitterIDint),
		Timestamp:               Timestamp(timestamp),
		EditedTimestamp:         Timestamp(editedTimestamp),
	}, true, nil
}

// Builds gets the information for all builds in the database
func (d *Database) Builds() ([]Build, error) {
	// Query the database
	rows, err := d.db.Query(`
		SELECT *
		FROM Builds
	`)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []Build{}
	var (
		idInt                   int
		verifiedInt             int
		verifierIDint           int
		verifiedTimestampString string
		reportedInt             int
		reporterIDint           int
		reportedTimestampString string
		updateRequestInt        int
		updateRequestBuildIDint int
		editionIDint            int
		buildClassIDint         int
		name                    string
		description             string
		creators                string
		creationTimestampString string
		width                   int
		height                  int
		depth                   int
		normalCloseDuration     int
		normalOpenDuration      int
		visibleCloseDuration    int
		visibleOpenDuration     int
		delayCloseDuration      int
		delayOpenDuration       int
		resetCloseDuration      int
		resetOpenDuration       int
		extensionDuration       int
		retractionDuration      int
		extensionDelayDuration  int
		retractionDelayDuration int
		imageURL                string
		youtubeURL              string
		worldDownloadURL        string
		serverIPAddress         string
		serverCoordinates       string
		serverCommand           string
		submitterIDint          int
		timestampString         string
		editedTimestampString   string
		verifiedTimestamp       time.Time
		reportedTimestamp       time.Time
		creationTimestamp       time.Time
		timestamp               time.Time
		editedTimestamp         time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&idInt, &verifiedInt, &verifierIDint, &verifiedTimestampString, &reportedInt, &reporterIDint,
			&reportedTimestampString, &updateRequestInt, &updateRequestBuildIDint, &editionIDint,
			&buildClassIDint, &name, &description, &creators, &creationTimestampString, &width,
			&height, &depth, &normalCloseDuration, &normalOpenDuration, &visibleCloseDuration,
			&visibleOpenDuration, &delayCloseDuration, &delayOpenDuration, &resetCloseDuration,
			&resetOpenDuration, &extensionDuration, &retractionDuration, &extensionDelayDuration,
			&retractionDelayDuration, &imageURL, &youtubeURL, &worldDownloadURL, &serverIPAddress,
			&serverCoordinates, &serverCommand, &submitterIDint, &timestampString, &editedTimestampString,
		); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Parse timestamps
		if verifiedTimestamp, err = time.Parse(timeLayout, verifiedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse verified timestamp")
		}
		if reportedTimestamp, err = time.Parse(timeLayout, reportedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse reported timestamp")
		}
		if creationTimestamp, err = time.Parse(timeLayout, creationTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse creation timestamp")
		}
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}
		if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse edited timestamp")
		}
		// Add to results
		results = append(results, Build{
			ID:                      strconv.Itoa(idInt),
			Verified:                verifiedInt != 0,
			VerifierID:              strconv.Itoa(verifierIDint),
			VerifiedTimestamp:       Timestamp(verifiedTimestamp),
			Reported:                reportedInt != 0,
			ReporterID:              strconv.Itoa(reporterIDint),
			ReportedTimestamp:       Timestamp(reportedTimestamp),
			UpdateRequest:           updateRequestInt != 0,
			UpdateRequestBuildID:    strconv.Itoa(updateRequestBuildIDint),
			EditionID:               strconv.Itoa(editionIDint),
			BuildClassID:            strconv.Itoa(buildClassIDint),
			Name:                    name,
			Description:             description,
			Creators:                creators,
			CreationTimestamp:       Timestamp(creationTimestamp),
			Width:                   width,
			Height:                  height,
			Depth:                   depth,
			NormalCloseDuration:     normalCloseDuration,
			NormalOpenDuration:      normalOpenDuration,
			VisibleCloseDuration:    visibleCloseDuration,
			VisibleOpenDuration:     visibleOpenDuration,
			DelayCloseDuration:      delayCloseDuration,
			DelayOpenDuration:       delayOpenDuration,
			ResetCloseDuration:      resetCloseDuration,
			ResetOpenDuration:       resetOpenDuration,
			ExtensionDuration:       extensionDuration,
			RetractionDuration:      retractionDuration,
			ExtensionDelayDuration:  extensionDelayDuration,
			RetractionDelayDuration: retractionDelayDuration,
			ImageURL:                imageURL,
			YoutubeURL:              youtubeURL,
			WorldDownloadURL:        worldDownloadURL,
			ServerIPAddress:         serverIPAddress,
			ServerCoordinates:       serverCoordinates,
			ServerCommand:           serverCommand,
			SubmitterID:             strconv.Itoa(submitterIDint),
			Timestamp:               Timestamp(timestamp),
			EditedTimestamp:         Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// BuildCreate creates a new build
func (d *Database) BuildCreate(b Build) (Build, error) {
	// Convert ids to ints
	verifierIDint, err := strconv.Atoi(b.VerifierID)
	if err != nil {
		return Build{}, errors.Wrap(err, "failed to convert verifier id to integer")
	}
	reporterIDint, err := strconv.Atoi(b.ReporterID)
	if err != nil {
		return Build{}, errors.Wrap(err, "failed to convert reporter id to integer")
	}
	updateRequestBuildIDint, err := strconv.Atoi(b.UpdateRequestBuildID)
	if err != nil {
		return Build{}, errors.Wrap(err, "failed to convert update request build id to integer")
	}
	editionIDint, err := strconv.Atoi(b.EditionID)
	if err != nil {
		return Build{}, errors.Wrap(err, "failed to convert edition id to integer")
	}
	buildClassIDint, err := strconv.Atoi(b.BuildClassID)
	if err != nil {
		return Build{}, errors.Wrap(err, "failed to convert build class id to integer")
	}
	submitterIDint, err := strconv.Atoi(b.SubmitterID)
	if err != nil {
		return Build{}, errors.Wrap(err, "failed to convert submitter id to integer")
	}
	// Edit build
	b.Timestamp = Timestamp(time.Now())
	b.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO Builds (Verified, VerifierID, VerifiedTimestamp, Reported,
			ReporterID, ReportedTimestamp, UpdateRequest, UpdateRequestBuildID,
			EditionID, BuildClassID, Name, Description, Creators,
			CreationTimestamp, Width, Height, Depth, NormalCloseDuration,
			NormalOpenDuration, VisibleCloseDuration, VisibleOpenDuration,
			DelayCloseDuration, DelayOpenDuration, ResetCloseDuration,
			ResetOpenDuration, ExtensionDuration, RetractionDuration,
			ExtensionDelayDuration, RetractionDelayDuration, ImageURL,
			YoutubeURL, WorldDownloadURL, ServerIPAddress, ServerCoordinates,
			ServerCommand, SubmitterID, Timestamp, EditedTimestamp
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?,
			?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?
		)
	`)
	if err != nil {
		return Build{}, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	res, err := s.Exec(
		d.btoi(b.Verified), verifierIDint, time.Time(b.VerifiedTimestamp).Format(timeLayout),
		d.btoi(b.Reported), reporterIDint, time.Time(b.ReportedTimestamp).Format(timeLayout),
		d.btoi(b.UpdateRequest), updateRequestBuildIDint, editionIDint, buildClassIDint,
		b.Name, b.Description, b.Creators, time.Time(b.CreationTimestamp).Format(timeLayout),
		b.Width, b.Height, b.Depth, b.NormalCloseDuration, b.NormalOpenDuration,
		b.VisibleCloseDuration, b.VisibleOpenDuration, b.DelayCloseDuration, b.DelayOpenDuration,
		b.ResetCloseDuration, b.ResetOpenDuration, b.ExtensionDuration, b.RetractionDuration,
		b.ExtensionDelayDuration, b.RetractionDelayDuration, b.ImageURL, b.YoutubeURL,
		b.WorldDownloadURL, b.ServerIPAddress, b.ServerCoordinates, b.ServerCommand,
		submitterIDint,
		time.Time(b.Timestamp).Format(timeLayout),
		time.Time(b.EditedTimestamp).Format(timeLayout),
	)
	if err != nil {
		return Build{}, errors.Wrap(err, "database query failed")
	}
	// Update build id
	idInt, err := res.LastInsertId()
	if err != nil {
		return Build{}, errors.Wrap(err, "couldn't update build id")
	}
	// strconv.Itoa(int(idInt)) would probably suffice
	// but doesn't give an error in case precision is lost
	// from int(idInt) int64 -> int
	b.ID = strconv.FormatInt(idInt, 10)
	return b, nil
}

// BuildDelete removes build information from the database
func (d *Database) BuildDelete(buildID string) (Build, bool, error) {
	// Convert buildID to int
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	// Get the build to return after the deletion and
	// to check if it exists
	b, ok, err := d.Build(buildID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to determine if build exists")
	} else if !ok {
		// Row doesn't exist
		return Build{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM Builds
		WHERE ID = ?
	`)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(buildIDint); err != nil {
		return Build{}, false, errors.Wrap(err, "database query failed")
	}
	return b, true, nil
}

// BuildEdit edits the information for a build in the database
func (d *Database) BuildEdit(buildID string, build Build) (Build, bool, error) {
	// Convert ids to int
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	verifierIDint, err := strconv.Atoi(build.VerifierID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to convert verifier id to integer")
	}
	reporterIDint, err := strconv.Atoi(build.ReporterID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to convert reporter id to integer")
	}
	updateRequestBuildIDint, err := strconv.Atoi(build.UpdateRequestBuildID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to convert update request build id to integer")
	}
	editionIDint, err := strconv.Atoi(build.EditionID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to convert edition id to integer")
	}
	buildClassIDint, err := strconv.Atoi(build.BuildClassID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to convert build class id to integer")
	}
	submitterIDint, err := strconv.Atoi(build.SubmitterID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to convert submitter id to integer")
	}
	// Get the build that is to be updated
	b, ok, err := d.Build(buildID)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to determine if build exists")
	} else if !ok {
		// Row doesn't exist
		return Build{}, false, nil
	}
	// Update information
	b.Verified = build.Verified
	b.VerifierID = build.VerifierID
	b.VerifiedTimestamp = build.VerifiedTimestamp
	b.Reported = build.Reported
	b.ReporterID = build.ReporterID
	b.ReportedTimestamp = build.ReportedTimestamp
	b.UpdateRequest = build.UpdateRequest
	b.UpdateRequestBuildID = build.UpdateRequestBuildID
	b.EditionID = build.EditionID
	b.BuildClassID = build.BuildClassID
	b.Name = build.Name
	b.Description = build.Description
	b.Creators = build.Creators
	b.CreationTimestamp = build.CreationTimestamp
	b.Width = build.Width
	b.Height = build.Height
	b.Depth = build.Depth
	b.NormalCloseDuration = build.NormalCloseDuration
	b.NormalOpenDuration = build.NormalOpenDuration
	b.VisibleCloseDuration = build.VisibleCloseDuration
	b.VisibleOpenDuration = build.VisibleOpenDuration
	b.DelayCloseDuration = build.DelayCloseDuration
	b.DelayOpenDuration = build.DelayOpenDuration
	b.ResetCloseDuration = build.ResetCloseDuration
	b.ResetOpenDuration = build.ResetOpenDuration
	b.ExtensionDuration = build.ExtensionDuration
	b.RetractionDuration = build.RetractionDuration
	b.ExtensionDelayDuration = build.ExtensionDelayDuration
	b.RetractionDelayDuration = build.RetractionDelayDuration
	b.ImageURL = build.ImageURL
	b.YoutubeURL = build.YoutubeURL
	b.WorldDownloadURL = build.WorldDownloadURL
	b.ServerIPAddress = build.ServerIPAddress
	b.ServerCoordinates = build.ServerCoordinates
	b.ServerCommand = build.ServerCommand
	b.SubmitterID = build.SubmitterID
	b.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE Builds
		SET Verified = ?, VerifierID = ?, VerifiedTimestamp = ?, Reported = ?,
			ReporterID = ?, ReportedTimestamp = ?, UpdateRequest = ?,
			UpdateRequestBuildID = ?, EditionID = ?, BuildClassID = ?,
			Name = ?, Description = ?, Creators = ?, CreationTimestamp = ?,
			Width = ?, Height = ?, Depth = ?, NormalCloseDuration = ?,
			NormalOpenDuration = ?, VisibleCloseDuration = ?, VisibleOpenDuration = ?,
			DelayCloseDuration = ?, DelayOpenDuration = ?, ResetCloseDuration = ?,
			ResetOpenDuration = ?, ExtensionDuration = ?, RetractionDuration = ?,
			ExtensionDelayDuration = ?, RetractionDelayDuration = ?, ImageURL = ?,
			YoutubeURL = ?, WorldDownloadURL = ?, ServerIPAddress = ?,
			ServerCoordinates = ?, ServerCommand = ?, SubmitterID = ?, EditedTimestamp = ?
		WHERE ID = ?
	`)
	if err != nil {
		return Build{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		d.btoi(b.Verified), verifierIDint, time.Time(b.VerifiedTimestamp).Format(timeLayout),
		d.btoi(b.Reported), reporterIDint, time.Time(b.ReportedTimestamp).Format(timeLayout),
		d.btoi(b.UpdateRequest), updateRequestBuildIDint, strconv.Itoa(editionIDint),
		strconv.Itoa(buildClassIDint), b.Name, b.Description, b.Creators,
		time.Time(b.CreationTimestamp).Format(timeLayout), b.Width, b.Height, b.Depth,
		b.NormalCloseDuration, b.NormalOpenDuration, b.VisibleCloseDuration, b.VisibleOpenDuration,
		b.DelayCloseDuration, b.DelayOpenDuration, b.ResetCloseDuration, b.ResetOpenDuration,
		b.ExtensionDuration, b.RetractionDuration, b.ExtensionDelayDuration, b.RetractionDelayDuration,
		b.ImageURL, b.YoutubeURL, b.WorldDownloadURL, b.ServerIPAddress, b.ServerCoordinates,
		b.ServerCommand, submitterIDint, time.Time(b.EditedTimestamp).Format(timeLayout), buildIDint,
	); err != nil {
		return Build{}, false, errors.Wrap(err, "database query failed")
	}
	return b, true, nil
}

// Version gets information for the specified version
func (d *Database) Version(versionID string) (Version, bool, error) {
	// Convert versionID to int
	versionIDint, err := strconv.Atoi(versionID)
	if err != nil {
		return Version{}, false, errors.Wrap(err, "failed to covnvert version id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT EditionID, MajorVersion, MinorVersion, Patch, Name,
			Description, VersionTimestamp, Timestamp, EditedTimestamp
		FROM Versions
		WHERE ID = ?
	`, versionIDint)
	if err != nil {
		return Version{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if version exists in database
	if !rows.Next() {
		return Version{}, false, nil
	}
	// Extract data
	var (
		editionIDint           int
		majorVersion           int
		minorVersion           int
		patch                  int
		name                   string
		description            string
		versionTimestampString string
		timestampString        string
		editedTimestampString  string
		versionTimestamp       time.Time
		timestamp              time.Time
		editedTimestamp        time.Time
	)
	if err = rows.Scan(
		&editionIDint, &majorVersion, &minorVersion, &patch,
		&name, &description, &versionTimestampString,
		&timestampString, &editedTimestampString,
	); err != nil {
		return Version{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if versionTimestamp, err = time.Parse(timeLayout, versionTimestampString); err != nil {
		return Version{}, false, errors.Wrap(err, "failed to parse version timestamp")
	}
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return Version{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return Version{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return Version{
		ID:               versionID,
		EditionID:        strconv.Itoa(editionIDint),
		MajorVersion:     majorVersion,
		MinorVersion:     minorVersion,
		Patch:            patch,
		Name:             name,
		Description:      description,
		VersionTimestamp: Timestamp(versionTimestamp),
		Timestamp:        Timestamp(timestamp),
		EditedTimestamp:  Timestamp(editedTimestamp),
	}, true, nil
}

// Versions gets information for all versions
func (d *Database) Versions() ([]Version, error) {
	// Query the database
	rows, err := d.db.Query(`
		SELECT *
		FROM Versions
	`)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []Version{}
	var (
		idInt                  int
		editionIDint           int
		majorVersion           int
		minorVersion           int
		patch                  int
		name                   string
		description            string
		versionTimestampString string
		timestampString        string
		editedTimestampString  string
		versionTimestamp       time.Time
		timestamp              time.Time
		editedTimestamp        time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&idInt, &editionIDint, &majorVersion, &minorVersion, &patch,
			&name, &description, &versionTimestampString,
			&timestampString, &editedTimestampString,
		); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Parse timestamps
		if versionTimestamp, err = time.Parse(timeLayout, versionTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse version timestamp")
		}
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}
		if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse edited timestamp")
		}
		// Add to results
		results = append(results, Version{
			ID:               strconv.Itoa(idInt),
			EditionID:        strconv.Itoa(editionIDint),
			MajorVersion:     majorVersion,
			MinorVersion:     minorVersion,
			Patch:            patch,
			Name:             name,
			Description:      description,
			VersionTimestamp: Timestamp(versionTimestamp),
			Timestamp:        Timestamp(timestamp),
			EditedTimestamp:  Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// VersionCreate creates a new version in the database
func (d *Database) VersionCreate(version Version) (Version, error) {
	// Convert ids to ints
	editionIDint, err := strconv.Atoi(version.EditionID)
	if err != nil {
		return Version{}, errors.Wrap(err, "failed to convert edition id to integer")
	}
	// Edit version
	version.Timestamp = Timestamp(time.Now())
	version.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO Versions (EditionID, MajorVersion, MinorVersion, Patch,
			Name, Description, VersionTimestamp, Timestamp, EditedTimestamp
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return Version{}, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	res, err := s.Exec(
		editionIDint, version.MajorVersion, version.MinorVersion, version.Patch,
		version.Name, version.Description,
		time.Time(version.VersionTimestamp).Format(timeLayout),
		time.Time(version.Timestamp).Format(timeLayout),
		time.Time(version.EditedTimestamp).Format(timeLayout),
	)
	if err != nil {
		return Version{}, errors.Wrap(err, "database query failed")
	}
	// Update version id
	idInt, err := res.LastInsertId()
	if err != nil {
		return Version{}, errors.Wrap(err, "couldn't update version id")
	}
	// strconv.Itoa(int(idInt)) would probably suffice
	// but doesn't give an error in case precision is lost
	// from int(idInt) int64 -> int
	version.ID = strconv.FormatInt(idInt, 10)
	return version, nil
}

// VersionDelete removes a version from the database
func (d *Database) VersionDelete(versionID string) (Version, bool, error) {
	// Convert version id to int
	versionIDint, err := strconv.Atoi(versionID)
	if err != nil {
		return Version{}, false, errors.Wrap(err, "failed to convert version id to integer")
	}
	// Get the version to return after deletion and
	// to check if it exists
	v, ok, err := d.Version(versionID)
	if err != nil {
		return Version{}, false, errors.Wrap(err, "failed to determine if version exists")
	} else if !ok {
		// Row doesn't exist
		return Version{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM Versions
		WHERE ID = ?
	`)
	if err != nil {
		return Version{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err := s.Exec(versionIDint); err != nil {
		return Version{}, false, errors.Wrap(err, "database query failed")
	}
	return v, true, nil
}

// VersionEdit edits the version information for a specified version
func (d *Database) VersionEdit(versionID string, version Version) (Version, bool, error) {
	// Convert ids to ints
	versionIDint, err := strconv.Atoi(versionID)
	if err != nil {
		return Version{}, false, errors.Wrap(err, "failed to convert version id to integer")
	}
	editionIDint, err := strconv.Atoi(version.EditionID)
	if err != nil {
		return Version{}, false, errors.Wrap(err, "failed to convert edition id to integer")
	}
	// Get the version that is to be updated
	v, ok, err := d.Version(versionID)
	if err != nil {
		return Version{}, false, errors.Wrap(err, "failed to determine if version exists")
	} else if !ok {
		// Row doesn't exist
		return Version{}, false, nil
	}
	// Update information
	v.EditionID = version.EditionID
	v.MajorVersion = version.MajorVersion
	v.MinorVersion = version.MinorVersion
	v.Patch = version.Patch
	v.Name = version.Name
	v.Description = version.Description
	v.VersionTimestamp = version.VersionTimestamp
	v.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE Versions
		SET EditionID = ?, MajorVersion = ?, MinorVersion = ?, Patch = ?, Name = ?,
			Description = ?, VersionTimestamp = ?, EditedTimestamp = ?
		WHERE ID = ?
	`)
	if err != nil {
		return Version{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		editionIDint, v.MajorVersion, v.MinorVersion, v.Patch,
		v.Name, v.Description,
		time.Time(v.VersionTimestamp).Format(timeLayout),
		time.Time(v.EditedTimestamp).Format(timeLayout),
		versionIDint,
	); err != nil {
		return Version{}, false, errors.Wrap(err, "database query failed")
	}
	return v, true, nil
}

// Record gets the information for a specified record
func (d *Database) Record(recordID string) (Record, bool, error) {
	// Convert record id to int
	recordIDint, err := strconv.Atoi(recordID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to convert record id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT Verified, VerifierID, VerifiedTimestap, UpdateRequest, UpdateRequestRecordID,
			EditionID, BuildClassID, RecordTypeID, Name, Description, SubmitterID,
			Timestamp, EditedTimestamp
		FROM Records
		WHERE ID = ?
	`, recordIDint)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if record exists in database
	if !rows.Next() {
		return Record{}, false, nil
	}
	// Extract data
	var (
		verifiedInt              int
		verifierIDint            int
		verifiedTimestampString  string
		updateRequestInt         int
		updateRequestRecordIDint int
		editionIDint             int
		buildClassIDint          int
		recordTypeIDint          int
		name                     string
		description              string
		submitterIDint           int
		timestampString          string
		editedTimestampString    string
		verifiedTimestamp        time.Time
		timestamp                time.Time
		editedTimestamp          time.Time
	)
	if err = rows.Scan(
		&verifiedInt, &verifierIDint, &verifiedTimestampString, &updateRequestInt,
		&updateRequestRecordIDint, &editionIDint, &buildClassIDint, &recordTypeIDint,
		&name, &description, &submitterIDint, &timestampString, &editedTimestampString,
	); err != nil {
		return Record{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if verifiedTimestamp, err = time.Parse(timeLayout, verifiedTimestampString); err != nil {
		return Record{}, false, errors.Wrap(err, "failed to parse verified timestamp")
	}
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return Record{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return Record{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return Record{
		ID:                    recordID,
		Verified:              verifiedInt != 0,
		VerifierID:            strconv.Itoa(verifierIDint),
		VerifiedTimestamp:     Timestamp(verifiedTimestamp),
		UpdateRequest:         updateRequestInt != 0,
		UpdateRequestRecordID: strconv.Itoa(updateRequestRecordIDint),
		EditionID:             strconv.Itoa(editionIDint),
		BuildClassID:          strconv.Itoa(buildClassIDint),
		RecordTypeID:          strconv.Itoa(recordTypeIDint),
		Name:                  name,
		Description:           description,
		SubmitterID:           strconv.Itoa(submitterIDint),
		Timestamp:             Timestamp(timestamp),
		EditedTimestamp:       Timestamp(editedTimestamp),
	}, true, nil
}

// Records gets information for all records in the database
func (d *Database) Records() ([]Record, error) {
	// Query the database
	rows, err := d.db.Query(`
		SELECT *
		FROM Records
	`)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []Record{}
	var (
		idInt                    int
		verifiedInt              int
		verifierIDint            int
		verifiedTimestampString  string
		updateRequestInt         int
		updateRequestRecordIDint int
		editionIDint             int
		buildClassIDint          int
		recordTypeIDint          int
		name                     string
		description              string
		submitterIDint           int
		timestampString          string
		editedTimestampString    string
		verifiedTimestamp        time.Time
		timestamp                time.Time
		editedTimestamp          time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&idInt, &verifiedInt, &verifierIDint, &verifiedTimestampString,
			&updateRequestInt, &updateRequestRecordIDint, &editionIDint,
			&buildClassIDint, &recordTypeIDint, &name, &description,
			&submitterIDint, &timestampString, &editedTimestampString,
		); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Parse timestamps
		if verifiedTimestamp, err = time.Parse(timeLayout, verifiedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse verified timestamp")
		}
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}
		if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse edited timestamp")
		}
		// Add to results
		results = append(results, Record{
			ID:                    strconv.Itoa(idInt),
			Verified:              verifiedInt != 0,
			VerifierID:            strconv.Itoa(verifierIDint),
			VerifiedTimestamp:     Timestamp(verifiedTimestamp),
			UpdateRequest:         updateRequestInt != 0,
			UpdateRequestRecordID: strconv.Itoa(updateRequestRecordIDint),
			EditionID:             strconv.Itoa(editionIDint),
			BuildClassID:          strconv.Itoa(buildClassIDint),
			RecordTypeID:          strconv.Itoa(recordTypeIDint),
			Name:                  name,
			Description:           description,
			SubmitterID:           strconv.Itoa(submitterIDint),
			Timestamp:             Timestamp(timestamp),
			EditedTimestamp:       Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// RecordCreate creates a new record
func (d *Database) RecordCreate(record Record) (Record, error) {
	// Convert ids to ints
	verifierIDint, err := strconv.Atoi(record.VerifierID)
	if err != nil {
		return Record{}, errors.Wrap(err, "failed to convert verifier id to integer")
	}
	updateRequestRecordIDint, err := strconv.Atoi(record.UpdateRequestRecordID)
	if err != nil {
		return Record{}, errors.Wrap(err, "failed to convert update request record id to integer")
	}
	editionIDint, err := strconv.Atoi(record.EditionID)
	if err != nil {
		return Record{}, errors.Wrap(err, "failed to convert edition id to integer")
	}
	buildClassIDint, err := strconv.Atoi(record.BuildClassID)
	if err != nil {
		return Record{}, errors.Wrap(err, "failed to convert build class id to integer")
	}
	recordTypeIDint, err := strconv.Atoi(record.RecordTypeID)
	if err != nil {
		return Record{}, errors.Wrap(err, "failed to convert record type id to integer")
	}
	submitterIDint, err := strconv.Atoi(record.SubmitterID)
	if err != nil {
		return Record{}, errors.Wrap(err, "failed to convert submitter id to integer")
	}
	// Edit record
	record.Timestamp = Timestamp(time.Now())
	record.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO Records (Verified, VerifierID, VerifiedTimestamp, UpdateRequest,
			UpdateRequestRecordID, EditionID, BuildClassID, RecordTypeID, Name,
			Description, SubmitterID, Timestamp, Editedtimestamp
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return Record{}, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	res, err := s.Exec(
		d.btoi(record.Verified), verifierIDint,
		time.Time(record.VerifiedTimestamp).Format(timeLayout),
		d.btoi(record.UpdateRequest), updateRequestRecordIDint,
		editionIDint, buildClassIDint, recordTypeIDint, record.Name,
		record.Description, submitterIDint,
		time.Time(record.Timestamp).Format(timeLayout),
		time.Time(record.EditedTimestamp).Format(timeLayout),
	)
	if err != nil {
		return Record{}, errors.Wrap(err, "database query failed")
	}
	// Update record id
	idInt, err := res.LastInsertId()
	if err != nil {
		return Record{}, errors.Wrap(err, "couldn't update record id")
	}
	// strconv.Itoa(int(idInt)) would probably suffice
	// but doesn't give an error in case precision is lost
	// from int(idInt) int64 -> int
	record.ID = strconv.FormatInt(idInt, 10)
	return record, nil
}

// RecordDelete removes a specified record from the database
func (d *Database) RecordDelete(recordID string) (Record, bool, error) {
	// Convert recordID to int
	recordIDint, err := strconv.Atoi(recordID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to convert record id to integer")
	}
	// Get the record to return after deletion and
	// to check if it exists
	r, ok, err := d.Record(recordID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to determine if record exists")
	} else if !ok {
		// Row doesn't exits
		return Record{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM Records
		WHERE ID = ?
	`)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err := s.Exec(recordIDint); err != nil {
		return Record{}, false, errors.Wrap(err, "database query failed")
	}
	return r, true, nil
}

// RecordEdit edits the information for a record in the database
func (d *Database) RecordEdit(recordID string, record Record) (Record, bool, error) {
	// Convert ids to ints
	recordIDint, err := strconv.Atoi(recordID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to convert record id to integer")
	}
	verifierIDint, err := strconv.Atoi(record.VerifierID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to convert verifier id to integer")
	}
	updateRequestRecordIDint, err := strconv.Atoi(record.UpdateRequestRecordID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to convert update request record id to integer")
	}
	editionIDint, err := strconv.Atoi(record.EditionID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to convert edition id to integer")
	}
	buildClassIDint, err := strconv.Atoi(record.BuildClassID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to convert build class id to integer")
	}
	recordTypeIDint, err := strconv.Atoi(record.RecordTypeID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to convert record type id to integer")
	}
	submitterIDint, err := strconv.Atoi(record.SubmitterID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to convert submitter id to integer")
	}
	// Get the record that is to be updated
	r, ok, err := d.Record(recordID)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to determine if record exists")
	} else if !ok {
		return Record{}, false, nil
	}
	// Update information
	r.Verified = record.Verified
	r.VerifierID = record.VerifierID
	r.VerifiedTimestamp = record.VerifiedTimestamp
	r.UpdateRequest = record.UpdateRequest
	r.UpdateRequestRecordID = record.UpdateRequestRecordID
	r.EditionID = record.EditionID
	r.BuildClassID = record.BuildClassID
	r.RecordTypeID = record.RecordTypeID
	r.Name = record.Name
	r.Description = record.Description
	r.SubmitterID = record.SubmitterID
	r.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE Records
		SET Verified = ?, VerifierID = ?, VerifiedTimestamp = ?, UpdateRequest = ?,
			UpdateRequestRecordID = ?, EditionID = ?, BuildClassID = ?, RecordTypeID = ?,
			Name = ?, Description = ?, SubmitterID = ?, EditedTimestamp = ?
		WHERE ID = ?
	`)
	if err != nil {
		return Record{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		d.btoi(r.Verified), verifierIDint, time.Time(r.VerifiedTimestamp).Format(timeLayout),
		d.btoi(r.UpdateRequest), updateRequestRecordIDint, editionIDint, buildClassIDint,
		recordTypeIDint, r.Name, r.Description, submitterIDint,
		time.Time(r.EditedTimestamp).Format(timeLayout), recordIDint,
	); err != nil {
		return Record{}, false, errors.Wrap(err, "database query failed")
	}
	return r, true, nil
}

// GuildBuildMessage gets the guild build message information for a specified
// guild and build
func (d *Database) GuildBuildMessage(guildID, buildID string) (GuildBuildMessage, bool, error) {
	// Convert guildID and buildID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT ChannelID, MessageID, Timestamp, EditedTimestamp
		FROM GuildBuildMessages
		WHERE GuildID = ? AND BuildID = ?
	`, guildIDint, buildIDint)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if the guild build message exists
	if !rows.Next() {
		return GuildBuildMessage{}, false, nil
	}
	// Extract data
	var (
		channelIDint          int
		messageIDint          int
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	if err = rows.Scan(
		&channelIDint, &messageIDint, &timestampString, &editedTimestampString,
	); err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parsing timestamps
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return GuildBuildMessage{
		GuildID:         guildID,
		BuildID:         buildID,
		ChannelID:       strconv.Itoa(channelIDint),
		MessageID:       strconv.Itoa(messageIDint),
		Timestamp:       Timestamp(timestamp),
		EditedTimestamp: Timestamp(editedTimestamp),
	}, true, nil
}

// GuildBuildMessages get the guild build message information for a
// specified guild
func (d *Database) GuildBuildMessages(guildID string) ([]GuildBuildMessage, error) {
	// Query the database
	rows, err := d.db.Query(`
		SELECT *
		FROM GuildBuildMessages
	`)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []GuildBuildMessage{}
	var (
		guildIDint            int
		buildIDint            int
		channelIDint          int
		messageIDint          int
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&guildIDint, &buildIDint, &channelIDint, &messageIDint,
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
		results = append(results, GuildBuildMessage{
			GuildID:         strconv.Itoa(guildIDint),
			BuildID:         strconv.Itoa(buildIDint),
			ChannelID:       strconv.Itoa(channelIDint),
			MessageID:       strconv.Itoa(messageIDint),
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// GuildBuildMessageCreate creates guild build message information in the database
func (d *Database) GuildBuildMessageCreate(guildID, buildID, channelID, messageID string) (GuildBuildMessage, bool, error) {
	// Convert ids to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	channelIDint, err := strconv.Atoi(channelID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert channel id to integer")
	}
	messageIDint, err := strconv.Atoi(messageID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert message id to integer")
	}
	// Check if guild build message already exists
	if _, ok, err := d.GuildBuildMessage(guildID, buildID); err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to determine if guild build message exists")
	} else if ok {
		// Row already exists
		return GuildBuildMessage{}, false, nil
	}
	// Create guild build message
	gbm := GuildBuildMessage{
		GuildID:         guildID,
		BuildID:         buildID,
		ChannelID:       channelID,
		MessageID:       messageID,
		Timestamp:       Timestamp(time.Now()),
		EditedTimestamp: Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO GuildBuildMessages
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		guildIDint, buildIDint, channelIDint, messageIDint,
		time.Time(gbm.Timestamp).Format(timeLayout),
		time.Time(gbm.EditedTimestamp).Format(timeLayout),
	); err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "database query failed")
	}
	return gbm, true, nil
}

// GuildBuildMessageDelete removes guild build message information from the database
func (d *Database) GuildBuildMessageDelete(guildID, buildID string) (GuildBuildMessage, bool, error) {
	// Convert guildID and buildID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	// Get the guild build message to return after deletion
	// and to check if it exists
	gbm, ok, err := d.GuildBuildMessage(guildID, buildID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to determine if guild build message exists")
	} else if !ok {
		// Row doesn't exits
		return GuildBuildMessage{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM GuildBuildMessages
		WHERE GuildID = ? AND BuildID = ?
	`)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(guildIDint, buildIDint); err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "database query failed")
	}
	return gbm, true, nil
}

// GuildBuildMessageEdit edits the build build message information for a specified
// guild and build
func (d *Database) GuildBuildMessageEdit(guildID, buildID, channelID, messageID string) (GuildBuildMessage, bool, error) {
	// Convert ids to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	channelIDint, err := strconv.Atoi(channelID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert channel id to integer")
	}
	messageIDint, err := strconv.Atoi(messageID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to convert message id to integer")
	}
	// Get the guild build message that is to be updated
	gbm, ok, err := d.GuildBuildMessage(guildID, buildID)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to determine if build build message exists")
	} else if !ok {
		// Row doesn't exist
		return GuildBuildMessage{}, false, nil
	}
	// Update values
	gbm.ChannelID = channelID
	gbm.MessageID = messageID
	gbm.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE GuildBuildMessages
		SET ChannelID = ?, MessageID = ?, EditedTimestamp = ?
		WHERE GuildID = ? AND BuildID = ?
	`)
	if err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "failed to prepare query")
	}
	// Execute query
	if _, err = s.Exec(
		channelIDint, messageIDint,
		time.Time(gbm.EditedTimestamp).Format(timeLayout),
		guildIDint, buildIDint,
	); err != nil {
		return GuildBuildMessage{}, false, errors.Wrap(err, "database query failed")
	}
	return gbm, true, nil
}

// BuildVersion gets specified build version information for
// a build and a version
func (d *Database) BuildVersion(buildID, versionID string) (BuildVersion, bool, error) {
	// Convert buildID and versionID to ints
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	versionIDint, err := strconv.Atoi(versionID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to convert version id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT StatusID, Notes, Timestamp, EditedTimestamp
		FROM BuildVersions
		WHERE BuildID = ? AND VersionID = ?
	`, buildIDint, versionIDint)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if build version exists
	if !rows.Next() {
		return BuildVersion{}, false, nil
	}
	// Extract data
	var (
		statusIDint           int
		notes                 string
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	if err = rows.Scan(&statusIDint, &notes, &timestampString, &editedTimestampString); err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parsing timestamps
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return BuildVersion{
		BuildID:         buildID,
		VersionID:       versionID,
		StatusID:        strconv.Itoa(statusIDint),
		Notes:           notes,
		Timestamp:       Timestamp(timestamp),
		EditedTimestamp: Timestamp(editedTimestamp),
	}, true, nil
}

// BuildVersionCreate creates information in the database for a specified
// build and version
func (d *Database) BuildVersionCreate(buildID, versionID, statusID, notes string) (BuildVersion, bool, error) {
	// Convert buildID, versionID and statusID to ints
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	versionIDint, err := strconv.Atoi(versionID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to convert version id to integer")
	}
	statusIDint, err := strconv.Atoi(statusID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to convert status id to integer")
	}
	// Check if the build version already exists
	if _, ok, err := d.BuildVersion(buildID, versionID); err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to determine if build version exists")
	} else if ok {
		// Row already exists
		return BuildVersion{}, false, nil
	}
	// Create build version
	bv := BuildVersion{
		BuildID:         buildID,
		VersionID:       versionID,
		StatusID:        statusID,
		Notes:           notes,
		Timestamp:       Timestamp(time.Now()),
		EditedTimestamp: Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO BuildVersions
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		buildIDint, versionIDint, statusIDint, notes,
		time.Time(bv.Timestamp).Format(timeLayout),
		time.Time(bv.EditedTimestamp).Format(timeLayout),
	); err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "database query failed")
	}
	return bv, true, nil
}

// BuildVersionDelete removes build version information from the database
// for a specified build and version
func (d *Database) BuildVersionDelete(buildID, versionID string) (BuildVersion, bool, error) {
	// Convert buildID and versionID to ints
	// Convert buildID, versionID and statusID to ints
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	versionIDint, err := strconv.Atoi(versionID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to convert version id to integer")
	}
	// Get the build version to return after deletion
	// and to check if it exists
	bv, ok, err := d.BuildVersion(buildID, versionID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to determine if build version exists")
	} else if !ok {
		// Row doesn't exist
		return BuildVersion{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM BuildVersions
		WHERE BuildID = ? AND VersionID = ?
	`)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(buildIDint, versionIDint); err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "database query failed")
	}
	return bv, true, nil
}

// BuildVersionEdit edits build version information from the database
// for a specified build and version
func (d *Database) BuildVersionEdit(buildID, versionID, statusID, notes string) (BuildVersion, bool, error) {
	// Convert buildID, versionID and statusID to ints
	buildIDint, err := strconv.Atoi(buildID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	versionIDint, err := strconv.Atoi(versionID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to convert version id to integer")
	}
	statusIDint, err := strconv.Atoi(statusID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to convert status id to integer")
	}
	// Get the build version that is to be updated
	bv, ok, err := d.BuildVersion(buildID, versionID)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to determine if build version exists")
	} else if !ok {
		// Row doesn't exist
		return BuildVersion{}, false, nil
	}
	// Update values
	bv.StatusID = statusID
	bv.Notes = notes
	bv.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE BuildVersions
		SET StatusID = ?, Notes = ?, EditedTimestamp = ?
		WHERE BuildID = ? AND VersionID = ?
	`)
	if err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		statusIDint, notes, time.Time(bv.EditedTimestamp).Format(timeLayout),
		buildIDint, versionIDint,
	); err != nil {
		return BuildVersion{}, false, errors.Wrap(err, "database query failed")
	}
	return bv, true, nil
}

// Status gets a specified status's information
func (d *Database) Status(statusID string) (Status, bool, error) {
	// Convert statusID to int
	statusIDint, err := strconv.Atoi(statusID)
	if err != nil {
		return Status{}, false, errors.Wrap(err, "failed to convert status id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT Name, Description, Timestamp, EditedTimestamp
		FROM Statuses
		WHERE ID = ?
	`, statusIDint)
	if err != nil {
		return Status{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if status exists in database
	if !rows.Next() {
		return Status{}, false, nil
	}
	// Extract data
	var (
		name                  string
		description           string
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	if err = rows.Scan(&name, &description, &timestampString, &editedTimestampString); err != nil {
		return Status{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return Status{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return Status{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return Status{
		ID:              statusID,
		Name:            name,
		Description:     description,
		Timestamp:       Timestamp(timestamp),
		EditedTimestamp: Timestamp(editedTimestamp),
	}, true, nil
}

// Statuses gets all statuses and their information
func (d *Database) Statuses() ([]Status, error) {
	// Query the database
	rows, err := d.db.Query(`
		SELECT *
		FROM Statuses
	`)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []Status{}
	var (
		idInt                 int
		name                  string
		description           string
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(&idInt, &name, &description, &timestampString, &editedTimestampString); err != nil {
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
		results = append(results, Status{
			ID:              strconv.Itoa(idInt),
			Name:            name,
			Description:     description,
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// StatusCreate creates a new status
func (d *Database) StatusCreate(name, description string) (Status, error) {
	// Create status
	status := Status{
		Name:            name,
		Description:     description,
		Timestamp:       Timestamp(time.Now()),
		EditedTimestamp: Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO Statuses (Name, Description, Timestamp, EditedTimestamp)
		VALUES (?, ?, ?, ?)
	`)
	if err != nil {
		return Status{}, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	res, err := s.Exec(
		name, description,
		time.Time(status.Timestamp).Format(timeLayout),
		time.Time(status.EditedTimestamp).Format(timeLayout),
	)
	if err != nil {
		return Status{}, errors.Wrap(err, "database query failed")
	}
	// Update status id
	idInt, err := res.LastInsertId()
	if err != nil {
		return Status{}, errors.Wrap(err, "couldn't update status id")
	}
	// strconv.Itoa(int(idInt)) would probably suffice
	// but doesn't give an error in case precision is lost
	// from int(idInt) int64 -> int
	status.ID = strconv.FormatInt(idInt, 10)
	return status, nil
}

// StatusDelete removes a status
func (d *Database) StatusDelete(statusID string) (Status, bool, error) {
	// Convert status id to int
	statusIDint, err := strconv.Atoi(statusID)
	if err != nil {
		return Status{}, false, errors.Wrap(err, "failed to convert status id to integer")
	}
	// Get the status to return after deletion and
	// to check if it exists
	status, ok, err := d.Status(statusID)
	if err != nil {
		return Status{}, false, errors.Wrap(err, "failed to determine if status exists")
	} else if !ok {
		// Row doesn't exist
		return Status{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM Statuses
		WHERE ID = ?
	`)
	if err != nil {
		return Status{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err := s.Exec(statusIDint); err != nil {
		return Status{}, false, errors.Wrap(err, "database query failed")
	}
	return status, true, nil
}

// StatusEdit edits a status
func (d *Database) StatusEdit(statusID, name, description string) (Status, bool, error) {
	// Convert statusID to int
	statusIDint, err := strconv.Atoi(statusID)
	if err != nil {
		return Status{}, false, errors.Wrap(err, "failed to convert status id to integer")
	}
	// Get the status that is to be updated
	status, ok, err := d.Status(statusID)
	if err != nil {
		return Status{}, false, errors.Wrap(err, "failed to determine if status exists")
	} else if !ok {
		// Row doesn't exist
		return Status{}, false, nil
	}
	// Update information
	status.Name = name
	status.Description = description
	status.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE Statuses
		SET Name = ?, Description = ?, EditedTimestamp = ?
		WHERE ID = ?
	`)
	if err != nil {
		return Status{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		name, description,
		time.Time(status.EditedTimestamp).Format(timeLayout),
		statusIDint,
	); err != nil {
		return Status{}, false, errors.Wrap(err, "database query failed")
	}
	return status, true, nil
}

// BuildRecord gets build record information for a build record
func (d *Database) BuildRecord(buildRecordID string) (BuildRecord, bool, error) {
	// Convert buildRecordID to int
	buildRecordIDint, err := strconv.Atoi(buildRecordID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert build record id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT BuildID, RecordID, Verified, VerifierID, VerifiedTimestamp, Reported, ReporterID,
			ReportedTimestamp, JointBuildRecord, JointBuildRecordID, SubmitterID, Timestamp, EditedTimestamp
		FROM BuildRecords
		WHERE ID = ?
	`, buildRecordIDint)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if build record exists
	if !rows.Next() {
		return BuildRecord{}, false, nil
	}
	// Extract data
	var (
		buildIDint              int
		recordIDint             int
		verifiedInt             int
		verifierIDint           int
		verifiedTimestampString string
		reportedInt             int
		reporterIDint           int
		reportedTimestampString string
		jointBuildRecordInt     int
		jointBuildRecordIDint   int
		submitterIDint          int
		timestampString         string
		editedTimestampString   string
		verifiedTimestamp       time.Time
		reportedTimestamp       time.Time
		timestamp               time.Time
		editedTimestamp         time.Time
	)
	if err = rows.Scan(
		&buildIDint, &recordIDint, &verifiedInt, &verifierIDint, &verifiedTimestampString,
		&reportedInt, &reporterIDint, &reportedTimestampString, &jointBuildRecordInt,
		&jointBuildRecordIDint, &submitterIDint, &timestampString, &editedTimestampString,
	); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if verifiedTimestamp, err = time.Parse(timeLayout, verifiedTimestampString); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to parse verified timestamp")
	}
	if reportedTimestamp, err = time.Parse(timeLayout, reportedTimestampString); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to parse reported timestamp")
	}
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return BuildRecord{
		ID:                 buildRecordID,
		BuildID:            strconv.Itoa(buildIDint),
		RecordID:           strconv.Itoa(recordIDint),
		Verified:           verifiedInt != 0,
		VerifierID:         strconv.Itoa(verifierIDint),
		VerifiedTimestamp:  Timestamp(verifiedTimestamp),
		Reported:           reportedInt != 0,
		ReporterID:         strconv.Itoa(reporterIDint),
		ReportedTimestamp:  Timestamp(reportedTimestamp),
		JointBuildRecord:   jointBuildRecordInt != 0,
		JointBuildRecordID: strconv.Itoa(jointBuildRecordIDint),
		SubmitterID:        strconv.Itoa(submitterIDint),
		Timestamp:          Timestamp(timestamp),
		EditedTimestamp:    Timestamp(editedTimestamp),
	}, true, nil
}

// BuildRecordCreate creates new build record information
func (d *Database) BuildRecordCreate(br BuildRecord) (BuildRecord, error) {
	// Convert ids to ints
	buildIDint, err := strconv.Atoi(br.BuildID)
	if err != nil {
		return BuildRecord{}, errors.Wrap(err, "failed to convert build id to integer")
	}
	recordIDint, err := strconv.Atoi(br.RecordID)
	if err != nil {
		return BuildRecord{}, errors.Wrap(err, "failed to convert record id to integer")
	}
	verifierIDint, err := strconv.Atoi(br.VerifierID)
	if err != nil {
		return BuildRecord{}, errors.Wrap(err, "failed to convert verifier id to integer")
	}
	reporterIDint, err := strconv.Atoi(br.ReporterID)
	if err != nil {
		return BuildRecord{}, errors.Wrap(err, "failed to convert reporter id to integer")
	}
	jointBuildRecordIDint, err := strconv.Atoi(br.JointBuildRecordID)
	if err != nil {
		return BuildRecord{}, errors.Wrap(err, "failed to convert joint build record id to integer")
	}
	submitterIDint, err := strconv.Atoi(br.SubmitterID)
	if err != nil {
		return BuildRecord{}, errors.Wrap(err, "failed to convert submitter id to integer")
	}
	// Edit build record
	br.Timestamp = Timestamp(time.Now())
	br.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO BuildRecords (BuildID, RecordID, Verified, VerifierID, VerifiedTimestamp,
			Reported, ReporterID, ReportedTimestamp, JointBuildRecord, JointBuildRecordID,
			SubmitterID, Timestamp, EditedTimestamp
		)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return BuildRecord{}, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	res, err := s.Exec(
		buildIDint, recordIDint, d.btoi(br.Verified), verifierIDint,
		time.Time(br.VerifiedTimestamp).Format(timeLayout), d.btoi(br.Reported),
		reporterIDint, time.Time(br.ReportedTimestamp).Format(timeLayout),
		d.btoi(br.JointBuildRecord), jointBuildRecordIDint, submitterIDint,
		time.Time(br.Timestamp).Format(timeLayout),
		time.Time(br.EditedTimestamp).Format(timeLayout),
	)
	if err != nil {
		return BuildRecord{}, errors.Wrap(err, "database query failed")
	}
	// Update build record id
	idInt, err := res.LastInsertId()
	if err != nil {
		return BuildRecord{}, errors.Wrap(err, "couldn't update build record id")
	}
	// strconv.Itoa(int(idInt)) would probably suffice
	// but doesn't give an error in case precision is lost
	// from int(idInt) int64 -> int
	br.ID = strconv.FormatInt(idInt, 10)
	return br, nil
}

// BuildRecordDelete removes build record information from the database
func (d *Database) BuildRecordDelete(buildRecordID string) (BuildRecord, bool, error) {
	// Convert build record id to int
	buildRecordIDint, err := strconv.Atoi(buildRecordID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert build record id to integer")
	}
	// Get the build record to return after deletion and
	// to check if it exists
	br, ok, err := d.BuildRecord(buildRecordID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to determine if build record exists")
	} else if !ok {
		// Row doesn't exist
		return BuildRecord{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM BuildRecords
		WHERE ID = ?
	`)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err := s.Exec(buildRecordIDint); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "database query failed")
	}
	return br, true, nil
}

// BuildRecordEdit edits build record information within the database
func (d *Database) BuildRecordEdit(buildRecordID string, br BuildRecord) (BuildRecord, bool, error) {
	// Convert ids to ints
	buildRecordIDint, err := strconv.Atoi(buildRecordID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert build record id to integer")
	}
	buildIDint, err := strconv.Atoi(br.BuildID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert build id to integer")
	}
	recordIDint, err := strconv.Atoi(br.RecordID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert record id to integer")
	}
	verifiedIDint, err := strconv.Atoi(br.VerifierID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert verifier id to integer")
	}
	reporterIDint, err := strconv.Atoi(br.ReporterID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert reporter id to integer")
	}
	jointBuildRecordIDint, err := strconv.Atoi(br.JointBuildRecordID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert joint build record id to integer")
	}
	submitterIDint, err := strconv.Atoi(br.SubmitterID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to convert submitter id to integer")
	}
	// Get the build record that is to be updated
	_, ok, err := d.BuildRecord(buildRecordID)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to determine if build record exists")
	} else if !ok {
		// Row doesn't exist
		return BuildRecord{}, false, nil
	}
	// Update information
	br.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE BuildRecords
		SET BuildID = ?, RecordID = ?, Verified = ?, VerifierID = ?, VerifiedTimestamp = ?,
			Reported = ?, ReporterID = ? ReportedTimestamp = ?, JointBuildRecord = ?,
			JointBuildRecordID = ?, SubmitterID = ?, EditedTimestamp = ?
		WHERE ID = ?
	`)
	if err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err := s.Exec(
		buildIDint, recordIDint, d.btoi(br.Verified), verifiedIDint,
		time.Time(br.VerifiedTimestamp).Format(timeLayout), d.btoi(br.Reported),
		reporterIDint, time.Time(br.ReportedTimestamp).Format(timeLayout),
		d.btoi(br.JointBuildRecord), jointBuildRecordIDint, submitterIDint,
		time.Time(br.EditedTimestamp).Format(timeLayout), buildRecordIDint,
	); err != nil {
		return BuildRecord{}, false, errors.Wrap(err, "database query failed")
	}
	return br, true, nil
}

// GuildRecordMessage gets the guild record message information for a specified
// guild and record
func (d *Database) GuildRecordMessage(guildID, recordID string) (GuildRecordMessage, bool, error) {
	// Convert guildID and recordID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	recordIDint, err := strconv.Atoi(recordID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert record id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT ChannelID, MessageID, Timestamp, EditedTimestamp
		FROM GuildRecordMessages
		WHERE GuildID = ? AND RecordID = ?
	`, guildIDint, recordIDint)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if guild record message exists
	if !rows.Next() {
		return GuildRecordMessage{}, false, nil
	}
	// Extract data
	var (
		channelIDint          int
		messageIDint          int
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	if err = rows.Scan(&channelIDint, &messageIDint, &timestampString, &editedTimestampString); err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to extract data")
	}
	// Parse timestamps
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	if editedTimestamp, err = time.Parse(timeLayout, editedTimestampString); err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to parse edited timestamp")
	}
	return GuildRecordMessage{
		GuildID:         guildID,
		RecordID:        recordID,
		ChannelID:       strconv.Itoa(channelIDint),
		MessageID:       strconv.Itoa(messageIDint),
		Timestamp:       Timestamp(timestamp),
		EditedTimestamp: Timestamp(editedTimestamp),
	}, true, nil
}

// GuildRecordMessages gets the guild record message information for a specified guild
func (d *Database) GuildRecordMessages(guildID string) ([]GuildRecordMessage, error) {
	// Convert guildID to int
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert guild id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT RecordID, ChannelID, MessageID, Timestamp, EditedTimestamp
		FROM GuildRecordMessages
		WHERE GuildID = ?
	`, guildIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []GuildRecordMessage{}
	var (
		recordIDint           int
		channelIDint          int
		messageIDint          int
		timestampString       string
		editedTimestampString string
		timestamp             time.Time
		editedTimestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&recordIDint, &channelIDint, &messageIDint,
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
		results = append(results, GuildRecordMessage{
			GuildID:         guildID,
			RecordID:        strconv.Itoa(recordIDint),
			ChannelID:       strconv.Itoa(channelIDint),
			MessageID:       strconv.Itoa(messageIDint),
			Timestamp:       Timestamp(timestamp),
			EditedTimestamp: Timestamp(editedTimestamp),
		})
	}
	return results, nil
}

// GuildRecordMessageCreate creates guild record message information for a specified
// guild and record
func (d *Database) GuildRecordMessageCreate(guildID, recordID, channelID, messageID string) (GuildRecordMessage, bool, error) {
	// Convert ids to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert guild id integer")
	}
	recordIDint, err := strconv.Atoi(recordID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert record id integer")
	}
	channelIDint, err := strconv.Atoi(channelID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert channel id integer")
	}
	messageIDint, err := strconv.Atoi(messageID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert message id integer")
	}
	// Check if guild record message already exists
	if _, ok, err := d.GuildRecordMessage(guildID, recordID); err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to determine if guild record message exists")
	} else if ok {
		// Row already exist
		return GuildRecordMessage{}, false, nil
	}
	// Create build record message
	grm := GuildRecordMessage{
		GuildID:         guildID,
		RecordID:        recordID,
		ChannelID:       channelID,
		MessageID:       messageID,
		Timestamp:       Timestamp(time.Now()),
		EditedTimestamp: Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO GuildRecordMessages
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err := s.Exec(
		guildIDint, recordIDint, channelIDint, messageIDint,
		time.Time(grm.Timestamp).Format(timeLayout),
		time.Time(grm.EditedTimestamp).Format(timeLayout),
	); err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "database query failed")
	}
	return grm, true, nil
}

// GuildRecordMessageDelete removes guild record message information for a specified
// guild and record
func (d *Database) GuildRecordMessageDelete(guildID, recordID string) (GuildRecordMessage, bool, error) {
	// Convert guildID and recordID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	recordIDint, err := strconv.Atoi(recordID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert record id to integer")
	}
	// Get the guild record message to return after deletion
	// and to check if it exists
	grm, ok, err := d.GuildRecordMessage(guildID, recordID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to determine if guild record message exists")
	} else if !ok {
		// Row doesn't exist
		return GuildRecordMessage{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM GuildRecordMessages
		WHERE GuildID = ? AND RecordID = ?
	`)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(guildIDint, recordIDint); err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "database query failed")
	}
	return grm, true, nil
}

// GuildRecordMessageEdit edits guild record message information for a specified
// guild and record
func (d *Database) GuildRecordMessageEdit(guildID, recordID, channelID, messageID string) (GuildRecordMessage, bool, error) {
	// Convert ids to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	recordIDint, err := strconv.Atoi(recordID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert record id to integer")
	}
	channelIDint, err := strconv.Atoi(channelID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert channel id to integer")
	}
	messageIDint, err := strconv.Atoi(messageID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to convert message id to integer")
	}
	// Get the guild record message that is to be updated
	grm, ok, err := d.GuildRecordMessage(guildID, recordID)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to determine if guild record message exists")
	} else if !ok {
		// Row doesn't exist
		return GuildRecordMessage{}, false, nil
	}
	// Update values
	grm.ChannelID = channelID
	grm.MessageID = messageID
	grm.EditedTimestamp = Timestamp(time.Now())
	// Prepare query
	s, err := d.db.Prepare(`
		UPDATE GuildRecordMessages
		SET ChannelID = ?, MessageID = ?, EditedTimestamp = ?
		WHERE GuildID = ?, RecordID = ?
	`)
	if err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		channelIDint, messageIDint,
		time.Time(grm.EditedTimestamp).Format(timeLayout),
		guildIDint, recordIDint,
	); err != nil {
		return GuildRecordMessage{}, false, errors.Wrap(err, "database query failed")
	}
	return grm, true, nil
}

// GuildTicketChannel gets information for a specified ticket within a guild
func (d *Database) GuildTicketChannel(guildID, channelID string) (GuildTicketChannel, bool, error) {
	// Convert guildID and channelID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	channelIDint, err := strconv.Atoi(channelID)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to convert channel id to integer")
	}
	// Query database
	rows, err := d.db.Query(`
		SELECT TicketID, TicketType, CreatorID, Timestamp
		FROM GuildTicketChannels
		WHERE GuildID = ? AND ChannelID = ?
	`, guildIDint, channelIDint)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Check if guild ticket channel exists
	if !rows.Next() {
		return GuildTicketChannel{}, false, nil
	}
	// Extract data
	var (
		ticketIDint     int
		ticketType      int
		creatorIDint    int
		timestampString string
		timestamp       time.Time
	)
	if err = rows.Scan(&ticketIDint, &ticketType, &creatorIDint, &timestampString); err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "database query failed")
	}
	// Parse timestamp
	if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to parse timestamp")
	}
	return GuildTicketChannel{
		GuildID:    guildID,
		ChannelID:  channelID,
		TicketID:   strconv.Itoa(ticketIDint),
		TicketType: TicketType(ticketType),
		CreatorID:  strconv.Itoa(creatorIDint),
		Timestamp:  Timestamp(timestamp),
	}, true, nil
}

// GuildTicketChannels gets information for all tickets within a guild
func (d *Database) GuildTicketChannels(guildID string) ([]GuildTicketChannel, error) {
	// Convert guildID to int
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return nil, errors.Wrap(err, "failed to convert guild id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT ChannelID, TicketID, TicketType, CreatorID, Timestamp
		FROM GuildTicketChannels
		WHERE GuildID = ?
	`, guildIDint)
	if err != nil {
		return nil, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// Create space to store results
	results := []GuildTicketChannel{}
	var (
		channelIDint    int
		ticketIDint     int
		ticketType      int
		creatorIDint    int
		timestampString string
		timestamp       time.Time
	)
	// For each row
	for rows.Next() {
		// Extract data
		if err = rows.Scan(
			&channelIDint, &ticketIDint, &ticketType,
			&creatorIDint, &timestampString,
		); err != nil {
			return nil, errors.Wrap(err, "failed to extract data")
		}
		// Parse timestamp
		if timestamp, err = time.Parse(timeLayout, timestampString); err != nil {
			return nil, errors.Wrap(err, "failed to parse timestamp")
		}
		// Add to results
		results = append(results, GuildTicketChannel{
			GuildID:    guildID,
			ChannelID:  strconv.Itoa(channelIDint),
			TicketID:   strconv.Itoa(ticketIDint),
			TicketType: TicketType(ticketType),
			CreatorID:  strconv.Itoa(creatorIDint),
			Timestamp:  Timestamp(timestamp),
		})
	}
	return results, nil
}

// GuildTicketChannelCreate creates a new ticket channel within the database
func (d *Database) GuildTicketChannelCreate(guildID, channelID string, ticketType TicketType, creatorID string) (GuildTicketChannel, bool, error) {
	// Convert ids to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	channelIDint, err := strconv.Atoi(channelID)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to convert channel id to integer")
	}
	creatorIDint, err := strconv.Atoi(creatorID)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to convert creator id to integer")
	}
	// Check if guild ticket channel already exists
	if _, ok, err := d.GuildTicketChannel(guildID, channelID); err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to determine if guild ticket channel exists")
	} else if ok {
		// Row already exists
		return GuildTicketChannel{}, false, nil
	}
	// Get next ticket id
	ticketIDint, err := d.nextTicketID(guildID)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to get next ticket id")
	}
	// Create guild ticket channel
	gtc := GuildTicketChannel{
		GuildID:    guildID,
		ChannelID:  channelID,
		TicketID:   strconv.Itoa(ticketIDint),
		TicketType: ticketType,
		CreatorID:  creatorID,
		Timestamp:  Timestamp(time.Now()),
	}
	// Prepare query
	s, err := d.db.Prepare(`
		INSERT INTO GuildTicketChannels
		VALUES (?, ?, ?, ?, ?, ?)
	`)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(
		guildIDint, channelIDint, ticketIDint, ticketType,
		creatorIDint, time.Time(gtc.Timestamp).Format(timeLayout),
	); err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "database query failed")
	}
	return gtc, true, nil
}

// GuildTicketChannelDelete removes an existing ticket channel from the database
func (d *Database) GuildTicketChannelDelete(guildID, channelID string) (GuildTicketChannel, bool, error) {
	// Convert guildID and channelID to ints
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to convert guild id to integer")
	}
	channelIDint, err := strconv.Atoi(channelID)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to convert channel id to integer")
	}
	// Get the guild ticket channel to return after deletion
	// and to check if it exists
	gtc, ok, err := d.GuildTicketChannel(guildID, channelID)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to determine if guild ticket channel exists")
	} else if !ok {
		// Row doesn't exist
		return GuildTicketChannel{}, false, nil
	}
	// Prepare query
	s, err := d.db.Prepare(`
		DELETE FROM GuildTicketChannels
		WHERE GuildID = ? AND ChannelID = ?
	`)
	if err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "failed to prepare query")
	}
	defer s.Close()
	// Execute query
	if _, err = s.Exec(guildIDint, channelIDint); err != nil {
		return GuildTicketChannel{}, false, errors.Wrap(err, "database query failed")
	}
	return gtc, true, nil
}

// Private functions

// nextStrikeID gets the next strike id for a specified user
func (d *Database) nextStrikeID(userID string) (int, error) {
	// Convert userID to int
	userIDint, err := strconv.Atoi(userID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert user id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT COALESCE (
			(
				SELECT MAX(StrikeID) + 1
				FROM UserStrikes
				WHERE UserID = ?
			),
			0
		)
	`, userIDint)
	if err != nil {
		return 0, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// The query should always return a value
	// Therefore the if block shouldn't be ran
	if !rows.Next() {
		return 0, errors.New("query didn't return a value")
	}
	// Extract data
	var strikeID int
	if err = rows.Scan(&strikeID); err != nil {
		return 0, errors.Wrap(err, "failed to extract data")
	}
	return strikeID, nil
}

func (d *Database) nextTicketID(guildID string) (int, error) {
	// Convert guildID to int
	guildIDint, err := strconv.Atoi(guildID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to convert guild id to integer")
	}
	// Query the database
	rows, err := d.db.Query(`
		SELECT COALESCE (
			(
				SELECT MAX(TicketID) + 1
				FROM GuildTicketChannels
				WHERE GuildID = ?
			),
			0
		)
	`, guildIDint)
	if err != nil {
		return 0, errors.Wrap(err, "database query failed")
	}
	defer rows.Close()
	// The query should always return a value
	// Therefore the if block shouldn't be ran
	if !rows.Next() {
		return 0, errors.New("query didn't return a value")
	}
	var ticketID int
	if err = rows.Scan(&ticketID); err != nil {
		return 0, errors.Wrap(err, "failed to extract data")
	}
	return ticketID, nil
}

// btoi converts a bool to an int
// true 	-> 1
// false	-> 0
func (d *Database) btoi(b bool) int {
	if b {
		return 1
	}
	return 0
}
