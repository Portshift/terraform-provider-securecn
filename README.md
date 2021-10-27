# SecureCN Terraform Provider

- [Provider Documentation Website](https://securecn.readme.io/docs/terraform-provider)
- [Provider Terraform Registry Page](https://registry.terraform.io/providers/Portshift/securecn/latest)

<img src="https://cdn.rawgit.com/hashicorp/terraform-website/master/content/source/assets/images/logo-hashicorp.svg" width="600px">

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
