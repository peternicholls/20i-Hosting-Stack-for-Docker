# Feature Specification: Secrets Management Integration

**Feature Branch**: `013-secrets-management`  
**Created**: 2025-12-28  
**Status**: Draft  
**Priority**: 🟡 High  
**Input**: User description: "Avoid plain-text credentials in .env with support for external secret providers"

## User Scenarios & Testing *(mandatory)*

### User Story 1 - Reference Secrets from External Provider (Priority: P1)

As a security-conscious developer, I want to reference secrets from an external provider so that credentials are never stored in plain text in my project files.

**Why this priority**: This is the core security value - eliminating plain-text credentials from project files.

**Independent Test**: Configure 1Password reference in config, run `20i start`, verify database connects using secret fetched from 1Password.

**Acceptance Scenarios**:

1. **Given** config references `op://vault/item/field`, **When** stack starts, **Then** actual secret value is fetched and used
2. **Given** external secret is updated, **When** stack restarts, **Then** new secret value is used
3. **Given** secret reference is invalid, **When** stack starts, **Then** clear error message indicates which secret failed

---

### User Story 2 - Support Multiple Secret Providers (Priority: P2)

As a developer on a team, I want to use my team's preferred secret provider so that I can integrate with existing security infrastructure.

**Why this priority**: Flexibility in providers enables adoption across different team setups.

**Independent Test**: Configure AWS Secrets Manager reference, run `20i start`, verify secret is fetched from AWS.

**Acceptance Scenarios**:

1. **Given** 1Password CLI is configured, **When** secrets reference 1Password, **Then** secrets are fetched via `op` CLI
2. **Given** AWS credentials are configured, **When** secrets reference AWS Secrets Manager, **Then** secrets are fetched via AWS SDK
3. **Given** provider is not available, **When** stack starts, **Then** error suggests installing/configuring the provider

---

### User Story 3 - Fallback to Local .env (Priority: P3)

As a developer working offline, I want the system to fall back to local .env values so that I can work without network access to secret providers.

**Why this priority**: Fallback ensures developers aren't blocked when providers are unavailable.

**Independent Test**: Disconnect from network, configure fallback in config, run `20i start`, verify local .env values are used.

**Acceptance Scenarios**:

1. **Given** external provider is unavailable, **When** fallback is configured, **Then** local .env values are used with warning
2. **Given** no fallback configured, **When** provider is unavailable, **Then** stack fails to start with clear error
3. **Given** `--offline` flag, **When** stack starts, **Then** provider is not contacted and local values are used

---

### User Story 4 - Encrypted Local Secrets (Priority: P4)

As a solo developer, I want to encrypt my .env file so that credentials are protected without needing an external provider.

**Why this priority**: Local encryption provides security benefits for users without external infrastructure.

**Independent Test**: Encrypt .env with passphrase, run `20i start`, enter passphrase when prompted, verify stack starts with decrypted values.

**Acceptance Scenarios**:

1. **Given** encrypted .env file, **When** stack starts, **Then** user is prompted for decryption passphrase
2. **Given** correct passphrase entered, **When** decryption succeeds, **Then** stack starts with decrypted values
3. **Given** `20i secrets encrypt`, **When** command runs, **Then** plain-text .env is encrypted and original is removed

---

### Edge Cases

- What happens when secret provider rate limits are hit?
- How does the system handle secrets with special characters (newlines, quotes)?
- What happens when encryption passphrase is forgotten?
- How does the system handle concurrent access to secrets from multiple stacks?

## Requirements *(mandatory)*

### Functional Requirements

- **FR-001**: System MUST support 1Password CLI (`op`) as a secret provider
- **FR-002**: System MUST support AWS Secrets Manager as a secret provider
- **FR-003**: System MUST support encrypted local .env files (using age or similar)
- **FR-004**: Secret references MUST use clear syntax (e.g., `op://vault/item/field`)
- **FR-005**: System MUST provide fallback to plain .env when providers unavailable
- **FR-006**: System MUST display clear errors when secret fetching fails
- **FR-007**: System MUST support `--offline` flag to skip provider communication
- **FR-008**: System MUST provide `20i secrets encrypt` command for local encryption
- **FR-009**: No plain-text passwords MUST be stored in git-tracked files

### Key Entities

- **Secret Reference**: Pointer to a secret in external provider (URI format like `op://vault/item/field`)
- **Secret Provider**: External service or tool that stores and retrieves secrets (1Password, AWS, local encrypted file)
- **Encrypted Envelope**: Local .env file encrypted with a symmetric key or passphrase

## Success Criteria *(mandatory)*

### Measurable Outcomes

- **SC-001**: Zero plain-text passwords in git-tracked configuration files
- **SC-002**: Secret fetching adds less than 5 seconds to stack startup
- **SC-003**: Developers can configure secrets in under 10 minutes following documentation
- **SC-004**: Fallback to local values works 100% of the time when configured
- **SC-005**: Support for at least 2 external providers (1Password, AWS)

## Assumptions

- Developers have accounts/access to their preferred secret provider
- CLI tools for providers are installable on macOS and Linux
- Secret providers have rate limits that won't be exceeded by normal usage
- Encryption keys/passphrases can be securely stored or remembered by users
