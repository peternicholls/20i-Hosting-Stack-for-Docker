# SpecKit Research Agent Bundle

**Version**: 1.0.0  
**Schema**: 1.0.0  
**Date**: 2025-12-28

## What is this?

The SpecKit Research Agent is a project-agnostic AI agent that conducts technical research, evaluates options, and produces actionable findings with proper documentation and decision tracking.

This bundle contains everything you need to add research capabilities to your SpecKit project.

## Features

✅ **Project-Agnostic**: Works with any project type (web apps, libraries, infrastructure, CLI tools)  
✅ **Configuration-Driven**: Adapts to your project priorities via YAML config  
✅ **Quality Gates**: Ensures research completeness with validation  
✅ **Integration Ready**: Feeds findings into specs, plans, and tasks  
✅ **Template-Based**: Consistent research structure  
✅ **Optional Constitution**: Works with or without guiding principles  

## What's Included

### Core Agent
- `speckit.research.agent.md` - AI agent definition with research workflow

### Templates
- `research-template.yaml` - Structure for research areas
- `research-feedback-template.md` - Integration guide template

### Scripts
- `research.sh` - Helper script for managing research
- `check-research-prerequisites.sh` - Validation script

### Documentation
- `research-framework.md` - Generic methodology guide
- `research-agent-setup.md` - Complete setup guide
- `agent-integration.md` - How it works with other SpecKit agents

### Configuration Examples
- `default.yaml` - Generic balanced configuration
- `web-app.yaml` - Web application (security-focused)
- `library.yaml` - Library/SDK (API-focused)
- `infrastructure.yaml` - Infrastructure/DevOps (reliability-focused)

### Context
- `research-agent-context.yaml` - Agent context and reference

## Quick Install

### Option 1: Automated Install (Recommended)

```bash
# From your project root
curl -fsSL https://raw.githubusercontent.com/peternicholls/SafeDownload/develop/.specify/bundles/research-agent/install.sh | bash
```

Or download and run:

```bash
wget https://raw.githubusercontent.com/peternicholls/SafeDownload/develop/.specify/bundles/research-agent/install.sh
chmod +x install.sh
./install.sh
```

### Option 2: Manual Install

1. **Download the bundle**:
   ```bash
   git clone --depth 1 --filter=blob:none --sparse https://github.com/peternicholls/SafeDownload.git
   cd SafeDownload
   git sparse-checkout set .specify/bundles/research-agent
   ```

2. **Run the install script**:
   ```bash
   cd .specify/bundles/research-agent
   ./install.sh
   ```

3. **Or copy files manually** (see Manual Installation below)

## What Gets Installed

```
your-project/
├── .github/
│   └── agents/
│       └── speckit.research.agent.md          # Agent definition
├── .specify/
│   ├── config/
│   │   └── research.yaml                      # Your configuration
│   ├── context/
│   │   └── research-agent-context.yaml        # Agent context
│   ├── docs/
│   │   ├── research-framework.md              # Methodology
│   │   ├── research-agent-setup.md            # Setup guide
│   │   └── agent-integration.md               # Integration guide
│   ├── examples/
│   │   └── research-configs/                  # Config templates
│   │       ├── default.yaml
│   │       ├── web-app.yaml
│   │       ├── library.yaml
│   │       └── infrastructure.yaml
│   ├── scripts/
│   │   └── bash/
│   │       ├── research.sh                    # Helper script
│   │       └── check-research-prerequisites.sh # Validation
│   └── templates/
│       ├── research-template.yaml             # Research template
│       └── research-feedback-template.md      # Feedback template
```

## Post-Installation Setup

### 1. Initialize Configuration

Choose your project type and run:

```bash
# Interactive setup
.specify/scripts/bash/research.sh init

# Or copy a template
cp .specify/examples/research-configs/web-app.yaml .specify/config/research.yaml
```

### 2. Customize Configuration

Edit `.specify/config/research.yaml`:

```yaml
metadata:
  project_name: "YourProject"
  project_type: "web-app"  # web-app | library | infrastructure | cli

paths:
  research_root: "research"  # Where to store research
  constitution: null         # Path to constitution.md or null
  feature_specs: "specs"     # Your specs directory
  
decision_criteria:
  weights:
    security: 30        # Adjust weights for your priorities
    scalability: 25
    # ...
```

### 3. Verify Installation

```bash
# Check prerequisites
.specify/scripts/bash/check-research-prerequisites.sh

# Should show all green checks
```

### 4. Create First Research

```bash
# Using the agent (in VS Code with GitHub Copilot)
@speckit.research "Evaluate React vs Vue vs Svelte for our dashboard"

# Or using helper script
.specify/scripts/bash/research.sh create "React vs Vue"
```

## Configuration Templates

### Generic/Default
Balanced weights for any project:
```bash
cp .specify/examples/research-configs/default.yaml .specify/config/research.yaml
```

### Web Application
Security and scalability focused:
```bash
cp .specify/examples/research-configs/web-app.yaml .specify/config/research.yaml
```

### Library/SDK
API design and compatibility focused:
```bash
cp .specify/examples/research-configs/library.yaml .specify/config/research.yaml
```

### Infrastructure
Reliability and operational simplicity focused:
```bash
cp .specify/examples/research-configs/infrastructure.yaml .specify/config/research.yaml
```

## Usage Examples

### Library Evaluation
```bash
@speckit.research "Compare date-fns vs dayjs vs luxon"
```

### Pattern Research
```bash
@speckit.research "Research error handling patterns in React"
```

### Architecture Design
```bash
@speckit.research "Design microservices communication architecture"
```

### Tool Selection
```bash
@speckit.research "Evaluate monitoring solutions: Datadog vs New Relic vs Grafana"
```

## Integration with SpecKit

The research agent integrates with other SpecKit agents:

```
constitution → specify → clarify → research → plan → tasks → implement
```

**Research provides**:
- Library recommendations
- Architecture patterns
- Technical constraints
- Performance benchmarks
- Implementation guidance

**Research consumes**:
- Constitution principles (optional)
- Feature specifications
- Sprint blockers
- Technical questions

## Requirements

- **Git**: For version control
- **yq**: YAML processor (`brew install yq` on macOS, `apt install yq` on Linux)
- **Bash 4+**: For scripts
- **VS Code + GitHub Copilot**: For AI agent execution

## Troubleshooting

### Prerequisites Check Fails

```bash
# Run validation
.specify/scripts/bash/check-research-prerequisites.sh --verbose

# Install missing tools
brew install yq  # macOS
apt install yq   # Linux
```

### Agent Not Found

Make sure file exists:
```bash
ls -la .github/agents/speckit.research.agent.md
```

If using VS Code, check Chat Participants settings.

### Configuration Errors

Validate your config:
```bash
yq validate .specify/config/research.yaml
yq '.decision_criteria.weights' .specify/config/research.yaml
```

## Manual Installation

If the automated script doesn't work for your setup:

1. **Create directories**:
   ```bash
   mkdir -p .github/agents
   mkdir -p .specify/{config,context,docs,examples/research-configs,scripts/bash,templates}
   ```

2. **Copy agent definition**:
   ```bash
   cp bundle/files/agents/speckit.research.agent.md .github/agents/
   ```

3. **Copy templates**:
   ```bash
   cp bundle/files/templates/research-template.yaml .specify/templates/
   cp bundle/files/templates/research-feedback-template.md .specify/templates/
   ```

4. **Copy scripts**:
   ```bash
   cp bundle/files/scripts/*.sh .specify/scripts/bash/
   chmod +x .specify/scripts/bash/*.sh
   ```

5. **Copy documentation**:
   ```bash
   cp bundle/files/docs/*.md .specify/docs/
   ```

6. **Copy examples**:
   ```bash
   cp bundle/files/examples/*.yaml .specify/examples/research-configs/
   ```

7. **Copy context**:
   ```bash
   cp bundle/files/context/research-agent-context.yaml .specify/context/
   ```

8. **Choose and customize config**:
   ```bash
   cp .specify/examples/research-configs/default.yaml .specify/config/research.yaml
   # Edit .specify/config/research.yaml for your project
   ```

## Documentation

- **Setup Guide**: `.specify/docs/research-agent-setup.md`
- **Framework Guide**: `.specify/docs/research-framework.md`
- **Integration Guide**: `.specify/docs/agent-integration.md`

## Support

- **Issues**: https://github.com/peternicholls/SafeDownload/issues
- **Documentation**: See `.specify/docs/` after installation
- **Examples**: See `.specify/examples/research-configs/`

## Version History

### 1.0.0 (2025-12-28)
- Initial release
- Project-agnostic configuration system
- 4 project type templates (default, web-app, library, infrastructure)
- Prerequisites validation
- Quality gates and completion criteria
- Integration feedback templates
- Comprehensive documentation

## License

Same as SafeDownload project (see LICENSE file in repository)

## Credits

Developed as part of the SafeDownload project's SpecKit framework.
