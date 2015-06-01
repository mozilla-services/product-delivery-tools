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

// Empty returns true if the bucket contains zero keys
func (b *BucketLister) Empty() (bool, error) {
	listParams := &s3.ListObjectsInput{
		Bucket:  aws.String(b.Bucket),
		MaxKeys: aws.Long(1),
	}

	s3Service := s3.New(b.AWSConfig)
	res, err := s3Service.ListObjects(listParams)
	if err != nil {
		return true, fmt.Errorf("listing %s err: %s", b.Bucket, err)
	}

	return len(res.Contents) <= 0, nil
}

func objectToListFileInfo(obj *s3.Object) *listFileInfo {
	size := *obj.Size
	sizeStr := ""
	if size < 1024 {
		sizeStr = fmt.Sprintf("%d B", size)
	} else {
		sizeStr = fmt.Sprintf("%d KB", size/1024)
	}

	return &listFileInfo{
		Key:          *obj.Key,
		LastModified: (*obj.LastModified).Format("02-Jan-2006 15:04"),
		Size:         sizeStr,
	}
}

// ServeHTTP implements http.Handler
func (b *BucketLister) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	reqPath := strings.Trim(strings.TrimPrefix(req.URL.Path, "/"+b.mountedAt), "/")
	prefix := path.Join(b.prefix, reqPath) + "/"

	prefixes := []*s3.CommonPrefix{}
	objects := []*s3.Object{}

	s3Service := s3.New(b.AWSConfig)

	objects, prefixes, err := listObjects(s3Service, b.Bucket, prefix)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal Server Error."))
		log.Printf("Error %s", err)
		return
	}

	if len(objects) == 0 && len(prefixes) == 0 {
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Not Found"))
		return
	}

	tmplParams := &listTemplateInput{
		Path:        prefix,
		Directories: make([]string, len(prefixes)),
		Files:       make([]*listFileInfo, len(objects)),
	}

	for i, p := range prefixes {
		tmplParams.Directories[i] = *p.Prefix
	}

	for i, o := range objects {
		tmplParams.Files[i] = objectToListFileInfo(o)
	}

	w.Header().Set("Content-Type", "text/html")
	err = listTemplate.Execute(w, tmplParams)
	if err != nil {
		log.Printf("Error executing template err: %s", err)
	}
}
