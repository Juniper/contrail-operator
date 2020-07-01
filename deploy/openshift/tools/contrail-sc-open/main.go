package main

import (
	"flag"
	"log"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

type port struct {
	From, To int64
	Protocol string
}

const (
	protocolTCP = "tcp"
	protocolUDP = "udp"
)

// ports is an array of all port ranges that should be open in order to make contrail work properly
var ports = [...]port{
	{5995, 5995, protocolTCP},
	{6379, 6379, protocolTCP},
	{5920, 5920, protocolTCP},
	{5921, 5921, protocolTCP},
	{4369, 4369, protocolTCP},
	{5673, 5673, protocolTCP},
	{25672, 25672, protocolTCP},
	{514, 514, protocolTCP},
	{6343, 6343, protocolTCP},
	{4739, 4739, protocolTCP},
	{5269, 5269, protocolTCP},
	{53, 53, protocolTCP},
	{179, 179, protocolTCP},
	{5672, 5672, protocolTCP},
	{10250, 10250, protocolTCP},
	{10256, 10256, protocolTCP},
	{80, 80, protocolTCP},
	{443, 443, protocolTCP},
	{1936, 1936, protocolTCP},
	{7000, 10000, protocolTCP},
	{2000, 3888, protocolTCP},
	{8053, 8053, protocolUDP},
	{4789, 4789, protocolUDP},
}

func setRules(svc *ec2.EC2, group *string) error {
	ruleCIDR := "0.0.0.0/0"
	for portIndex := range ports {
		_, err := svc.AuthorizeSecurityGroupIngress(&ec2.AuthorizeSecurityGroupIngressInput{
			GroupId: group,
			IpPermissions: []*ec2.IpPermission{
				&ec2.IpPermission{
					IpProtocol: &ports[portIndex].Protocol,
					FromPort:   &ports[portIndex].From,
					ToPort:     &ports[portIndex].To,
					IpRanges: []*ec2.IpRange{
						&ec2.IpRange{
							CidrIp: &ruleCIDR,
						},
					},
				},
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func main() {
	clusterName := flag.String("cluster-name", "", "Openshift cluster name.")
	region := flag.String("region", "eu-central-1", "AWS region where security groups are located.")
	flag.Parse()

	if *clusterName == "" {
		log.Fatal("No cluster name has been specified.")
		os.Exit(1)
	}

	clusterRegex := *clusterName + "-*"
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(*region)},
	)
	if err != nil {
		log.Fatal("Couldn't create new session to AWS region: ", region, "\n", err)
		os.Exit(1)
	}

	svc := ec2.New(sess)
	tagValueString := "tag-value"
	vpc, err := svc.DescribeVpcs(&ec2.DescribeVpcsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: &tagValueString,
				Values: []*string{
					&clusterRegex,
				},
			},
		},
	})
	if err != nil {
		log.Fatal("Couldn't get VPCs in region. Error:", err)
		os.Exit(1)
	}

	if len(vpc.Vpcs) != 1 {
		log.Fatal("Unable to get VPCs - there's ", len(vpc.Vpcs), " VPCs (should be 1)")
		os.Exit(1)
	}

	vpcID := *vpc.Vpcs[0].VpcId
	vpcIdstring := "vpc-id"
	result, err := svc.DescribeSecurityGroups(&ec2.DescribeSecurityGroupsInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name: &vpcIdstring,
				Values: []*string{
					&vpcID,
				},
			},
		},
	})
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case "InvalidGroupId.Malformed":
				fallthrough
			case "InvalidGroup.NotFound":
				log.Fatal(aerr.Message())
				os.Exit(1)
			}
		}
		log.Fatal("Unable to get descriptions for security groups, ", err)
		os.Exit(1)
	}

	log.Print("Security Group:")
	for _, group := range result.SecurityGroups {
		if *group.GroupName == "default" {
			continue
		}
		log.Print("Adding rules for security group: ", *group.GroupName)
		err := setRules(svc, group.GroupId)
		if err != nil {
			if aerr, ok := err.(awserr.Error); ok {
				if aerr.Code() != "InvalidPermission.Duplicate" {
					log.Fatal("Unable to set rules for security groups, ", *group.GroupName, "\nError: ", err)
					os.Exit(1)
				}
			}
		}
		log.Print("Rules added for scurity group: ", *group.GroupName)
	}
}
