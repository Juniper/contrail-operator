package job_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	batch "k8s.io/api/batch/v1"
	v1 "k8s.io/api/core/v1"

	"github.com/Juniper/contrail-operator/pkg/job"
)

func TestStatus_Pending(t *testing.T) {
	tests := map[string]struct {
		jobStatus       batch.JobStatus
		expectedPending bool
	}{
		"no conditions": {
			jobStatus: batch.JobStatus{
				Conditions: []batch.JobCondition{},
			},
			expectedPending: true,
		},
		"nil conditions": {
			jobStatus: batch.JobStatus{
				Conditions: nil,
			},
			expectedPending: true,
		},
		"condition type Failed, status True": {
			jobStatus: batch.JobStatus{
				Conditions: []batch.JobCondition{{Type: batch.JobFailed, Status: v1.ConditionTrue}},
			},
			expectedPending: false,
		},
		"condition type Failed, status False": {
			jobStatus: batch.JobStatus{
				Conditions: []batch.JobCondition{{Type: batch.JobFailed, Status: v1.ConditionFalse}},
			},
			expectedPending: true,
		},
		"condition type fake, status true": {
			jobStatus: batch.JobStatus{
				Conditions: []batch.JobCondition{{Type: "fake", Status: v1.ConditionTrue}},
			},
			expectedPending: true,
		},
		"condition type fake, status false": {
			jobStatus: batch.JobStatus{
				Conditions: []batch.JobCondition{{Type: "fake", Status: v1.ConditionFalse}},
			},
			expectedPending: true,
		},
		"condition type Complete, status True": {
			jobStatus: batch.JobStatus{
				Conditions: []batch.JobCondition{{Type: batch.JobComplete, Status: v1.ConditionTrue}},
			},
			expectedPending: false,
		},
		"condition type Complete, status False": {
			jobStatus: batch.JobStatus{
				Conditions: []batch.JobCondition{{Type: batch.JobComplete, Status: v1.ConditionFalse}},
			},
			expectedPending: true,
		},
		"multiple conditions": {
			jobStatus: batch.JobStatus{
				Conditions: []batch.JobCondition{
					{Type: batch.JobComplete, Status: v1.ConditionFalse},
					{Type: batch.JobFailed, Status: v1.ConditionTrue},
				},
			},
			expectedPending: false,
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// when
			pending := job.Status(test.jobStatus).JobPending()
			// then
			assert.Equal(t, test.expectedPending, pending)
		})
	}
}
