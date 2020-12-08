package job

import (
	batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"
)

type Job batch.Job

func (j Job) JobPending() bool {
	if len(j.Status.Conditions) == 0 {
		return true
	}
	for _, condition := range j.Status.Conditions {
		if condition.Status == v1.ConditionTrue && (condition.Type == batch.JobComplete || condition.Type == batch.JobFailed) {
			return false
		}
	}
	return true
}

func (j Job) JobCompleted() bool {
	if len(j.Status.Conditions) == 0 {
		return false
	}
	for _, condition := range j.Status.Conditions {
		if condition.Status == v1.ConditionTrue && condition.Type == batch.JobComplete {
			return true
		}
	}
	return false
}

func (j Job) JobFailed() bool {
	for _, condition := range j.Status.Conditions {
		if condition.Status == v1.ConditionTrue && condition.Type == batch.JobFailed {
			return true
		}
	}
	// This is a workaround for k8s bug that Status of Job may be updated with a delay
	// and number of pods ran can exceed BackOffLimit specified in Job.Spec
	if j.Spec.BackoffLimit != nil && *j.Spec.BackoffLimit <= j.Status.Failed {
		return true
	}
	return false
}
