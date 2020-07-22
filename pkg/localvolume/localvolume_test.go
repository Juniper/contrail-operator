package localvolume_test

import (
	"context"
	"testing"

	"github.com/Juniper/contrail-operator/pkg/localvolume"
	"github.com/stretchr/testify/assert"
	core "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestNew(t *testing.T) {
	t.Run("should return localvolumes", func(t *testing.T) {
		cl := fake.NewFakeClient()
		lvs := localvolume.New(cl)
		assert.NotNil(t, lvs)
	})
}

func TestLocalVolumesNew(t *testing.T) {
	cl := fake.NewFakeClient()
	lvs := localvolume.New(cl)
	label := map[string]string{"app": "test"}
	nodeSelector := map[string]string{"node": ""}
	quantity5Gi := resource.MustParse("5Gi")
	path := "/mnt/storage"

	t.Run("should return error when empty path is provided", func(t *testing.T) {
		_, err := lvs.New("volume", label, nodeSelector, "", quantity5Gi)
		assert.Error(t, err)
	})

	t.Run("should return error when empty nodeSelector is provided", func(t *testing.T) {
		_, err := lvs.New("volume", label, nil, path, quantity5Gi)
		assert.Error(t, err)
	})

	t.Run("should return error when empty label is provided", func(t *testing.T) {
		_, err := lvs.New("volume", nil, nodeSelector, path, quantity5Gi)
		assert.Error(t, err)
	})

	t.Run("should return error when empty quantity is provided", func(t *testing.T) {
		_, err := lvs.New("volume", nil, nodeSelector, path, resource.Quantity{})
		assert.Error(t, err)
	})

	t.Run("should create valid correct logical volumes", func(t *testing.T) {
		lvs, err := lvs.New("volume", label, nodeSelector, path, quantity5Gi)
		assert.NoError(t, err)
		assert.NotNil(t, lvs)
	})

}

func TestLocalVolumeEnsureExists(t *testing.T) {
	cl := fake.NewFakeClient()
	lvs := localvolume.New(cl)
	name := "volume"
	label := map[string]string{"app": "test"}
	nodeSelector := map[string]string{"node": ""}
	path := "/mnt/storage"
	quantity5Gi := resource.MustParse("5Gi")

	t.Run("should create valid correct PV when all values are provided", func(t *testing.T) {
		volumeMode := core.PersistentVolumeMode("Filesystem")
		expPvSpec := core.PersistentVolumeSpec{
			Capacity:   core.ResourceList{core.ResourceStorage: quantity5Gi},
			VolumeMode: &volumeMode,
			AccessModes: []core.PersistentVolumeAccessMode{
				"ReadWriteOnce",
			},
			StorageClassName: "local-storage",
			NodeAffinity: &core.VolumeNodeAffinity{
				Required: &core.NodeSelector{
					NodeSelectorTerms: []core.NodeSelectorTerm{{
						MatchExpressions: []core.NodeSelectorRequirement{
							core.NodeSelectorRequirement{Key: "node", Operator: "In", Values: []string{""}},
						},
					}},
				},
			},
			PersistentVolumeSource: core.PersistentVolumeSource{
				Local: &core.LocalVolumeSource{
					Path: path,
				},
			},
		}
		lv, err := lvs.New(name, label, nodeSelector, path, quantity5Gi)
		assert.NoError(t, err)
		err = lv.EnsureExist()
		assert.NoError(t, err)
		pv := &core.PersistentVolume{}
		cl.Get(context.Background(), types.NamespacedName{Name: name}, pv)
		assert.Equal(t, expPvSpec, pv.Spec)
	})
}
