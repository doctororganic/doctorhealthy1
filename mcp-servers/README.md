# Roo Automation MCP Server

This MCP (Model Context Protocol) server provides autonomous task execution capabilities to the Roo AI Assistant, enabling seamless integration with system operations, file management, package installation, and AI collaboration features.

## 🚀 Features

- **Shell Command Execution**: Run system commands with timeout and error handling
- **Software Installation**: Install packages via multiple package managers (brew, npm, pip, apt)
- **Redis Operations**: Direct Redis database operations for shared memory
- **File Operations**: Read, write, and manage files and directories
- **AI Collaboration**: Manage shared memory between AI assistants (Roo, Code, Kilo)
- **Docker Operations**: Container management and orchestration
- **Git Operations**: Version control operations

## 📋 Prerequisites

- Node.js 18+
- Redis server running (configured with authentication)
- Access to system package managers

## 🛠️ Installation

```bash
cd nutrition-platform/mcp-servers
npm install
```

## ⚙️ Configuration

### Environment Variables

Create a `.env` file or set environment variables:

```bash
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=secure_redis_password_2025
```

### VSCode Integration

Add to your `.vscode/settings.json`:

```json
{
  "roo.mcpServers": {
    "roo-automation": {
      "command": "node",
      "args": ["/absolute/path/to/nutrition-platform/mcp-servers/index.js"],
      "env": {
        "REDIS_HOST": "localhost",
        "REDIS_PORT": "6379",
        "REDIS_PASSWORD": "secure_redis_password_2025"
      }
    }
  }
}
```

## 🎯 Available Tools

### execute_shell_command
Execute shell commands with full control.

**Parameters:**
- `command` (string): Shell command to execute
- `cwd` (string, optional): Working directory
- `timeout` (number, optional): Timeout in milliseconds (default: 30000)

**Example:**
```javascript
{
  "command": "npm install express",
  "cwd": "/path/to/project",
  "timeout": 60000
}
```

### install_software
Install packages using various package managers.

**Parameters:**
- `package` (string): Package name
- `manager` (string): Package manager (brew, npm, pip, apt, yarn)
- `global` (boolean, optional): Install globally (npm/yarn only)

**Example:**
```javascript
{
  "package": "typescript",
  "manager": "npm",
  "global": true
}
```

### redis_operations
Perform Redis database operations.

**Parameters:**
- `operation` (string): Operation type (get, set, keys, del, config, ping, info)
- `key` (string, optional): Redis key
- `value` (string, optional): Value for set operations
- `ttl` (number, optional): Time-to-live in seconds (default: 3600)

**Example:**
```javascript
{
  "operation": "set",
  "key": "my_key",
  "value": "my_value",
  "ttl": 3600
}
```

### file_operations
Manage files and directories.

**Parameters:**
- `operation` (string): Operation type (read, write, append, list, delete, mkdir, exists)
- `path` (string): File or directory path
- `content` (string, optional): Content for write/append
- `encoding` (string, optional): File encoding (default: utf8)

**Example:**
```javascript
{
  "operation": "write",
  "path": "/path/to/file.txt",
  "content": "Hello, World!",
  "encoding": "utf8"
}
```

### ai_collaboration
Manage AI assistant coordination and shared memory.

**Parameters:**
- `action` (string): Action type (set_task, get_task, complete_task, wait_for_task, list_active, update_progress)
- `agentId` (string): AI agent identifier
- `taskId` (string, optional): Task identifier
- `data` (object, optional): Task data payload

**Example:**
```javascript
{
  "action": "set_task",
  "agentId": "roo",
  "taskId": "api_development_123",
  "data": {
    "type": "code_generation",
    "description": "Create user authentication API",
    "priority": "high",
    "status": "in_progress"
  }
}
```

### docker_operations
Manage Docker containers and images.

**Parameters:**
- `operation` (string): Operation type (build, run, stop, remove, logs, ps, images)
- `image` (string, optional): Docker image name
- `container` (string, optional): Container name
- `command` (string, optional): Command to run
- `ports` (array, optional): Port mappings

### git_operations
Perform Git version control operations.

**Parameters:**
- `operation` (string): Git operation (status, add, commit, push, pull, clone, branch, checkout)
- `repo` (string, optional): Repository URL for clone
- `files` (array, optional): Files to add/commit
- `message` (string, optional): Commit message
- `branch` (string, optional): Branch name

## 🚀 Usage Examples

### Basic Command Execution
```javascript
// Install a package
await useTool("install_software", {
  package: "lodash",
  manager: "npm"
});

// Run a build command
await useTool("execute_shell_command", {
  command: "npm run build",
  cwd: "/path/to/project"
});
```

### AI Collaboration Workflow
```javascript
// Roo starts a task
await useTool("ai_collaboration", {
  action: "set_task",
  agentId: "roo",
  taskId: "feature_xyz",
  data: {
    type: "feature_development",
    description: "Implement user dashboard",
    priority: "high"
  }
});

// Code reviews the work
await useTool("ai_collaboration", {
  action: "update_progress",
  agentId: "code",
  taskId: "code_review_xyz",
  data: {
    status: "in_progress",
    findings: "Architecture looks good, minor optimizations needed"
  }
});

// Kilo waits for completion
const result = await useTool("ai_collaboration", {
  action: "wait_for_task",
  agentId: "roo",
  taskId: "feature_xyz"
});
```

### File Management
```javascript
// Create a new file
await useTool("file_operations", {
  operation: "write",
  path: "/path/to/newfile.js",
  content: "console.log('Hello from MCP!');"
});

// Read existing file
const content = await useTool("file_operations", {
  operation: "read",
  path: "/path/to/config.json"
});
```

## 🔒 Security Features

- **Command Validation**: All shell commands are validated before execution
- **Timeout Protection**: Commands automatically timeout to prevent hanging
- **Redis Authentication**: Secure Redis connections with password protection
- **Output Sanitization**: Command outputs are sanitized for safety
- **Audit Logging**: All operations are logged for monitoring

## 📊 Monitoring & Diagnostics

The MCP server provides comprehensive logging:

```
🚀 Starting Roo Automation MCP Server...
✅ Redis client connected
🔧 Executing: npm install express
✅ Command completed successfully
📝 Task set: roo:api_development_123
```

## 🐛 Troubleshooting

### Common Issues

1. **Redis Connection Failed**
   - Ensure Redis is running: `redis-cli ping`
   - Check password configuration
   - Verify network connectivity

2. **Command Timeout**
   - Increase timeout parameter for long-running commands
   - Check system resources

3. **Permission Denied**
   - Ensure proper file permissions
   - Use sudo for system-level operations when appropriate

### Debug Mode

Enable verbose logging by setting the environment variable:

```bash
DEBUG=mcp:* node index.js
```

## 🤝 Integration with Claude Desktop

To integrate with Claude Desktop:

1. Install Claude Desktop
2. Add the MCP server configuration to your Claude Desktop config file
3. Restart Claude Desktop
4. The Roo Automation tools will be available in your conversations

## 📈 Performance Considerations

- **Connection Pooling**: Redis connections are pooled for efficiency
- **Timeout Management**: Prevents resource exhaustion
- **Error Recovery**: Automatic retry mechanisms for transient failures
- **Memory Management**: Efficient handling of large command outputs

## 🔄 Version History

- **v1.0.0**: Initial release with core automation features
  - Shell command execution
  - Software installation
  - Redis operations
  - File operations
  - AI collaboration
  - Docker and Git operations

---

This MCP server transforms Roo from a passive assistant into an autonomous development partner capable of executing complex workflows, managing infrastructure, and coordinating with other AI assistants through shared memory.