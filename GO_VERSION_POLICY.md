# Go Version Policy

The OpenAI Go SDK supports the current stable Go release and the immediately
preceding stable Go release. The oldest supported release is declared by the
`go` directive in [`go.mod`](go.mod) and is tested on every pull request.

The SDK team may retain the most recently retired Go release for up to six
months when the dependency graph and security posture allow it. This grace
period is discretionary, is not an LTS commitment, and may end early because of
security, dependency, platform, or toolchain requirements. During a grace
period, CI tests the retired minimum in addition to the current and preceding
stable Go releases.

Minimum Go version increases:

- ship in an SDK minor release, not a patch release;
- are documented in the README and release notes;
- require approval from the SDK CODEOWNERS; and
- do not require a new SDK major version when exported APIs and the module
  import path remain compatible.

The complete approved release note and SDK compatibility boundary are stored in
`.github/go-support-policy.json`. A pull request that changes the minimum Go
version copies that text into a `## Release note` section in its description.
When the generated release pull request opens, the release manager copies the
same text into the matching `CHANGELOG.md` release section. CI blocks the
generated release pull request until the changelog contains the approved text.

The SDK team reviews this policy within 30 days of each scheduled February and
August Go release. Automation compares the official Go release feed with the
repository's module, CI, and documentation settings each month, but never
changes the minimum version automatically.

### Current compatibility

| SDK version | Go requirement |
| --- | --- |
| v3.45.0 through current | Go 1.25 or later |
| v3.44.0 | Final release that builds with Go 1.22–1.24 |

Previously published SDK versions remain available. Unsupported Go releases
and older SDK versions receive no guaranteed fixes or security backports. Users
who need current security fixes must use a supported Go toolchain and SDK
release.

For the upstream lifecycle and toolchain-selection rules, see the [Go release
policy](https://go.dev/doc/devel/release#policy) and [Go toolchain
documentation](https://go.dev/doc/toolchain).
