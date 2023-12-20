package monitor_test

import (
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"

	"github.com/brewinski/home-lab-server/src/metro-volleyball-bot/monitor"
)

func TestDS_MonitorPage(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testServer := httptest.NewServer(tt.fields.handler)
			defer testServer.Close()

			w := monitor.New(monitor.Config{
				Client: tt.fields.client,
			})

			got, err := w.Monitor(testServer.URL)
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

func TestDS_CheckPageForChanges(t *testing.T) {
	t.Parallel()
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
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			testServer := httptest.NewServer(tt.fields.handler)
			defer testServer.Close()

			w := monitor.New(monitor.Config{
				Client:   tt.fields.client,
				InitData: tt.fields.prevResponse,
			})
			got, prev, err := w.CheckForChanges(testServer.URL)
			if (err != nil) != tt.wantErr {
				t.Errorf("Web.CheckPageForChanges() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Web.CheckPageForChanges() = %v, want %v", got, tt.want)
			}
			if prev != tt.wantLastPageResponse {
				t.Errorf("Web.CheckPageForChanges() = %v, want %v", prev, tt.wantLastPageResponse)
			}
		})
	}
}

func TestDS_CheckPageForChangesRace(t *testing.T) {
	t.Parallel()

	testServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(time.Now().String()))
	}))
	defer testServer.Close()

	w := monitor.New(monitor.Config{
		Client: &http.Client{},
	})
	workers := 100
	var wg sync.WaitGroup
	wg.Add(workers)

	for i := 0; i < workers; i++ {
		go func() {
			w.CheckForChanges(testServer.URL)
			wg.Done()
		}()
	}

	wg.Wait()
}
