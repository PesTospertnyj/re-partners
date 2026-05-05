package utils

import (
	"math"
	"re-partners/internal/dto"
)

const unreachablePackCount = math.MaxInt

func MinimumTotalItems(quantity int, sizesAsc []int) int {
	total := 0
	remaining := quantity
	// Greedily take as many items as possible from each pack size.
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
	// Edge case: all packs are smaller than the remainder (single pack size).
	smallest := sizesAsc[0]
	return total + ((remaining+smallest-1)/smallest)*smallest
}

// MinimumPackBreakdown finds the minimum number of packs for exactly totalItems.
func MinimumPackBreakdown(totalItems int, sizesAsc []int) []dto.Pack {
	minPackCountByAmount := initializedMinPackCounts(totalItems)
	lastPackSizeByAmount := make([]int, totalItems+1)

	for amount := 1; amount <= totalItems; amount++ {
		for _, packSize := range sizesAsc {
			if packSize > amount {
				break
			}

			previousAmount := amount - packSize
			if minPackCountByAmount[previousAmount] == unreachablePackCount {
				continue
			}

			candidatePackCount := minPackCountByAmount[previousAmount] + 1
			if candidatePackCount < minPackCountByAmount[amount] {
				minPackCountByAmount[amount] = candidatePackCount
				lastPackSizeByAmount[amount] = packSize
			}
		}
	}

	if minPackCountByAmount[totalItems] == unreachablePackCount {
		return nil
	}

	packCounts := countPacksBySize(totalItems, lastPackSizeByAmount)

	return convertToDto(packCounts, sizesAsc)
}

func initializedMinPackCounts(totalItems int) []int {
	minPackCountByAmount := make([]int, totalItems+1)
	for amount := 1; amount <= totalItems; amount++ {
		minPackCountByAmount[amount] = unreachablePackCount
	}

	return minPackCountByAmount
}

func countPacksBySize(totalItems int, lastPackSizeByAmount []int) map[int]int {
	packCounts := make(map[int]int)
	for amount := totalItems; amount > 0; {
		packSize := lastPackSizeByAmount[amount]
		packCounts[packSize]++
		amount -= packSize
	}

	return packCounts
}

func convertToDto(packCounts map[int]int, sizesAsc []int) []dto.Pack {
	packs := make([]dto.Pack, 0, len(packCounts))

	for i := len(sizesAsc) - 1; i >= 0; i-- {
		packSize := sizesAsc[i]
		quantity := packCounts[packSize]
		if quantity == 0 {
			continue
		}

		packs = append(packs, dto.Pack{
			Size:     packSize,
			Quantity: quantity,
		})
	}

	return packs
}
