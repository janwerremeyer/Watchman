package aws

import (
	"context"
	"encoding/base64"
	"fmt"
	"io"
	"os"
	"path"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	"github.com/docker/docker/api/types/image"

	"github.com/docker/docker/client"
)

func BaseImageName(full string) string {
	if full == "" {
		return ""
	}

	// Remove digest if present (everything after @)
	if i := strings.Index(full, "@"); i != -1 {
		full = full[:i]
	}

	// Remove tag if present (everything after :)
	if i := strings.LastIndex(full, ":"); i != -1 {
		full = full[:i]
	}

	// Extract the last segment after /
	return path.Base(full)
}

func (c *Client) ListTags(image string) ([]string, error) {
	svc := ecr.NewFromConfig(c.config)

	input := &ecr.ListImagesInput{
		RepositoryName: aws.String(BaseImageName(image)),
	}

	paginator := ecr.NewListImagesPaginator(svc, input)

	tags := make([]string, 0)

	for paginator.HasMorePages() {
		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			return nil, fmt.Errorf("failed to describe images: %w", err)
		}
		for _, img := range page.ImageIds {
			tags = append(tags, *img.ImageTag)
		}
	}
	return tags, nil
}
func (c *Client) PullImage(imageWithTag string) error {
	svc := ecr.NewFromConfig(c.config)
	auth, err := svc.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return fmt.Errorf("failed to get ECR auth token: %w", err)
	}

	if len(auth.AuthorizationData) == 0 {
		return fmt.Errorf("no authorization data received")
	}

	// Decode the base64 token (format: "AWS:password")
	decoded, err := base64.StdEncoding.DecodeString(*auth.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		return fmt.Errorf("failed to decode auth token: %w", err)
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid authorization token format")
	}

	username, password := parts[0], parts[1]
	endpoint := strings.TrimPrefix(*auth.AuthorizationData[0].ProxyEndpoint, "https://")
	repository := BaseImageName(imageWithTag)

	split := strings.Split(imageWithTag, ":")
	tag := split[1]

	imageName := fmt.Sprintf("%s/%s:%s", endpoint, repository, tag)

	cli, err := client.NewClientWithOpts(client.FromEnv)
	if err != nil {
		return fmt.Errorf("failed to create Docker client: %w", err)
	}

	encodedAuth := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "{\"username\":\"%s\",\"password\":\"%s\"}\n", username, password))

	rc, err := cli.ImagePull(context.TODO(), imageName, image.PullOptions{RegistryAuth: encodedAuth})
	if err != nil {
		return fmt.Errorf("docker image pull failed: %w", err)
	}
	defer rc.Close()

	// Important: drain the stream so the pull actually completes.
	// You can also io.Copy(os.Stdout, rc) to see progress JSON.
	if _, err := io.Copy(os.Stdout, rc); err != nil {
		return fmt.Errorf("pulling: %w", err)

	}

	return nil
}

//func PullECRImageWithDockerClient(cfg aws.Config, repository, tag string) error {
//	svc := ecr.NewFromConfig(cfg)
//	auth, err := svc.GetAuthorizationToken(context.TODO(), &ecr.GetAuthorizationTokenInput{})
//	if err != nil {
//		return fmt.Errorf("failed to get ECR auth token: %w", err)
//	}
//
//	if len(auth.AuthorizationData) == 0 {
//		return fmt.Errorf("no authorization data received")
//	}
//
//	// Decode the base64 token (format: "AWS:password")
//	decoded, err := base64.StdEncoding.DecodeString(*auth.AuthorizationData[0].AuthorizationToken)
//	if err != nil {
//		return fmt.Errorf("failed to decode auth token: %w", err)
//	}
//
//	parts := strings.SplitN(string(decoded), ":", 2)
//	if len(parts) != 2 {
//		return fmt.Errorf("invalid authorization token format")
//	}
//
//	username, password := parts[0], parts[1]
//	endpoint := strings.TrimPrefix(*auth.AuthorizationData[0].ProxyEndpoint, "https://")
//	imageName := fmt.Sprintf("%s/%s:%s", endpoint, repository, tag)
//
//	cli, err := client.NewClientWithOpts(client.FromEnv)
//	if err != nil {
//		return fmt.Errorf("failed to create Docker client: %w", err)
//	}
//
//	encodedAuth := base64.StdEncoding.EncodeToString(fmt.Appendf(nil, "{\"username\":\"%s\",\"password\":\"%s\"}\n", username, password))
//
//	rc, err := cli.ImagePull(context.TODO(), imageName, image.PullOptions{RegistryAuth: encodedAuth})
//	if err != nil {
//		return fmt.Errorf("docker image pull failed: %w", err)
//	}
//	defer rc.Close()
//
//	// Important: drain the stream so the pull actually completes.
//	// You can also io.Copy(os.Stdout, rc) to see progress JSON.
//	if _, err := io.Copy(os.Stdout, rc); err != nil {
//		return fmt.Errorf("pulling: %w", err)
//
//	}
//
//	return nil
//}

//func InitAWS(awsConfig watchmanConfig.AWSConfig) (aws.Config, error) {
//	creds := aws.NewCredentialsCache(
//		credentials.NewStaticCredentialsProvider(awsConfig.Credentials.AWSAccessKeyID, awsConfig.Credentials.AWSSecretAccessKey, awsConfig.Credentials.AWSSessionToken),
//	)
//	cfg, err := config.LoadDefaultConfig(context.TODO(),
//		config.WithCredentialsProvider(creds),
//		config.WithRegion(awsConfig.Credentials.AWSRegion),
//	)
//	if err != nil {
//		return aws.Config{}, err
//	}
//
//	assumedRole := stscreds.NewAssumeRoleProvider(
//		sts.NewFromConfig(cfg),
//		awsConfig.Credentials.AWSAssumedRole,
//	)
//
//	cfg.Credentials = aws.NewCredentialsCache(assumedRole)
//
//	slog.Info("AWS Config loaded", "Region", cfg.Region)
//	return cfg, nil
//}
