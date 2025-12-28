Thanks for considering contributing! Please follow these guidelines:

- Fork the repository and create a feature branch.
- Make small, focused commits with clear messages.
- Run `docker compose up -d` and ensure the stack starts successfully before opening a PR.
- Open a pull request against the `main` branch and describe your change.

Development tips:
- Copy `.env.example` to `.env` and edit values locally (`cp .env.example .env`).
- To test with a different compose file, set `STACK_FILE` to the full path of the `docker-compose.yml` you want to use.

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
