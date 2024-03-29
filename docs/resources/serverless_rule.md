---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "securecn_serverless_rule Resource - terraform-provider-securecn"
subcategory: ""
description: |-
  A SecureCN serverless rule
---

# securecn_serverless_rule (Resource)

A SecureCN serverless rule



<!-- schema generated by tfplugindocs -->
## Schema

### Required

- `rule_name` (String)

### Optional

- `action` (String)
- `match_by_function_any` (Boolean) The rule will match on any function
- `match_by_function_arn` (Block List, Max: 1) The rule will match using function arns (see [below for nested schema](#nestedblock--match_by_function_arn))
- `match_by_function_name` (Block List, Max: 1) The rule will match using function names (see [below for nested schema](#nestedblock--match_by_function_name))
- `scope` (Block List) Scope defines the scope of this rule (see [below for nested schema](#nestedblock--scope))
- `serverless_function_validation` (Block List, Max: 1) Define function security validations (see [below for nested schema](#nestedblock--serverless_function_validation))
- `status` (String)

### Read-Only

- `id` (String) The ID of this resource.

<a id="nestedblock--match_by_function_arn"></a>
### Nested Schema for `match_by_function_arn`

Required:

- `arns` (Map of String)


<a id="nestedblock--match_by_function_name"></a>
### Nested Schema for `match_by_function_name`

Required:

- `names` (List of String)


<a id="nestedblock--scope"></a>
### Nested Schema for `scope`

Optional:

- `cloud_account` (String)
- `regions` (List of String)


<a id="nestedblock--serverless_function_validation"></a>
### Nested Schema for `serverless_function_validation`

Optional:

- `data_access_risk` (String)
- `function_permission_risk` (String)
- `is_unused_function` (Boolean)
- `publicly_accessible_risk` (String)
- `risk` (String)
- `secrets_risk` (String)
- `vulnerability` (String)
