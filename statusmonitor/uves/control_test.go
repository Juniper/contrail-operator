package uves

import "testing"

func Test_statusData_Status(t *testing.T) {
	tests := []struct {
		name    string
		s       statusData
		want    string
		wantErr bool
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
			want:    "5",
			wantErr: false,
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
			want:    "0",
			wantErr: false,
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
			want:    "5",
			wantErr: false,
		},
		{
			name:    "Empty status list",
			s:       [][]interface{}{},
			want:    "0",
			wantErr: false,
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
			want:    "",
			wantErr: true,
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
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.s.Status()
			if (err != nil) != tt.wantErr {
				t.Errorf("statusData.Status() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("statusData.Status() = %v, want %v", got, tt.want)
			}
		})
	}
}
