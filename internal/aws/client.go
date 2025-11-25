package aws

import (
	"context"
	"log/slog"
	watchmanConfig "watchman/config"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/credentials/stscreds"
	"github.com/aws/aws-sdk-go-v2/service/sts"
)

type Client struct {
	config aws.Config
}

func NewClient(awsConfig watchmanConfig.AWSConfig) (*Client, error) {
	creds := aws.NewCredentialsCache(
		credentials.NewStaticCredentialsProvider(awsConfig.Credentials.AWSAccessKeyID, awsConfig.Credentials.AWSSecretAccessKey, awsConfig.Credentials.AWSSessionToken),
	)
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(awsConfig.Credentials.AWSRegion),
	)
	if err != nil {
		return nil, err
	}

	assumedRole := stscreds.NewAssumeRoleProvider(
		sts.NewFromConfig(cfg),
		awsConfig.Credentials.AWSAssumedRole,
	)

	cfg.Credentials = aws.NewCredentialsCache(assumedRole)
	slog.Info("AWS Config loaded", "Region", cfg.Region)

	awsClient := Client{
		config: cfg,
	}

	return &awsClient, nil
}
