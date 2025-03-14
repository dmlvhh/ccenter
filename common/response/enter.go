package response

import (
	"net/http"

	"github.com/zeromicro/go-zero/rest/httpx"
)

type Body struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
	Msg  string      `json:"msg"`
}

// Response http返回
func Response(r *http.Request, w http.ResponseWriter, resp interface{}, err error, statusCode ...int) {
	var errCode int
	if len(statusCode) > 0 {
		errCode = statusCode[0] // 如果传入了状态码，则使用传入的状态码
	} else {
		errCode = 7
	}

	if err == nil {
		if resp == nil {
			resp = map[string]interface{}{}
		}
		// 成功返回
		httpx.WriteJson(w, http.StatusOK, &Body{
			Code: 0,
			Data: resp,
			Msg:  "成功",
		})
		return
	}

	httpx.WriteJson(w, http.StatusOK, &Body{
		Code: errCode,
		Data: map[string]interface{}{},
		Msg:  err.Error(),
	})
}
