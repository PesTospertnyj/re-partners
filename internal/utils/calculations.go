package utils

import "re-partners/internal/dto"

func MinimumTotalItems(quantity int, sizesAsc []int) int {
	total := 0
	remaining := quantity
	// Greedily take as many items as possible from each pack size (largest to smallest).
	for i := len(sizesAsc) - 1; i >= 0; i-- {
		size := sizesAsc[i]
		if remaining < size {
			continue
		}
		count := remaining / size
		total += count * size
		remaining -= count * size
	}
	if remaining == 0 {
		return total
	}
	// Cover remaining items with the smallest pack that is enough.
	for _, size := range sizesAsc {
		if size >= remaining {
			return total + size
		}
	}
	// Edge case: all packs are smaller than remainder (typically single pack size).
	smallest := sizesAsc[0]
	return total + ((remaining+smallest-1)/smallest)*smallest
}

// MinimumPackBreakdown finds minimum number of packs for exactly totalItems (coin change DP).
func MinimumPackBreakdown(totalItems int, sizesAsc []int) []dto.Pack {
	const inf = int(^uint(0) >> 1)
	dp := make([]int, totalItems+1)
	choice := make([]int, totalItems+1) // pack size chosen at amount i.
	for i := 1; i <= totalItems; i++ {
		dp[i] = inf
	}
	for amount := 1; amount <= totalItems; amount++ {
		for _, size := range sizesAsc {
			if size > amount || dp[amount-size] == inf {
				continue
			}
			if candidate := dp[amount-size] + 1; candidate < dp[amount] {
				dp[amount] = candidate
				choice[amount] = size
			}
		}
	}
	// Reconstruct selected pack counts.
	counts := map[int]int{}
	for cur := totalItems; cur > 0; {
		size := choice[cur]
		counts[size]++
		cur -= size
	}
	// Build result from largest pack size to smallest.
	var packs []dto.Pack
	for i := len(sizesAsc) - 1; i >= 0; i-- {
		size := sizesAsc[i]
		if qty := counts[size]; qty > 0 {
			packs = append(packs, dto.Pack{Size: size, Quantity: qty})
		}
	}
	return packs
}
