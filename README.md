# README
## Run the container with
docker run -dit -v /var/lib/cfn_signaler/:/app/config -p 8080:8080 marjamis/cfn-signaler

## Manual call value
curl -X POST -d "send=FAILURE" <ip>:<port>/signal/
curl -X POST -d "send=SUCCESS" <ip>:<port>/signal/

## Sample Cloud Formation template
{
    "Resources": {
        "SimpleConfig": {
            "Type": "AWS::AutoScaling::LaunchConfiguration",
            "Properties": {
                "ImageId": "ami-9ff7e8af",
                "SecurityGroups": [
                    "sg-0c986c69",
                    "sg-0e986c6b"
                ],
                "KeyName" : "172.31.x.x-testing",
                "InstanceType": "t2.small",
                "IamInstanceProfile": "arn:aws:iam::<accountId>:instance-profile/cfn-signaler",
                "UserData": {
                    "Fn::Base64": {
                        "Fn::Join": [
                            "",
                            [
                                "#!/bin/bash -xe\n",
                                "yum update -y aws-cfn-bootstrap\n",
                                "yum install docker -y\n",
                                "chkconfig docker on\n",
                                "service docker start\n",
                                "mkdir /var/lib/cfn-signaler/ \n",
                                "echo '{\"LogicalResourceId\": \"ASG1\",\"StackName\":\"",
                                {
                                    "Ref": "AWS::StackName"
                                },
                                "\"}' > /var/lib/cfn-signaler/cfn-signaler.json \n",
                                "docker run -dit -v /var/lib/cfn-signaler/:/app/config -p 8080:8080 marjamis/cfn-signaler\n"
                            ]
                        ]
                    }
                }
            }
        },
        "ASG1": {
            "CreationPolicy": {
                "ResourceSignal": {
                    "Count": 1,
                    "Timeout": "PT15M"
                }
            },
            "UpdatePolicy": {
                "AutoScalingRollingUpdate": {
                    "MinInstancesInService": "1",
                    "MaxBatchSize": "1",
                    "PauseTime": "PT15M"
                }
            },
            "Type": "AWS::AutoScaling::AutoScalingGroup",
            "Properties": {
                "AvailabilityZones": {
                    "Fn::GetAZs": {
                        "Ref": "AWS::Region"
                    }
                },
                "LaunchConfigurationName": {
                    "Ref": "SimpleConfig"
                },
                "MaxSize": "1",
                "MinSize": "1"
            }
        }
    }
}

