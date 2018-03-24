# cfn-signaler
This is a basic method to signal back to the Creation/Update Policy of an ASG within a CFN(CloudFormation) Stack. In this way you can test how this works in certain circumstances against a test application.

This program is essentially a web server which allows you to signal a success or failure, depending on what you're attempting to test, via a website or a POST request to the IP address of the EC2 Instance that is apart of the ASG.

Issues or Feature Requests are welcome.

**NOTE:** The master branch is automatically built into a docker image which can be accessed publicly at: https://hub.docker.com/r/marjamis/cfn-signaler/

## Run the application directly with go
```bash
make run -e LOGICALID=<logical_id_of_ASG> -e STACKNAME=<stack_name>
```

## Build and run the application via a docker container
```bash
make dbuild
make drun -e LOGICALID=<logical_id_of_ASG> -e STACKNAME=<stack_name> -e PUBLICPORT=<public_port>
```

## Manual curl against the endpoint
```bash
curl -X POST -d "send=FAILURE" <ip>:<port>/signal/
curl -X POST -d "send=SUCCESS" <ip>:<port>/signal/
```

## Sample CloudFormation Template
```yaml
Parameters:
  AMI:
    Type: AWS::EC2::Image::Id
  SecurityGroups:
    Description: Make sure port used is open.
    Type: CommaDelimitedList
  SSHKeyPair:
    Type: AWS::EC2::KeyPair::KeyName
  InstanceType:
    Type: String
    Default: t2.micro
  InstanceProfile:
    Type: String
  PublicPort:
    Type: String
  InstanceName:
    Type: String
Resources:
  SimpleConfig:
    Type: AWS::AutoScaling::LaunchConfiguration
    Properties:
      ImageId: !Ref AMI
      SecurityGroups: !Ref SecurityGroups
      KeyName: !Ref SSHKeyPair
      InstanceType: !Ref InstanceType
      IamInstanceProfile: !Ref InstanceProfile
      UserData:
        Fn::Base64:
          !Sub |
            #!/bin/bash -xe
            yum update -y aws-cfn-bootstrap
            yum install docker -y
            chkconfig docker on
            service docker start
            docker run -dit -e LOGICALID=ASG1 -e STACKNAME=${AWS::StackName} -p ${PublicPort}:8080 marjamis/cfn-signaler
  ASG1:
    CreationPolicy:
      ResourceSignal:
        Count: 1
        Timeout: PT15M
    UpdatePolicy:
      AutoScalingRollingUpdate:
        MinInstancesInService: 1
        MaxBatchSize: 1
        PauseTime: PT15M
    Type: AWS::AutoScaling::AutoScalingGroup
    Properties:
      AvailabilityZones:
        Fn::GetAZs: !Ref AWS::Region
      LaunchConfigurationName: !Ref SimpleConfig
      MaxSize: 2
      MinSize: 1
      Tags:
        - Key: Name
          Value: !Ref InstanceName
          PropagateAtLaunch: true
```
