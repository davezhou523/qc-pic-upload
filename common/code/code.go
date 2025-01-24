package code

type ResponseList struct {
	BarCode string `json:"bar_code"`
}
type HttpResponse struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Success bool   `json:"success"`
}
