package bucketlister

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func listMirror(objects []*s3.Object, prefixes []*s3.CommonPrefix, err error) func(*s3.S3, string, string) ([]*s3.Object, []*s3.CommonPrefix, error) {
	return func(svc *s3.S3, bucket, prefix string) ([]*s3.Object, []*s3.CommonPrefix, error) {
		return objects, prefixes, err
	}
}

func TestBucketPrefix(t *testing.T) {
	now := time.Now()
	listObjects = listMirror(
		[]*s3.Object{
			&s3.Object{
				Key:          aws.String("key1"),
				LastModified: &now,
				Size:         aws.Long(100),
			},
		},
		[]*s3.CommonPrefix{
			&s3.CommonPrefix{
				Prefix: aws.String("prefix1"),
			},
		},
		nil,
	)
	bl := New("bucket", "/prefix/", nil)

	assert.Equal(t, bl.prefix, "prefix/")

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/", nil)
	assert.NoError(t, err)
	bl.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "/key1")
	assert.Contains(t, recorder.Body.String(), "/prefix1")
}
