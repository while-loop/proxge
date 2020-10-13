package ge

import (
	"net/http"
	"testing"
	"time"
)

func Test_osrsGe_PriceById(t *testing.T) {
	tests := []struct {
		name    string
		id      int
		want    int
		wantErr bool
	}{
		{"cannonball", 2, 1000, false},
		{"negative id", -1, 0, true},
	}
	ge := &osrsGe{
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := ge.PriceById(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("PriceById() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr {
				return
			}

			if got > tt.want {
				t.Errorf("PriceById() got > %v, want %v", got, tt.want)
			}
			if got <= 0 {
				t.Errorf("PriceById() got <= %v, want %v", got, 0)
			}
		})
	}
}
