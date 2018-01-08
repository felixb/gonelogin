# gonelogin

A simple onelogin -> aws cli client.

## build

Copy the `example/config.go` into your main folder and change the constants to your organisations onelogin settings.
Then run `make` to build `gonelogin`.

## use

You can list available AWS roles by leaving the `-role-arn` parameter empty:
```bash
$ gonelogin \
    -username "<your onelogin username or email>" \
    -password "<your onelogin password>" \
    -mfa-code "<your onelogin mfa code>"
Successfully logged into onelogin
Available roles:
arn:aws:iam::111111111111:role/some-role
arn:aws:iam::222222222222:role/some-role
arn:aws:iam::222222222222:role/some-other-role
```

Assume into a single role by stating so:
```bash
$ gonelogin \
    -username "<your onelogin username or email>" \
    -password "<your onelogin password>" \
    -mfa-code "<your onelogin mfa code>" \
    -role-arn arn:aws:iam::111111111111:role/some-role
Successfully logged into onelogin
Wrote session token for profile some-role
Token is valid until: 2018-01-08 11:38:03 +0000 UTC
```
Use the AWS profile as always e.g. be running `aws sts get-caller-identity --profile some-role`.

You may specify the target role name by setting the `-profile` parameter as well.

Assume into multiple roles ar onc by stating multiple role arns like this:
```bash
$ gonelogin \
    -username "<your onelogin username or email>" \
    -password "<your onelogin password>" \
    -mfa-code "<your onelogin mfa code>" \
    -role-arn arn:aws:iam::111111111111:role/some-role,arn:aws:iam::222222222222:role/some-other-role
Successfully logged into onelogin
Wrote session token for profile some-role
Token is valid until: 2018-01-08 11:40:39 +0000 UTC
Wrote session token for profile some-other-role
Token is valid until: 2018-01-08 11:40:39 +0000 UTC
```

You may specify the target role names by setting the `-profile` parameter as well with multiple roles:
```bash
$ gonelogin \
    -username "<your onelogin username or email>" \
    -password "<your onelogin password>" \
    -mfa-code "<your onelogin mfa code>" \
    -role-arn arn:aws:iam::111111111111:role/some-role,arn:aws:iam::222222222222:role/some-role,arn:aws:iam::222222222222:role/some-other-role \
    -profile some-profile,some-other-profile,yet-another-profile
Successfully logged into onelogin
Wrote session token for profile some-profile
Token is valid until: 2018-01-08 11:40:39 +0000 UTC
Wrote session token for profile some-other-profile
Token is valid until: 2018-01-08 11:40:39 +0000 UTC
Wrote session token for profile yet-another-profile
Token is valid until: 2018-01-08 11:40:39 +0000 UTC
```

