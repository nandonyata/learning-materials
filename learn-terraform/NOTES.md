# Terraform
We will be learning terraform with LocalStack so we don't use any money for real AWS services here.

Youtube source - https://www.youtube.com/watch?v=pnzlqoYNuQc

# Install LocalStack
cloud: https://app.localstack.cloud/getting-started
run the docker compose proivded and then:
```
npm install -g @localstack/lstk
lstk login
```

# Install Terraform
MacOS:
```
brew tap hashicorp/tap
brew install hashicorp/tap/terraform
```
doc source - https://developer.hashicorp.com/terraform/install

---

# For the Web UI Endpoint
Once restarted, in the LocalStack web app at app.localstack.cloud/instances, use:
```
http://<your-wsl-ip>:4566
```

# Find your WSL IP by running:
```
hostname -I
```
This will give you an IP like 172.x.x.x that your Windows browser can reach.
Alternatively, if host.docker.internal doesn't work in your setup, the WSL IP approach should connect properly.

---

# Add user IAM in the web or cli
CLI:
```
lstk aws iam create-user --user-name terraform_user

# OR

aws iam create-user \
  --user-name terraform_user \
  --endpoint-url http://172.24.34.66:4566 \
  --region us-east-1

```
Web - https://app.localstack.cloud/inst/default/resources/iam

# Run some commands accordingly
```
lstk terraform init # might be not needed for existing project
lstk terraform fmt
lstk terraform validate
lstk terraform plan
lstk terraform destroy
lstk terraform apply
```