package e2e

import "time"

var (
	retryInterval        = time.Second * 20
	waitTimeout          = time.Second * 480
	cleanupRetryInterval = time.Second * 4
	cleanupTimeout       = time.Second * 20
)
