package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/codegangsta/cli"

	"github.com/awslabs/aws-sdk-go/aws"
	"github.com/awslabs/aws-sdk-go/aws/awserr"
	"github.com/awslabs/aws-sdk-go/aws/credentials"
	"github.com/awslabs/aws-sdk-go/service/ec2"
	"github.com/awslabs/aws-sdk-go/service/elb"
)

const (
	ERR_SUCCESS = iota
	ERR_GENERAL
	_
	_
	_
	ERR_API_FAILURE
	ERR_BAD_INPUT
)

var (
	creds = credentials.NewChainCredentials([]credentials.Provider{
		&credentials.EC2RoleProvider{},
		&credentials.EnvProvider{},
		&credentials.SharedCredentialsProvider{},
	})
	inDebug bool
)

type config struct {
	*aws.Config
}

func main() {
	app := cli.NewApp()
	app.Name = "elb-discovery"
	app.Usage = "Return the IP Address of one or more healthy instance(s) behind a load balancer."
	app.Version = "1.0"
	app.Authors = []cli.Author{
		cli.Author{
			Name:  "Peter Olds",
			Email: "polds@kyanicorp.com",
		},
	}
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name:   "region, r",
			Value:  "us-east-1",
			Usage:  "aws region to query",
			EnvVar: "AWS_REGION",
		},
		cli.StringFlag{
			Name:  "load-balancer-name",
			Usage: "name of elb to query",
		},
		cli.IntFlag{
			Name:  "count",
			Value: 1,
			Usage: "Number of results to return. 0 returns all. Note: Any number that is *NOT* 1 returns as a json array.",
		},
		cli.BoolFlag{
			Name:  "private-ip-only",
			Usage: "If set, application will only return private IP Addresses.",
		},
		cli.BoolFlag{
			Name:  "debug",
			Usage: "Verbose: Turn on logging.",
		},
	}

	app.Action = func(c *cli.Context) {
		var err error
		inDebug = c.Bool("debug")

		// Validate the arguments
		err = validation(c)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(ERR_BAD_INPUT)
		}
		log("Arguments valid.")

		// Verify we have credentials. If we don't we die.
		_, err = creds.Get()
		if err != nil {
			fmt.Fprintf(os.Stderr, "[error] [credentials] %s\n\n", err.Error())
			os.Exit(ERR_GENERAL)
		}
		log("Credentials valid.")

		var cfg config
		cfg.Config = &aws.Config{
			Region:      c.String("region"),
			Credentials: creds,
		}

		// Get a list of all instances behind a load balancer
		instances, err := cfg.ELBInstances(c.String("load-balancer-name"))
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(ERR_API_FAILURE)
		}
		log(fmt.Sprintf("Instances Available: %d", len(instances.InstanceStates)))

		// Filter those instances to only healthy instances
		healthy := cfg.ELBHealthyInstances(instances)
		if len(healthy) == 0 {
			fmt.Fprintln(os.Stderr, "[error] No healthy instances available in the ELB (1).")
			os.Exit(ERR_API_FAILURE)
		}

		// Get the list of IP Addresses based on the healthy instances
		count := c.Int("count")
		ips := cfg.FetchIP(healthy, count, c.Bool("private-ip-only"))
		if len(ips) == 0 {
			fmt.Fprintln(os.Stderr, "[error] No healthy instances available in the ELB (2).")
			os.Exit(ERR_API_FAILURE)
		}

		// Return a string if they wanted one IP (Default)
		if len(ips) == 1 && count == 1 {
			fmt.Fprintln(os.Stdout, ips[0])
			os.Exit(ERR_SUCCESS)
		}

		// Return a JSON Object otherwise.
		bips, err := json.Marshal(ips)
		if err != nil {
			fmt.Fprintf(os.Stderr, "[error] %s\n", err.Error())
			os.Exit(ERR_GENERAL)
		}

		fmt.Fprintf(os.Stdout, "%s\n", bips)
		os.Exit(ERR_SUCCESS)
	}

	app.Run(os.Args)
}

// FetchIP determines how many healthy instances we have versus what the user defined as the flag
// and then loops over each healthy instance. If a call fails it attempts to go to the next one.
func (cfg config) FetchIP(instances []*elb.InstanceState, count int, isPrivate bool) []string {
	var res []string

	// Determine exactly how many times we need to iterate.
	num := count
	if count == 0 || count > len(instances) {
		num = len(instances)
	}
	num--

	for i := 0; i <= num; i++ {
		if i == len(instances) {
			break
		}

		ip, err := cfg.getIP(instances[i].InstanceID, isPrivate)
		if err != nil {
			log(fmt.Sprint(err.Error()))
			// We threw an error so we'll recurse again.

			num++
			continue
		}

		res = append(res, *ip)
	}

	return res
}

func (cfg config) getIP(instanceID *string, isPrivate bool) (*string, error) {
	svc := ec2.New(cfg.Config)

	params := &ec2.DescribeInstancesInput{
		InstanceIDs: []*string{
			instanceID,
		},
	}

	resp, err := svc.DescribeInstances(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				err = fmt.Errorf("[error] [RequestFailure] %d: %s %s", reqErr.StatusCode(), reqErr.Code(), reqErr.Message())
				return nil, err
			}

			err = fmt.Errorf("[error] [aws.Error] %s %s", awsErr.Code(), awsErr.Message())
		} else {
			err = fmt.Errorf("[error] [UNKNOWN] %s", err.Error())
		}
	}

	instance := resp.Reservations[0].Instances[0]
	if isPrivate {
		return instance.PrivateIPAddress, err
	}

	return instance.PublicIPAddress, err
}

// validation validates arguments before continuing.
func validation(c *cli.Context) error {
	if c.String("region") == "" {
		return fmt.Errorf("[error] [validation] region cannot be empty.\n")
	}
	log(fmt.Sprintf("region = %s", c.String("region")))

	if c.String("load-balancer-name") == "" {
		return fmt.Errorf("[error] [validation] load-balancer-name cannot be empty.\n")
	}
	log(fmt.Sprintf("load-balancer-name = %s", c.String("load-balancer-name")))

	log(fmt.Sprintf("count = %d", c.Int("count")))
	log(fmt.Sprintf("private ip addresses = %t", c.Bool("private-ip-only")))

	return nil
}

// ELBInstances takes an ELB name and returns all instances attached to the ELB.
func (cfg config) ELBInstances(name string) (*elb.DescribeInstanceHealthOutput, error) {
	svc := elb.New(cfg.Config)

	params := &elb.DescribeInstanceHealthInput{
		LoadBalancerName: aws.String(name),
	}

	resp, err := svc.DescribeInstanceHealth(params)
	if err != nil {
		if awsErr, ok := err.(awserr.Error); ok {
			if reqErr, ok := err.(awserr.RequestFailure); ok {
				err = fmt.Errorf("[error] [RequestFailure] %d: %s %s\n", reqErr.StatusCode(), reqErr.Code(), reqErr.Message())
				return resp, err
			}

			err = fmt.Errorf("[error] [aws.Error] %s %s\n", awsErr.Code(), awsErr.Message())
		} else {
			err = fmt.Errorf("[error] [UNKNOWN] %s\n", err.Error())
		}
	}

	return resp, err
}

// ELBHealthyInstances takes the attached instances to the ELB and returns only "InService" instances.
func (config) ELBHealthyInstances(instances *elb.DescribeInstanceHealthOutput) []*elb.InstanceState {
	var (
		healthy = "InService"
		out     []*elb.InstanceState
	)

	for _, instance := range instances.InstanceStates {
		if *instance.State == healthy {
			out = append(out, instance)
		}
	}

	return out
}

func log(msg string) {
	if inDebug {
		fmt.Fprintf(os.Stdout, "[debug] %s\n", msg)
	}
}
