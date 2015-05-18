package postupload

import (
	"errors"
	"fmt"
	"time"
)

// BuildID represents a release build id
type BuildID struct {
	id   string
	time time.Time
}

// NewBuildID returns a new *BuildID
//
// id must be at least 14 characters long
func NewBuildID(id string) (*BuildID, error) {
	buildID := new(BuildID)
	err := buildID.Set(id)
	return buildID, err
}

// Set sets a *BuildID from a string
func (b *BuildID) Set(id string) error {
	if len(id) < 14 {
		return errors.New("id must be at least 14 characters")
	}

	l, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		return fmt.Errorf("NewBuildID/LoadLocation: %s", err)
	}

	t, err := time.ParseInLocation("20060102150405", id, l)
	if err != nil {
		return fmt.Errorf("NewBuildID/Parse: %s", err)
	}

	b.id = id
	b.time = t
	return nil
}

// String returns a BuildID as a string
func (b *BuildID) String() string {
	return b.id
}

// Year returns BuildID's year
func (b *BuildID) Year() string {
	return b.id[0:4]
}

// Month returns id's month
func (b *BuildID) Month() string {
	return b.id[4:6]
}

// Day Returns id's Day
func (b *BuildID) Day() string {
	return b.id[6:8]
}

// Hour Returns BuildID's Hour
func (b *BuildID) Hour() string {
	return b.id[8:10]
}

// Minute returns BuildID's Minute
func (b *BuildID) Minute() string {
	return b.id[10:12]
}

// Second returns BuildID's Second
func (b *BuildID) Second() string {
	return b.id[12:14]
}

// Time returns BuildID's time.Time
func (b *BuildID) Time() time.Time {
	return b.time
}
