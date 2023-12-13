package common

import (
	"encoding/base64"
	"fmt"
	"strings"
)

func GetVirtualNodeID(clusterID, queueFQDN string) (string, error) {
	if clusterID == "" {
		return "", fmt.Errorf("invalid clusterID, clusterID cannot be empty")
	}

	if len(clusterID) > 32 {
		return "", fmt.Errorf("invalid clusterID, clusterID cannot exceed 32 chars")
	}

	if strings.Contains(clusterID, "-root") {
		return "", fmt.Errorf("invalid clusterID, the clusterID cannot contain '-root' chars")
	}

	nodeStr := fmt.Sprintf("%s-%s", clusterID, queueFQDN)
	return base64.StdEncoding.EncodeToString([]byte(nodeStr)), nil
}

func ParseVirtualNodeID(virtualNodeID string) (string, string, error) {
	decodedID, err := base64.StdEncoding.DecodeString(virtualNodeID)
	if err != nil {
		return "", "", err
	}

	parts := strings.Split(string(decodedID), "-root")
	if len(parts) != 2 {
		return "", "", fmt.Errorf("incorrect virtual node ID format %s", virtualNodeID)
	}

	return parts[0], "root" + parts[1], nil
}
