package bot

import (
	"testing"
)

func TestMonitorPage(t *testing.T) {
	t.Parallel()
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
		{
			name:    "TestMonitorPage will return an error if the page does not exist",
			args:    args{pageUrl: "https://grumbo.dev/this-page-does-not-exist"},
			wantErr: true,
		},
		{
			name:    "TestMonitorPage will return an error if the page is not a valid url",
			args:    args{pageUrl: "this-is-not-a-valid-url"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := MonitorPage(tt.args.pageUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("MonitorPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}
