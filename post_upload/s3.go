package main

import (
	"fmt"
	"os"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/service/s3"
)

var awsConfig = &aws.Config{
	MaxRetries: 5,
	Region:     "us-east-1",
}

func s3CopyFile(src, bucket, key string) error {
	file, err := os.Open(src)
	if err != nil {
		return fmt.Errorf("s3CopyFile: %s", err)
	}
	defer file.Close()

	s3Service := s3.New(awsConfig)
	putObjectInput := &s3.PutObjectInput{
		Body:   file,
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}
	_, err = s3Service.PutObject(putObjectInput)
	if err != nil {
		return fmt.Errorf("s3CopyFile: %s", err)
	}
	return nil
}
