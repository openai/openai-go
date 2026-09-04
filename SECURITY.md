# Security Policy

## Reporting a vulnerability

To report a security issue affecting this SDK or OpenAI's services, please
contact disclosure@openai.com.

Do not report security vulnerabilities through public GitHub issues, pull requests, or discussions.

## What to include

- The affected package or product and version.
- A clear description of the potential impact.
- Sanitized reproduction steps or a minimal proof of concept.

For the Go SDK, also include the affected module path, release or commit, Go
version, operating system, and any known mitigations when relevant.

Do not include live credentials, API keys, customer data, or unredacted sensitive logs.
Remove access tokens, private keys, credential-bearing headers and URLs, and
other sensitive diagnostics before sharing reproduction details.

This policy covers the source code in this repository, official OpenAI Go SDK
modules such as `github.com/openai/openai-go/v3`, and official tagged SDK release
artifacts. Security issues affecting other OpenAI services may also be reported
through the same private channel.

For supported Go versions and SDK release compatibility, see
[GO_VERSION_POLICY.md](GO_VERSION_POLICY.md).

## Coordinated disclosure

Please give the maintainers a reasonable opportunity to investigate and address the issue before public disclosure.

Please follow OpenAI's
[Coordinated Vulnerability Disclosure Policy](https://openai.com/policies/coordinated-vulnerability-disclosure-policy)
when reporting and discussing the issue.
