package postupload

// BuildID represents a release build id
type BuildID string

// Validate ensures BuildID is long enough
func (b BuildID) Validate() bool {
	return len(b) >= 14
}

// Year returns BuildID's year
func (b BuildID) Year() string {
	return string(b[0:4])
}

// Month returns buildID's month
func (b BuildID) Month() string {
	return string(b[4:6])
}

// Day Returns buildID's Day
func (b BuildID) Day() string {
	return string(b[6:8])
}

// Hour Returns BuildID's Hour
func (b BuildID) Hour() string {
	return string(b[8:10])
}

// Minute returns BuildID's Minute
func (b BuildID) Minute() string {
	return string(b[10:12])
}

// Second returns BuildID's Second
func (b BuildID) Second() string {
	return string(b[12:14])
}
