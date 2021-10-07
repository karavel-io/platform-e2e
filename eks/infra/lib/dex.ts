import { Construct } from "@aws-cdk/core";
import { Secret } from "@aws-cdk/aws-secretsmanager";

export class Dex extends Construct {
  constructor(scope: Construct, id: string) {
    super(scope, id);

    new Secret(this, 'Secret', {
      secretName: 'eks-e2e-cluster/dex-secret',
      generateSecretString: {
        generateStringKey: 'GITHUB_CLIENT_SECRET',
        secretStringTemplate: '{"GITHUB_CLIENT_ID":"change-me"}'
      }
    })
  }
}

