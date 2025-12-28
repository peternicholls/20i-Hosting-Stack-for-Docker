Thanks for considering contributing! Please follow these guidelines:

- Fork the repository and create a feature branch.
- Make small, focused commits with clear messages following [Conventional Commits](https://www.conventionalcommits.org/).
- Run `docker compose up -d` and ensure the stack starts successfully before opening a PR.
- Open a pull request against the `main` branch and describe your change.

Development tips:
- Copy `.env.example` to `.env` and edit values locally (`cp .env.example .env`).
- To test with a different compose file, set `STACK_FILE` to the full path of the `docker-compose.yml` you want to use.

## Conventional Commits

This project uses conventional commits for automated changelog generation and semantic versioning. All commit messages should follow this format:

```
<type>: <description>

[optional body]

[optional footer]
```

**Common types:**
- `feat:` - New feature (triggers minor version bump)
- `fix:` - Bug fix (triggers patch version bump)
- `docs:` - Documentation changes
- `refactor:` - Code refactoring
- `perf:` - Performance improvements
- `test:` - Adding or updating tests
- `chore:` - Maintenance tasks
- `ci:` - CI/CD changes

**Breaking changes:**
```
feat!: change default PHP version

BREAKING CHANGE: Projects using PHP 7.x must update their configuration.
```

Breaking changes trigger a major version bump.

## Release Process

### For Maintainers

The release process is fully automated using [release-please](https://github.com/googleapis/release-please):

1. **Merge PRs** with conventional commit messages to `main`
2. **Review the Release PR** that release-please creates automatically
   - Title: `chore(main): release X.Y.Z`
   - Contains auto-generated CHANGELOG entries
3. **Merge the Release PR** when ready to release
4. **Automatic actions:**
   - Git tag is created
   - GitHub Release is published
   - Release artifacts are built and attached
   - CHANGELOG is updated

### Pre-release Versions

For alpha, beta, or release candidate versions, use the workflow_dispatch trigger:

1. Go to **Actions** → **Release** workflow
2. Click **Run workflow**
3. Enter version (e.g., `2.0.0-alpha.1`)
4. Select pre-release option
5. Click **Run workflow**

### Release Artifacts

Each release includes:
- `20i-stack-vX.Y.Z.tar.gz` - Complete distribution archive
- `install.sh` - Standalone quick installer
- `checksums.sha256` - SHA256 verification hashes

### Quick Install for Users

Users can install the latest release with:

```bash
curl -fsSL https://github.com/peternicholls/20i-stack/releases/latest/download/install.sh | bash
```

Or install a specific version:

```bash
curl -fsSL https://github.com/peternicholls/20i-stack/releases/download/v2.0.0/install.sh | bash -s 2.0.0
```

## Feature Branch Workflow

For maintainers and AI agents working on features:

### Pre-PR Checklist

Before creating a final PR for your feature:

1. **Complete your feature work** on the feature branch (e.g., `004-release-workflow`)
2. **Merge latest changes from main/master**:
   ```bash
   git checkout main
   git pull origin main
   git checkout <your-feature-branch>
   git merge main
   # Resolve any conflicts
   git push origin <your-feature-branch>
   ```
3. **Sync specs across all branches** before final PR:
   ```bash
   # From your feature branch with updated specs
   for branch in $(git branch | grep -E "^  [0-9]" | sed 's/^  //'); do
     git checkout "$branch"
     git checkout <source-branch> -- specs/
     git add specs/
     git commit -m "chore: sync all spec folders for cross-reference"
     git push origin "$branch"
   done
   ```
4. **Return to your feature branch** and create the PR

### Why This Workflow?

- **Merging from main first**: Ensures your feature includes any upstream changes, reducing merge conflicts later
- **Spec synchronization**: All feature branches maintain a complete copy of the `specs/` directory for cross-reference, maintaining visibility across the project roadmap

Thanks — maintainers will review PRs and request changes as needed.
