package deliverytools

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWSConfig is the global *aws.Config for all tools
var AWSSession = session.Must(session.NewSession(
	&aws.Config{
		MaxRetries: aws.Int(5),
		Region:     aws.String("us-east-1"),
	}))
