package utils

import (
	"reflect"
	"testing"

	"re-partners/internal/dto"
)

func TestMinimumTotalItems(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name     string
		quantity int
		sizesAsc []int
		want     int
	}{
		{
			name:     "1_with_250_500",
			quantity: 1,
			sizesAsc: []int{250, 500},
			want:     250,
		},
		{
			name:     "250_with_250_500",
			quantity: 250,
			sizesAsc: []int{250, 500},
			want:     250,
		},
		{
			name:     "251_prefers_fewer_packs",
			quantity: 251,
			sizesAsc: []int{250, 500},
			want:     500,
		},
		{
			name:     "501_with_250_500_1000",
			quantity: 501,
			sizesAsc: []int{250, 500, 1000},
			want:     750,
		},
		{
			name:     "12001_with_250_2000_5000",
			quantity: 12001,
			sizesAsc: []int{250, 2000, 5000},
			want:     12250,
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			got := MinimumTotalItems(tc.quantity, tc.sizesAsc)
			if got != tc.want {
				t.Fatalf("MinimumTotalItems(%d, %v) = %d, want %d", tc.quantity, tc.sizesAsc, got, tc.want)
			}
		})
	}
}

func TestMinimumPackBreakdown(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		quantity  int
		sizesAsc  []int
		wantTotal int
		wantPacks map[int]int
	}{
		{
			name:      "1_with_250_500",
			quantity:  1,
			sizesAsc:  []int{250, 500},
			wantTotal: 250,
			wantPacks: map[int]int{250: 1},
		},
		{
			name:      "250_with_250_500",
			quantity:  250,
			sizesAsc:  []int{250, 500},
			wantTotal: 250,
			wantPacks: map[int]int{250: 1},
		},
		{
			name:      "251_prefers_1x500_over_2x250",
			quantity:  251,
			sizesAsc:  []int{250, 500},
			wantTotal: 500,
			wantPacks: map[int]int{500: 1},
		},
		{
			name:      "501_prefers_500_plus_250",
			quantity:  501,
			sizesAsc:  []int{250, 500, 1000},
			wantTotal: 750,
			wantPacks: map[int]int{500: 1, 250: 1},
		},
		{
			name:      "12001_prefers_2x5000_1x2000_1x250",
			quantity:  12001,
			sizesAsc:  []int{250, 2000, 5000},
			wantTotal: 12250,
			wantPacks: map[int]int{5000: 2, 2000: 1, 250: 1},
		},
	}

	for _, tc := range tests {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			total := MinimumTotalItems(tc.quantity, tc.sizesAsc)
			if total != tc.wantTotal {
				t.Fatalf("MinimumTotalItems(%d, %v) = %d, want %d", tc.quantity, tc.sizesAsc, total, tc.wantTotal)
			}

			packs := MinimumPackBreakdown(total, tc.sizesAsc)
			got := packsToMap(packs)
			if !reflect.DeepEqual(got, tc.wantPacks) {
				t.Fatalf("MinimumPackBreakdown(%d, %v) = %v, want %v", total, tc.sizesAsc, got, tc.wantPacks)
			}
		})
	}
}

func TestMinimumPackBreakdownEdgeCaseConfiguration(t *testing.T) {
	t.Parallel()

	sizesAsc := []int{23, 31, 53}
	totalItems := 500000
	want := map[int]int{
		23: 2,
		31: 7,
		53: 9429,
	}

	got := packsToMap(MinimumPackBreakdown(totalItems, sizesAsc))
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("MinimumPackBreakdown(%d, %v) = %v, want %v", totalItems, sizesAsc, got, want)
	}
}

func packsToMap(packs []dto.Pack) map[int]int {
	res := make(map[int]int, len(packs))
	for _, p := range packs {
		res[p.Size] = p.Quantity
	}
	return res
}
