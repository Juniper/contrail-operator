package ring_test

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"k8s.io/apimachinery/pkg/types"

	core "k8s.io/api/core/v1"

	"github.com/Juniper/contrail-operator/pkg/swift/ring"
)

func TestNew(t *testing.T) {
	t.Run("should return error when ring type is empty", func(t *testing.T) {
		_, err := ring.New(types.NamespacedName{Name: "rings"}, "", "service-account")
		assert.Error(t, err)
	})
	t.Run("should return error when service account name is empty", func(t *testing.T) {
		_, err := ring.New(types.NamespacedName{Name: "rings"}, "account", "")
		assert.Error(t, err)
	})
	t.Run("should create a ring", func(t *testing.T) {
		accountRing, err := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		require.NoError(t, err)
		assert.NotNil(t, accountRing)
	})
}

func TestRing_BuildJob(t *testing.T) {
	nodeSelector := map[string]string{"node-role.kubernetes.io/master": ""}
	jobName := types.NamespacedName{
		Namespace: "contrail",
		Name:      "ring-account-job",
	}
	device := ring.Device{
		Region: "1",
		Zone:   "2",
		IP:     "192.168.0.1",
		Port:   5000,
		Device: "d1",
	}
	t.Run("should return error when no devices have been added", func(t *testing.T) {
		account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		_, err := account.BuildJob(jobName, nodeSelector)
		assert.Error(t, err)
	})
	t.Run("should create a job with given name", func(t *testing.T) {
		account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		_ = account.AddDevice(device)
		// when
		job, err := account.BuildJob(jobName, nodeSelector)
		require.NoError(t, err)
		assert.Equal(t, jobName.Name, job.Name)
		assert.Equal(t, jobName.Namespace, job.Namespace)
	})
	t.Run("should use default namespace when job namespace not given", func(t *testing.T) {
		account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		_ = account.AddDevice(device)
		// when
		job, err := account.BuildJob(types.NamespacedName{
			Namespace: "",
			Name:      "ring-account-job",
		}, nodeSelector)
		require.NoError(t, err)
		assert.Equal(t, "default", job.Namespace)
	})
	t.Run("should return error when job name not given", func(t *testing.T) {
		account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		_ = account.AddDevice(device)
		// when
		_, err := account.BuildJob(types.NamespacedName{
			Namespace: "contrail",
			Name:      "",
		}, nodeSelector)
		assert.Error(t, err)
	})
	t.Run("should return error when config map name not given", func(t *testing.T) {
		_, err := ring.New(types.NamespacedName{Name: ""}, "account", "service-account")
		assert.Error(t, err)
	})
	t.Run("should pass namespace name of config map", func(t *testing.T) {
		tests := map[string]struct {
			configMap            types.NamespacedName
			expectedRingFileName string
		}{
			"1": {
				configMap:            types.NamespacedName{Namespace: "contrail", Name: "swift"},
				expectedRingFileName: "contrail/swift",
			},
			"2": {
				configMap:            types.NamespacedName{Namespace: "", Name: "swift"},
				expectedRingFileName: "default/swift",
			},
			"3": {
				configMap:            types.NamespacedName{Namespace: "contrail", Name: "contrail"},
				expectedRingFileName: "contrail/contrail",
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				account, _ := ring.New(test.configMap, "account", "service-account")
				_ = account.AddDevice(device)
				// when
				job, _ := account.BuildJob(jobName, nodeSelector)
				// then
				containers := job.Spec.Template.Spec.Containers
				require.Len(t, containers, 1)
				args := containers[0].Args
				require.NotEmpty(t, args)
				assert.Equal(t, test.expectedRingFileName, args[0])
			})
		}
	})
	t.Run("should return error when ringType not given", func(t *testing.T) {
		_, err := ring.New(types.NamespacedName{Name: "rings"}, "", "service-account")
		assert.Error(t, err)
	})
	t.Run("should pass ringType", func(t *testing.T) {
		account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		_ = account.AddDevice(device)
		// when
		job, _ := account.BuildJob(jobName, nodeSelector)
		// then
		containers := job.Spec.Template.Spec.Containers
		require.Len(t, containers, 1)
		args := containers[0].Args
		require.NotEmpty(t, args)
		assert.Equal(t, "account", args[1])
	})
	t.Run("should pass formatted devices arguments", func(t *testing.T) {
		tests := map[string]struct {
			devices         []ring.Device
			expectedDevices []string
		}{
			"one device": {
				devices: []ring.Device{
					{
						Region: "1",
						Zone:   "2",
						IP:     "192.168.0.2",
						Port:   6000,
						Device: "d2",
					},
				},
				expectedDevices: []string{
					"r1z2-192.168.0.2:6000/d2",
				},
			},
			"one device, different values": {
				devices: []ring.Device{
					{
						Region: "2",
						Zone:   "3",
						IP:     "192.168.0.3",
						Port:   5000,
						Device: "d3",
					},
				},
				expectedDevices: []string{
					"r2z3-192.168.0.3:5000/d3",
				},
			},
			"two devices": {
				devices: []ring.Device{
					{
						Region: "2",
						Zone:   "3",
						IP:     "192.168.0.3",
						Port:   5000,
						Device: "d3",
					},
					{
						Region: "5",
						Zone:   "6",
						IP:     "192.168.0.8",
						Port:   6000,
						Device: "d9",
					},
				},
				expectedDevices: []string{
					"r2z3-192.168.0.3:5000/d3",
					"r5z6-192.168.0.8:6000/d9",
				},
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
				for _, device := range test.devices {
					_ = account.AddDevice(device)
				}
				// when
				job, err := account.BuildJob(jobName, nodeSelector)
				// then
				require.NoError(t, err)
				containers := job.Spec.Template.Spec.Containers
				require.Len(t, containers, 1)
				args := containers[0].Args
				assert.Equal(t, test.expectedDevices, args[2:])
			})
		}
	})
	t.Run("should specify container's required properties", func(t *testing.T) {
		account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		_ = account.AddDevice(device)
		// when
		job, _ := account.BuildJob(jobName, nodeSelector)
		// then
		containers := job.Spec.Template.Spec.Containers
		require.NotEmpty(t, containers)
		container := containers[0]

		t.Run("command", func(t *testing.T) {
			assert.Empty(t, container.Command)
		})

		t.Run("name", func(t *testing.T) {
			assert.NotEmpty(t, container.Name)
		})
	})
	t.Run("should specify container's image from local registry", func(t *testing.T) {
		account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		_ = account.AddDevice(device)
		// when
		job, _ := account.BuildJob(jobName, nodeSelector)
		// then
		containers := job.Spec.Template.Spec.Containers
		require.NotEmpty(t, containers)
		image := containers[0].Image
		assert.True(t, strings.HasPrefix(image, "localhost:5000/"))
	})

	t.Run("should specify restartPolicy (default Always is not supported in jobs)", func(t *testing.T) {
		account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		_ = account.AddDevice(device)
		// when
		job, _ := account.BuildJob(jobName, nodeSelector)
		// then
		assert.Equal(t, core.RestartPolicyNever, job.Spec.Template.Spec.RestartPolicy)
	})

	t.Run("should specify nodeSelector", func(t *testing.T) {
		account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		_ = account.AddDevice(device)
		// when
		job, _ := account.BuildJob(jobName, nodeSelector)
		// then
		assert.Equal(t, nodeSelector, job.Spec.Template.Spec.NodeSelector)
	})

	t.Run("should pass service account name", func(t *testing.T) {
		account, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
		_ = account.AddDevice(device)
		// when
		job, _ := account.BuildJob(jobName, nodeSelector)
		// then
		assert.Equal(t, "service-account", job.Spec.Template.Spec.ServiceAccountName)
	})
}

func TestRing_AddDevice(t *testing.T) {
	t.Run("should return error when region is empty", func(t *testing.T) {
		tests := map[string]ring.Device{
			"empty region": {
				Region: "",
				Zone:   "1",
				IP:     "192.168.0.1",
				Port:   300,
				Device: "d1",
			},
			"empty zone": {
				Region: "1",
				Zone:   "",
				IP:     "192.168.0.1",
				Port:   300,
				Device: "d1",
			},
			"empty IP": {
				Region: "1",
				Zone:   "1",
				IP:     "",
				Port:   300,
				Device: "d1",
			},
			"empty device": {
				Region: "1",
				Zone:   "1",
				IP:     "192.168.0.1",
				Port:   300,
				Device: "",
			},
		}
		for name, device := range tests {
			t.Run(name, func(t *testing.T) {
				theRing, _ := ring.New(types.NamespacedName{Name: "rings"}, "account", "service-account")
				// when
				err := theRing.AddDevice(device)
				// then
				require.Error(t, err)
				// and when
				_, err = theRing.BuildJob(types.NamespacedName{
					Namespace: "t",
					Name:      "t",
				}, map[string]string{})
				// then
				require.Error(t, err)
			})
		}
	})
}
