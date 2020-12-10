package ha

import "time"

var (
	retryInterval          = time.Second * 10
	waitTimeout            = time.Minute * 20
	cleanupRetryInterval   = time.Second * 1
	cleanupTimeout         = time.Second * 5
	waitForOperatorTimeout = time.Minute * 10
)
