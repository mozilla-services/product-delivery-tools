package bucketlister

import (
	"fmt"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/s3"
)

// BucketLister is a directory listing service for S3
type BucketLister struct {
	Bucket    string
	mountedAt string
	prefix    string

	AWSConfig *aws.Config
}

// New returns a *BucketLister
//
// prefix is the starting point for this lister
func New(bucket, prefix string, awsConfig *aws.Config) *BucketLister {
	trimmedPrefix := strings.Trim(prefix, "/") + "/"
	return &BucketLister{
		AWSConfig: awsConfig,
		Bucket:    bucket,
		mountedAt: trimmedPrefix,
		prefix:    trimmedPrefix,
	}
}

// ServeHTTP implements http.Handler
func (b *BucketLister) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqPath := strings.Trim(strings.TrimPrefix(req.URL.Path, "/"+b.mountedAt), "/")
	prefix := path.Join(b.prefix, reqPath) + "/"
	s3Service := s3.New(b.AWSConfig)

	prefixes := []*s3.CommonPrefix{}
	objects := []*s3.Object{}

	listParams := &s3.ListObjectsInput{
		Bucket:    aws.String(b.Bucket),
		Delimiter: aws.String("/"),
		Prefix:    aws.String(prefix),
	}

	for {
		res, err := s3Service.ListObjects(listParams)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Internal Server Error."))
			log.Printf("Error listing %s/%s err: %s", b.Bucket, prefix, err)
			return
		}
		prefixes = append(prefixes, res.CommonPrefixes...)
		objects = append(objects, res.Contents...)
		if res.NextMarker != nil {
			continue
		}
		break
	}

	if len(objects) == 0 && len(prefixes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("Prefixes<br>"))
	for _, p := range prefixes {
		w.Write([]byte(fmt.Sprintf("<a href=\"/%s\">/%s</a>", *p.Prefix, *p.Prefix)))
		w.Write([]byte("<br>"))
	}

	w.Write([]byte("Objects<br>"))
	for _, p := range objects {
		w.Write([]byte(*p.Key))
		w.Write([]byte("<br>"))
	}
}
