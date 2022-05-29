package persist

import (
	"database/sql"
	"encoding/json"
	"fmt"

	"Webex.API.Integration.And.Visualization/types"
)

// The db instance is sharable between multiple goroutines.
// It can therefore be initialized once in the main and used for the server's lifetime.
type Persist struct {
	db *sql.DB
}

// NewPersist create a new instance of Persits provided a database pointer.
func NewPersist(db *sql.DB) (*Persist, error) {
	// create if not exists the table
	_, err := db.Exec(
		"CREATE TABLE IF NOT EXISTS meeting_qualities (meeting_id TEXT PRIMARY KEY, client_id TEXT, data_dump TEXT)",
	)
	if err != nil {
		return nil, err
	}

	return &Persist{db}, nil
}

// Save the Webex analytics data to persitence storage.
// One can only make one call per 5 min for analytics data for single ID.
// This function assumes the successful authorization happened for client_id.
func (p *Persist) SaveAnalyticsData(meetingID, clientID, dataDump string) error {
	// validate that the data dump is non-empty
	if len(dataDump) == 0 {
		return fmt.Errorf("data dump is empty")
	}

	// save data to db, replace if already exists
	_, err := p.db.Exec("REPLACE INTO meeting_qualities (meeting_id, client_id, data_dump) VALUES (?, ?, ?)",
		meetingID, clientID, dataDump)
	return err
}

// RetieveAnalyticsData retrieves the analytics data for a given meeting if present.
// This function assumes the successful authorization happened for client_id.
func (p *Persist) RetriveAnalyticsData(clientID, meetingID string) (*types.MeetingQualities, error) {
	var data types.MeetingQualities
	var dataDump string
	if err := p.db.QueryRow(
		"SELECT data_dump FROM meeting_qualities WHERE client_id = ? AND meeting_id = ?",
		clientID, meetingID,
	).Scan(&dataDump); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}

	if err := json.Unmarshal([]byte(dataDump), &data); err != nil {
		return nil, err
	}

	data.MeetingID = meetingID
	return &data, nil
}
