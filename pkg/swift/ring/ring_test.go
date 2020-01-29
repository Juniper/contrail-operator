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
	t.Run("should return error when claim name is empty", func(t *testing.T) {
		_, err := ring.New("", "/etc/swift", "account")
		assert.Error(t, err)
	})
	t.Run("should return error when path is empty", func(t *testing.T) {
		_, err := ring.New("rings", "", "account")
		assert.Error(t, err)
	})
	t.Run("should return error when ring type is empty", func(t *testing.T) {
		_, err := ring.New("rings", "/etc/swift", "")
		assert.Error(t, err)
	})
	t.Run("should create a ring", func(t *testing.T) {
		accountRing, err := ring.New("rings", "/etc/swift", "account")
		require.NoError(t, err)
		assert.NotNil(t, accountRing)
	})
}

func TestRing_BuildJob(t *testing.T) {
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
		account, _ := ring.New("rings", "/etc/swift", "account")
		_, err := account.BuildJob(jobName)
		assert.Error(t, err)
	})
	t.Run("should create a job with given name", func(t *testing.T) {
		account, _ := ring.New("rings", "/etc/swift", "account")
		_ = account.AddDevice(device)
		// when
		job, err := account.BuildJob(jobName)
		require.NoError(t, err)
		assert.Equal(t, jobName.Name, job.Name)
		assert.Equal(t, jobName.Namespace, job.Namespace)
	})
	t.Run("should use default namespace when job namespace not given", func(t *testing.T) {
		account, _ := ring.New("rings", "/etc/swift", "account")
		_ = account.AddDevice(device)
		// when
		job, err := account.BuildJob(types.NamespacedName{
			Namespace: "",
			Name:      "ring-account-job",
		})
		require.NoError(t, err)
		assert.Equal(t, "default", job.Namespace)
	})
	t.Run("should return error when job name not given", func(t *testing.T) {
		account, _ := ring.New("rings", "/etc/swift", "account")
		_ = account.AddDevice(device)
		// when
		_, err := account.BuildJob(types.NamespacedName{
			Namespace: "contrail",
			Name:      "",
		})
		assert.Error(t, err)
	})
	t.Run("should add volume with given claimName to pod template spec", func(t *testing.T) {
		tests := map[string]string{
			"claimName=rings":   "rings",
			"claimName=another": "another",
		}
		for name, claimName := range tests {
			t.Run(name, func(t *testing.T) {
				account, _ := ring.New(claimName, "/etc/swift", "account")
				_ = account.AddDevice(device)
				// when
				job, _ := account.BuildJob(jobName)
				// then
				volumes := job.Spec.Template.Spec.Volumes
				require.Len(t, volumes, 1)
				expectedVolume := core.Volume{
					Name: "rings",
					VolumeSource: core.VolumeSource{
						PersistentVolumeClaim: &core.PersistentVolumeClaimVolumeSource{
							ClaimName: claimName,
							ReadOnly:  false,
						},
					},
				}
				assert.Equal(t, expectedVolume, volumes[0])
			})
		}
	})
	t.Run("should mount volume to container with specified path", func(t *testing.T) {
		tests := map[string]string{
			"path=/etc/swift": "/etc/swift",
			"path=/etc/other": "/etc/other",
		}
		for name, path := range tests {
			t.Run(name, func(t *testing.T) {
				account, _ := ring.New("rings", path, "account")
				_ = account.AddDevice(device)
				// when
				job, _ := account.BuildJob(jobName)
				// then
				containers := job.Spec.Template.Spec.Containers
				require.Len(t, containers, 1)
				volumeMounts := containers[0].VolumeMounts
				require.Len(t, volumeMounts, 1)
				expectedVolumeMount := core.VolumeMount{
					Name:      "rings",
					MountPath: path,
				}
				assert.Equal(t, expectedVolumeMount, volumeMounts[0])
			})
		}
	})
	t.Run("should pass a ring file name as first argument", func(t *testing.T) {
		tests := map[string]struct {
			path                 string
			ringType             string
			expectedRingFileName string
		}{
			"1": {
				path:                 "/etc/swift",
				ringType:             "account",
				expectedRingFileName: "/etc/swift/account",
			},
			"2": {
				path:                 "/etc/other",
				ringType:             "another",
				expectedRingFileName: "/etc/other/another",
			},
			"path /": {
				path:                 "/",
				ringType:             "account",
				expectedRingFileName: "/account",
			},
			"path ends with /": {
				path:                 "/etc/swift/",
				ringType:             "account",
				expectedRingFileName: "/etc/swift/account",
			},
		}
		for name, test := range tests {
			t.Run(name, func(t *testing.T) {
				account, _ := ring.New("rings", test.path, test.ringType)
				_ = account.AddDevice(device)
				// when
				job, _ := account.BuildJob(jobName)
				// then
				containers := job.Spec.Template.Spec.Containers
				require.Len(t, containers, 1)
				args := containers[0].Args
				require.NotEmpty(t, args)
				assert.Equal(t, test.expectedRingFileName, args[0])
			})
		}
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
				account, _ := ring.New("rings", "/etc/swift", "account")
				for _, device := range test.devices {
					_ = account.AddDevice(device)
				}
				// when
				job, err := account.BuildJob(jobName)
				// then
				require.NoError(t, err)
				containers := job.Spec.Template.Spec.Containers
				require.Len(t, containers, 1)
				args := containers[0].Args
				assert.Equal(t, test.expectedDevices, args[1:])
			})
		}
	})
	t.Run("should specify container's required properties", func(t *testing.T) {
		account, _ := ring.New("rings", "/etc/swift", "account")
		_ = account.AddDevice(device)
		// when
		job, _ := account.BuildJob(jobName)
		// then
		containers := job.Spec.Template.Spec.Containers
		require.NotEmpty(t, containers)
		container := containers[0]

		t.Run("command", func(t *testing.T) {
			assert.NotEmpty(t, container.Command)
		})

		t.Run("name", func(t *testing.T) {
			assert.NotEmpty(t, container.Name)
		})
	})
	t.Run("should specify container's image from local registry", func(t *testing.T) {
		account, _ := ring.New("rings", "/etc/swift", "account")
		_ = account.AddDevice(device)
		// when
		job, _ := account.BuildJob(jobName)
		// then
		containers := job.Spec.Template.Spec.Containers
		require.NotEmpty(t, containers)
		image := containers[0].Image
		assert.True(t, strings.HasPrefix(image, "registry:5000/"))
	})

	t.Run("should specify restartPolicy (default Always is not supported in jobs)", func(t *testing.T) {
		account, _ := ring.New("rings", "/etc/swift", "account")
		_ = account.AddDevice(device)
		// when
		job, _ := account.BuildJob(jobName)
		// then
		assert.Equal(t, core.RestartPolicyNever, job.Spec.Template.Spec.RestartPolicy)
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
				theRing, _ := ring.New("a", "/etc/swift", "account")
				// when
				err := theRing.AddDevice(device)
				// then
				require.Error(t, err)
				// and when
				_, err = theRing.BuildJob(types.NamespacedName{
					Namespace: "t",
					Name:      "t",
				})
				// then
				require.Error(t, err)
			})
		}
	})
}
