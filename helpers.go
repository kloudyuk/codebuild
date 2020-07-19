package main

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
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
