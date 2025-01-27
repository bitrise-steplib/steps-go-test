# Go test

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-go-test?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-go-test/releases)

Runs Go test

<details>
<summary>Description</summary>

Runs Go test on the given packages one-by-one:

`go test -v <package>`
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://devcenter.bitrise.io/steps-and-workflows/steps-and-workflows-index/).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `packages` | Newline separated list of Go packages, to run the Go test command against.  __Example:__  ``` github.com/my/step github.com/bitrise/step/tool ``` |  | `$BITRISE_GO_PACKAGES` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `GO_CODE_COVERAGE_REPORT_PATH` | Path to the code coverage report file, which contains each package's code coverage report. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-go-test/pulls) and [issues](https://github.com/bitrise-steplib/steps-go-test/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://devcenter.bitrise.io/bitrise-cli/run-your-first-build/).

Learn more about developing steps:

- [Create your own step](https://devcenter.bitrise.io/contributors/create-your-own-step/)
- [Testing your Step](https://devcenter.bitrise.io/contributors/testing-and-versioning-your-steps/)
