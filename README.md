# SecureCN Terraform Provider

[![Build and Test](https://github.com/Portshift/terraform-provider-securecn/actions/workflows/test.yml/badge.svg)](https://github.com/Portshift/terraform-provider-securecn/actions/workflows/test.yml)

- [Provider Documentation Website](https://securecn.readme.io/docs/terraform-provider)
- [Provider Terraform Registry Page](https://registry.terraform.io/providers/Portshift/securecn/latest)

<img src="https://raw.githubusercontent.com/hashicorp/terraform-website/master/public/img/logo-hashicorp.svg" width="600px">

## Development

### Building

```bash
go install
```

### Update the docs after code changes

```bash
go install go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
go generate
```
