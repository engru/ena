// MIT License

// Copyright (c) 2018 soren yang

// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:

// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.

// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

package cronjob

import (
	"time"

	"github.com/robfig/cron"

	"github.com/lsytj0413/ena/logger"
)

var (
	cronjob = cron.New()
)

// Job is the cronjob invoke object interface
type Job interface {
	Name() string
	Spec() string
	Run()
}

// TemplateJob for Job interface implement
type TemplateJob struct {
	name string
	spec string
	fn   func()
}

// Name return the job name
func (j *TemplateJob) Name() string {
	return j.name
}

// Spec return the job spec
func (j *TemplateJob) Spec() string {
	return j.spec
}

// Run the job with fn
func (j *TemplateJob) Run() {
	if j.fn != nil {
		j.fn()
	}
}

// NewJob constrct Job object
func NewJob(name string, spec string, fn func()) (Job, error) {
	if err := ValidateSpec(spec); err != nil {
		return nil, err
	}

	return &TemplateJob{
		name: name,
		spec: spec,
		fn:   fn,
	}, nil
}

// ValidateSpec validate the spec string
func ValidateSpec(spec string) error {
	_, err := cron.Parse(spec)
	return err
}

// RegisterFunc a cronjob from func
func RegisterFunc(name string, spec string, fn func()) error {
	job, err := NewJob(name, spec, fn)
	if err != nil {
		return err
	}

	return Register(job)
}

// Register a cronjob
func Register(job Job) error {
	return cronjob.AddFunc(job.Spec(), func() {
		start := time.Now()

		logger.Infof("Job[%v] Start", job.Name())
		defer func() {
			logger.Infof("Job[%v] Done, Cost[%v]", job.Name(), time.Now().Sub(start))
		}()

		job.Run()
	})
}

// Start will start the cronjob scheduler
func Start() {
	cronjob.Start()
}

// Stop will stop the cronjob scheduler
func Stop() {
	cronjob.Stop()
}
