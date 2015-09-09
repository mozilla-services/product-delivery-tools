package services

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
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
				Size:         aws.Int64(2048),
			},
			&s3.Object{
				Key:          aws.String("/pub/firebird/MozillaFirebird-i686-linux-gtk2+xft.tar.gz"),
				LastModified: &now,
				Size:         aws.Int64(2048),
			},
		},
		[]*s3.CommonPrefix{
			&s3.CommonPrefix{
				Prefix: aws.String("prefix+1"),
			},
		},
		nil,
	)
	bl := NewBucketLister("bucket", "/prefix/", nil)

	assert.Equal(t, bl.basePrefix, "prefix/")

	recorder := httptest.NewRecorder()
	req, err := http.NewRequest("GET", "/pre+fix/", nil)
	assert.NoError(t, err)
	bl.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "text/html", recorder.Header().Get("Content-Type"))
	assert.Contains(t, recorder.Body.String(), "/pre%2Bfix/key1")
	assert.Contains(t, recorder.Body.String(), "/pre&#43;fix/prefix&#43;1/")
	assert.Contains(t, recorder.Body.String(), "2K")
	assert.Contains(t, recorder.Body.String(), "/pre%2Bfix/MozillaFirebird-i686-linux-gtk2%2Bxft.tar.gz")

	recorder = httptest.NewRecorder()
	req, err = http.NewRequest("GET", "/", nil)
	req.Header.Add("Accept", "application/json")
	assert.NoError(t, err)
	bl.ServeHTTP(recorder, req)

	assert.Equal(t, 200, recorder.Code)
	assert.Equal(t, "application/json", recorder.Header().Get("Content-Type"))
	res := new(PrefixListing)
	err = json.Unmarshal(recorder.Body.Bytes(), res)
	assert.NoError(t, err)

	assert.Equal(t, "prefix+1/", res.Prefixes[0])
}
