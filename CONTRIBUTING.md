# Contributing Guidelines

The following is a set of guidelines for contributing to the NGINX Agent. We really appreciate that you are considering contributing!

#### Table Of Contents

[Getting Started](#getting-started)

[Contributing](#contributing)

[Code Guidelines](#code-guidelines)

[Code of Conduct](https://github.com/nginx/agent/blob/main/docs/CODE_OF_CONDUCT.md)

## Getting Started

Follow our [Installation Guide](https://github.com/nginx/agent/blob/main/README.md#Installation) to get the NGINX Agent up and running.

<!-- ### Project Structure (OPTIONAL) -->

## Contributing

### Report a Bug

To report a bug, open an issue on GitHub with the label `bug` using the available bug report issue template. Please ensure the bug has not already been reported. **If the bug is a potential security vulnerability, please report using our security policy.**

### Suggest a Feature or Enhancement

To suggest a feature or enhancement, please create an issue on GitHub with the label `feature` or `enhancement` using the available feature request issue template. Please ensure the feature or enhancement has not already been suggested.

### Open a Pull Request

* Fork the repo, create a branch, implement your changes, add any relevant tests, submit a PR when your changes are **tested** and ready for review.
* Fill in [our pull request template](https://github.com/nginx/agent/blob/main/.github/pull_request_template.md).

Note: if you'd like to implement a new feature, please consider creating a feature request issue first to start a discussion about the feature.

### F5 Contributor License Agreement (CLA)

F5 requires all external contributors to agree to the terms of the F5 CLA (available [here](https://github.com/f5/.github/blob/main/CLA/cla-markdown.md))
before any of their changes can be incorporated into an F5 Open Source repository.

If you have not yet agreed to the F5 CLA terms and submit a PR to this repository, a bot will prompt you to view and
agree to the F5 CLA. You will have to agree to the F5 CLA terms through a comment in the PR before any of your changes
can be merged. Your agreement signature will be safely stored by F5 and no longer be required in future PRs.

## Code Guidelines

<!-- ### Go Guidelines Here Linter-->

* when writing [table-driven tests](https://go.dev/wiki/TableDrivenTests), add a numbered prefix to each test case name (`Test 1:`, `Test 2:`, and so on) to make it easier to see which cases failed.

### Git Guidelines

* Keep a clean, concise and meaningful git commit history on your branch (within reason), rebasing locally and squashing before submitting a PR.
* If possible and/or relevant, use the [Conventional Commits](https://www.conventionalcommits.org/en/v1.0.0/) format when writing a commit message, so that changelogs can be automatically generated
* Follow the guidelines of writing a good commit message as described here <https://chris.beams.io/posts/git-commit/> and summarised in the next few points:
  * In the subject line, use the present tense ("Add feature" not "Added feature").
  * In the subject line, use the imperative mood ("Move cursor to..." not "Moves cursor to...").
  * Limit the subject line to 72 characters or less.
  * Reference issues and pull requests liberally after the subject line.
  * Add more detailed description in the body of the git message (`git commit -a` to give you more space and time in your text editor to write a good message instead of `git commit -am`).
