# Quick Start Guide

## Install in 60 Seconds

### 1. Download and Install

```bash
# Clone the bundle
git clone --depth 1 --filter=blob:none --sparse https://github.com/peternicholls/SafeDownload.git
cd SafeDownload
git sparse-checkout set .specify/bundles/research-agent

# Run installer
cd .specify/bundles/research-agent
./install.sh
```

### 2. Choose Your Project Type

When prompted:
- **1** = Generic/Default (balanced)
- **2** = Web App (security & scalability)
- **3** = Library/SDK (API design)
- **4** = Infrastructure (reliability)

### 3. Verify

```bash
.specify/scripts/bash/check-research-prerequisites.sh
```

### 4. Start Researching

```bash
@speckit.research "Compare React vs Vue vs Svelte"
```

## What You Get

```
your-project/
├── .github/agents/speckit.research.agent.md  ← AI agent
├── .specify/
│   ├── config/research.yaml                  ← Your config
│   ├── templates/                            ← Templates
│   ├── scripts/bash/                         ← Helper scripts
│   ├── docs/                                 ← Documentation
│   └── examples/research-configs/            ← Config templates
```

## Common Commands

```bash
# Initialize configuration
.specify/scripts/bash/research.sh init

# Check prerequisites
.specify/scripts/bash/check-research-prerequisites.sh

# Get research status
.specify/scripts/bash/research.sh status

# Create new research area
@speckit.research "Your question"
```

## Need Help?

- **Setup Guide**: `.specify/docs/research-agent-setup.md`
- **Framework**: `.specify/docs/research-framework.md`
- **Full README**: `README.md`

## One-Line Install (Automated)

```bash
curl -fsSL https://raw.githubusercontent.com/peternicholls/SafeDownload/develop/.specify/bundles/research-agent/install.sh | bash
```

## Configuration Quick Reference

Edit `.specify/config/research.yaml`:

```yaml
metadata:
  project_name: "YourProject"      # Change this
  project_type: "web-app"          # Your type

paths:
  research_root: "research"        # Where to store research
  constitution: null               # Path or null

decision_criteria:
  weights:
    security: 30                   # Adjust for your priorities
    scalability: 25
    maintainability: 20
```

## Next Steps

1. ✅ Customize `.specify/config/research.yaml`
2. ✅ Run prerequisites check
3. ✅ Read the setup guide
4. ✅ Start your first research

That's it! You're ready to research. 🎉
