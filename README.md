# Run go test

[![Step changelog](https://shields.io/github/v/release/bitrise-steplib/steps-go-test?include_prereleases&label=changelog&color=blueviolet)](https://github.com/bitrise-steplib/steps-go-test/releases)

Runs go test command and exports the test results.

<details>
<summary>Description</summary>

Runs go test command and exports the test results.

`go test -v <package>`
</details>

## 🧩 Get started

Add this step directly to your workflow in the [Bitrise Workflow Editor](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/steps/adding-steps-to-a-workflow.html).

You can also run this step directly with [Bitrise CLI](https://github.com/bitrise-io/bitrise).

## ⚙️ Configuration

<details>
<summary>Inputs</summary>

| Key | Description | Flags | Default |
| --- | --- | --- | --- |
| `package` | The package argument to be passed to the go test command.  For example 'go test math', 'go test ./...', and even 'go test .' | required |  |
| `covermode` | Set the mode for coverage analysis for the package[s] being tested.  Possible values are: - 'none': no coverage analysis - 'set': bool: does this statement run? - 'count': int: how many times does this statement run? - 'atomic': int: count, but correct in multithreaded tests; significantly more expensive. | required | `none` |
| `test_options` | Additional options to be added to the executed go test command. |  |  |
| `test_report_name` | Name of the generated test report JUnit xml to be used in Bitrise Test Reports.  If not specified the provided package name will be used. |  |  |
| `output_dir` | This directory will contain the generated artifacts. | required | `$BITRISE_DEPLOY_DIR` |
</details>

<details>
<summary>Outputs</summary>

| Environment Variable | Description |
| --- | --- |
| `GO_TEST_RUN_LOG_PATH` | Path to the generated test run log file. |
| `GO_CODE_COVERAGE_REPORT_PATH` | Path to the generated code coverage profile file.  This output will be available only if the covermode is not 'none'. |
</details>

## 🙋 Contributing

We welcome [pull requests](https://github.com/bitrise-steplib/steps-go-test/pulls) and [issues](https://github.com/bitrise-steplib/steps-go-test/issues) against this repository.

For pull requests, work on your changes in a forked repository and use the Bitrise CLI to [run step tests locally](https://docs.bitrise.io/en/bitrise-ci/bitrise-cli/running-your-first-local-build-with-the-cli.html).

Learn more about developing steps:

- [Create your own step](https://docs.bitrise.io/en/bitrise-ci/workflows-and-pipelines/developing-your-own-bitrise-step/developing-a-new-step.html)
