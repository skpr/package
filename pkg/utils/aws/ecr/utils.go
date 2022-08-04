package ecr

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

const (
	// RegionUsEast1 identifies the US region.
	RegionUsEast1 = "us-east-1"
	// RegionApSouthEast2 identifies the Australia region.
	RegionApSouthEast2 = "ap-southeast-2"
)

// Helper function to convert a base64 token to a string.
func decodeAuthorizationToken(auth string) (string, error) {
	decoded, err := base64.StdEncoding.DecodeString(auth)
	if err != nil {
		return "", err
	}

	parts := strings.SplitN(string(decoded), ":", 2)
	if len(parts) != 2 {
		return "", errors.New("auth data contains invalid payload")
	}

	return parts[1], nil
}

// Helper function to derive a region for a URL.
func extractRegionFromURL(url string) (string, error) {
	regions := []string{
		RegionUsEast1,
		RegionApSouthEast2,
	}

	for _, region := range regions {
		if strings.Contains(url, region) {
			return region, nil
		}
	}

	return "", fmt.Errorf("region not found for url: %s", url)
}
