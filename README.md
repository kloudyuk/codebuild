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
    	Override environment variable (can be provided multiple times e.g. -e NAME=value -e ANOTHER_NAME=value)
  -role string
    	Assume the given role before making the request to CodeBuild
  -service-role string
    	Override the service role
  -src-location string
    	Override the source location
  -src-type string
    	Override the source type
  -src-version string
    	Override the source version
  -tail
    	Tail the logs (implies -wait)
  -wait
    	Wait for the build to complete
```
