resource "securecn_k8s_cluster" "local" {
  name                       = "local"
  kubernetes_cluster_context = "kind-kind"
  orchestration_type         = "KUBERNETES"
}

resource "securecn_environment" "env1" {
  name        = "env1"
  description = "desc"
  risk        = "MEDIUM"

  kubernetes_environment {
    cluster_name = securecn_k8s_cluster.local.name

    namespaces_by_labels = {
      key11 = "value11"
      key22 = "value22"
    }
  }
}

resource "securecn_deployer" "vault" {
  name = "vault"
  operator_deployer {
    cluster_id      = securecn_k8s_cluster.local.id
    service_account = "vault"
    namespace       = "default"
    rule_creation   = false
    security_check  = true
  }
}

resource "securecn_ci_policy" "vault" {
  name = "vault"
  vulnerability_policy {
    permissible_vulnerability_level = "MEDIUM"
    enforcement_option              = "FAIL"
  }
}

resource "securecn_cd_policy" "vault" {
  name      = "vault"
  deployers = [
    securecn_deployer.vault.id,
  ]
  secret_policy {
    permissible_vulnerability_level = "NO_KNOWN_RISK"
    enforcement_option              = "FAIL"
  }
}

resource "securecn_deployment_rule" "rule1" {
  rule_name = "terraform deployment rule"

  match_by_pod_name {
    names                             = ["Finance"]
    vulnerability_severity_level      = "HIGH"
    vulnerability_on_violation_action = "BLOCK"
    psp_profile                       = "Baseline"
    psp_on_violation_action           = "ENFORCE"
  }

  scope = "any"
}
