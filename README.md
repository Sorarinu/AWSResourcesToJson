# AWSResourcesToJson

This outputs EC2 existing in AWS and ALB information associated with EC2 in json.

## Requirements

* Docker
    * Docker version 19.03.11, build 42e35e61f3
* docker-compose
    * docker-compose version 1.25.0, build 0a186604
* Golang
    * go version go1.14.4 linux/amd64

## Build

```
$ make
```

## Usage

1. Register the `REGION` in the environment variable.
2. Register AWS credentials with the `aws configure` command. This tool uses the default profile.
3. Run `go run main.go` or `./main.handle` after doing the build.

## Output

```
[
    :
    :
  {
    "InstanceId": "i-xxxxxxxxxxxxx",
    "InstanceType": "t3.small",
    "Placement": "ap-northeast-1a",
    "PrivateIP": "xxx.xxx.xxx.xxx",
    "PublicIP": "yyy.yyy.yyy.yyy",
    "State": "running",
    "Tags": [
      {
        "Key": "Name",
        "Value": "ec2_instance_hoge"
      }
    ],
    "Name": "ec2_instance_hoge",
    "LoadBalancer": {
      "Arn": "arn:aws:elasticloadbalancing:ap-northeast-1:xxxxxxxxxx:loadbalancer/app/hoge-loadbalancer/zzzzzzzzzzzz",
      "Name": "cnt-coi-alb01",
      "Tags": [
        {
          "Key": "hoge",
          "Value": "true"
        }
      ],
      "TargetGroups": [
        {
          "Arn": "arn:aws:elasticloadbalancing:ap-northeast-1:xxxxxxxxxxx:targetgroup/hoge-target/zzzzzzzzzzzzz",
          "Name": "hoge-target",
          "Targets": [
            {
              "InstanceId": "i-xxxxxxxxxxxxx",
              "State": "healthy"
            }
          ]
        }
      ]
    }
  },
    :
    :
]
  ```

## [License](LICENSE)

The MIT License (MIT)