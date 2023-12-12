package bot

import "testing"

func TestMonitorPage(t *testing.T) {
	type args struct {
		pageUrl string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "TestMonitorPage will return the page response",
			args:    args{pageUrl: "https://www.google.com"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := MonitorPage(tt.args.pageUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("MonitorPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
