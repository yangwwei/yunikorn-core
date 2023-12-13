package common

import (
	"encoding/base64"
	"gotest.tools/v3/assert"
	"testing"
)

func TestGetVirtualNodeID(t *testing.T) {
	tests := []struct {
		clusterID    string
		queueFQDN    string
		valid        bool
		errorMessage string
	}{
		{"compute-cluster-1", "root.a", true, ""},
		{"012345678901234567890123456789012",
			"root.a", false, "clusterID cannot exceed 32 chars"},
		{"01234567890123456789012345678901", "root.a", true, ""},
		{"compute-root-cluster", "root.a", false,
			"the clusterID cannot contain '-root' chars"},
	}
	for _, test := range tests {
		_, err := GetVirtualNodeID(test.clusterID, test.queueFQDN)
		if test.valid {
			assert.NilError(t, err)
		} else {
			assert.ErrorContains(t, err, test.errorMessage)
		}
	}
}

func TestParseVirtualNodeID(t *testing.T) {
	tests := []struct {
		virtualNodeID     string
		expectedClusterID string
		expectedQueue     string
		valid             bool
		errorMessage      string
	}{
		{base64.StdEncoding.EncodeToString([]byte("compute-cluster-1-root.abc")),
			"compute-cluster-1",
			"root.abc",
			true,
			""},
		{base64.StdEncoding.EncodeToString([]byte("compute-cluster-1024-root.a.b.c.d")),
			"compute-cluster-1024",
			"root.a.b.c.d",
			true,
			""},
		{base64.StdEncoding.EncodeToString([]byte("compute-cluster-1024")),
			"compute-cluster-1024",
			"",
			false,
			"incorrect virtual node ID format"},
	}
	for _, test := range tests {
		clusterID, queue, err := ParseVirtualNodeID(test.virtualNodeID)
		if test.valid {
			assert.Equal(t, clusterID, test.expectedClusterID)
			assert.Equal(t, queue, test.expectedQueue)
		} else {
			assert.ErrorContains(t, err, test.errorMessage)
		}
	}
}
