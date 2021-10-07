import { Construct } from "@aws-cdk/core";
import { Secret } from "@aws-cdk/aws-secretsmanager";

export class ArgoCD extends Construct {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new Secret(this, 'Secret', {
      secretName: 'eks-e2e-cluster/argocd-secret',
      generateSecretString: {
        generateStringKey: 'placeholder',
        secretStringTemplate: '{}'
      }
    });

    new Secret(this, 'RepoSecret', {
      secretName: 'eks-e2e-cluster/argocd-pull-creds',
      generateSecretString: {
        generateStringKey: 'sshPrivateKey',
        secretStringTemplate: '{}'
      }
    });
  }
}

