package model

type Kind struct {
	Kind string `json:"kind" yaml:"kind"`
}

const (
	KindTypeBuild = "build"
	KindTypeJob   = "job"
)

type KindBuild struct {
	Kind string    `json:"kind" yaml:"kind"`
	Spec SpecBuild `json:"spec" yaml:"spec"`
}

type SpecBuild struct {
	JobName    string            `json:"jobName" yaml:"job_name"`
	Parameters map[string]string `json:"parameters" yaml:"parameters"`
}

type KindJob struct {
	Kind string  `json:"kind" yaml:"kind"`
	Spec SpecJob `json:"spec" yaml:"spec"`
}

type SpecJob struct {
	JobName string `json:"jobName" yaml:"job_name"`
	XMLData string `json:"xmlData" yaml:"xml_data"`
}
