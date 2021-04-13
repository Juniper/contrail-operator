package e2e

import "time"

var (
	retryInterval          = time.Second * 15
	waitTimeout            = time.Second * 800
	cleanupRetryInterval   = time.Second * 1
	cleanupTimeout         = time.Second * 5
	waitForOperatorTimeout = time.Minute * 30
)
