package handler

import (
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"net/http"
	"re-partners/internal/dto"
	"re-partners/internal/service"
)

type OrderHandler struct {
	log         *zap.SugaredLogger
	packService service.PackService
}

func NewOrderHandler(log *zap.SugaredLogger, packService service.PackService) *OrderHandler {
	return &OrderHandler{
		log:         log,
		packService: packService,
	}
}

func (h *OrderHandler) Calculate(c echo.Context) error {
	calculateDto := dto.Calculate{}
	if err := c.Bind(&calculateDto); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	sizes, err := h.packService.GetPackSizes(c.Request().Context())
	if err != nil {
		h.log.Error(err)
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	ascSizes := make([]int, len(sizes))
	for i, s := range sizes {
		ascSizes[i] = s.Size
	}

	minTotalItems := minimumTotalItems(calculateDto.ItemsOrdered, ascSizes)
	minPacks := minimumPackBreakdown(minTotalItems, ascSizes)

	return c.JSON(http.StatusOK, dto.CalculateResult{
		ItemsOrdered: calculateDto.ItemsOrdered,
		TotalItems:   minTotalItems,
		Packs:        minPacks,
	})
}

func minimumTotalItems(quantity int, sizesAsc []int) int {
	total := 0
	remaining := quantity
	// Жадно берём максимум каждого пака (от большего к меньшему)
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
	// Остаток покрываем наименьшим подходящим паком
	for _, size := range sizesAsc {
		if size >= remaining {
			return total + size
		}
	}
	// Крайний случай: все паки меньше остатка (только один размер)
	smallest := sizesAsc[0]
	return total + ((remaining+smallest-1)/smallest)*smallest
}

// Pass2: минимум паков для ровно totalItems предметов (coin change DP)
func minimumPackBreakdown(totalItems int, sizesAsc []int) []dto.Pack {
	const inf = int(^uint(0) >> 1)
	dp := make([]int, totalItems+1)
	choice := make([]int, totalItems+1) // какой пак взяли на шаге i
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
	// Восстанавливаем какие паки взяли
	counts := map[int]int{}
	for cur := totalItems; cur > 0; {
		size := choice[cur]
		counts[size]++
		cur -= size
	}
	// Собираем результат от большего к меньшему
	var packs []dto.Pack
	for i := len(sizesAsc) - 1; i >= 0; i-- {
		size := sizesAsc[i]
		if qty := counts[size]; qty > 0 {
			packs = append(packs, dto.Pack{Size: size, Quantity: qty})
		}
	}
	return packs
}
