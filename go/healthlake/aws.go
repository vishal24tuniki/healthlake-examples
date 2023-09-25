package healthlake

import (
	"context"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
)

var awsConfigOnce sync.Once
var awsCfg aws.Config
var errFetchConfig error

func GetAWSConfig() (aws.Config, error) {
	awsConfigOnce.Do(func() {
		awsCfg, errFetchConfig = awsConfig.LoadDefaultConfig(context.TODO(),
			awsConfig.WithRegion(awsRegion),
		)
	})

	return awsCfg, errFetchConfig
}
