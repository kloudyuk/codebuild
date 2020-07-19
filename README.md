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
  -help -h
    	Show command help
  -e value
    	Override a CodeBuild environment variable (can be provided multiple times e.g. -e NAME=value -e ANOTHER_NAME=value)
  -src-location string
    	Override the CodeBuild source location
  -src-type string
    	Override the CodeBuild source type
  -src-version string
    	Override the CodeBuild source version
  -tail
    	Tail the logs via the CloudWatch log stream (implies -wait)
  -wait
    	Wait for the build to complete
```
