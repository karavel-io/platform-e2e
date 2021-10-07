#!/usr/bin/env node
import 'source-map-support/register';
import { App } from '@aws-cdk/core';
import { E2eStack } from '../lib/stack';
import { KubernetesVersion } from "@aws-cdk/aws-eks";

const app = new App();

new E2eStack(app, 'E2eStack', {
  kubernetesVersion: KubernetesVersion.V1_21,
});
