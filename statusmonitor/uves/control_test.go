package uves

import "testing"

func Test_statusData_Status(t *testing.T) {
	tests := []struct {
		name string
		s    statusData
		want string
	}{
		{
			name: "Single filled status",
			s: [][]interface{}{
				{
					map[string]string{
						"@type": "u32",
						"#text": "5",
					},
				},
			},
			want: "5",
		},
		{
			name: "Single empty status",
			s: [][]interface{}{
				{
					map[string]string{
						"@type": "u32",
						"#text": "",
					},
				},
			},
			want: "0",
		},
		{
			name: "Filled status with second string",
			s: [][]interface{}{
				{
					map[string]string{
						"@type": "u32",
						"#text": "5",
					},
					"ip-10.0.12.13-ec2-instance.internal",
				},
			},
			want: "5",
		},
		{
			name: "Empty status list",
			s:    [][]interface{}{},
			want: "0",
		},
		{
			name: "Too long status list",
			s: [][]interface{}{
				{
					map[string]string{
						"@type": "u32",
						"#text": "5",
					},
				},
				{
					map[string]string{
						"@type": "u32",
						"#text": "6",
					},
				},
			},
			want: "0",
		},
		{
			name: "Missing status field",
			s: [][]interface{}{
				{
					map[string]string{
						"@type": "u32",
					},
				},
			},
			want: "0",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.s.Status()
			if got != tt.want {
				t.Errorf("statusData.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
