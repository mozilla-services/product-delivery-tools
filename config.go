package deliverytools

import "github.com/aws/aws-sdk-go/aws"

// AWSConfig is the global *aws.Config for all tools
var AWSConfig = &aws.Config{
	MaxRetries: 5,
	Region:     "us-east-1",
}
