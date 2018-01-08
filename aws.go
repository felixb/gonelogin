package main

import (
	"fmt"
	"path/filepath"
	"os"
	"os/user"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/go-ini/ini"
)

func AssumeRoleWithSaml(assertion *SamlAssertion, roleArn string, duration int64) (*sts.Credentials, error) {
	sess := session.Must(session.NewSession())
	svc := sts.New(sess)

	providerArn := assertion.GetProvider(roleArn)
	if providerArn == nil {
		return nil, fmt.Errorf("Error finding identity provider for role %s", roleArn)
	}

	output, err := svc.AssumeRoleWithSAML(&sts.AssumeRoleWithSAMLInput{
		RoleArn:         &roleArn,
		PrincipalArn:    providerArn,
		SAMLAssertion:   assertion.Assertion,
		DurationSeconds: &duration,
	})

	if err != nil {
		return nil, err
	}

	return output.Credentials, nil

}

func WriteProfile(cred *sts.Credentials, name, region string) error {
	usr, err := user.Current()
	if err != nil {
		return err
	}

	awsPath := filepath.Join(usr.HomeDir, ".aws")
	filename := filepath.Join(awsPath, "credentials")

	cfg, err := ini.Load(filename)
	if err != nil {
		fmt.Printf("Unable to find credentials file %s. Creating new file.\n", filename)

		if err := os.MkdirAll(awsPath, os.ModePerm); err != nil {
			return err
		}
		cfg = ini.Empty()
	}
	sec, err := cfg.GetSection(name)
	if err != nil {
		sec, err = cfg.NewSection(name)
		if err != nil {
			return err
		}
	}
	updateKey(sec, "aws_access_key_id", cred.AccessKeyId)
	updateKey(sec, "aws_secret_access_key", cred.SecretAccessKey)
	updateKey(sec, "aws_session_token", cred.SessionToken)
	if region != "" {
		updateKey(sec, "region", &region)
	}

	if err := cfg.SaveTo(filename); err != nil {
		return err
	}

	fmt.Printf("Wrote session token for profile %s\n", name)
	fmt.Printf("Token is valid until: %v\n", cred.Expiration)

	return nil
}

func updateKey(sec *ini.Section, name string, value *string) error {
	key, err := sec.GetKey(name)
	if err != nil {
		_, err := sec.NewKey(name, *value)
		if err != nil {
			return err
		}
	} else {
		key.SetValue(*value)
	}
	return nil
}