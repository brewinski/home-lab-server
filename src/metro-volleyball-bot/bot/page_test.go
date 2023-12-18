package bot

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

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
		t.Run(tt.name, func(t *testing.T) {
			_, err := MonitorPage(tt.args.pageUrl)
			if (err != nil) != tt.wantErr {
				t.Errorf("MonitorPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
		})
	}
}

func TestWeb_MonitorPage(t *testing.T) {
	type fields struct {
		client  *http.Client
		handler http.HandlerFunc
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "TestWeb_MonitorPage will return the page response",
			fields: fields{
				client: &http.Client{},
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("OK"))

					w.WriteHeader(http.StatusOK)
				}),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "TestWeb_MonitorPage will return an error if the page does not exist",
			fields: fields{
				client: &http.Client{},
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(tt.fields.handler)
			defer testServer.Close()

			w := &Web{
				client: tt.fields.client,
			}

			got, err := w.MonitorPage(testServer.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Web.MonitorPage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Web.MonitorPage() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestWeb_CheckPageForChanges(t *testing.T) {
	type fields struct {
		client       *http.Client
		prevResponse string
		handler      http.HandlerFunc
	}

	tests := []struct {
		name                 string
		fields               fields
		wantLastPageResponse string
		want                 bool
		wantErr              bool
	}{
		{
			name: "TestWeb_CheckPageForChanges set the page response to the first response if it's empty",
			fields: fields{
				client: &http.Client{},
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("OK"))

					w.WriteHeader(http.StatusOK)
				}),
			},
			wantLastPageResponse: "OK",
			want:                 false,
			wantErr:              false,
		},
		{
			name: "TestWeb_CheckPageForChanges will return false if the page response is the same as the last response",
			fields: fields{
				client: &http.Client{},
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("OK"))

					w.WriteHeader(http.StatusOK)
				}),
				prevResponse: "OK",
			},
			wantLastPageResponse: "OK",
			want:                 false,
			wantErr:              false,
		},
		{
			name: "TestWeb_CheckPageForChanges will return true if the page response is different than the last response",
			fields: fields{
				client: &http.Client{},
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.Write([]byte("OK"))

					w.WriteHeader(http.StatusOK)
				}),
				prevResponse: "NOT OK",
			},
			wantLastPageResponse: "OK",
			want:                 true,
			wantErr:              false,
		},
		{
			name: "TestWeb_CheckPageForChanges will return an error if the page response is not a valid url",
			fields: fields{
				client: &http.Client{},
				handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
					w.WriteHeader(http.StatusNotFound)
				}),
				prevResponse: "NOT OK",
			},
			wantLastPageResponse: "NOT OK",
			want:                 false,
			wantErr:              true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			testServer := httptest.NewServer(tt.fields.handler)
			defer testServer.Close()

			w := &Web{
				client:       tt.fields.client,
				prevResponse: tt.fields.prevResponse,
			}
			got, err := w.CheckPageForChanges(testServer.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Web.CheckPageForChanges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Web.CheckPageForChanges() = %v, want %v", got, tt.want)
			}
			if w.prevResponse != tt.wantLastPageResponse {
				t.Errorf("Web.CheckPageForChanges() = %v, want %v", w.prevResponse, tt.wantLastPageResponse)
			}
		})
	}
}
