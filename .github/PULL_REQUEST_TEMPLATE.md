## Description

<!-- Please describe the changes in this pull request. Explain what is being changed and why. -->

## Motivation and Context

<!-- Why is this change required? What problem does it solve? Link any relevant issues (e.g., "Fixes #123"). -->

## Conventional Commit Format

This project uses [Conventional Commits](https://www.conventionalcommits.org/) for automated changelog generation and semantic versioning.

**Please ensure your PR title follows this format:**

```
<type>: <description>

Examples:
feat: add new configuration option
fix: resolve port detection issue
docs: update README with ARM instructions
```

**Common types:**
- `feat:` - New feature (minor version bump)
- `fix:` - Bug fix (patch version bump)
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks

**Breaking changes:**
```
feat!: change default PHP version

BREAKING CHANGE: Projects using PHP 7.x must update their configuration.
```

## How Has This Been Tested?

<!-- Describe the tests that you ran to verify your changes. Provide instructions so we can reproduce. -->
<!-- Examples: unit tests, integration tests, manual tests, screenshots, etc. -->

## Screenshots (if appropriate)

<!-- Add screenshots or animated GIFs to help review the UI/UX changes, if applicable. -->

## Types of Changes

<!-- Check all that apply: -->
- [ ] Bug fix (non-breaking change that fixes an issue)
- [ ] New feature (non-breaking change that adds functionality)
- [ ] Breaking change (fix or feature that would cause existing functionality to change)
- [ ] Documentation update
- [ ] Refactor / code quality improvement
- [ ] Other (please describe)

## Checklist

<!-- Go through this checklist and mark completed items with an "x" (e.g., [x]). -->
- [ ] I have read the contributing guidelines (if available).
- [ ] I have updated the documentation as needed.
- [ ] I have added or updated tests to cover my changes.
- [ ] All new and existing tests pass locally.
- [ ] I have considered security and privacy implications of these changes.

## Additional Context

<!-- Add any other context or information that reviewers should know about this pull request. -->
