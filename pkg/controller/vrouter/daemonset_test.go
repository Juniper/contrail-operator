package vrouter_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	v1 "k8s.io/api/core/v1"

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
	pathType := v1.HostPathType("")
	assert.Contains(t, ds.Spec.Template.Spec.Volumes, v1.Volume{Name: "cni-bin",
		VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: testBinariesPath, Type: &pathType}}})
	assert.Contains(t, ds.Spec.Template.Spec.Volumes, v1.Volume{Name: "cni-config-files",
		VolumeSource: v1.VolumeSource{HostPath: &v1.HostPathVolumeSource{Path: testConfigPath, Type: &pathType}}})
}
