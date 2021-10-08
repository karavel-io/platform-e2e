# EKS end-to-end testing environment

This folder contains deployment code for testing the [Karavel Container Platform] on [Amazon EKS].

## Deploying the infrastructure

The cloud infrastructure is deployed using [CDK]. NodeJS 14 and [pnpm] are required. You should have the AWS CLI already
configured to connect to the target AWS account. Verify it by running `aws sts get-caller-identity`.

```bash
$ pwd
/path/to/repo/eks

$ cd infra && pnpm install
$ pnpm cdk bootstrap # only needed once per AWS account

$ pnpm cdk diff # verify that everything looks ok

$ pnpm cdk deploy # wait for a very long time for the EKS cluster to come up.
```

When the CloudFormation stack is green, fill in the secrets information on [AWS Secrets Manager]. The required secrets
should already be available with placeholder information.

Currently, you need the following credential sets:
- an SSH key pair for ArgoCD to pull this repository
- a GitHub OAuth app for Dex to authenticate services

### ArgoCD key pair

Generate a new SSH key pair with the following command: `ssh-keygen -t ed25519 -C "Argo CD"`.
Create the keys somewhere handy (we'll be deleting them in a few moments), then copy the secret key to AWS Secrets Manager
in the `eks-e2e-cluster/argocd-pull-creds` secret. The key should be base64-encoded, on UNIX systems this is as easy as:
`cat argo-secret-key.pem | base64 -w0`.

The public key (`argo-secret-key.pem.pub` in our example) should be added to the GitHub repository as a [deploy key](https://github.com/karavel-io/platform-e2e/settings/keys).

### Dex OAuth App

Create a new Client Secret (or ask an administrator to create one for you) for [this OAuth App](https://github.com/organizations/karavel-io/settings/applications/1728020).
Copy the Client ID and Client Secret to AWS Secrets Manager in the `eks-e2e-cluster/dex-secret` secret.

### Connecting to the cluster

Navigate the AWS console to the CloudFormation stack you just created, and check its outputs. Look for an output called `ClusterConfigCommand[random characters]`
and run the CLI command it contains. This will configure `kubectl` to access the test cluster. Validate that you have access by
running `kubectl get nodes`.

## Deploying the platform

Choose the platform version you want to test (or create a new folder to add a new version), like [unstable](unstable), and enter it.

From this point on, follow the regular platform install procedures highlighted in the [documentation](https://platform.karavel.io/docs).

Long story short:

```bash
$ cd unstable # or any other target version
$ karavel render
$ kustomize build | kubectl apply -f -
$ kustomize build | kubectl apply -f - # run it a second time to work around some race conditions between resources in kubectl
```

Give the platform few minutes to bootstrap, and you should be able to access ArgoCD on https://deploy.eks.e2e.karavel.io. 
The rest of the components should come up soon after.

## Teardown

When you're done, simply delete the CloudFormation stack from the AWS console or run `pnpm cdk destroy` from the [infra](infra)
folder. Some resources, like S3 buckets, may fail to delete because they're busy or not empty. Mark them down and delete them
manually from the AWS console. For S3 buckets, remember to first empty them, then delete them. Thanks Amazon!

[Karavel Container Platform]: https://platform.karavel.io
[Amazon EKS]: https://aws.amazon.com/eks/
[CDK]: https://docs.aws.amazon.com/cdk/latest/guide/home.html
[pnpm]: https://pnpm.io
[AWS Secrets Manager]: https://aws.amazon.com/secrets-manager
