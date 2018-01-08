package main

import (
	"errors"
	"fmt"
	"os"
	"flag"
	"strings"
)

const (
	DEFAULT_SESSION_TOKEN_DURATION = int64(60 * 60)
	DEFAULT_AWS_REGION = "eu-central-1"
	DEFAULT_ONELOGIN_REGION = "us"
)

func dieOnError(err error, message string) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: %v\n", message, err)
		os.Exit(1)
	}
}

func checkStringFlagNotEmpty(name string, f *string) {
	if f == nil || *f == "" {
		fmt.Fprintf(os.Stderr, "Missing mandatory parameter: %s\n\n", name)
		flag.Usage()
		os.Exit(1)
	}
}

func printAvailableRoles(assertion *SamlAssertion) {
	availableRoles := assertion.Parse()
	if availableRoles == nil {
		dieOnError(errors.New("No role avaiable to assume into"), "Error assuming roles")
	}
	fmt.Println("Available roles:")
	for role := range availableRoles {
		fmt.Println(role)
	}
}

func main() {
	// set up command line flags
	roleArns := flag.String("role-arn", "", "AWS role arn to assume into, separate by comma to assume multiple roles at once")
	profileName := flag.String("profile", "", "Write this AWS CLI profile, defaults to role name, separate by comma to assume multiple roles at once")
	duration := flag.Int64("duration", DEFAULT_SESSION_TOKEN_DURATION, "Token duration in seconds for target profile")
	awsRegion := flag.String("region", DEFAULT_AWS_REGION, "AWS region")
	username := flag.String("username", "", "onelogin user name")
	password := flag.String("password", "", "onelogin password")
	oneloginRegion := flag.String("onelogin-region", DEFAULT_ONELOGIN_REGION, "onelogin region, us/eu")
	mfaCode := flag.String("mfa-code", "", "MFA code")
	flag.Parse()

	checkStringFlagNotEmpty("username", username)
	checkStringFlagNotEmpty("password", password)
	checkStringFlagNotEmpty("mfa-code", mfaCode)

	client, err := NewOneloginClient(CLIENT_ID, CLIENT_SECRET, *oneloginRegion)
	dieOnError(err, "Error creating onelogin client")

	assertion, err := client.GetSamlAssertion(APP_ID, SUBDOMAIN, *username, *password, *mfaCode)
	dieOnError(err, "Error getting SAML assertion")

	fmt.Println("Successfully logged into onelogin")

	if *roleArns == "" {
		printAvailableRoles(assertion)
	} else {
		profileNames := strings.Split(*profileName, ",")
		for i, roleArn := range strings.Split(*roleArns, ",") {
			creds, err := AssumeRoleWithSaml(assertion, roleArn, *duration)
			dieOnError(err, "Error assuming AWS role")
			if *profileName == "" || i >= len(profileNames) {
				WriteProfile(creds, strings.SplitN(roleArn, "/", 2)[1], *awsRegion)
			} else {
				WriteProfile(creds, profileNames[i], *awsRegion)
			}
		}
	}
}