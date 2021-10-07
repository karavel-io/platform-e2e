version="unstable"

component "cert-manager" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "cert-manager"

  # Params

  letsencrypt = {
    email = "tech@neosperience.com"
  }
}

component "external-secrets" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "external-secrets"

  # Params

  defaultProvider = "aws"

  aws = {
    enable = true
    region = "eu-west-1"
    eksRole = "arn:aws:iam::568706034652:role/KaravelE2eExternalSecretsRole"
  }
}

component "external-dns" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "external-dns"

  # Params

  domainFilter = "eks.e2e.karavel.io"

  provider = "route53"
  route53 = {
    zoneId = "Z02598523M9WR611ST687"
    eksRole = "arn:aws:iam::568706034652:role/KaravelE2eExternalDnsRole"
  }
}

component "ingress-nginx" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "ingress-nginx"
}

component "goldpinger" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "monitoring"
}

component "dex" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "dex"

  # Params
  publicURL = "https://auth.eks.e2e.karavel.io"

  connectors = [
    {
      type = "github"
      id = "github"
      name = "GitHub"
      config = {
        clientID = "$GITHUB_CLIENT_ID"
        clientSecret = "$GITHUB_CLIENT_SECRET"
        redirectURI = "https://auth.eks.e2e.karavel.io/callback"
        orgs = [
          {
            name = "karavel-io",
            teams = ["platform"]
          }
        ]
        teamNameField = "slug",
        useLoginAsId = false
      }
    }
  ]

  secret = {
    key = "eks-e2e-cluster/dex-secret"
  }
}

component "argocd" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "argocd"

  # Params
  publicURL = "https://deploy.eks.e2e.karavel.io"

  adminGroup = "platform"

  git = {
    repo = "git@github.com:karavel-io/platform-e2e.git"
  }

  secret = {
    key = "eks-e2e-cluster/argocd-secret"
  }

  credentialsSecret = {
    key = "eks-e2e-cluster/argocd-pull-creds"
  }
}

component "grafana" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "monitoring"

  # Params

  publicURL = "https://grafana.eks.e2e.karavel.io"

  dex = {
    issuer = "https://auth.eks.e2e.karavel.io"
  }
}

component "prometheus" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "monitoring"

  # Params

  store = "s3"
  s3 = {
    bucket = "karavel-eks-e2e-cluster-metrics-prometheus"
    region = "eu-west-1"
    eksRole = "arn:aws:iam::568706034652:role/KaravelE2ePrometheusRole"
  }
}

component "loki" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "monitoring"

  # Params

  store = "s3"
  s3 = {
    bucket = "karavel-eks-e2e-cluster-logging-loki"
    region = "eu-west-1"
    eksRole = "arn:aws:iam::568706034652:role/KaravelE2eLokiRole"
  }
}

//component "tempo" {
//  version = "unstable:0.1.0-SNAPSHOT"
//  namespace = "monitoring"
//
//  # Params
//
//  store = "s3"
//  s3 = {
//    bucket = "karavel-eks-e2e-cluster-tracing-tempo"
//    region = "eu-west-1"
//    eksRole = "arn:aws:iam::568706034652:role/KaravelE2eTempoRole"
//  }
//}

component "velero" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "velero"

  backups = {
    provider = "aws"
    s3 = {
      bucket = "karavel-eks-e2e-cluster-backups"
      region = "eu-west-1"
    }
  }

  snapshots = {
    provider = "aws"
    region = "eu-west-1"
  }

  aws = {
    eksRole = "arn:aws:iam::568706034652:role/KaravelE2eVeleroRole"
  }
}

component "olm" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "olm"
}
