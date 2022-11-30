# SecureCN Terraform Provider

[![Build and Test](https://github.com/Portshift/terraform-provider-securecn/actions/workflows/test.yml/badge.svg)](https://github.com/Portshift/terraform-provider-securecn/actions/workflows/test.yml)

- [Provider Documentation Website](https://panoptica.readme.io/docs/terraform-provider)
- [Provider Terraform Registry Page](https://registry.terraform.io/providers/Portshift/securecn/latest)

<img src="https://raw.githubusercontent.com/hashicorp/terraform-website/master/public/img/logo-hashicorp.svg" width="600px">

## Development

### Building

```bash
go test ./...
go install
```

### Update the docs after changing resources

```bash
go install github.com/hashicorp/terraform-plugin-docs/cmd/tfplugindocs
go generate
```

### CI and testing

An acceptance test is running for all submitted PRs in the repository with GitHub Actions [test.yml](.github/workflows/test.yml).
It is compiling the provider, setting up a [kind](https://kind.sigs.k8s.io/) cluster
and performs the registration of this cluster in a separate account of the staging environment,
see the [examples/main.tf](https://github.com/Portshift/terraform-provider-securecn/blob/main/examples/main.tf) file for all the resources that are getting created during this test. After a successful test the provider destroys these resources and cleans up the account.

If the CI test fails because of Escher authentication errors, just try restarting the build (this happens time to time). The root cause of this is that Escher auth is a time based authentication method and the time on the management server and on GitHub Actions can differ.

### Releasing

The release process is fully automated by GitHub Actions and GoReleaser. To execute the process
you only need to tag the repository at the target commit with a semantically versioned git commit like: `v1.1.6`
After the release is compiled GoReleaser will upload the binaries next to the GitHub release it creates
and those are getting grabbed by the Terraform Registry automatically.
