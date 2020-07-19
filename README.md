# codebuild

Start a build in AWS CodeBuild

Includes options to tail the logs and wait for a build to complete

## Usage

```text
Usage:
  codebuild [flags] project

Args:
  project string
        The name of the CodeBuild project

Flags:
  -e value
    	Environment variable override (can be provided multiple times)
    	e.g. -e NAME=value -e ANOTHER_NAME=value
  -location string
    	Source location override
  -tail
    	Tail the logs via the CloudWatch log stream (implies -wait)
  -wait
    	Wait for the build to complete
```
