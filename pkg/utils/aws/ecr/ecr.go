package ecr

import (
	"context"
	"fmt"
	"strings"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ecr"
	docker "github.com/fsouza/go-dockerclient"
	"github.com/pkg/errors"
)

// Username to pass to the Docker registry.
// https://docs.aws.amazon.com/cli/latest/reference/ecr/get-authorization-token.html
const Username = "AWS"

// IsRegistry managed by AWS ECR.
func IsRegistry(registry string) bool {
	return strings.Contains(registry, ".ecr.")
}

// UpgradeAuth to use an AWS IAM token for authentication..
// https://docs.aws.amazon.com/cli/latest/reference/ecr/get-login.html
func UpgradeAuth(url string, auth docker.AuthConfiguration) (docker.AuthConfiguration, error) {
	region, err := extractRegionFromURL(url)
	if err != nil {
		return auth, errors.Wrap(err, "failed to determine registry region")
	}
	ctx := context.TODO()
	cfg, err := config.LoadDefaultConfig(
		ctx,
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(auth.Username, auth.Password, "")),
	)
	if err != nil {
		return auth, fmt.Errorf("failed to get session: %w", err)
	}
	ecrClient := ecr.NewFromConfig(cfg)

	res, err := ecrClient.GetAuthorizationToken(ctx, &ecr.GetAuthorizationTokenInput{})
	if err != nil {
		return auth, err
	}

	if len(res.AuthorizationData) == 0 {
		return auth, errors.New("failed get authorization data")
	}
	if res.AuthorizationData[0].AuthorizationToken == nil {
		return auth, errors.New("failed get authorization token")
	}

	password, err := decodeAuthorizationToken(*res.AuthorizationData[0].AuthorizationToken)
	if err != nil {
		return auth, errors.Wrap(err, "failed to decode authorization token")
	}

	auth.Username = Username
	auth.Password = password

	return auth, nil
}
