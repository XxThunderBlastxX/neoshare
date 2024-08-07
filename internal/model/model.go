package model

import "strconv"

type WebResponse struct {
	Message    string `json:"data,omitempty"`
	Error      string `json:"error,omitempty"`
	StatusCode int    `json:"status_code"`
	Success    bool   `json:"success"`
}

func (w *WebResponse) ConvertToMap() map[string]interface{} {
	return map[string]interface{}{
		"message":       w.Message,
		"error_message": w.Error,
		"status_code":   w.StatusCode,
		"success":       w.Success,
	}
}

func (w *WebResponse) ConvertToStruct(m map[string]interface{}) {
	statusCode, _ := strconv.Atoi(m["status_code"].(string))
	success, _ := strconv.ParseBool(m["success"].(string))

	w.Message = m["message"].(string)
	w.Error = m["error_message"].(string)
	w.StatusCode = statusCode
	w.Success = success
}
