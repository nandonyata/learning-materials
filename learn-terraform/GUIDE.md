# Terraform + AWS EC2: Beginner Guide

## Part 1: Key Terms in Plain Language

| Term | Simple explanation |
|---|---|
| **AMI (Amazon Machine Image)** | A ready-made template for a server. It contains the operating system (Ubuntu, Windows, etc.) and sometimes pre-installed software. Like a phone's factory-reset image. |
| **EC2 instance** | A virtual server (computer) running in AWS. `aws_instance` in Terraform creates one. |
| **Instance type** | How big/powerful the server is (e.g. `t3.micro` is small and cheap). |
| **VPC** | Your own private network inside AWS. Like the walls of an office building. |
| **Subnet** | A section of that network. Like a floor or room in the building. |
| **Security group** | A firewall for your server. Decides who can connect (e.g. allow web traffic on port 80). |
| **Region** | A geographic location for AWS, like `us-east-1` or `ap-southeast-1`. |
| **Availability zone (AZ)** | A separate data center within a region, so if one fails, others keep working. |
| **Provider** | A plugin that lets Terraform talk to a service (AWS, Azure, Google Cloud, etc.). |
| **State** | A file where Terraform remembers what it has created, so it knows what to change next time. |

AMI IDs look like `ami-0abc123...`. They differ per region and change often, which is why it is better to look them up with a data source instead of hard-coding them.

---

## Part 2: What Is a Data Source?

A **data source** lets you **read information that already exists** (outside your configuration, or in another Terraform state) and use it in your resources. It does not create or manage anything.

### Resource vs. data source

| | Resource | Data source |
|---|---|---|
| Keyword | `resource` | `data` |
| Purpose | Create, update, delete infrastructure | Read existing information |
| Managed by Terraform | Yes | No (read-only) |
| Reference syntax | `aws_instance.web.id` | `data.aws_ami.ubuntu.id` |

### Example

Look up the latest Ubuntu AMI and use it for an EC2 instance:

```hcl
data "aws_ami" "ubuntu" {
  most_recent = true
  owners      = ["099720109477"] # Canonical

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-jammy-22.04-amd64-server-*"]
  }
}

resource "aws_instance" "web" {
  ami           = data.aws_ami.ubuntu.id
  instance_type = "t3.micro"
}
```

Terraform first **looks up** the newest Ubuntu AMI (data source), then **creates a server** from it (resource).

### Common use cases

- Fetching the latest image/AMI ID
- Referencing an existing VPC, subnet, or security group created by someone else
- Reading availability zones for a region
- Getting the current account ID or region (`aws_caller_identity`, `aws_region`)
- Reading outputs from another Terraform project (`terraform_remote_state`)
- Pulling secrets or parameters (e.g. `aws_secretsmanager_secret_version`)

### Key points

- Data sources are read during `terraform plan` (or deferred to apply if they depend on values not yet known).
- If the thing you look up doesn't exist, the plan fails with an error.
- Each provider defines its own data sources, so check the provider docs.

---

## Part 3: Creating an EC2 Instance with Terraform (Step by Step)

### Step 1: Install what you need

- **Terraform**: download from terraform.io and check with `terraform -version`
- **AWS CLI**: to store your AWS login credentials
- An **AWS account**

### Step 2: Give Terraform access to AWS

- In AWS, create an **IAM user** with permission to create EC2 (for learning, `AmazonEC2FullAccess`)
- Create an **Access Key** for that user
- On your computer, run `aws configure` and enter the access key, secret key, and region
- Terraform will automatically use these credentials

### Step 3: Create a project folder

```bash
mkdir my-first-ec2
cd my-first-ec2
```

Create a file called `main.tf`.

### Step 4: Write the configuration

You need 3 basic parts.

**a) Provider**: tells Terraform to use AWS and which region

```hcl
provider "aws" {
  region = "ap-southeast-1"
}
```

**b) Resource**: the thing you want to create

```hcl
resource "aws_instance" "my_server" {
  ami           = "ami-xxxxxxxx"   # the OS template
  instance_type = "t3.micro"       # the size

  tags = {
    Name = "my-first-server"
  }
}
```

**c) (Optional) Output**: prints useful info after creation

```hcl
output "server_ip" {
  value = aws_instance.my_server.public_ip
}
```

### Step 5: Run the 4 main commands

1. `terraform init`: downloads the AWS provider plugin
2. `terraform plan`: shows what *will* be created (a preview, nothing is created yet)
3. `terraform apply`: actually creates it (type `yes` to confirm)
4. `terraform destroy`: deletes everything when you're done, so you don't get charged

---

## Part 4: Where to Check Which Fields You Need

The **Terraform Registry documentation** is your main reference.

1. Go to **registry.terraform.io**
2. Click **Providers**, then **hashicorp/aws**, then **Documentation**
3. Search for the resource name, e.g. `aws_instance`
4. On that page, look at these sections:
   - **Example Usage**: copy and adapt this first
   - **Argument Reference**: the fields *you write*. Each is marked **Required** or **Optional**
   - **Attribute Reference**: the values *you can read* afterwards (like `id`, `public_ip`)

**Tip:** For `aws_instance`, the only truly required fields are usually `ami` and `instance_type`. Everything else is optional.

The same pattern works for any resource: search the name in the docs, read the example, then check the required arguments.

### Where to find values like the AMI ID

- Hard-coding is risky since AMI IDs differ per region and change over time
- Better: use the `data "aws_ami"` data source (see Part 2) to look up the latest one automatically

---

## Beginner Tips

- Always run `terraform plan` before `terraform apply`
- Always run `terraform destroy` after practicing to avoid unexpected AWS bills
- `terraform fmt` auto-formats your code and `terraform validate` checks for mistakes
- Never put your AWS keys inside `.tf` files