package model

// CustomResponse struct untuk format respons kustom
type CustomResponse struct {
	StatusCode int    `json:"status_code"`
	Message    string `json:"message"`
}

// CustomDataResponse struct untuk format respons kustom dengan data
type CustomDataResponse struct {
	StatusCode int         `json:"status_code"`
	Message    string      `json:"message"`
	Data       interface{} `json:"data"`
}
