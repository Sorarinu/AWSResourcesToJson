/*
Package resources は、AWS に存在する各種リソースの情報を取得したりするパッケージです。
*/
package resources

import (
	"strings"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/elbv2"
)

// Instance .
type Instance struct {
	InstanceID   string `json:"InstanceId"`
	InstanceType string `json:"InstanceType"`
	Placement    string `json:"Placement"`
	PrivateIP    string `json:"PrivateIP"`
	PublicIP     string `json:"PublicIP"`
	State        string `json:"State"`
	Tags         []Tag  `json:"Tags"`
	Name         string `json:"Name"`
	LoadBalancer LoadBalancer `json:"LoadBalancer"`
}

// LoadBalancer .
type LoadBalancer struct {
	Arn string `json:"Arn"`
	Name string `json:"Name"`
	Tags []Tag `json:"Tags"`
	TargetGroups []TargetGroup `json:"TargetGroups"`
}

// TargetGroup .
type TargetGroup struct {
	Arn string `json:"Arn"`
	Name string `json:"Name"`
	Targets []Target `json:"Targets"`
}

// Target .
type Target struct {
	InstanceID string `json:"InstanceId"`
	State string `json:"State"`
}

// Tag .
type Tag struct {
	Key   string `json:"Key"`
	Value string `json:"Value"`
}

// GetEC2Instances は、与えられたリージョンに存在する EC2 インスタンスの一覧を取得して返します。
func GetEC2Instances(region string) ([]Instance, error) {
	sess := session.Must(session.NewSession())
	svc := ec2.New(
		sess,
		aws.NewConfig().WithRegion(region),
	)

	res, err := svc.DescribeInstances(nil)
	if err != nil {
		return nil, err
	}

	var instances []Instance

	for _, r := range res.Reservations {
		for _, i := range r.Instances {
			var name string
			var tags []Tag

			for _, t := range i.Tags {
				if *t.Key == "Name" {
					name = *t.Value
				}

				tags = append(tags, Tag{
					*t.Key,
					*t.Value,
				})
			}

			if i.PublicIpAddress == nil {
				i.PublicIpAddress = aws.String("")
			}

			if i.PrivateIpAddress == nil {
				i.PrivateIpAddress = aws.String("")
			}

			instances = append(instances, Instance{
				*i.InstanceId,
				*i.InstanceType,
				*i.Placement.AvailabilityZone,
				*i.PrivateIpAddress,
				*i.PublicIpAddress,
				*i.State.Name,
				tags,
				name,
				LoadBalancer{},
			})
		}
	}

	return instances, nil
}

// GetALBResources は、与えられたリージョンに存在する ALB リソースの一覧を取得して返します。
func GetALBResources(region string) ([]LoadBalancer, error) {
	sess := session.Must(session.NewSession())
	svc := elbv2.New(
		sess,
		aws.NewConfig().WithRegion(region),
	)

	res, err := svc.DescribeLoadBalancers(nil)
	if err != nil {
		return nil, err
	}

	var loadbalancers []LoadBalancer

	for _, l := range res.LoadBalancers {
		var tags []Tag
		var targetGroups []TargetGroup

		tags, err = getALBTags(region, *l.LoadBalancerArn)
		if err != nil {
			return nil, err
		}

		targetGroups, err = getTargetGroups(region, *l.LoadBalancerArn)
		if err != nil {
			return nil, err
		}

		loadbalancers = append(loadbalancers, LoadBalancer{
			*l.LoadBalancerArn,
			*l.LoadBalancerName,
			tags,
			targetGroups,
		})
	}

	return loadbalancers, nil
}

// getALBTags は、与えられたロードバランサに付与されているタグの情報を取得して返します。
func getALBTags(region string, arn string) ([]Tag, error) {
	sess := session.Must(session.NewSession())
	svc := elbv2.New(
		sess,
		aws.NewConfig().WithRegion(region),
	)

	input := &elbv2.DescribeTagsInput{
		ResourceArns: []*string{
			aws.String(arn),
		},
	}

	res, err := svc.DescribeTags(input)
	if err != nil {
		return nil, err
	}

	var tags []Tag

	for _, t := range res.TagDescriptions {
		if t.Tags == nil {
			return tags, nil
		}

		for _, v := range t.Tags {
			tags = append(tags, Tag{
				*v.Key,
				*v.Value,
			})
		}
	}

	return tags, nil
}

// getTargetGroups は、与えられたロードバランサに紐付いているターゲットグループを取得して返します。
func getTargetGroups(region string, arn string) ([]TargetGroup, error) {
	sess := session.Must(session.NewSession())
	svc := elbv2.New(
		sess,
		aws.NewConfig().WithRegion(region),
	)

	input := &elbv2.DescribeTargetGroupsInput{
		LoadBalancerArn: aws.String(arn),
	}

	res, err := svc.DescribeTargetGroups(input)
	if err != nil {
		return nil, err
	}

	var targetGroups []TargetGroup

	for _, t := range res.TargetGroups {
		var targets []Target

		targets, err = getTargets(region, *t.TargetGroupArn)
		if err != nil {
			return nil, err
		}

		targetGroups = append(targetGroups, TargetGroup{
			*t.TargetGroupArn,
			*t.TargetGroupName,
			targets,
		})
	}

	return targetGroups, nil
}

// getTargets は、与えられたターゲットグループに紐付いているターゲット（インスタンス）の情報を取得して返します。
func getTargets(region string, arn string) ([]Target, error) {
	sess := session.Must(session.NewSession())
	svc := elbv2.New(
		sess,
		aws.NewConfig().WithRegion(region),
	)

	input := &elbv2.DescribeTargetHealthInput{
		TargetGroupArn: aws.String(arn),
	}

	res, err := svc.DescribeTargetHealth(input)
	if err != nil {
		return nil, err
	}

	var targets []Target

	for _, t := range res.TargetHealthDescriptions {
		targets = append(targets, Target{
			*t.Target.Id,
			*t.TargetHealth.State,
		})
	}

	return targets, nil
}

// MergeResources は、与えられたそれぞれの EC2、ALB をもとに情報を集約して返します。
func MergeResources(instances []Instance, loadbalancers []LoadBalancer) []Instance{
	for n, i := range instances {
		for _, l := range loadbalancers {
			for _, tg := range l.TargetGroups {
				for _, t := range tg.Targets {
					if strings.Index(t.InstanceID, i.InstanceID) == -1 {
						continue
					}
					instances[n].LoadBalancer = l
				}
			}
		}
	}
	
	return instances
}