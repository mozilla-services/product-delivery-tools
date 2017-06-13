package deliverytools

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

// AWSSession is the global *session.Session for all tools
var AWSSession = session.Must(session.NewSession(
	&aws.Config{
		MaxRetries: aws.Int(5),
		Region:     aws.String("us-east-1"),
	}))
