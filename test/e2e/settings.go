package e2e

import "time"

var (
	retryInterval        = time.Second * 5
	waitTimeout          = time.Second * 180
	cleanupRetryInterval = time.Second * 1
	cleanupTimeout       = time.Second * 5
)
