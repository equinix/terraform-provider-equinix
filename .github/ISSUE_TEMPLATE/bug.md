---
name: Bug
about: For when something is there, but doesn't work how it should
title: ''
labels: ''
assignees: ''

---

<!--- Please keep this note for the community --->

### Community Note

* Please vote on this issue by adding a 👍 [reaction](https://blog.github.com/2016-03-10-add-reactions-to-pull-requests-issues-and-comments/) to the original issue to help the community and maintainers prioritize this request.
* Please do not leave _+1_ or _me too_ comments, they generate extra noise for issue followers and do not help prioritize the request.
* If you are interested in working on this issue or have submitted a pull request, please leave a comment.

<!--- Thank you for keeping this note for the community --->

### Terraform Version

<!--- Please run `terraform -v` to show the Terraform core version and provider version(s). If you are not running the latest version of Terraform or the provider, please upgrade because your issue may have already been fixed. [Terraform documentation on provider versioning](https://www.terraform.io/docs/configuration/providers.html#provider-versions). --->

### Affected Resource(s)

<!--- Please list the affected resources and data sources. --->

* equinix_XXXXX

### Terraform Configuration Files

<!--- Information about code formatting: https://help.github.com/articles/basic-writing-and-formatting-syntax/#quoting-code --->

```tf
# Copy-paste your Terraform configurations here.
#
# For large Terraform configs, please use a service like Dropbox and share a link to the ZIP file.
#
# If reproducing the bug involves modifying the config file (e.g., apply a config,
# change a value, apply the config again, see the bug), then please include both:
# * the version of the config before the change, and
# * the version of the config after the change.
```

### Debug Output

<!---
Please provide a link to a GitHub Gist containing the complete debug output. Please do NOT paste the debug output in the issue; just paste a link to the Gist.

Note: the debug log may contain your API key, please search the log for `X-Auth-Token` and `Authorization` (HTTP request header holding the API token) and remove it.

To obtain the debug output, run `terraform apply` with the environment variable `TF_LOG=DEBUG`. See the [Terraform documentation on debugging](https://www.terraform.io/docs/internals/debugging.html) for more information.
--->

### Panic Output

<!--- If Terraform produced a panic, please provide a link to a GitHub Gist containing the output of the `crash.log`. --->

### Expected Behavior

<!--- What should have happened? --->

### Actual Behavior

<!--- What actually happened? --->

### Steps to Reproduce

<!--- Please list the steps required to reproduce the issue. --->

1. `terraform apply`

### References

<!---
Information about referencing Github Issues: https://help.github.com/articles/basic-writing-and-formatting-syntax/#referencing-issues-and-pull-requests

Are there any other GitHub issues (open or closed) or pull requests that should be linked here? Vendor documentation? For example:
--->

* #0000
