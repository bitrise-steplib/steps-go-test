# Run Go Tests

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-go-test?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-go-test/releases)

Executes tests using Go's test runner.

<details>
<summary>Description</summary>

This step will run tests and generate [coverage profiles][profiles] for
specified packages.  In addition, the [data race detector][drr] is activated to
locate potential runtime issues.

[profiles]: https://go.dev/blog/cover
[drr]: https://go.dev/doc/articles/race_detector

#### Related steps

[Go list](https://github.com/bitrise-steplib/steps-go-list)

</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `package` | Newline-deliminated list of packages to test | required | `true` |
| `BITRISE_DEPLOY_DIR` | A path to a location to save the aggregated code coverage reports | required | `true` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `GO_CODE_COVERAGE_REPORT_PATH` | Path to file containing code coverage results |

</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/bitrise-step-restore-gradle-cache/pulls) and [issues](https://github.com/bitrise-steplib/bitrise-step-restore-gradle-cache/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

**Note:** this step's end-to-end tests (defined in `e2e/bitrise.yml`) are working with secrets which are intentionally not stored in this repo. External contributors won't be able to run those tests. Don't worry, if you open a PR with your contribution, we will help with running tests and make sure that they pass.


Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)