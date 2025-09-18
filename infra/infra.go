package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsiam"
	"github.com/aws/aws-cdk-go/awscdk/v2/awslogs"
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/joho/godotenv"
)

type MCPStackProps struct {
	awscdk.StackProps
}

func NewInfraStack(scope constructs.Construct, id string, props *MCPStackProps) awscdk.Stack {
	var sprops awscdk.StackProps
	if props != nil {
		sprops = props.StackProps
	}
	stack := awscdk.NewStack(scope, &id, &sprops)

	// 1. Obtain the Default VPC
	vpc := awsec2.Vpc_FromLookup(stack, jsii.String("DefaultVPC"), &awsec2.VpcLookupOptions{
		IsDefault: jsii.Bool(true),
	})

	// 2. Create the security group and ingress rule on 8080 from anywhere
	securityGroup := awsec2.NewSecurityGroup(stack, jsii.String("SecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc: vpc,
	})

	securityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(8080)), jsii.String("Allow access on port 8080 on all IP"), jsii.Bool(true))
	
	// 3. Create Log Group
	logGroup := awslogs.NewLogGroup(stack, jsii.String("MCPLogGroup"), &awslogs.LogGroupProps{
		LogGroupName: jsii.String("/ecs/mcp-service"),
		Retention: awslogs.RetentionDays_ONE_WEEK,
	})

	// 4. Create Task Role with access to CloudWatch Logs
	taskRole := awsiam.NewRole(stack, jsii.String("MCPServerTaskRole"), &awsiam.RoleProps{
		AssumedBy: awsiam.NewServicePrincipal(jsii.String("ecs-tasks.amazonaws.com"), nil),
	})

	taskRole.AddToPolicy(awsiam.NewPolicyStatement(&awsiam.PolicyStatementProps{
		Resources: &[]*string {
			logGroup.LogGroupArn(),
		},
		Actions: &[]*string{
			jsii.String("logs:CreateLogStream"),
			jsii.String("logs:PutLogEvents"),
		},
	}))

	// 5. Create the ECS cluster
	cluster := awsecs.NewCluster(stack, jsii.String("MCPCluster"), &awsecs.ClusterProps{
		ClusterName: jsii.String("MCPCluster"),
		Vpc: vpc,
	})

	// 6. Create Fargate Task Definition
	fargateTaskDefintion := awsecs.NewFargateTaskDefinition(stack, jsii.String("MCPFargateTaskDefinition"), &awsecs.FargateTaskDefinitionProps {
		Cpu: jsii.Number(256),
		MemoryLimitMiB: jsii.Number(512),
	})

	// 7. Create Container using fargateTaskDefintion and add port mapping to 8080
	container := fargateTaskDefintion.AddContainer(jsii.String("MCPContainer"), &awsecs.ContainerDefinitionOptions {
		Image: awsecs.AssetImage_FromAsset(jsii.String("../"), &awsecs.AssetImageProps{
			File: jsii.String("Dockerfile"),
		}),
		MemoryLimitMiB: jsii.Number(512),
		Logging: awsecs.LogDrivers_AwsLogs(&awsecs.AwsLogDriverProps{
			StreamPrefix: jsii.String("mcp-server"),
			LogGroup: logGroup,
		}),
	})

	container.AddPortMappings(&awsecs.PortMapping{
		ContainerPort: jsii.Number(8080),
	})

	// 8. Fargate Service
	awsecs.NewFargateService(stack, jsii.String("MCPFargateService"), &awsecs.FargateServiceProps{
		Cluster: cluster,
		TaskDefinition: fargateTaskDefintion,
		DesiredCount: jsii.Number(1),
		AssignPublicIp: jsii.Bool(true),
		SecurityGroups: &[]awsec2.ISecurityGroup {
			securityGroup,
		},
	})

	return stack
}

func main() {
	defer jsii.Close()

	app := awscdk.NewApp(nil)

	NewInfraStack(app, "MCPStack", &MCPStackProps{
		awscdk.StackProps{
			Env: env(),
		},
	})

	app.Synth(nil)
}


func env() *awscdk.Environment {
	err := godotenv.Load()

	if err != nil {
		log.Fatal("Could not retrieve environment variables.")
	}

	return &awscdk.Environment{
		Account: jsii.String(os.Getenv("AWS_ACCOUNT_ID")),
		Region:  jsii.String(os.Getenv("REGION")),
	}
}