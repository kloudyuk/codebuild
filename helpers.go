package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/codebuild"
)

// getBuild gets info for a given CodebBuild build ID
func getBuild(id string) (*codebuild.Build, error) {
	svc := codebuild.New(sess)
	in := &codebuild.BatchGetBuildsInput{
		Ids: []*string{
			aws.String(id),
		},
	}
	out, err := svc.BatchGetBuilds(in)
	if err != nil {
		return nil, err
	}
	return out.Builds[0], nil
}

// waitForLogInfo waits until CodeBuild build info to contain the CloudWatch log group and log stream
func waitForLogInfo(id string) (string, string, error) {
	for {
		build, err := getBuild(id)
		if err != nil {
			return "", "", err
		}
		if build.Logs.GroupName != nil && build.Logs.StreamName != nil {
			return *build.Logs.GroupName, *build.Logs.StreamName, nil
		}
		time.Sleep(3 * time.Second)
	}
}

// buildURL returns the build URL for a given build ID
func buildURL(id string) string {
	return fmt.Sprintf("https://%s.console.aws.amazon.com/codesuite/codebuild/%s/projects/%s/build/%s/?region=%s",
		region, accountID, project, id, region)
}

var bitbucketEnvVars = map[string]string{
	"sourceLocation": "BITBUCKET_GIT_HTTP_ORIGIN",
	"sourceVersion":  "BITBUCKET_COMMIT",
}

// sourceFromEnv tries to get missing source info from environment variables based on the type
func sourceFromEnv(src *Source) (*Source, error) {
	if src.Type == "" {
		return src, nil
	}
	var lookup map[string]string
	switch src.Type {
	case "BITBUCKET":
		lookup = bitbucketEnvVars
	default:
		return nil, fmt.Errorf("Unknown source type: %s", src.Type)
	}
	if src.Location == "" {
		src.Location = os.Getenv(lookup["sourceLocation"])
	}
	if src.Version == "" {
		src.Version = os.Getenv(lookup["sourceVersion"])
	}
	return src, nil
}

// getSessionForRole returns a new session valid the given role ARN
func getSessionForRole(roleARN string) (*session.Session, error) {
	creds := stscreds.NewCredentials(sess, roleARN)
	return session.NewSession(&aws.Config{Credentials: creds})
}
