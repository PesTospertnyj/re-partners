package dto

type (
	Pack struct {
		Size     int `json:"size"`
		Quantity int `json:"quantity"`
	}

	PackSize struct {
		ID   int `json:"id"`
		Size int `json:"size"`
	}

	AddPackSize struct {
		Size int `json:"size"`
	}
)
