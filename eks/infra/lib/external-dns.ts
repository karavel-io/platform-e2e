import { Construct } from "@aws-cdk/core";
import { IOpenIdConnectProvider, PolicyDocument, PolicyStatement, Role } from "@aws-cdk/aws-iam";
import { IHostedZone } from "@aws-cdk/aws-route53";
import { clusterPrincipalForServiceAccount } from "./iam";

export interface ExternalDnsProps {
  clusterOidcProvider: IOpenIdConnectProvider,
  hostedZone: IHostedZone,
  roleName: string,
}

export class ExternalDns extends Construct {
  constructor(scope: Construct, id: string, props: ExternalDnsProps) {
    super(scope, id);

    const { clusterOidcProvider, hostedZone, roleName } = props;

    new Role(this, 'Role', {
      roleName,
      assumedBy: clusterPrincipalForServiceAccount(this, clusterOidcProvider, [{ namespace: 'external-dns', name: 'external-dns' }]),
      inlinePolicies: {
        Route53ReadWriteAccess: new PolicyDocument({
          statements: [
            new PolicyStatement({
              actions: [
                'route53:ChangeResourceRecordSets',
              ],
              resources: [
                hostedZone.hostedZoneArn,
              ],
            }),
            new PolicyStatement({
              actions: [
                'route53:ListHostedZones',
                'route53:ListResourceRecordSets',
              ],
              resources: [
                '*',
              ],
            }),
          ],
        }),
      },
    });
  }
}
