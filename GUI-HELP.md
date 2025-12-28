# 20i Stack GUI Manager

## 🚀 Usage

From any project directory, simply run:

```bash
20i-gui
```

This gives you an interactive menu with these options:

### 📋 Menu Options:

1. **🚀 Start Stack (current directory)**
   - Uses the current directory as your project root
   - Auto-detects project name from folder name
   - Prompts for custom web port (defaults to 80)
   - Lets you choose phpMyAdmin image type:
     - Cross‑platform (default): `phpmyadmin/phpmyadmin:latest`
     - ARM‑native (Apple Silicon): `arm64v8/phpmyadmin:latest`
   - Loads `.20i-local` file if present for project-specific settings

2. **🛑 Stop Stack**
   - Shows list of running 20i stacks
   - Choose specific project to stop or stop all
   - Clean shutdown of containers

3. **📊 View Status**
   - Overview of all running Docker containers
   - List of active 20i projects

4. **📋 View Logs**
   - Shows running 20i stacks
   - Follow real-time logs for selected project
   - Press Ctrl+C to stop following

## 🎯 Perfect For:

- **Quick project switching** without remembering commands
- **Beginners** who prefer menus over command line
- **Visual confirmation** of what's running
- **Project demos** with clean interface

## 🛠 Integration with Existing Workflow:

Your existing aliases still work perfectly:
- `dcu` - Start stack (command line)
- `dcd` - Stop stack (command line) 
- `20i` - View status (command line)
- `20i-gui` - Interactive menu (new!)

## 💡 Pro Tips:

- **Dialog Support**: Install `dialog` package for prettier menus:
  ```bash
  brew install dialog
  ```

- **Project Settings**: Create `.20i-local` in your project root:
  ```bash
  export HOST_PORT=8080
  export MYSQL_DATABASE=myproject_db
  ```

- **phpMyAdmin Architecture**: You can override selection via environment:
   ```bash
   # Cross‑platform (default)
   export PMA_IMAGE=phpmyadmin/phpmyadmin:latest

   # Apple Silicon / ARM (native)
   export PMA_IMAGE=arm64v8/phpmyadmin:latest
   ```

- **From Anywhere**: The `20i-gui` command works from any project directory

Perfect complement to your powerful shell-based workflow! 🎉
