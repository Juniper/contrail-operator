package vrouter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/Juniper/contrail-operator/pkg/controller/vrouter"
)

func TestGetDaemonset(t *testing.T) {
	assert.NotPanics(t, func() { _ = vrouter.GetDaemonset() }, "Daemonset got properly")
}
