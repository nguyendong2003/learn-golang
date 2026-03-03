package util

import "github.com/gin-gonic/gin"

const RequestIDKey = "request_id"

func GetRequestID(c *gin.Context) string {
	if v, ok := c.Get(RequestIDKey); ok {
		if rid, ok := v.(string); ok {
			return rid
		}
	}
	return ""
}
