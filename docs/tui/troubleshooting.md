# 20i Stack Manager TUI - Troubleshooting Guide

## Common Issues and Solutions

### Installation Issues

#### "command not found: 20i-stack-manager"

**Cause**: Binary not in PATH or not installed correctly.

**Solutions**:
1. Verify installation:
   ```bash
   ls -la ~/go/bin/20i-stack-manager
   ```

2. Check PATH includes `$GOPATH/bin`:
   ```bash
   echo $PATH | grep -q "$GOPATH/bin" && echo "✓ In PATH" || echo "✗ Not in PATH"
   ```

3. Add to PATH if needed (add to `~/.bashrc` or `~/.zshrc`):
   ```bash
   export PATH="$HOME/go/bin:$PATH"
   ```

4. Or run directly:
   ```bash
   cd /path/to/20i-stack/tui
   ./bin/20i-stack-manager
   ```

#### "cannot find package" during build

**Cause**: Missing Go dependencies.

**Solution**:
```bash
cd /path/to/20i-stack/tui
go mod download
go mod tidy
make build
```

### Docker Connection Issues

#### "Cannot connect to Docker daemon"

**Cause**: Docker is not running or socket permissions are incorrect.

**Solutions**:

1. **Check Docker is running**:
   ```bash
   docker ps  # Should list containers, not error
   ```

2. **Start Docker Desktop** (macOS/Windows):
   - Open Docker Desktop application
   - Wait for "Docker is running" indicator

3. **Check Docker service** (Linux):
   ```bash
   sudo systemctl status docker
   sudo systemctl start docker  # If not running
   ```

4. **Fix socket permissions** (Linux):
   ```bash
   sudo usermod -aG docker $USER
   # Log out and back in for changes to take effect
   ```

5. **Verify socket exists**:
   ```bash
   ls -la /var/run/docker.sock
   # Should show: srw-rw---- 1 root docker 0 ...
   ```

#### "Error response from daemon: dial unix /var/run/docker.sock: permission denied"

**Cause**: Your user doesn't have permissions to access Docker socket.

**Solution** (Linux):
```bash
sudo usermod -aG docker $USER
newgrp docker  # Or log out and back in
```

**Solution** (macOS/Windows):
- Reinstall Docker Desktop
- Ensure Docker Desktop is running with proper permissions

### Container Operation Issues

#### "Container not found" error

**Cause**: Container was removed outside the TUI or name changed.

**Solutions**:
1. The container list should auto-refresh after errors
2. Exit and restart the TUI to force refresh
3. Verify containers exist:
   ```bash
   docker ps -a
   ```

#### "Port already allocated" error

**Cause**: Another container or process is using the same port.

**Solutions**:
1. Check what's using the port:
   ```bash
   lsof -i :80      # Replace 80 with your port
   netstat -tulpn | grep :80
   ```

2. Stop conflicting container/process
3. Or modify port in docker-compose.yml

#### Container won't start

**Causes and solutions**:

1. **Check logs manually**:
   ```bash
   docker logs <container-name>
   ```

2. **Port conflicts** - See "Port already allocated" above

3. **Resource constraints**:
   ```bash
   docker system df  # Check disk space
   docker system prune  # Clean up if needed
   ```

4. **Image issues**:
   ```bash
   docker pull <image-name>  # Re-pull image
   ```

### UI Display Issues

#### Terminal rendering is broken/garbled

**Causes and solutions**:

1. **Terminal too small**:
   - Minimum size: 80x24
   - Recommended: 120x40 or larger
   - Resize terminal window

2. **Color support issues**:
   ```bash
   echo $TERM  # Should be xterm-256color or similar
   export TERM=xterm-256color
   ```

3. **Try different terminal**:
   - macOS: iTerm2, Alacritty
   - Linux: gnome-terminal, Alacritty, kitty
   - Windows: Windows Terminal, Alacritty

4. **Reset terminal**:
   ```bash
   reset
   clear
   ```

#### Text is unreadable or colors are wrong

**Solutions**:

1. **Check terminal supports 256 colors**:
   ```bash
   tput colors  # Should return 256
   ```

2. **Update terminal color scheme** to use a modern theme

3. **Disable custom terminal transparency** if text is hard to read

#### UI doesn't update after operations

**Cause**: Refresh not working correctly.

**Solutions**:
1. Wait a few seconds - some operations take time
2. Press `Esc` to return to dashboard (may trigger refresh)
3. Exit and restart the TUI

### Performance Issues

#### TUI is slow or laggy

**Solutions**:

1. **Check Docker daemon health**:
   ```bash
   docker info  # Should respond quickly
   ```

2. **Reduce Docker resource usage**:
   - Stop unused containers: `docker stop $(docker ps -q)`
   - Clean up: `docker system prune`

3. **Terminal emulator performance**:
   - Switch to a faster terminal (Alacritty recommended)
   - Disable unnecessary terminal features (transparency, blur)

4. **System resources**:
   - Close other applications
   - Check available RAM and CPU

### Testing Issues

#### Tests fail with "Docker not available"

**This is expected in CI** - the integration tests are designed to be CI-safe and skip tests when Docker is unavailable.

**For local testing without Docker**:
```bash
cd tui
go test ./...  # Should skip Docker-dependent tests
```

**To run with Docker**:
```bash
# First ensure Docker is running
docker ps

# Then run tests
cd tui
go test ./...
```

#### "make test" command not found

**Solution**:
```bash
# Use go test directly
cd tui
go test -v ./...

# Or install make (if needed)
# macOS:
brew install make

# Linux:
sudo apt-get install make  # Debian/Ubuntu
sudo yum install make       # RHEL/CentOS
```

## Getting Help

### Check logs and diagnostics

1. **Docker logs**:
   ```bash
   docker logs <container-name>
   ```

2. **Docker system info**:
   ```bash
   docker info
   docker version
   ```

3. **TUI build info** (if implemented):
   ```bash
   20i-stack-manager --version
   ```

### Reporting Issues

When reporting issues, include:

1. **Your environment**:
   - OS and version
   - Terminal emulator and version
   - Docker version: `docker --version`
   - Go version: `go version`

2. **Steps to reproduce**:
   - What you did
   - What you expected
   - What actually happened

3. **Relevant output**:
   - Error messages
   - Docker logs
   - Screenshots if UI-related

See [CONTRIBUTING.md](../../CONTRIBUTING.md) for how to file issues.

### Additional Resources

- [User Guide](user-guide.md) - Complete feature documentation
- [Architecture](architecture.md) - Technical design details
- [Main README](../../README.md) - Project overview
- [TUI README](../../tui/README.md) - Development guide
