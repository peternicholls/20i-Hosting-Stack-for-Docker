#!/usr/bin/osascript

# 20i Stack Manager - Standalone Application
# Double-click this file to manage your 20i Docker stack

try
    # Main menu dialog
    set menuChoice to choose from list {"🚀 Start Stack", "🛑 Stop Stack", "📊 View Status", "📋 View Logs", "❌ Cancel"} with title "20i Stack Manager" with prompt "What would you like to do?" default items {"🚀 Start Stack"}
    
    if menuChoice is false or menuChoice = {"❌ Cancel"} then
        return
    end if
    
    set action to item 1 of menuChoice
    
    if action = "🚀 Start Stack" then
        startStack()
    else if action = "🛑 Stop Stack" then
        stopStack()
    else if action = "📊 View Status" then
        viewStatus()
    else if action = "📋 View Logs" then
        viewLogs()
    end if
    
on error errMsg
    display alert "❌ Error" message errMsg buttons {"OK"} default button "OK"
end try

# Function to start the stack
on startStack()
    try
        # Get project directory
        set projectPath to choose folder with prompt "📁 Select your project directory:"
        set projectPath to POSIX path of projectPath
        
        # Get project name for display
        set projectName to basename(projectPath)
        
        # Ask for phpMyAdmin image architecture
        set pmaChoice to choose from list {"Cross-platform (default)", "ARM-native (Apple Silicon)"} with title "phpMyAdmin Architecture" with prompt "Choose phpMyAdmin image type:" default items {"Cross-platform (default)"}

        # Ask for custom settings
        set settingsDialog to display dialog "⚙️ Custom settings (optional):" default answer "HOST_PORT=80" with title "20i Stack Settings" buttons {"Skip", "Use Settings"} default button "Skip"
        
        set useCustomSettings to button returned of settingsDialog = "Use Settings"
        set customSettings to ""
        if useCustomSettings then
            set customSettings to text returned of settingsDialog
        end if
        
        # Build the command
        set shellScript to "cd '" & projectPath & "';" & return
        
        if customSettings is not "" then
            set shellScript to shellScript & "export " & customSettings & ";" & return
        end if

        # Define STACK_HOME for central YAML lookup
        set shellScript to shellScript & "STACK_HOME=\"${STACK_HOME:-$HOME/docker/20i-stack}\";" & return

        # Load defaults from central YAML if not set
        set shellScript to shellScript & "if [ -f \"$STACK_HOME/config/stack-vars.yml\" ]; then" & return
        set shellScript to shellScript & "  if [ -z \"$PHP_VERSION\" ]; then PHP_VERSION=\"$(awk -F': ' '/^PHP_VERSION:/ {print $2}' \"$STACK_HOME/config/stack-vars.yml\" | tr -d '\"' | tr -d " & quoted form of "'" & " | tr -d ' ')\"; export PHP_VERSION; fi;" & return
        set shellScript to shellScript & "  if [ -z \"$MYSQL_VERSION\" ]; then MYSQL_VERSION=\"$(awk -F': ' '/^MYSQL_VERSION:/ {print $2}' \"$STACK_HOME/config/stack-vars.yml\" | tr -d '\"' | tr -d " & quoted form of "'" & " | tr -d ' ')\"; export MYSQL_VERSION; fi;" & return
        set shellScript to shellScript & "  if [ -z \"$PMA_IMAGE\" ]; then PMA_IMAGE=\"$(awk -F': ' '/^PMA_IMAGE:/ {print $2}' \"$STACK_HOME/config/stack-vars.yml\" | tr -d '\"')\"; export PMA_IMAGE; fi;" & return
        set shellScript to shellScript & "fi;" & return

        # Apply phpMyAdmin image choice
        if pmaChoice is not false then
            set pmaSelection to item 1 of pmaChoice
            if pmaSelection = "ARM-native (Apple Silicon)" then
                set shellScript to shellScript & "export PMA_IMAGE=arm64v8/phpmyadmin:latest;" & return
            else
                set shellScript to shellScript & "export PMA_IMAGE=phpmyadmin/phpmyadmin:latest;" & return
            end if
        end if
        
        set shellScript to shellScript & "export COMPOSE_PROJECT_NAME='" & projectName & "';" & return
        set shellScript to shellScript & "export CODE_DIR='" & projectPath & "';" & return
        set shellScript to shellScript & "echo '🚀 Starting 20i stack for project: " & projectName & "';" & return
        set shellScript to shellScript & "echo '📁 Code directory: " & projectPath & "';" & return
        set shellScript to shellScript & "STACK_FILE=\"${STACK_FILE:-$HOME/docker/20i-stack/docker-compose.yml}\"; docker compose -f \"$STACK_FILE\" up -d;" & return
        set shellScript to shellScript & "echo '✅ Stack started! Access your site at: http://localhost';" & return
        set shellScript to shellScript & "echo '🔧 phpMyAdmin: http://localhost:8081';"
        
        # Execute in Terminal
        tell application "Terminal"
            activate
            do script shellScript
        end tell
        
        # Show success notification
        display notification "Stack starting for: " & projectName with title "🚀 20i Stack" subtitle "Check Terminal for details"
        
    on error errMsg
        display alert "❌ Error Starting Stack" message errMsg buttons {"OK"} default button "OK"
    end try
end startStack

# Function to stop the stack
on stopStack()
    try
        # Get list of running compose projects
        set shellScript to "docker ps --format 'table {{.Names}}' | grep -E '-.+-[0-9]+$' | sed 's/-[^-]*-[0-9]*$//' | sort -u"
        set runningProjects to do shell script shellScript
        
        if runningProjects = "" then
            display alert "ℹ️ No Running Stacks" message "No 20i stacks appear to be running." buttons {"OK"} default button "OK"
            return
        end if
        
        # Convert to list for dialog
        set projectList to paragraphs of runningProjects
        set selectedProject to choose from list projectList with title "🛑 Stop 20i Stack" with prompt "Select project to stop:" default items {item 1 of projectList}
        
        if selectedProject is false then
            return
        end if
        
        set projectName to item 1 of selectedProject
        
        set shellScript to "export COMPOSE_PROJECT_NAME='" & projectName & "';" & return
        set shellScript to shellScript & "echo '🛑 Stopping 20i stack: " & projectName & "';" & return
        set shellScript to shellScript & "STACK_FILE=\"${STACK_FILE:-$HOME/docker/20i-stack/docker-compose.yml}\" && docker compose -f \"$STACK_FILE\" down;" & return
        set shellScript to shellScript & "echo '✅ Stack stopped: " & projectName & "';"
        
        tell application "Terminal"
            activate
            do script shellScript
        end tell
        
        display notification "Stack stopped: " & projectName with title "🛑 20i Stack"
        
    on error errMsg
        display alert "❌ Error Stopping Stack" message errMsg buttons {"OK"} default button "OK"
    end try
end stopStack

# Function to view status
on viewStatus()
    try
        set shellScript to "echo '📊 20i Stack Status:';" & return
        set shellScript to shellScript & "STACK_FILE=\"${STACK_FILE:-$HOME/docker/20i-stack/docker-compose.yml}\" && docker compose -f \"$STACK_FILE\" ps;" & return
        set shellScript to shellScript & "echo '';" & return
        set shellScript to shellScript & "echo '🐳 All Docker containers:';" & return
        set shellScript to shellScript & "docker ps --format 'table {{.Names}}\\t{{.Status}}\\t{{.Ports}}' | head -20"
        
        tell application "Terminal"
            activate
            do script shellScript
        end tell
        
    on error errMsg
        display alert "❌ Error Viewing Status" message errMsg buttons {"OK"} default button "OK"
    end try
end viewStatus

# Function to view logs
on viewLogs()
    try
        # Get list of running compose projects
        set shellScript to "docker ps --format 'table {{.Names}}' | grep -E '.+-[a-z0-9]+-[0-9]+$' | sed 's/-[^-]*-[0-9]*$//' | sort -u"
        set runningProjects to do shell script shellScript
        
        if runningProjects = "" then
            display alert "ℹ️ No Running Stacks" message "No 20i stacks appear to be running." buttons {"OK"} default button "OK"
            return
        end if
        
        # Convert to list for dialog
        set projectList to paragraphs of runningProjects
        set selectedProject to choose from list projectList with title "📋 View 20i Stack Logs" with prompt "Select project to view logs:" default items {item 1 of projectList}
        
        if selectedProject is false then
            return
        end if
        
        set projectName to item 1 of selectedProject
        
        set shellScript to "export COMPOSE_PROJECT_NAME='" & projectName & "';" & return
        set shellScript to shellScript & "echo '📋 Viewing logs for: " & projectName & "';" & return
        set shellScript to shellScript & "echo 'Press Ctrl+C to stop following logs';" & return
        set shellScript to shellScript & "STACK_FILE=\"${STACK_FILE:-$HOME/docker/20i-stack/docker-compose.yml}\" && docker compose -f \"$STACK_FILE\" logs -f"
        
        tell application "Terminal"
            activate
            do script shellScript
        end tell
        
    on error errMsg
        display alert "❌ Error Viewing Logs" message errMsg buttons {"OK"} default button "OK"
    end try
end viewLogs

# Helper function to get basename
on basename(posixPath)
    set AppleScript's text item delimiters to "/"
    set pathItems to text items of posixPath
    set AppleScript's text item delimiters to ""
    
    # Remove trailing slash if present
    if item -1 of pathItems = "" then
        return item -2 of pathItems
    else
        return item -1 of pathItems
    end if
end basename
