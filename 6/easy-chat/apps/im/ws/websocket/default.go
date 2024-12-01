package websocket

import (
	"math"
	"time"
)

const (
	// 定义最大空闲时间
	defaultMaxConnectionIdle = time.Duration(math.MaxInt)
)
