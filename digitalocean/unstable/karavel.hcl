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

  defaultProvider = "vault"

  vault = {
    enable = true
  }
}

component "external-dns" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "external-dns"

  # Params

  provider = "route53"
  route53 = {
    zoneId = "Z02598523M9WR611ST687"
    iamRole = "karavel-e2e-testing-external-dns"
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
  publicURL = "https://auth.do.e2e.karavel.io"

#  connectors = [
#    {
#      type = "github"
#      id = "github"
#      name = "GitHub"
#      config = {
#        clientID = "$GITHUB_CLIENT_ID"
#        clientSecret = "$GITHUB_CLIENT_SECRET"
#        redirectURI = "https://auth.do.e2e.karavel.io/callback"
#        orgs = [
#          {
#            name = "karavel-io",
#            teams = ["platform"]
#          }
#        ]
#        teamNameField = "slug",
#        useLoginAsId = false
#      }
#    }
#  ]
#
#  secret = {
#    key = "do-e2e-cluster/dex-secret"
#  }
}

component "argocd" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "argocd"

  # Params
  publicURL = "https://deploy.do.e2e.karavel.io"

  adminGroup = "karavel-io:platform"

  git = {
    repo = "git@github.com:karavel-io/platform-e2e.git"
  }

  secret = {
    key = "do-e2e-cluster/argocd-secret"
  }

  credentialsSecret = {
    key = "do-e2e-cluster/argocd-pull-creds"
  }

  dex = {
    issuer = "https://auth.do.e2e.karavel.io"
  }
}

component "grafana" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "monitoring"

  # Params

  publicURL = "https://grafana.do.e2e.karavel.io"

  dex = {
    issuer = "https://auth.do.e2e.karavel.io"
  }
}

component "prometheus" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "monitoring"

  # Params

  store = "s3"
  s3 = {
    bucket = "karavel-do-e2e-metrics-prometheus"
    region = "eu-west-1"
    doRole = "arn:aws:iam::568706034652:role/KaravelE2ePrometheusRole"
  }
}

component "loki" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "monitoring"

  # Params

  store = "s3"
  s3 = {
    bucket = "karavel-do-e2e-logging-loki"
    region = "eu-west-1"
    doRole = "arn:aws:iam::568706034652:role/KaravelE2eLokiRole"
  }
}

component "tempo" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "monitoring"

  # Params

  store = "s3"
  s3 = {
    bucket = "karavel-do-e2e-cluster-tracing-tempo"
    region = "eu-west-1"
    doRole = "arn:aws:iam::568706034652:role/KaravelE2eTempoRole"
  }
}

component "velero" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "velero"

  backups = {
    provider = "aws"
    s3 = {
      bucket = "karavel-do-e2e-backups"
      region = "eu-west-1"
    }
  }

  snapshots = {
    provider = "aws"
    region = "eu-west-1"
  }

  aws = {
    doRole = "arn:aws:iam::568706034652:role/KaravelE2eVeleroRole"
  }
}

component "olm" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "olm"
}

component "vault" {
  version = "unstable:0.1.0-SNAPSHOT"
  namespace = "vault"

  publicURL = "https://vault.do.e2e.karavel.io"

  dex = {
    enable = false
    issuer = "https://auth.do.e2e.karavel.io"
    adminGroup = "karavel-io:platform"
  }
}
