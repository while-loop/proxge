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
		{"bgs unhumanize price", 11804, 30000000, false},
		{"sara brew", 6685, 20000, false},
		{"harmonised orb", 24511, 5000000000, false},
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
				t.Errorf("PriceById() got <= 0 %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_unhumanizeNumber(t *testing.T) {
	tests := []struct {
		num  string
		want int
	}{
		{"10.235k", 10235},
		{"10k", 10000},
		{"10,000k", 10000000},
		{"10000k", 10000000},
		{"9.2m", 9200000},
		{"9.2M", 9200000},
		{"54.235m", 54235000},
		{"2.147b", 2147000000},
		{"1b", 1000000000},
		{"10,732", 10732},
		{"10,732,321", 10732321},
	}

	for _, tt := range tests {
		t.Run(tt.num, func(t *testing.T) {
			if got := unhumanizeNumber(tt.num); got != tt.want {
				t.Errorf("unhumanizeNumber() = %v, want %v", got, tt.want)
			}
		})
	}
}
