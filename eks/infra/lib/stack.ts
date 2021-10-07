import { Construct, Stack, StackProps, Tags } from '@aws-cdk/core';
import {
  GatewayVpcEndpoint, GatewayVpcEndpointAwsService,
  InstanceClass,
  InstanceSize,
  InstanceType,
  NatProvider,
  PrivateSubnet,
  PublicSubnet,
  Vpc
} from "@aws-cdk/aws-ec2";
import { Cluster, DefaultCapacityType, KubernetesVersion } from "@aws-cdk/aws-eks";
import { Dex } from "./dex";
import { ArgoCD } from "./argocd";
import { ObservabilityStorage } from "./observability-storage";
import { ExternalSecrets } from "./external-secrets";
import { ExternalDns } from "./external-dns";
import { HostedZone } from "@aws-cdk/aws-route53";

export interface E2eStackProps extends StackProps {
  kubernetesVersion: KubernetesVersion,
}

export class E2eStack extends Stack {
  constructor(scope: Construct, id: string, props: E2eStackProps) {
    super(scope, id, props);

    const vpc = new Vpc(this, 'Vpc', {
      natGateways: 1,
      natGatewayProvider: NatProvider.gateway(),
    });

    vpc.selectSubnets({}).subnets.forEach((subnet) => {
      const tags = Tags.of(subnet);
      tags.add(`kubernetes.io/cluster/KaravelE2e`, '1');

      if (subnet instanceof PublicSubnet) {
        tags.add('kubernetes.io/role/elb', '1');
      }

      if (subnet instanceof PrivateSubnet) {
        tags.add('kubernetes.io/role/internal-elb', '1');
      }
    });

    new GatewayVpcEndpoint(this, "S3Endpoint", {
      service: GatewayVpcEndpointAwsService.S3,
      vpc: vpc,
    });

    new GatewayVpcEndpoint(this, "DynamoEndpoint", {
      service: GatewayVpcEndpointAwsService.DYNAMODB,
      vpc: vpc,
    });

    const cluster = new Cluster(this, 'Cluster', {
      clusterName: 'KaravelE2e',
      version: props.kubernetesVersion,
      vpc,
      defaultCapacity: 3,
      defaultCapacityInstance: InstanceType.of(InstanceClass.M5A, InstanceSize.LARGE),
      defaultCapacityType: DefaultCapacityType.NODEGROUP,
    });

    new ArgoCD(this, 'ArgoCD')

    new Dex(this, 'Dex');

    const clusterOidcProvider = cluster.openIdConnectProvider;

    new ObservabilityStorage(this, 'Prometheus', {
      clusterOidcProvider,
      bucketName: 'karavel-eks-e2e-metrics-prometheus',
      roleName: 'KaravelE2ePrometheusRole',
      serviceAccounts: [{ namespace: 'monitoring', name: 'prometheus-operated' }, {
        namespace: 'monitoring',
        name: 'thanos'
      }],
      ...props,
    });

    new ObservabilityStorage(this, 'Loki', {
      clusterOidcProvider,
      bucketName: 'karavel-eks-e2e-logging-loki',
      roleName: 'KaravelE2eLokiRole',
      serviceAccounts: [{ namespace: 'monitoring', name: 'loki' }],
      ...props,
    });

    new ObservabilityStorage(this, 'Tempo', {
      clusterOidcProvider,
      bucketName: 'karavel-eks-e2e-tracing-tempo',
      roleName: 'KaravelE2eTempoRole',
      serviceAccounts: [{ namespace: 'monitoring', name: 'tempo' }],
      ...props,
    });

    new ExternalSecrets(this, 'ExternalSecrets', {
      clusterOidcProvider,
      roleName: 'KaravelE2eExternalSecretsRole',
    });

    new ExternalDns(this, 'ExternalDns', {
      clusterOidcProvider,
      roleName: 'KaravelE2eExternalDnsRole',
      hostedZone: HostedZone.fromHostedZoneId(this, 'HostedZone', 'Z02598523M9WR611ST687'),
    });
  }
}
