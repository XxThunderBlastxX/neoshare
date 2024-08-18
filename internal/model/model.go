package model

import "strconv"

type WebResponse struct {
	Message    string `json:"data,omitempty"`
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
}

func (w *WebResponse) ConvertToMap() map[string]any {
	return map[string]any{
		"message":     w.Message,
		"status_code": w.StatusCode,
		"success":     w.Success,
	}
}

func (w *WebResponse) ConvertToStruct(m map[string]any) {
	statusCode, _ := strconv.Atoi(m["status_code"].(string))
	success, _ := strconv.ParseBool(m["success"].(string))

	w.Message = m["message"].(string)
	w.StatusCode = statusCode
	w.Success = success
}
