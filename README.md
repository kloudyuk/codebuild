# codebuild

Start a build in AWS CodeBuild

## Features

- flags or environment variables to configure the AWS CodeBuild options
- ability to assume an IAM role before executing the build
- can wait for a build to complete execution
- can tail the build logs from the CloudWatch log stream
- automatically populates source info from the environment if source-type is provided (currently supports type: BITBUCKET)

## Usage

```text
Usage:
  codebuild PROJECT [FLAGS]

Flags:
  -c, --compute-type string      Override the compute type
  -e, --env NAME=VALUE           Override environment variables
                                 Can be provided multiple times or as a comma separated list
  -f, --follow                   Tail the logs (implies --wait=true)
      --role-arn string          Assume the given role before making the request to CodeBuild
      --service-role string      Override the service role
  -l, --source-location string   Override the source location
  -t, --source-type string       Override the source type
  -v, --source-version string    Override the source version
  -w, --wait                     Wait for the build to complete
```
