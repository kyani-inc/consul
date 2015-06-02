# consul/discovery

Discovery is a utility that allows you to discover healthy instances that are behind an Elastic Load Balancer.

### Usage

Install the binary to `/usr/bin` or `/usr/local/bin`.

```
[#] elb-discovery --help
   --region, -r "us-east-1" aws region to query [$AWS_REGION]
   --load-balancer-name     name of elb to query
   --count "1"          Number of results to return. 0 returns all. Note: Any number that is *NOT* 1 returns as a json array.
   --private-ip-only        If set, application will only return private IP Addresses.
   --debug          Verbose: Turn on logging.
   --help, -h           show help
   --version, -v        print the version
```

Example response:

```
[#] elb-discovery --load-balancer-name consul --private-ip-only
10.0.3.119

[#] elb-discovery --load-balancer-name consul --count 0 --private-ip-only
["10.0.4.22","10.0.3.119"]

[#] elb-discovery --load-balancer-name consul
127.0.0.1
```

### Defaults

- `--region`: Defaults to "us-east-1". This can be overridden by the ENV Variable `$AWS_REGION` or the `--region` flag.
- `--count`: Defaults to 1. If Count equals 0 then it will return an "unlimited" number of results.
- `--private-ip-only`: Defaults to false. If not set the utility will return public ip addresses.


### Credentials

The utility attempts three methods of retrieving AWS Credentials and they are attempted in the order listed:

- `EC2RoleProvider`: This method uses an instance role, please see `iam-policy.json` for an IAM Policy that works with the instance role.
- `EnvProvider`: This method looks for the following environment variables: Access Key ID: AWS_ACCESS_KEY_ID or AWS_ACCESS_KEY - Secret Access Key: AWS_SECRET_ACCESS_KEY or AWS_SECRET_KEY
- `SharedCredentialsProvider`: This method looks for an ini file located at ` $HOME/.aws/credentials` See [here](http://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html) for more info.

### Why not use Atlas?

Atlas states that they will not always be free so we have developed this so there is no need to depend on Atlas.