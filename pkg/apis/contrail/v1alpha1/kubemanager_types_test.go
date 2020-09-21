package v1alpha1

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestKubemanagerJoinServerList(t *testing.T) {
	var testPort = 42
	tests := []struct {
		servers []string
		port    *int
		sep     string
		want    string
	}{
		{servers: []string{"1.1.1.1"}, port: nil, sep: ",", want: "1.1.1.1"},
		{servers: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}, port: nil, sep: ",", want: "1.1.1.1,2.2.2.2,3.3.3.3"},
		{servers: []string{"1.1.1.1"}, port: &testPort, sep: ",", want: "1.1.1.1:42"},
		{servers: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}, port: &testPort, sep: ",", want: "1.1.1.1:42,2.2.2.2:42,3.3.3.3:42"},
		{servers: []string{"1.1.1.1"}, port: nil, sep: " ", want: "1.1.1.1"},
		{servers: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}, port: nil, sep: " ", want: "1.1.1.1 2.2.2.2 3.3.3.3"},
		{servers: []string{"1.1.1.1"}, port: &testPort, sep: " ", want: "1.1.1.1:42"},
		{servers: []string{"1.1.1.1", "2.2.2.2", "3.3.3.3"}, port: &testPort, sep: " ", want: "1.1.1.1:42 2.2.2.2:42 3.3.3.3:42"},
	}

	for _, tc := range tests {
		got := joinServerList(tc.servers, tc.port, tc.sep)
		assert.Equal(t, tc.want, got)
	}
}
