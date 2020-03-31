package vrouter_test

import (
	"github.com/stretchr/testify/assert"
	"testing"

	"github.com/Juniper/contrail-operator/pkg/apis/contrail/v1alpha1"
	"github.com/Juniper/contrail-operator/pkg/controller/vrouter"
)

func TestGetDaemonset(t *testing.T) {
	testBinariesPath := "/test/cni/bin"
	testConfigPath := "/config/test/cni"
	testCNIDirs := v1alpha1.VrouterCNIDirectories{
		BinariesDirectory:    testBinariesPath,
		ConfigFilesDirectory: testConfigPath,
	}
	ds, err := vrouter.GetDaemonset(testCNIDirs)
	assert.NoError(t, err)
	var binariesHostPath string
	var configHostPath string
	for _, volume := range ds.Spec.Template.Spec.Volumes {
		if volume.Name == "cni-bin" {
			binariesHostPath = volume.HostPath.Path
		}
		if volume.Name == "cni-config-files" {
			configHostPath = volume.HostPath.Path
		}
	}
	assert.Equal(t, binariesHostPath, testBinariesPath)
	assert.Equal(t, configHostPath, testConfigPath)
}
