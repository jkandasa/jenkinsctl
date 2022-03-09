package jenkins

import (
	"time"

	"github.com/bndr/gojenkins"
)

// Parameter struct
type Parameter struct {
	Name  string `json:"name" yaml:"name" structs:"name"`
	Value string `json:"value" yaml:"value" structs:"value"`
}

// Artifact struct
type Artifact struct {
	FileName string `json:"fileName" yaml:"file_name" structs:"file name"`
	Path     string `json:"path" yaml:"path" structs:"path"`
}

// BuildResponse struct
type BuildResponse struct {
	QueueID         int64                    `json:"queueId" yaml:"queue_id" structs:"queue id"`
	Number          int64                    `json:"number" yaml:"number" structs:"number"`
	URL             string                   `json:"url" yaml:"url" structs:"url"`
	TriggeredBy     string                   `json:"triggeredBy" yaml:"triggered_by" structs:"triggered by"`
	Parameters      []Parameter              `json:"parameters" yaml:"parameters" structs:"parameters"`
	InjectedEnvVars map[string]string        `json:"injectedEnvVars" yaml:"injected_env_vars" structs:"injected env var"`
	Causes          []map[string]interface{} `json:"causes" yaml:"causes" structs:"causes"`
	Duration        time.Duration            `json:"duration" yaml:"duration" structs:"duration"`
	Console         interface{}              `json:"console" yaml:"console" structs:"console"`
	Result          string                   `json:"result" yaml:"result" structs:"result"`
	IsRunning       bool                     `json:"isRunning" yaml:"is_running" structs:"is running"`
	Revision        string                   `json:"revision" yaml:"revision" structs:"revision"`
	RevisionBranch  string                   `json:"revisionBranch" yaml:"revision_branch" structs:"revision branch"`
	Timestamp       time.Time                `json:"timestamp" yaml:"timestamp" structs:"timestamp"`
	TestResult      *gojenkins.TestResult    `json:"testResult" yaml:"test_result" structs:"test result"`
	Artifacts       []Artifact               `json:"artifacts" yaml:"artifacts" structs:"artifacts"`
}
