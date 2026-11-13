data "aws_caller_identity" "current" {}

data "aws_region" "current" {}

locals {
  account_id = data.aws_caller_identity.current.account_id
  region     = data.aws_region.current.region
  prefix     = "arn:aws:%s:${local.region}:${local.account_id}"
  role_name  = "${var.name}-github-deploy"
}

resource "aws_iam_openid_connect_provider" "github" {
  count = var.create_oidc_provider ? 1 : 0

  url            = "https://token.actions.githubusercontent.com"
  client_id_list = ["sts.amazonaws.com"]
}

data "aws_iam_openid_connect_provider" "github" {
  count = var.create_oidc_provider ? 0 : 1

  url = "https://token.actions.githubusercontent.com"
}

locals {
  oidc_provider_arn = var.create_oidc_provider ? aws_iam_openid_connect_provider.github[0].arn : data.aws_iam_openid_connect_provider.github[0].arn
}

data "aws_iam_policy_document" "trust" {
  statement {
    actions = ["sts:AssumeRoleWithWebIdentity"]

    principals {
      type        = "Federated"
      identifiers = [local.oidc_provider_arn]
    }

    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:aud"
      values   = ["sts.amazonaws.com"]
    }

    # Only workflows running on main can deploy.
    condition {
      test     = "StringEquals"
      variable = "token.actions.githubusercontent.com:sub"
      values   = ["repo:${var.github_repository}:ref:refs/heads/main"]
    }
  }
}

resource "aws_iam_role" "this" {
  name               = local.role_name
  assume_role_policy = data.aws_iam_policy_document.trust.json
}

resource "aws_iam_role_policy" "this" {
  name   = "deploy"
  role   = aws_iam_role.this.id
  policy = data.aws_iam_policy_document.deploy.json
}

data "aws_iam_policy_document" "deploy" {
  statement {
    sid       = "TerraformStateList"
    actions   = ["s3:ListBucket"]
    resources = ["arn:aws:s3:::${var.state_bucket}"]
  }

  statement {
    sid       = "TerraformStateObjects"
    actions   = ["s3:GetObject", "s3:PutObject", "s3:DeleteObject"]
    resources = ["arn:aws:s3:::${var.state_bucket}/${var.state_key}*"]
  }

  statement {
    sid       = "Lambda"
    actions   = ["lambda:*"]
    resources = ["${format(local.prefix, "lambda")}:function:${var.name}*"]
  }

  statement {
    sid     = "IAM"
    actions = ["iam:*"]
    resources = [
      "arn:aws:iam::${local.account_id}:role/${var.name}*",
      "arn:aws:iam::${local.account_id}:policy/${var.name}*",
    ]
  }

  statement {
    sid       = "ReadOIDCProvider"
    actions   = ["iam:GetOpenIDConnectProvider", "iam:ListOpenIDConnectProviders"]
    resources = ["*"]
  }

  # CI can refresh the deploy role and OIDC provider but not change them,
  # so changes to this module have to be applied locally.
  statement {
    sid         = "DenySelfModification"
    effect      = "Deny"
    not_actions = ["iam:Get*", "iam:List*"]
    resources = [
      "arn:aws:iam::${local.account_id}:role/${local.role_name}",
      "arn:aws:iam::${local.account_id}:oidc-provider/token.actions.githubusercontent.com",
    ]
  }

  statement {
    sid       = "Logs"
    actions   = ["logs:*"]
    resources = ["${format(local.prefix, "logs")}:log-group:/aws/lambda/${var.name}*"]
  }

  statement {
    sid       = "SSM"
    actions   = ["ssm:*"]
    resources = ["${format(local.prefix, "ssm")}:parameter/${var.name}/*"]
  }

  dynamic "statement" {
    for_each = length(var.read_parameter_arns) > 0 ? [1] : []

    content {
      sid       = "ReadSharedParameters"
      actions   = ["ssm:GetParameter"]
      resources = var.read_parameter_arns
    }
  }

  # API and certificate IDs are generated, so these can't be scoped by name.
  statement {
    sid     = "APIGateway"
    actions = ["apigateway:*"]
    resources = [
      "arn:aws:apigateway:${local.region}::/apis",
      "arn:aws:apigateway:${local.region}::/apis/*",
      "arn:aws:apigateway:${local.region}::/domainnames",
      "arn:aws:apigateway:${local.region}::/domainnames/*",
      "arn:aws:apigateway:${local.region}::/tags/*",
    ]
  }

  statement {
    sid       = "ACM"
    actions   = ["acm:*"]
    resources = ["${format(local.prefix, "acm")}:certificate/*"]
  }

  statement {
    sid       = "Describe"
    actions   = ["logs:DescribeLogGroups", "ssm:DescribeParameters", "acm:RequestCertificate", "acm:ListCertificates"]
    resources = ["*"]
  }
}
