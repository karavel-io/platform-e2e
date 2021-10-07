import { CfnOutput, Construct, Duration, RemovalPolicy } from "@aws-cdk/core";
import { IOpenIdConnectProvider, PolicyDocument, PolicyStatement, Role } from "@aws-cdk/aws-iam";
import { BlockPublicAccess, Bucket, BucketAccessControl, BucketEncryption, StorageClass } from "@aws-cdk/aws-s3";
import { clusterPrincipalForServiceAccount, ServiceAccountRef } from "./iam";

export interface ObservabilityStorageProps {
  clusterOidcProvider: IOpenIdConnectProvider,
  bucketName: string,
  roleName: string,
  serviceAccounts: ServiceAccountRef[],
  deleteAfter?: Duration,
}

export class ObservabilityStorage extends Construct {
  constructor(scope: Construct, id: string, props: ObservabilityStorageProps) {
    super(scope, id);

    const {
      clusterOidcProvider,
      bucketName,
      roleName,
      serviceAccounts,
      deleteAfter,
    } = props;

    const bucket = new Bucket(this, 'Bucket', {
      bucketName,
      accessControl: BucketAccessControl.PRIVATE,
      blockPublicAccess: BlockPublicAccess.BLOCK_ALL,
      encryption: BucketEncryption.S3_MANAGED,
      removalPolicy: RemovalPolicy.DESTROY,
      lifecycleRules: [
        {
          transitions: [
            {
              transitionAfter: Duration.days(0),
              storageClass: StorageClass.INTELLIGENT_TIERING,
            },
          ],
          expiration: deleteAfter || Duration.days(365),
        },
      ],
    });

    const role = new Role(this, 'Role', {
      roleName,
      assumedBy: clusterPrincipalForServiceAccount(this, clusterOidcProvider, serviceAccounts),
      inlinePolicies: {
        ObservabilityStorageBucketReadWriteAccess: new PolicyDocument({
          statements: [
            new PolicyStatement({
              actions: [
                's3:ListBucket',
                's3:PutObject',
                's3:GetObject',
                's3:DeleteObject',
              ],
              resources: [
                bucket.bucketArn,
                bucket.arnForObjects('*'),
              ],
            }),
          ],
        }),
      },
    });

    new CfnOutput(this, 'BucketArn', {
      description: `ARN of the ${id} bucket`,
      value: bucket.bucketArn,
    });

    new CfnOutput(this, 'RoleArn', {
      description: `ARN of the ${id} role`,
      value: role.roleArn,
    });
  }
}
