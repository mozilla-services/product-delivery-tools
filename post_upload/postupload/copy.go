package postupload

// Copier is an interface for copying file
type Copier interface {
	Copy(src, dest string) error
}

// S3Copier copies files from disk to S3
type S3Copier struct {
	Bucket string
}

// Copy copies src to dst on Bucket
func (s *S3Copier) Copy(src, dest string) error {
	return nil
}
