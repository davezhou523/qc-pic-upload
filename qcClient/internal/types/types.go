// Code generated by goctl. DO NOT EDIT.
package types

type OrderRequest struct {
	Sno   string `json:"sno"`            //流水号
	Model string `json:"model,optional"` //型号
	Num   int64  `json:"num,optional"`   //
}

type OrderResponse struct {
	Code int64  `json:"code"`
	Msg  string `json:"msg"`
}

type ResponseList struct {
	BarCode string `json:"bar_code"`
}

type OrderListReq struct {
	Sno      string `json:"sno,optional"`
	PageSize int64  `json:"pageSize,default=100"`
	Category int64  `json:"category,default=1"` //1:未打印流号号，2当日已打印的流水号
}

type OrderListResp struct {
	List []ResponseList `json:"list"`
}
