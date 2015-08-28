package deliverytools

import "github.com/aws/aws-sdk-go/aws"

// AWSConfig is the global *aws.Config for all tools
var AWSConfig = &aws.Config{
	MaxRetries: aws.Int(5),
	Region:     aws.String("us-east-1"),
}
