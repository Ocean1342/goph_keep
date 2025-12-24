package config

import (
	"os"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	tests := []struct {
		name    string
		cfgData []byte
		want    *Config
		wantErr bool
	}{
		{
			name: "positive case",
			cfgData: func() []byte {
				f, err := os.ReadFile("test.config.json")
				if err != nil {
					t.Fatalf("couldn't read test.config.json: %v", err)
				}
				return f
			}(),
			want: &Config{
				RegisterURL:     "register_url",
				AuthURL:         "auth_url",
				SendLogoPassURL: "send_logo_pass_url",
			},
			wantErr: false,
		},
		{
			name:    "empty cfg data case",
			cfgData: nil,
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := New(tt.cfgData)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() got = %v, want %v", got, tt.want)
			}
		})
	}
}
