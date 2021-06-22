# jenkinsctl
`jenkinsctl` is a command line tool to interact with jenkins server.<br>
Inspired by `kubectl` and `oc` (OpenShift) client.

### Download the client
* [Releases](https://github.com/jkandasa/jenkinsctl/releases/latest)

### Examples
To get login token,
> The API token is available in your personal configuration page. Click your name on the top right corner on every page, then click "Configure" to see your API token. (The URL `$root/me/configure` is a good shortcut.) You can also change your API token from here.

Source: [Jenkins doc](https://www.jenkins.io/doc/book/system-administration/authenticating-scripted-clients/)

#### Login to jenkins server
```bash
$ jenkinsctl login http://localhost:8080 --username jeeva --password 11f3d04172eb97b3f4c0287911e5832b00
Login successful.

# prints the version details
$ jenkinsctl version
Client Version: {version:master, buildDate:2021-06-21T17:09:18+00:00, gitCommit:5878ed77e8c28eeae04cb2311bbe44deabaff6b2, goLangVersion:go1.16.3, platform:linux/amd64}
Server Version: 2.289.1

```
#### Display jobs
```bash
$ go run cmd/main.go jobs
COLOR   	NAME                      	CLASS                                         	URL                                                   
        	JobFolder1/job/Foleder2   	com.cloudbees.hudson.plugins.folder.Folder    	http://localhost:8080/job/JobFolder1/job/Foleder2/   	
notbuilt	JobFolder1/job/folder-job1	hudson.model.FreeStyleProject                 	http://localhost:8080/job/JobFolder1/job/folder-job1/	
red     	pipeline job              	org.jenkinsci.plugins.workflow.job.WorkflowJob	http://localhost:8080/job/pipeline%20job/            	
blue    	test-job-1                	hudson.model.FreeStyleProject                 	http://localhost:8080/job/test-job-1/                	
notbuilt	test-job-2                	hudson.model.FreeStyleProject                 	http://localhost:8080/job/test-job-2/              	

$ go run cmd/main.go jobs --depth 2
COLOR   	NAME                                     	CLASS                                         	URL                                                                    
blue    	JobFolder1/job/Foleder2/job/hello-job 123	hudson.model.FreeStyleProject                 	http://localhost:8080/job/JobFolder1/job/Foleder2/job/hello-job%20123/	
notbuilt	JobFolder1/job/folder-job1               	hudson.model.FreeStyleProject                 	http://localhost:8080/job/JobFolder1/job/folder-job1/                 	
red     	pipeline job                             	org.jenkinsci.plugins.workflow.job.WorkflowJob	http://localhost:8080/job/pipeline%20job/                             	
blue    	test-job-1                               	hudson.model.FreeStyleProject                 	http://localhost:8080/job/test-job-1/                                 	
notbuilt	test-job-2                               	hudson.model.FreeStyleProject                 	http://localhost:8080/job/test-job-2/  
```

#### Switch to a job
```bash
$ jenkinsctl job test-job-1
Switched to 'test-job-1' at 'http://localhost:8080'
```
#### Display parameters of the job
```bash
$ jenkinsctl get parameters
NAME  	DEFAULT VALUE	TYPE                      	DESCRIPTION   
PARAM1	test         	StringParameterDefinition 	Param value 1	
PARAM2	test2        	StringParameterDefinition 	Param value 2	
PARAM3	true         	BooleanParameterDefinition	Param value 3	
```
#### Build a job
```bash
# sample build config file
$ cat example_build.yaml
kind: build
spec:
  job_name: test-job-1
  parameters:
    PARAM1: custom value 1
    PARAM2: custom value 2
    PARAM3: false

# trigger a build
$ jenkinsctl create --file ./example_build.yaml 
build created on the job 'test-job-1', build queue id:11
```
#### Display past builds
```bash
$ jenkinsctl get builds --limit 2
NUMBER	TRIGGERED BY          	RESULT 	IS RUNNING	DURATION	TIMESTAMP                        	REVISION 
11    	Jeeva Kandasamy(jeeva)	SUCCESS	false     	8ms     	2021-06-21 22:26:09.756 +0530 IST	        	
10    	Jeeva Kandasamy(jeeva)	SUCCESS	false     	19ms    	2021-06-21 22:25:55.894 +0530 IST	        	
```
#### display details of a build
```
$ jenkinsctl get build 11
KEY            	VALUE                                    
URL            	http://localhost:8080/job/test-job-1/11/	
Build Number   	11                                      	
Triggered By   	Jeeva Kandasamy(jeeva)                  	
Result         	SUCCESS                                 	
Is Running     	false                                   	
Duration       	0s                                      	
Revision       	                                        	
Revision Branch	                                        	
Timestamp      	2021-06-21 22:26:09.756 +0530 IST       	

Parameters:
PARAM1	custom value 1	
PARAM2	custom value 2	
PARAM3	              	

Test Result:
KEY     	VALUE 
Passed  	0    	
Failed  	0    	
Skipped 	0    	
Duration	0s   	

Artifacts:
PATH 
```
#### Display a build details in yaml format
```bash
$ jenkinsctl get build 11 --output yaml
url: http://localhost:8080/job/test-job-1/11/
number: 11
triggered_by: Jeeva Kandasamy(jeeva)
parameters:
- name: PARAM1
  value: custom value 1
- name: PARAM2
  value: custom value 2
- name: PARAM3
  value: ""
injected_env_vars: {}
causes:
- _class: hudson.model.Cause$UserIdCause
  shortDescription: Started by user Jeeva Kandasamy
  userId: jeeva
  userName: Jeeva Kandasamy
duration: 8ms
console: null
result: SUCCESS
is_running: false
revision: ""
revision_branch: ""
timestamp: 2021-06-21T22:26:09.756+05:30
test_result:
  duration: 0
  empty: false
  failcount: 0
  passcount: 0
  skipcount: 0
  suites: []
artifacts: []
```
#### Display a build details in json format
```bash
$ jenkinsctl get build 11 --output json --pretty
{
 "url": "http://localhost:8080/job/test-job-1/11/",
 "number": 11,
 "triggeredBy": "Jeeva Kandasamy(jeeva)",
 "parameters": [
  {
   "name": "PARAM1",
   "value": "custom value 1"
  },
  {
   "name": "PARAM2",
   "value": "custom value 2"
  },
  {
   "name": "PARAM3",
   "value": ""
  }
 ],
 "injectedEnvVars": null,
 "causes": [
  {
   "_class": "hudson.model.Cause$UserIdCause",
   "shortDescription": "Started by user Jeeva Kandasamy",
   "userId": "jeeva",
   "userName": "Jeeva Kandasamy"
  }
 ],
 "duration": 8000000,
 "console": null,
 "result": "SUCCESS",
 "isRunning": false,
 "revision": "",
 "revisionBranch": "",
 "timestamp": "2021-06-21T22:26:09.756+05:30",
 "testResult": {
  "duration": 0,
  "empty": false,
  "failCount": 0,
  "passCount": 0,
  "skipCount": 0,
  "suites": null
 },
 "artifacts": null
```
#### Display the console log of a build
```bash
$ jenkinsctl get console 11
Started by user Jeeva Kandasamy
Running as SYSTEM
Building in workspace /var/jenkins_home/workspace/test-job-1
[test-job-1] $ /bin/sh -xe /tmp/jenkins2269490024621194842.sh
...
WORKSPACE_TMP=/var/jenkins_home/workspace/test-job-1@tmp
+ date
Mon 21 Jun 2021 04:56:09 PM UTC
+ hostname
33501f943b12
Finished: SUCCESS

# watch running builds console logs
$ jenkinsctl get console 12 --watch
Started by user Jeeva Kandasamy
Running as SYSTEM
Building in workspace /var/jenkins_home/workspace/test-job-1
[test-job-1] $ /bin/sh -xe /tmp/jenkins2269490024621194842.sh
...
```