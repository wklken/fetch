package config_test

// func TestReadFromFile(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		wantErr bool
// 	}{
// 		{
// 			name: "read from file",
// 			args: args{
// 				path: filepath.Join("testdata", "config.yaml"),
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "read from non-existent file",
// 			args: args{
// 				path: filepath.Join("testdata", "non-existent.yaml"),
// 			},
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			_, err := config.ReadFromFile(tt.args.path)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ReadFromFile() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 		})
// 	}
// }

// func TestReadCasesFromFile(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    []*config.Case
// 		wantErr bool
// 	}{
// 		{
// 			name: "read cases from file",
// 			args: args{
// 				path: filepath.Join("testdata", "cases.yaml"),
// 			},
// 			want: []*config.Case{
// 				{
// 					Request: config.Request{
// 						Method: "GET",
// 						URL:    "http://localhost:8080",
// 						Header: map[string]string{
// 							"Content-Type": "application/json",
// 						},
// 						Body: "test",
// 					},
// 				},
// 				{
// 					Request: config.Request{
// 						Method: "POST",
// 						URL:    "http://localhost:8080",
// 						Header: map[string]string{
// 							"Content-Type": "application/json",
// 						},
// 						Body: "test",
// 					},
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "read cases from non-existent file",
// 			args: args{
// 				path: filepath.Join("testdata", "non-existent.yaml"),
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := config.ReadCasesFromFile(tt.args.path)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ReadCasesFromFile() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }

// func TestReadLines(t *testing.T) {
// 	type args struct {
// 		path string
// 	}
// 	tests := []struct {
// 		name    string
// 		args    args
// 		want    map[int]map[int]string
// 		wantErr bool
// 	}{
// 		{
// 			name: "read lines from file",
// 			args: args{
// 				path: filepath.Join("testdata", "config.yaml"),
// 			},
// 			want: map[int]map[int]string{
// 				1: {
// 					1:  "name: test",
// 					2:  "request:",
// 					3:  "  method: GET",
// 					4:  "  url: http://localhost:8080",
// 					5:  "  headers:",
// 					6:  "    content-type: application/json",
// 					7:  "  body: test",
// 					8:  "response:",
// 					9:  "  status: 200",
// 					10: "  headers:",
// 					11: "    content-type: application/json",
// 					12: "  body: test",
// 				},
// 			},
// 			wantErr: false,
// 		},
// 		{
// 			name: "read lines from non-existent file",
// 			args: args{
// 				path: filepath.Join("testdata", "non-existent.yaml"),
// 			},
// 			want:    nil,
// 			wantErr: true,
// 		},
// 	}
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got, err := config.ReadLines(tt.args.path)
// 			if (err != nil) != tt.wantErr {
// 				t.Errorf("ReadLines() error = %v, wantErr %v", err, tt.wantErr)
// 				return
// 			}
// 			assert.Equal(t, tt.want, got)
// 		})
// 	}
// }
