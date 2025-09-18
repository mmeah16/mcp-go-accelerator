# Welcome to your CDK Go project!

This is a blank project for CDK development with Go.

The `cdk.json` file tells the CDK toolkit how to execute your app.

## Useful commands

- `cdk deploy` deploy this stack to your default AWS account/region
- `cdk diff` compare deployed stack with current state
- `cdk synth` emits the synthesized CloudFormation template
- `go test` run unit tests

# ECS Fargate Prototype Infrastructure

This project provisions an **ECS Fargate service** and supporting infrastructure using the **AWS CDK for Go**.  
It is designed as a prototype to demonstrate a containerized MCP Server deployment with minimal setup.

---

## What It Deploys

- **VPC & Networking**

  - Uses default VPC and subnets from context
  - Security groups configured for task networking

- **ECS Fargate Service**

  - Task definition and service
  - Task role with CloudWatch Logs permissions

- **CloudWatch Logs**
  - Automatic log groups created for ECS tasks

---

## Getting Started

1. Install dependencies:

   ```bash
   go mod tidy
   ```

2. Bootstap CDK:

   ```bash
   cdk bootstrap
   ```

3. Synthesize CloudFormation template:

   ```bash
   cdk synth
   ```

4. Deploy the stack:
   ```bash
   cdk deploy
   ```

## Files

- `infra/infra.go` – Main CDK stack definition
- `cdk.json` – CDK configuration (entry point)
- `go.mod`, `go.sum` – Go module dependencies
- `infra_test.go` – Basic CDK unit test (auto-generated)
- `README.md` – Documentation for the project

## Notes

- This stack is a **prototype**: it gets the MCP server running on ECS quickly.
- Currently, the ECS tasks use a **public IP** for access.
- Security groups are open to all IPs on port 8080 — suitable for testing but **not production-ready**.
- The task role has minimal CloudWatch permissions scoped to the log group.

---

## 📖 Roadmap (Next Steps)

1. **Secure Networking**

   - Restrict security groups to only trusted IPs or ALB.
   - Consider private subnets + NAT for outbound access.

2. **Load Balancer**

   - Add an Application Load Balancer (ALB) in public subnets.
   - Configure target group and health checks for ECS tasks.

3. **Custom Domain & HTTPS**

   - Set up Route 53 hosted zone and ACM certificate.
   - Attach ALB listener on HTTPS (443) with certificate.
   - Create DNS alias to the ALB.
