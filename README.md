# README
This is a basic method to signal back to the Creation or Update Policy of an SG within a CFN(Cloud Formations) Stack. In this way you can test how this works in certain cirumstances against a test application.

This program is essentially a webserver which allows you to signal a success or failure depending on what you're attempting to test  via a website or a POST request to the IP address of the EC2 Instance that is apart of the ASG.

Issues or Feature Requests are welcome.

	NOTE: The master branch is automatically created into a docker image which can be accessed publicly. This image is located at: https://hub.docker.com/r/marjamis/cfn-signaler/

## Run the container
    docker run -dit -e STACKNAME=<stackname> -e LOGICALID=<logicalid> -p <PublicPort>:8080 marjamis/cfn-signaler

## Run the binary directly
    LOGICALID=<logical_id_of_ASG> STACKNAME=<stack_name> go run main.go

## Manual curl against the endpoint
    curl -X POST -d "send=FAILURE" <ip>:<port>/signal/
    curl -X POST -d "send=SUCCESS" <ip>:<port>/signal/

## Sample Cloud Formation Template
    {
    "Parameters": {
        "AMI": {
            "Type": "AWS::EC2::Image::Id"
        },
        "SecurityGroups": {
            "Description": "Make sure port used is open.",
            "Type": "CommaDelimitedList"
        },
        "SSHKeyPair": {
            "Type": "AWS::EC2::KeyPair::KeyName"
        },
        "InstanceType": {
            "Type": "String",
            "Default": "t2.micro"
        },
        "InstanceProfile": {
            "Type": "String"
        },
        "PublicPort": {
            "Type": "String"
        },
        "InstanceName": {
            "Type": "String"
        }
    },
    "Resources": {
        "SimpleConfig": {
            "Type": "AWS::AutoScaling::LaunchConfiguration",
            "Properties": {
                "ImageId": {
                    "Ref": "AMI"
                },
                "SecurityGroups": {
                    "Ref": "SecurityGroups"
                },
                "KeyName": {
                    "Ref": "SSHKeyPair"
                },
                "InstanceType": {
                    "Ref": "InstanceType"
                },
                "IamInstanceProfile": {
                    "Ref": "InstanceProfile"
                },
                "UserData": {
                    "Fn::Base64": {
                        "Fn::Join": [
                            "", [
                                "#!/bin/bash -xe\n",
                                "yum update -y aws-cfn-bootstrap\n",
                                "yum install docker -y\n",
                                "chkconfig docker on\n",
                                "service docker start\n",
                                "docker run -dit -e LOGICALID=ASG1 -e STACKNAME=", {
                                    "Ref": "AWS::StackName"
                                },
                                " -p ", {
                                    "Ref": "PublicPort"
                                }, ":8080 marjamis/cfn-signaler\n"
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
                "MaxSize": "2",
                "MinSize": "1",
                "Tags": [{
                    "Key": "Name",
                    "Value": {
                        "Ref": "InstanceName"
                    },
                    "PropagateAtLaunch": "true"
                }]
            }
        }
    }
    }
