package jobs

import (
	"encoding/json"
	"time"
)

// Listen handles job execution.
type Handler func(id string, j *Job) error

// Listen handles job execution.
type ErrorHandler func(id string, j *Job, err error)

// Job carries information about single job.
type Job struct {
	// Job contains name of job broker (usually PHP class).
	Job string `json:"job"`

	// Payload is string data (usually JSON) passed to Job broker.
	Payload string `json:"payload"`

	// Options contains set of PipelineOptions specific to job execution. Can be empty.
	Options Options `json:"options,omitempty"`
}

// Body packs job payload into binary payload.
func (j *Job) Body() []byte {
	return []byte(j.Payload)
}

// Context pack job context (job, id) into binary payload.
func (j *Job) Context(id string) ([]byte, error) {
	return json.Marshal(
		struct {
			ID  string `json:"id"`
			Job string `json:"job"`
		}{ID: id, Job: j.Job},
	)
}

// Options carry information about how to handle given job.
type Options struct {
	// Pipeline manually specified pipeline.
	Pipeline string `json:"pipeline,omitempty"`

	// Delay defines time duration to delay execution for. Defaults to none.
	Delay int `json:"delay,omitempty"`

	// Maximum job retries. Defaults to none.
	MaxAttempts int `json:"maxAttempts,omitempty"`

	// RetryDelay defines for how long job should be waiting until next retry. Defaults to none.
	RetryDelay int `json:"retryDelay,omitempty"`

	// Reserve defines for how broker should wait until treating job are failed. Defaults to 30 min.
	Timeout int `json:"timeout,omitempty"`
}

// CanRetry must return true if broker is allowed to re-run the job.
func (o *Options) CanRetry(attempts int) bool {
	return o.MaxAttempts > attempts
}

// RetryDuration returns retry delay duration in a form of time.Duration.
func (o *Options) RetryDuration() time.Duration {
	return time.Second * time.Duration(o.RetryDelay)
}

// DelayDuration returns delay duration in a form of time.Duration.
func (o *Options) DelayDuration() time.Duration {
	return time.Second * time.Duration(o.Delay)
}

// DelayDuration returns timeout duration in a form of time.Duration.
func (o *Options) TimeoutDuration() time.Duration {
	if o.Timeout == 0 {
		return 30 * time.Minute
	}

	return time.Second * time.Duration(o.Timeout)
}
