import { CfnJson, Construct } from "@aws-cdk/core";
import { IOpenIdConnectProvider } from "@aws-cdk/aws-iam/lib/oidc-provider";
import { IPrincipal, OpenIdConnectPrincipal } from "@aws-cdk/aws-iam";

export interface ServiceAccountRef {
  namespace: string;
  name: string;
}

export function clusterPrincipalForServiceAccount(scope: Construct, clusterOidcProvider: IOpenIdConnectProvider, serviceAccounts: ServiceAccountRef[], id = 'ProviderIssuer'): IPrincipal {
  const iss = clusterOidcProvider.openIdConnectProviderIssuer;
  return new OpenIdConnectPrincipal(clusterOidcProvider, {
    StringEquals: new CfnJson(scope, `${id}RoleEqualsConditions`, {
      value: {

        [`${iss}:aud`]: 'sts.amazonaws.com',
      },
    }),
    'ForAnyValue:StringLike': new CfnJson(scope, `${id}RoleLikeConditions`, {
      value: {
        [`${iss}:sub`]: serviceAccounts.map(({ namespace, name }) => `system:serviceaccount:${namespace}:${name}`),
      },
    }),
  });
}
