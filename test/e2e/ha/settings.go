package ha

import "time"

var (
	RetryInterval          = time.Second * 5
	WaitTimeout            = time.Second * 240
	CleanupRetryInterval   = time.Second * 1
	CleanupTimeout         = time.Second * 5
	WaitForOperatorTimeout = time.Minute * 10
)
