# 20i Stack - Docker Development Environment

[![Version](https://img.shields.io/github/v/release/peternicholls/20i-Hosting-Stack-for-Docker)](https://github.com/peternicholls/20i-Hosting-Stack-for-Docker/releases)
[![License](https://img.shields.io/github/license/peternicholls/20i-Hosting-Stack-for-Docker)](LICENSE)

**A reusable, multi-platform Docker stack based on [20i's](https://www.20i.com) shared hosting environment for PHP projects with Nginx, PHP-FPM, MariaDB, and phpMyAdmin.**

Works great on Intel/AMD and Apple Silicon (ARM64). Choose ARM-native phpMyAdmin for best performance on Apple Silicon, or the cross-platform image for universal compatibility.

**Perfect for**: PHP development, Laravel projects, WordPress development, prototyping, and any web project needing a quick, reliable development environment.

## Topics
docker, docker-compose, php, php-fpm, nginx, mariadb, phpmyadmin, apple-silicon, arm64, cross-platform, development-environment, macos

## Overview
 A reusable, centralized Docker development stack for PHP projects using:
 - **PHP 8.5** with FPM on Alpine Linux
- **Nginx** as reverse proxy
- **MariaDB** for database
- **phpMyAdmin** for database management

## Quick Start

### Shell Commands (Recommended)
```bash
# From any project directory:
20i-gui      # Interactive menu
# Aliases:
dcu          # Start stack (uses current directory)
dcd          # Stop stack
```

### Manual Usage
```bash
cd /path/to/your/project
export CODE_DIR=$(pwd)
export COMPOSE_PROJECT_NAME=$(basename "$(pwd)")
cp .env.example .env   # edit .env if needed
docker compose up -d
```

## Features

✅ **Centralized Stack** - One stack serves any project  
✅ **Project Isolation** - Each project gets isolated containers  
✅ **Environment Variables** - Fully configurable via .env or .20i-local  
✅ **Shell Integration** - Convenient aliases and functions  
✅ **GUI Interface** - Interactive menu system  
✅ **Live Reloading** - Volume mounting for development  

## Access Points

- **Website**: http://localhost (or custom HOST_PORT)
- **phpMyAdmin**: http://localhost:8081
- **Database**: localhost:3306

## Default Credentials

- **MySQL Root**: `root` / `root`
- **MySQL User**: `devuser` / `devpass`
- **Default DB**: `devdb`

## Configuration

### Global Settings (.env.example)
```bash
HOST_PORT=80
PHP_VERSION=8.5
MYSQL_VERSION=10.6
MYSQL_PORT=3306
PMA_PORT=8081
```

### Central Variables (config/stack-vars.yml)
- Define shared defaults in `config/stack-vars.yml`. Scripts read these values and export them if not already set by `.env` or `.20i-local`.
- Supported keys:
	- `PHP_VERSION` (e.g., `PHP_VERSION: "8.5"`)
	- `MYSQL_VERSION` (e.g., `MYSQL_VERSION: "10.6"`)
	- `PMA_IMAGE` (e.g., `PMA_IMAGE: "phpmyadmin/phpmyadmin:latest"`)
- Override via `.env`, `.20i-local`, or environment variables as needed.

### phpMyAdmin Architecture (ARM vs Cross‑Platform)
- **Default (Cross‑Platform):** Uses `phpmyadmin/phpmyadmin:latest` and works on all architectures.
- **Apple Silicon/ARM (Recommended):** Set `PMA_IMAGE=arm64v8/phpmyadmin:latest` for native ARM performance.

You can set this in your `.env`:
```bash
# Cross‑platform (default)
PMA_IMAGE=phpmyadmin/phpmyadmin:latest

# Apple Silicon/ARM (native)
# PMA_IMAGE=arm64v8/phpmyadmin:latest
```

Or select it via the GUI when starting the stack.

### Per-Project Settings (.20i-local)
Create in your project root:
```bash
export HOST_PORT=8080
export MYSQL_DATABASE=myproject_db
export MYSQL_USER=projectuser
export MYSQL_PASSWORD=projectpass
```

## Architecture

- **Nginx (Port 80)**: Front-end web server and reverse proxy
- **Apache/PHP-FPM (Port 9000)**: PHP processing engine
- **MariaDB (Port 3306)**: Database server
- **phpMyAdmin (Port 8081)**: Database management interface

## Files Structure

```
20i-stack/
├── docker/
│   ├── apache/
│   │   ├── Dockerfile          # PHP 8.5 + extensions
│   │   └── php.ini            # PHP configuration
│   └── nginx.conf.tmpl        # Nginx reverse proxy config
├── docker-compose.yml         # Main stack definition
├── 20i-gui                   # Interactive CLI menu
├── .env.example              # Default configuration
└── README.md                 # This file
```

## Shell Integration

The easiest way to integrate with your shell is to source the provided example script:

```bash
# Add to your .zshrc or .bashrc:
source /path/to/20i-stack/zsh-example-script.zsh
```

**For zsh users**: Copy the script to your home directory for easier access:
```bash
cp /path/to/20i-stack/zsh-example-script.zsh ~/.20i-stack.zsh
echo "source ~/.20i-stack.zsh" >> ~/.zshrc
```

This provides convenient commands:
- `20i-up` - Start stack for current directory
- `20i-down` - Stop stack
- `20i-status` (or just `20i`) - View status
- `20i-logs` - Follow logs
- `20i-destroy` - Destroy stack and volumes (with confirmation)
- `20i-gui` - Launch interactive menu
- `dcu` / `dcd` - Aliases for up/down

**Optional**: Customize the stack location by setting environment variables before sourcing:
```bash
export STACK_FILE=/custom/path/docker-compose.yml  # Or set STACK_HOME
source ~/.20i-stack.zsh
```

## Workflow Examples

### Start New Project
```bash
cd /path/to/new-project
dcu                    # Starts stack for this project
# Site available at http://localhost
```

### Switch Projects
```bash
dcd                    # Stop current stack
cd /path/to/other-project
dcu                    # Start stack for new project
```

### Interactive Management
```bash
20i-gui               # Opens menu with start/stop/status/logs options
```

## Troubleshooting

### Project name normalization

If your project folder name contains spaces, uppercase letters, or other characters that Docker Compose disallows, the CLI will automatically normalize it into a safe project name. Names are converted to lowercase, invalid characters are replaced with hyphens, multiple separators are collapsed, and a leading letter is ensured if needed. For example:

- `DEV BirminghamFilms` → `dev-birminghamfilms`

The normalized project name will be shown when the stack starts.

### Port Conflicts
```bash
# Use custom port
export HOST_PORT=8080
dcu
```

### Database Issues
```bash
# Reset database (or use 20i-destroy for complete cleanup)
dcd
docker volume rm $(docker volume ls -q | grep db_data)
dcu
```

### Destroy Stack Completely
```bash
# WARNING: This removes all data including database volumes!
20i-destroy           # Prompts for confirmation
# Or from GUI: Choose option 6
```

### View Logs
```bash
20i-logs              # Follow all logs
20i-gui               # Menu option for specific service logs
```

## Requirements

- Docker Desktop for Mac
- Bash/Zsh shell
- Optional: `dialog` package for prettier GUI menus

## License

MIT License - Use freely for development purposes.

---


