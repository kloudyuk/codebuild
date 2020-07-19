package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/cloudwatchlogs"
	"github.com/aws/aws-sdk-go/service/codebuild"
	"github.com/aws/aws-sdk-go/service/sts"
)

// Flags
var envVars env
var roleARN string
var serviceRole string
var sourceLocation string
var sourceType string
var sourceVersion string
var tail bool
var wait bool

// Args
var project string

// Source represents the CodeBuild source options
type Source struct {
	Type     string
	Location string
	Version  string
}

var accountID string
var region string
var sess *session.Session

func usage() {
	fmt.Printf(`
Usage:
  codebuild [flags] project

Args:
  project string
	The name of the CodeBuild project

Flags:
  -help -h
	Show command help
`)
	flag.PrintDefaults()
	fmt.Println()
}

func init() {
	envVars = make(map[string]string)
	flag.Usage = usage
	flag.Var(&envVars, "e", "Override environment variable (can be provided multiple times e.g. -e NAME=value -e ANOTHER_NAME=value)")
	flag.StringVar(&roleARN, "role", "", "Assume the given role before making the request to CodeBuild")
	flag.StringVar(&serviceRole, "service-role", "", "Override the service role")
	flag.StringVar(&sourceLocation, "src-location", "", "Override the source location")
	flag.StringVar(&sourceType, "src-type", "", "Override the source type")
	flag.StringVar(&sourceVersion, "src-version", "", "Override the source version")
	flag.BoolVar(&tail, "tail", false, "Tail the logs (implies -wait)")
	flag.BoolVar(&wait, "wait", false, "Wait for the build to complete")
	flag.Parse()
	sess = session.Must(session.NewSession())
	region = *sess.Config.Region
}

func main() {

	// Vaildate args
	if len(flag.Args()) != 1 {
		flag.Usage()
		fmt.Println("ERROR: Missing required argument: project")
		os.Exit(2)
	}
	project = flag.Args()[0]

	// Ensure wait is true if tail is true
	if tail {
		wait = true
	}

	// if we've been given a role ARN, get a new session based on the assumed role
	var err error
	if roleARN != "" {
		sess, err = getSessionForRole(roleARN)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	// Get the AWS account ID
	svc := sts.New(sess)
	callerID, err := svc.GetCallerIdentity(&sts.GetCallerIdentityInput{})
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	accountID = *callerID.Account

	// Create a source object from the flags
	src := &Source{
		sourceType,
		sourceLocation,
		sourceVersion,
	}

	// Attempt to fill out missing source info from environment variables
	src, err = sourceFromEnv(src)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Start the AWS CodeBuild build
	fmt.Printf("Starting AWS CodeBuild for project: %s\n", project)
	out, err := StartBuild(project, serviceRole, src, envVars)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	buildID := *out.Build.Id
	fmt.Printf("Build URL: %s\n", buildURL(buildID))

	// Tail the CloudWatch log stream
	if tail {
		fmt.Println("Waiting for CloudWatch log info...")
		logGroup, logStream, err := waitForLogInfo(buildID)
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
		// Start tailing the logs
		go func() {
			fmt.Printf("Tailing logs from CloudWatch: %s/%s\n", logGroup, logStream)
			fmt.Println("--------------------------------------------------------------------------------")
			if err := Tail(context.Background(), logGroup, logStream); err != nil {
				fmt.Println(err)
				os.Exit(1)
			}
		}()
	}

	// Wait for the CodeBuild build to complete
	if wait {
		if err := Wait(context.Background(), buildID); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

}

// StartBuild starts an AWS CodeBuild build
func StartBuild(project string, serviceRole string, src *Source, envVars env) (*codebuild.StartBuildOutput, error) {
	svc := codebuild.New(sess)
	in := &codebuild.StartBuildInput{
		ProjectName: aws.String(project),
	}
	if serviceRole != "" {
		in.ServiceRoleOverride = aws.String(serviceRole)
	}
	if src.Type != "" {
		in.SourceTypeOverride = aws.String(src.Type)
	}
	if src.Location != "" {
		in.SourceLocationOverride = aws.String(src.Location)
	}
	if src.Version != "" {
		in.SourceVersion = aws.String(src.Version)
	}
	for k, v := range envVars {
		in.EnvironmentVariablesOverride = append(in.EnvironmentVariablesOverride, &codebuild.EnvironmentVariable{
			Name:  aws.String(k),
			Value: aws.String(v),
		})
	}
	return svc.StartBuild(in)
}

// Wait for a given AWS CodeBuild build to complete
func Wait(ctx context.Context, id string) error {
	for {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		out, err := getBuild(id)
		if err != nil {
			return err
		}
		if *out.BuildComplete {
			// the build may complete before we have the full set of logs so wait a bit before we exit
			time.Sleep(10 * time.Second)
			return nil
		}
		time.Sleep(3 * time.Second)
	}
}

// Tail tails logs from a CloudWatch log stream
func Tail(ctx context.Context, logGroup string, logStream string) error {
	svc := cloudwatchlogs.New(sess)
	in := &cloudwatchlogs.GetLogEventsInput{
		LogGroupName:  aws.String(logGroup),
		LogStreamName: aws.String(logStream),
		StartFromHead: aws.Bool(true),
	}
	for {
		err := svc.GetLogEventsPagesWithContext(ctx, in, func(page *cloudwatchlogs.GetLogEventsOutput, lastPage bool) bool {
			for _, p := range page.Events {
				fmt.Print(*p.Message)
			}
			in.NextToken = page.NextForwardToken
			return lastPage
		})
		if err != nil {
			return err
		}
		time.Sleep(3 * time.Second)
	}
}
