package models

type ExpressionRequest struct {
	Expression string `json:"expression"`
}

type ExpressionResponse struct {
	ExpressionID string `json:"id"`
}

type Expression struct {
	ID     uint32  `json:"id"`
	Status string  `json:"status"`
	Result float64 `json:"result"`
}
