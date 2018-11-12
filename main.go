package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
)

const (
	SourceableTemplate string = `AWS_ACCESS_KEY_ID=%s
AWS_SECRET_ACCESS_KEY=%s
AWS_SESSION_TOKEN=%s
AWS_ASSUMED_ROLE_ARN=%s
AWS_ASSUMED_ROLE_ID=%s`
)

var (
	roleArn     string
	sessionName string
	format      string
)

func init() {
	flag.StringVar(&roleArn, "role-arn", "", "ARN of the IAM role to assume")
	flag.StringVar(&sessionName, "session-name", "", "Session name for assumed credentials")
	flag.StringVar(&format, "format", "sourceable", "Set the output format")
	flag.Parse()

	if roleArn == "" || sessionName == "" {
		log.Fatalln("role-arn and session-name must be set")
	}

	if format != "sourceable" && format != "json" {
		log.Fatalln("Only 'json' and 'sourceable' output formats are supported")
	}
}

func renderJson(id, key, token, arn, rid string) {

}

func renderSourceable(id, key, token, arn, rid string) {
	fmt.Printf(SourceableTemplate, id, key, token, arn, rid)
}

func render(r *sts.AssumeRoleOutput) {
	keyId := aws.StringValue(r.Credentials.AccessKeyId)
	key := aws.StringValue(r.Credentials.SecretAccessKey)
	tok := aws.StringValue(r.Credentials.SessionToken)
	assumedArn := aws.StringValue(r.AssumedRoleUser.Arn)
	assumedId := aws.StringValue(r.AssumedRoleUser.AssumedRoleId)

	switch format {
	case "json":
		renderJson(keyId, key, tok, assumedArn, assumedId)
	case "sourceable":
		renderSourceable(keyId, key, tok, assumedArn, assumedId)
	}
}

func main() {
	svc := sts.New(session.New())
	input := &sts.AssumeRoleInput{
		RoleArn:         aws.String(roleArn),
		RoleSessionName: aws.String(sessionName),
	}

	result, err := svc.AssumeRole(input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case sts.ErrCodeMalformedPolicyDocumentException:
				log.Fatalln(sts.ErrCodeMalformedPolicyDocumentException, aerr.Error())
			case sts.ErrCodePackedPolicyTooLargeException:
				log.Fatalln(sts.ErrCodePackedPolicyTooLargeException, aerr.Error())
			case sts.ErrCodeRegionDisabledException:
				log.Fatalln(sts.ErrCodeRegionDisabledException, aerr.Error())
			default:
				log.Fatalln(aerr.Error())
			}
		}

		log.Fatalln(err.Error())
	}

	render(result)
}
