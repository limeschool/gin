package ext

import (
	"github.com/google/uuid"
	"github.com/limeschool/gin"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		trace := ctx.GetHeader(TraceID)
		if trace == "" {
			trace = uuid.New().String()
		}
		ctx.TraceID = trace
		ctx.Set(TraceID, trace)
		ctx.Log = newLog(trace)
		ctx.Config = newConfig(ctx.Log)
	}
}
