package main

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/mozilla-services/product-delivery-tools"
)

func s3Service() *s3.S3 {
	return s3.New(deliverytools.AWSConfig)
}

var s3FileCache = map[string]string{}

var keyExpiresPatterns = []struct {
	Pattern  *regexp.Regexp
	Duration time.Duration
}{
	{
		Pattern:  regexp.MustCompile("^pub/.*/nightly/latest.*"),
		Duration: 1 * time.Hour,
	},
}

func keyCacheControl(key string) *string {
	for _, p := range keyExpiresPatterns {
		if p.Pattern.MatchString(key) {
			return aws.String(fmt.Sprintf("max-age=%d", p.Duration/time.Second))
		}
	}
	return nil
}

func s3CopyObject(src, bucket, key string) error {
	copyInput := &s3.CopyObjectInput{
		Bucket:       aws.String(bucket),
		CacheControl: keyCacheControl(key),
		ContentType:  aws.String(ContentType(key)),
		CopySource:   aws.String(src),
		Key:          aws.String(key),
	}

	// Special case for .txt.gz
	if strings.HasSuffix(key, ".txt.gz") {
		copyInput.ContentType = aws.String("text/plain; charset=UTF-8")
		copyInput.ContentEncoding = aws.String("gzip")
	}

	_, err := s3Service().CopyObject(copyInput)

	if err != nil {
		return fmt.Errorf("copying %s to %s/%s, err: %s", src, bucket, key, err)
	}
	return nil
}

func s3PutFile(src, bucket, key string) error {
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("opening %s: err, %s", src, err)
	}
	defer file.Close()

	putObjectInput := &s3.PutObjectInput{
		Body:         file,
		Bucket:       aws.String(bucket),
		CacheControl: keyCacheControl(key),
		ContentType:  aws.String(ContentType(key)),
		Key:          aws.String(key),
	}

	// Special case for .txt.gz
	if strings.HasSuffix(key, ".txt.gz") {
		putObjectInput.ContentType = aws.String("text/plain; charset=UTF-8")
		putObjectInput.ContentEncoding = aws.String("gzip")
	}
	_, err = s3Service().PutObject(putObjectInput)
	if err != nil {
		return fmt.Errorf("putting %s to %s/%s err: %s", src, bucket, key, err)
	}
	return nil
}

func s3CopyFile(src, bucket, key string) error {
	destKey := "/" + bucket + "/" + key
	if cpSrc, ok := s3FileCache[src]; ok {
		// File has already been copied, so move on.
		if cpSrc == destKey {
			return nil
		}
		return s3CopyObject(cpSrc, bucket, key)
	}

	if err := s3PutFile(src, bucket, key); err != nil {
		return err
	}

	s3FileCache[src] = destKey
	return nil
}
