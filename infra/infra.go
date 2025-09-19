package main

import (
	"log"
	"os"

	"github.com/aws/aws-cdk-go/awscdk/v2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsec2"
	"github.com/aws/aws-cdk-go/awscdk/v2/awsecs"
	"github.com/aws/aws-cdk-go/awscdk/v2/awselasticloadbalancingv2"
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

	// 2. Create the security groups for application load balancer with ingress rule on 8080 from anywhere and for Fargate service with ingress on ALB SG
	albSecurityGroup := awsec2.NewSecurityGroup(stack, jsii.String("ALBSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc: vpc,
		SecurityGroupName: jsii.String("ALB-SG"),
		Description: jsii.String("Security group enabling traffic from HTTP and HTTPS into MCP Server on Fargate"),
	})

	albSecurityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(80)), jsii.String("Allow access on port 80 (HTTP) from all IP"), jsii.Bool(true))
	albSecurityGroup.AddIngressRule(awsec2.Peer_AnyIpv4(), awsec2.Port_Tcp(jsii.Number(443)), jsii.String("Allow access on port 443 (HTTPS) from all IP"), jsii.Bool(true))
	
	fargateSecurityGroup := awsec2.NewSecurityGroup(stack, jsii.String("FargateSecurityGroup"), &awsec2.SecurityGroupProps{
		Vpc: vpc,
		SecurityGroupName: jsii.String("Fargate-MCP-Server-SG"),
		Description: jsii.String("Security group enabling traffic from albSecurityGroup"),
	})

	fargateSecurityGroup.AddIngressRule(albSecurityGroup, awsec2.Port_Tcp(jsii.Number(8080)), jsii.String("Allow from ALB-SG"), jsii.Bool(true))

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
	fargateService := awsecs.NewFargateService(stack, jsii.String("MCPFargateService"), &awsecs.FargateServiceProps{
		Cluster: cluster,
		TaskDefinition: fargateTaskDefintion,
		DesiredCount: jsii.Number(1),
		AssignPublicIp: jsii.Bool(true),
		SecurityGroups: &[]awsec2.ISecurityGroup {
			fargateSecurityGroup,
		},
	})

	// 9. Create ALB
	alb := awselasticloadbalancingv2.NewApplicationLoadBalancer(stack, jsii.String("ALB"), &awselasticloadbalancingv2.ApplicationLoadBalancerProps{
		Vpc: vpc,
		InternetFacing: jsii.Bool(true),
	})

	listener := alb.AddListener(jsii.String("Listener"), &awselasticloadbalancingv2.BaseApplicationListenerProps{
		Port: jsii.Number(80),
		Open: jsii.Bool(true),
	})

	listener.AddTargets(jsii.String("ApplicationFleet"), &awselasticloadbalancingv2.AddApplicationTargetsProps{
		Port: jsii.Number(8080),
		Targets: &[]awselasticloadbalancingv2.IApplicationLoadBalancerTarget{
			fargateService,
		},
		HealthCheck: &awselasticloadbalancingv2.HealthCheck{
			Path: jsii.String("/health"),
			Interval: awscdk.Duration_Seconds(jsii.Number(10)),
			Timeout: awscdk.Duration_Seconds(jsii.Number(5)),
			HealthyThresholdCount: jsii.Number(2),
			UnhealthyThresholdCount: jsii.Number(5),
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