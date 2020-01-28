package e2e

import "time"

var (
	retryInterval        = time.Second * 10
	waitTimeout          = time.Second * 240
	cleanupRetryInterval = time.Second * 2
	cleanupTimeout       = time.Second * 10
)
