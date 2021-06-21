package client

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/bndr/gojenkins"
	"github.com/jkandasa/jenkinsctl/pkg/model/config"
	jenkinsML "github.com/jkandasa/jenkinsctl/pkg/model/jenkins"
)

// Client type
type Client struct {
	api *gojenkins.Jenkins
	ctx context.Context
}

// NewClient function to get client instance
func NewClient(cfg *config.Config) *Client {
	httpClient := http.DefaultClient
	httpClient.Transport = http.DefaultTransport
	httpClient.Transport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: cfg.InsecureSkipTLSVerify}

	ctx := context.TODO()
	jenkins := gojenkins.CreateJenkins(httpClient, cfg.URL, cfg.Username, cfg.GetPassword())
	_, err := jenkins.Init(ctx)
	if err != nil {
		log.Fatalf("error on login, %s", err.Error())
	}
	return &Client{api: jenkins, ctx: ctx}
}

func (jc *Client) Version() string {
	return jc.api.Version
}

// ListJobs details
func (jc *Client) ListJobs() ([]gojenkins.InnerJob, error) {
	return jc.api.GetAllJobNames(jc.ctx)
}

// Status of the server
func (jc *Client) Status() (*gojenkins.ExecutorResponse, error) {
	return jc.api.Info(jc.ctx)
}

// GetBuild returns build details of a job
func (jc *Client) GetBuild(jobName string, buildNumber int, withConsole bool) (*jenkinsML.BuildResponse, error) {
	buildRaw, err := jc.api.GetBuild(jc.ctx, jobName, int64(buildNumber))
	if err != nil {
		return nil, err
	}

	build := jenkinsML.BuildResponse{
		URL:    buildRaw.GetUrl(),
		Number: buildRaw.GetBuildNumber(),
	}
	// update triggered by
	if causes, err := buildRaw.GetCauses(jc.ctx); err == nil {
		build.Causes = causes
		for _, c := range causes {
			if c != nil && c["userName"] != nil {
				build.TriggeredBy, _ = c["userName"].(string)
				if userIdRaw, ok := c["userId"]; ok {
					userId, _ := userIdRaw.(string)
					if userId != "" {
						build.TriggeredBy = fmt.Sprintf("%s(%s)", build.TriggeredBy, userId)
					}
				}
			}
		}
	}

	build.Parameters = make([]jenkinsML.Parameter, 0)
	for _, p := range buildRaw.GetParameters() {
		build.Parameters = append(build.Parameters, jenkinsML.Parameter{Name: p.Name, Value: fmt.Sprintf("%v", p.Value)})
	}

	if injectedEnvVars, err := buildRaw.GetInjectedEnvVars(jc.ctx); err != nil {
		build.InjectedEnvVars = injectedEnvVars
	}

	build.Duration = time.Duration((buildRaw.GetDuration() / 1000) * 1000000000)

	build.Result = buildRaw.GetResult()
	build.IsRunning = buildRaw.IsRunning(jc.ctx)
	build.Revision = buildRaw.GetRevision()
	//build.RevisionBranch = bu.GetRevisionBranch()
	build.Timestamp = buildRaw.GetTimestamp()
	if testResult, err := buildRaw.GetResultSet(jc.ctx); err == nil {
		// build.TestResult = testResult // loading testResult keeps lot of data
		build.TestResult = &gojenkins.TestResult{
			Duration:  testResult.Duration,
			Empty:     testResult.Empty,
			FailCount: testResult.FailCount,
			PassCount: testResult.PassCount,
			SkipCount: testResult.SkipCount,
		}
	}

	for _, artifact := range buildRaw.GetArtifacts() {
		build.Artifacts = append(build.Artifacts, jenkinsML.Artifact{FileName: artifact.FileName, Path: artifact.Path})
	}

	if withConsole {
		consoleLog := strings.Split(buildRaw.GetConsoleOutput(jc.ctx), "\n")
		for _, line := range consoleLog {
			if strings.Contains(line, "Login to the console with user") {
				build.Console = line
				break
			}
		}
	}

	return &build, nil
}

// ListBuilds details
func (jc *Client) ListBuilds(jobName string, limit int, withConsole bool) ([]jenkinsML.BuildResponse, error) {
	builds := make([]jenkinsML.BuildResponse, 0)
	buildIds, err := jc.api.GetAllBuildIds(jc.ctx, jobName)
	if err != nil {
		return nil, err
	}
	for _, b := range buildIds {
		build, err := jc.GetBuild(jobName, int(b.Number), withConsole)
		if err != nil {
			return nil, err
		}
		builds = append(builds, *build)
		if len(builds) >= limit {
			return builds, nil
		}
	}
	return builds, nil
}

// GetConsole returns/prints build console log
func (jc *Client) GetConsole(jobName string, buildNumber int, watch bool, out io.Writer) (string, error) {
	build, err := jc.api.GetBuild(jc.ctx, jobName, int64(buildNumber))
	if err != nil {
		return "", err
	}

	if !watch {
		return build.GetConsoleOutput(jc.ctx), nil
	}

	ticker := time.NewTicker(3 * time.Second)
	defer ticker.Stop()

	startID := int64(0)
	for {
		console, err := build.GetConsoleOutputFromIndex(jc.ctx, startID)
		if err != nil {
			return "", err
		}
		if len(console.Content) > 0 {
			fmt.Fprint(out, console.Content)
		}
		startID = console.Offset
		if !console.HasMoreText {
			return "", nil
		}
		<-ticker.C
	}

}

// ListParameters of a job
func (jc *Client) ListParameters(jobName string) ([]gojenkins.ParameterDefinition, error) {
	job, err := jc.api.GetJob(jc.ctx, jobName)
	if err != nil {
		return nil, err
	}
	return job.GetParameters(jc.ctx)
}

// DownloadArtifacts of a build
func (jc *Client) DownloadArtifacts(jobName string, buildNumber int, toDirectory string) (string, error) {
	directoryFinal := filepath.Join(toDirectory, jobName, strconv.Itoa(buildNumber))
	build, err := jc.api.GetBuild(jc.ctx, jobName, int64(buildNumber))
	if err != nil {
		return directoryFinal, err
	}

	dirSplitter := fmt.Sprintf("%d/artifact/", buildNumber)
	for _, a := range build.GetArtifacts() {
		subDir := ""
		if dirs := strings.SplitAfterN(a.Path, dirSplitter, 2); len(dirs) > 1 {
			subDir = filepath.Dir(dirs[1])
		}
		err = os.MkdirAll(filepath.Join(directoryFinal, subDir), os.ModePerm)
		if err != nil {
			return filepath.Join(directoryFinal, subDir), err
		}
		_, err := a.SaveToDir(jc.ctx, filepath.Join(directoryFinal, subDir))
		if err != nil {
			return directoryFinal, err
		}
	}
	return directoryFinal, nil
}

// Build a job with parameters
func (jc *Client) Build(name string, parameters map[string]string) (int64, error) {
	return jc.api.BuildJob(jc.ctx, name, parameters)
}

// Build a job with parameters
func (jc *Client) CreateJob(jobName string, xmlData string) (string, error) {
	job, err := jc.api.CreateJob(jc.ctx, xmlData, jobName)
	if err != nil {
		return "", err
	}
	return job.GetName(), nil
}
