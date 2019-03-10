package assocentity

import "testing"

func Test_avg(t *testing.T) {
	tests := []struct {
		name string
		xs   []float64
		want float64
	}{
		{"3 numbers", []float64{1, 2, 3}, 2},
		{"3 numbers", []float64{1, 2, 3}, 2},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := avg(tt.xs); got != tt.want {
				t.Errorf("avg() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_round(t *testing.T) {
	tests := []struct {
		name string
		x    float64
		want float64
	}{
		{"5.12", 5.12, 5.12},
		{"5.4999999", 5.4999999, 5.5},
		{"5.6", 5.61111111, 5.61},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := round(tt.x); got != tt.want {
				t.Errorf("round() = %v, want %v", got, tt.want)
			}
		})
	}
}
