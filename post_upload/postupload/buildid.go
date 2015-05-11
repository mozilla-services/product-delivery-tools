package postupload

type buildID string

func (b buildID) Validate() bool {
	return len(b) >= 14
}

func (b buildID) Year() string {
	return string(b[0:4])
}

func (b buildID) Month() string {
	return string(b[4:6])
}
func (b buildID) Day() string {
	return string(b[6:8])
}
func (b buildID) Hour() string {
	return string(b[8:10])
}
func (b buildID) Minute() string {
	return string(b[10:12])
}
func (b buildID) Second() string {
	return string(b[12:14])
}
