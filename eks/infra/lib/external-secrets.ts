import { Arn, Construct, Stack } from "@aws-cdk/core";
import { IOpenIdConnectProvider, PolicyDocument, PolicyStatement, Role } from "@aws-cdk/aws-iam";
import { clusterPrincipalForServiceAccount } from "./iam";

export interface ExternalSecretsProps {
  clusterOidcProvider: IOpenIdConnectProvider;
  roleName: string;
}

export class ExternalSecrets extends Construct {
  constructor(scope: Construct, id: string, props: ExternalSecretsProps) {
    super(scope, id);

    const { clusterOidcProvider, roleName } = props;

    new Role(this, 'Role', {
      roleName,
      assumedBy: clusterPrincipalForServiceAccount(this, clusterOidcProvider, [{
        namespace: 'external-secrets',
        name: 'external-secrets',
      }]),
      inlinePolicies: {
        SecretsManagerReadOnlyAccess: new PolicyDocument({
          statements: [
            new PolicyStatement({
              actions: [
                'secretsmanager:GetResourcePolicy',
                'secretsmanager:GetSecretValue',
                'secretsmanager:DescribeSecret',
                'secretsmanager:ListSecretVersionIds',
              ],
              resources: [
                Arn.format({
                  service: 'secretsmanager',
                  resource: 'secret',
                  sep: ':',
                  resourceName: '*',
                }, Stack.of(this)),
              ],
            }),
          ],
        }),
      },
    });
  }
}
