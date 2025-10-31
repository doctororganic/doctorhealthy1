/**
 * Roo Automation MCP Server
 * Provides autonomous task execution capabilities to Roo AI Assistant
 */

const { Server } = require("@modelcontextprotocol/sdk/server/index.js");
const { StdioServerTransport } = require("@modelcontextprotocol/sdk/server/stdio.js");
const { CallToolRequestSchema, ListToolsRequestSchema } = require("@modelcontextprotocol/sdk/types.js");
const { execSync, exec } = require('child_process');
const util = require('util');
const execAsync = util.promisify(exec);
const fs = require('fs').promises;
const path = require('path');

class RooAutomationServer {
  constructor() {
    this.server = new Server(
      {
        name: "roo-automation-server",
        version: "1.0.0",
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    this.redis = null;
    this.setupToolHandlers();
  }

  // Initialize Redis connection for AI collaboration
  async initializeRedis() {
    if (!this.redis) {
      const Redis = require('ioredis');
      this.redis = new Redis({
        host: process.env.REDIS_HOST || 'localhost',
        port: process.env.REDIS_PORT || 6379,
        password: process.env.REDIS_PASSWORD || 'secure_redis_password_2025',
        retryDelayOnFailover: 100,
        enableReadyCheck: false,
        maxRetriesPerRequest: 3,
        lazyConnect: true
      });

      try {
        await this.redis.connect();
        console.error('✅ MCP Server Redis connected');
      } catch (error) {
        console.error('❌ MCP Server Redis connection failed:', error.message);
      }
    }
    return this.redis;
  }

  setupToolHandlers() {
    // List available tools
    this.server.setRequestHandler(ListToolsRequestSchema, async () => {
      return {
        tools: [
          {
            name: "execute_shell_command",
            description: "Execute shell commands on the system with timeout and error handling",
            inputSchema: {
              type: "object",
              properties: {
                command: {
                  type: "string",
                  description: "Shell command to execute"
                },
                cwd: {
                  type: "string",
                  description: "Working directory for command execution",
                  default: "."
                },
                timeout: {
                  type: "number",
                  description: "Command timeout in milliseconds",
                  default: 30000
                }
              },
              required: ["command"]
            }
          },
          {
            name: "install_software",
            description: "Install software packages using various package managers",
            inputSchema: {
              type: "object",
              properties: {
                package: {
                  type: "string",
                  description: "Package name to install"
                },
                manager: {
                  type: "string",
                  enum: ["brew", "npm", "pip", "apt", "yarn"],
                  description: "Package manager to use"
                },
                global: {
                  type: "boolean",
                  description: "Install globally (for npm/yarn)",
                  default: false
                }
              },
              required: ["package", "manager"]
            }
          },
          {
            name: "redis_operations",
            description: "Perform Redis database operations for shared memory",
            inputSchema: {
              type: "object",
              properties: {
                operation: {
                  type: "string",
                  enum: ["get", "set", "keys", "del", "config", "ping", "info"],
                  description: "Redis operation type"
                },
                key: {
                  type: "string",
                  description: "Redis key (for get/set/del operations)"
                },
                value: {
                  type: "string",
                  description: "Value for set operations"
                },
                ttl: {
                  type: "number",
                  description: "Time-to-live in seconds for set operations",
                  default: 3600
                }
              },
              required: ["operation"]
            }
          },
          {
            name: "file_operations",
            description: "Read, write, and manage files and directories",
            inputSchema: {
              type: "object",
              properties: {
                operation: {
                  type: "string",
                  enum: ["read", "write", "append", "list", "delete", "mkdir", "exists"],
                  description: "File operation type"
                },
                path: {
                  type: "string",
                  description: "File or directory path"
                },
                content: {
                  type: "string",
                  description: "Content for write/append operations"
                },
                encoding: {
                  type: "string",
                  description: "File encoding",
                  default: "utf8"
                }
              },
              required: ["operation", "path"]
            }
          },
          {
            name: "ai_collaboration",
            description: "Manage AI assistant shared memory and coordination",
            inputSchema: {
              type: "object",
              properties: {
                action: {
                  type: "string",
                  enum: ["set_task", "get_task", "complete_task", "wait_for_task", "list_active", "update_progress"],
                  description: "AI collaboration action type"
                },
                agentId: {
                  type: "string",
                  description: "AI agent identifier (e.g., 'roo', 'code', 'kilo')"
                },
                taskId: {
                  type: "string",
                  description: "Unique task identifier"
                },
                data: {
                  type: "object",
                  description: "Task data payload with flexible schema",
                  properties: {
                    type: { type: "string" },
                    description: { type: "string" },
                    status: { type: "string" },
                    priority: { type: "string" },
                    progress: { type: "string" },
                    dependencies: { type: "array" },
                    result: { type: "object" }
                  }
                }
              },
              required: ["action", "agentId"]
            }
          },
          {
            name: "docker_operations",
            description: "Manage Docker containers and images",
            inputSchema: {
              type: "object",
              properties: {
                operation: {
                  type: "string",
                  enum: ["build", "run", "stop", "remove", "logs", "ps", "images"],
                  description: "Docker operation type"
                },
                image: {
                  type: "string",
                  description: "Docker image name"
                },
                container: {
                  type: "string",
                  description: "Container name"
                },
                command: {
                  type: "string",
                  description: "Command to run in container"
                },
                ports: {
                  type: "array",
                  description: "Port mappings"
                }
              },
              required: ["operation"]
            }
          },
          {
            name: "git_operations",
            description: "Perform Git version control operations",
            inputSchema: {
              type: "object",
              properties: {
                operation: {
                  type: "string",
                  enum: ["status", "add", "commit", "push", "pull", "clone", "branch", "checkout"],
                  description: "Git operation type"
                },
                repo: {
                  type: "string",
                  description: "Repository URL for clone operations"
                },
                files: {
                  type: "array",
                  description: "Files to add/commit"
                },
                message: {
                  type: "string",
                  description: "Commit message"
                },
                branch: {
                  type: "string",
                  description: "Branch name"
                }
              },
              required: ["operation"]
            }
          }
        ]
      };
    });

    // Handle tool calls
    this.server.setRequestHandler(CallToolRequestSchema, async (request) => {
      const { name, arguments: args } = request.params;

      try {
        await this.initializeRedis();

        switch (name) {
          case "execute_shell_command":
            return await this.executeShellCommand(args);

          case "install_software":
            return await this.installSoftware(args);

          case "redis_operations":
            return await this.redisOperations(args);

          case "file_operations":
            return await this.fileOperations(args);

          case "ai_collaboration":
            return await this.aiCollaboration(args);

          case "docker_operations":
            return await this.dockerOperations(args);

          case "git_operations":
            return await this.gitOperations(args);

          default:
            throw new Error(`Unknown tool: ${name}`);
        }
      } catch (error) {
        console.error(`Tool execution error (${name}):`, error.message);
        return {
          content: [{ type: "text", text: `❌ Error executing ${name}: ${error.message}` }],
          isError: true
        };
      }
    });
  }

  async executeShellCommand(args) {
    const { command, cwd = process.cwd(), timeout = 30000 } = args;

    try {
      console.error(`🔧 Executing: ${command} (cwd: ${cwd})`);

      const result = execSync(command, {
        cwd,
        encoding: 'utf8',
        timeout,
        maxBuffer: 1024 * 1024 * 10, // 10MB buffer
        env: { ...process.env, FORCE_COLOR: '1' }
      });

      console.error(`✅ Command completed successfully`);
      return {
        content: [{ type: "text", text: result || "Command executed successfully (no output)" }]
      };
    } catch (error) {
      console.error(`❌ Command failed: ${error.message}`);
      return {
        content: [{ type: "text", text: `Command failed: ${error.message}\nStderr: ${error.stderr || 'N/A'}` }],
        isError: true
      };
    }
  }

  async installSoftware(args) {
    const { package: packageName, manager, global = false } = args;

    let command;
    switch (manager) {
      case "brew":
        command = `brew install ${packageName}`;
        break;
      case "npm":
        command = global ? `npm install -g ${packageName}` : `npm install ${packageName}`;
        break;
      case "yarn":
        command = global ? `yarn global add ${packageName}` : `yarn add ${packageName}`;
        break;
      case "pip":
        command = `pip install ${packageName}`;
        break;
      case "apt":
        command = `sudo apt-get update && sudo apt-get install -y ${packageName}`;
        break;
      default:
        throw new Error(`Unsupported package manager: ${manager}`);
    }

    return await this.executeShellCommand({ command });
  }

  async redisOperations(args) {
    const { operation, key, value, ttl = 3600 } = args;

    const redis = await this.initializeRedis();

    try {
      switch (operation) {
        case "ping":
          const pingResult = await redis.ping();
          return { content: [{ type: "text", text: `PONG: ${pingResult}` }] };

        case "get":
          const getResult = await redis.get(key);
          return { content: [{ type: "text", text: getResult || "Key not found" }] };

        case "set":
          await redis.setex(key, ttl, value);
          return { content: [{ type: "text", text: `Key set: ${key}` }] };

        case "keys":
          const keys = await redis.keys(key || "*");
          return { content: [{ type: "text", text: keys.join('\n') || "No keys found" }] };

        case "del":
          const delResult = await redis.del(key);
          return { content: [{ type: "text", text: `Deleted ${delResult} key(s)` }] };

        case "config":
          const configResult = await redis.config("GET", key || "*");
          return { content: [{ type: "text", text: JSON.stringify(configResult, null, 2) }] };

        case "info":
          const infoResult = await redis.info();
          return { content: [{ type: "text", text: infoResult }] };

        default:
          throw new Error(`Unsupported Redis operation: ${operation}`);
      }
    } catch (error) {
      return {
        content: [{ type: "text", text: `Redis operation failed: ${error.message}` }],
        isError: true
      };
    }
  }

  async fileOperations(args) {
    const { operation, path: filePath, content, encoding = 'utf8' } = args;

    try {
      switch (operation) {
        case "read":
          const data = await fs.readFile(filePath, encoding);
          return { content: [{ type: "text", text: data }] };

        case "write":
          await fs.writeFile(filePath, content || '', encoding);
          return { content: [{ type: "text", text: `File written: ${filePath}` }] };

        case "append":
          await fs.appendFile(filePath, content || '', encoding);
          return { content: [{ type: "text", text: `Content appended: ${filePath}` }] };

        case "list":
          const files = await fs.readdir(filePath);
          return { content: [{ type: "text", text: files.join('\n') }] };

        case "delete":
          await fs.unlink(filePath);
          return { content: [{ type: "text", text: `File deleted: ${filePath}` }] };

        case "mkdir":
          await fs.mkdir(filePath, { recursive: true });
          return { content: [{ type: "text", text: `Directory created: ${filePath}` }] };

        case "exists":
          try {
            await fs.access(filePath);
            return { content: [{ type: "text", text: `Path exists: ${filePath}` }] };
          } catch {
            return { content: [{ type: "text", text: `Path does not exist: ${filePath}` }] };
          }

        default:
          throw new Error(`Unsupported file operation: ${operation}`);
      }
    } catch (error) {
      return {
        content: [{ type: "text", text: `File operation failed: ${error.message}` }],
        isError: true
      };
    }
  }

  async aiCollaboration(args) {
    const { action, agentId, taskId, data = {} } = args;

    const redis = await this.initializeRedis();

    try {
      const key = taskId ? `ai_assistants:${agentId}:${taskId}` : `ai_assistants:${agentId}:*`;

      switch (action) {
        case "set_task":
          if (!taskId) throw new Error("taskId required for set_task");
          const taskData = {
            agentId,
            taskId,
            timestamp: Date.now(),
            status: data.status || 'active',
            ...data
          };
          await redis.setex(key, 3600, JSON.stringify(taskData));
          return { content: [{ type: "text", text: `Task set: ${agentId}:${taskId}` }] };

        case "get_task":
          if (!taskId) throw new Error("taskId required for get_task");
          const task = await redis.get(key);
          return { content: [{ type: "text", text: task ? JSON.parse(task) : "Task not found" }] };

        case "complete_task":
          if (!taskId) throw new Error("taskId required for complete_task");
          const currentTask = await redis.get(key);
          if (currentTask) {
            const updatedTask = { ...JSON.parse(currentTask), status: 'completed', completedAt: Date.now(), ...data };
            await redis.setex(key, 3600, JSON.stringify(updatedTask));
          }
          return { content: [{ type: "text", text: `Task completed: ${agentId}:${taskId}` }] };

        case "update_progress":
          if (!taskId) throw new Error("taskId required for update_progress");
          const existingTask = await redis.get(key);
          if (existingTask) {
            const updatedTask = { ...JSON.parse(existingTask), ...data, lastUpdated: Date.now() };
            await redis.setex(key, 3600, JSON.stringify(updatedTask));
          }
          return { content: [{ type: "text", text: `Progress updated: ${agentId}:${taskId}` }] };

        case "list_active":
          const allKeys = await redis.keys('ai_assistants:*:*:*');
          const activeTasks = [];

          for (const taskKey of allKeys) {
            const taskData = await redis.get(taskKey);
            if (taskData) {
              const parsed = JSON.parse(taskData);
              if (parsed.status === 'active') {
                activeTasks.push(parsed);
              }
            }
          }

          return { content: [{ type: "text", text: JSON.stringify(activeTasks, null, 2) }] };

        default:
          throw new Error(`Unsupported AI collaboration action: ${action}`);
      }
    } catch (error) {
      return {
        content: [{ type: "text", text: `AI collaboration failed: ${error.message}` }],
        isError: true
      };
    }
  }

  async dockerOperations(args) {
    const { operation, image, container, command, ports = [] } = args;

    let dockerCommand;
    switch (operation) {
      case "build":
        if (!image) throw new Error("Image name required for build");
        dockerCommand = `docker build -t ${image} .`;
        break;

      case "run":
        if (!image || !container) throw new Error("Image and container names required for run");
        const portMappings = ports.map(p => `-p ${p}`).join(' ');
        dockerCommand = `docker run -d --name ${container} ${portMappings} ${image} ${command || ''}`.trim();
        break;

      case "stop":
        if (!container) throw new Error("Container name required for stop");
        dockerCommand = `docker stop ${container}`;
        break;

      case "remove":
        if (!container) throw new Error("Container name required for remove");
        dockerCommand = `docker rm ${container}`;
        break;

      case "logs":
        if (!container) throw new Error("Container name required for logs");
        dockerCommand = `docker logs ${container}`;
        break;

      case "ps":
        dockerCommand = `docker ps -a`;
        break;

      case "images":
        dockerCommand = `docker images`;
        break;

      default:
        throw new Error(`Unsupported Docker operation: ${operation}`);
    }

    return await this.executeShellCommand({ command: dockerCommand });
  }

  async gitOperations(args) {
    const { operation, repo, files = [], message, branch } = args;

    let gitCommand;
    switch (operation) {
      case "status":
        gitCommand = "git status --porcelain";
        break;

      case "add":
        gitCommand = files.length > 0 ? `git add ${files.join(' ')}` : "git add .";
        break;

      case "commit":
        if (!message) throw new Error("Commit message required");
        gitCommand = `git commit -m "${message}"`;
        break;

      case "push":
        gitCommand = `git push origin ${branch || 'main'}`;
        break;

      case "pull":
        gitCommand = `git pull origin ${branch || 'main'}`;
        break;

      case "clone":
        if (!repo) throw new Error("Repository URL required for clone");
        gitCommand = `git clone ${repo}`;
        break;

      case "branch":
        gitCommand = branch ? `git checkout -b ${branch}` : "git branch -a";
        break;

      case "checkout":
        if (!branch) throw new Error("Branch name required for checkout");
        gitCommand = `git checkout ${branch}`;
        break;

      default:
        throw new Error(`Unsupported Git operation: ${operation}`);
    }

    return await this.executeShellCommand({ command: gitCommand });
  }

  async start() {
    console.error("🚀 Starting Roo Automation MCP Server...");
    console.error("Available tools:");
    console.error("  - execute_shell_command: Run shell commands");
    console.error("  - install_software: Install packages");
    console.error("  - redis_operations: Redis database operations");
    console.error("  - file_operations: File system operations");
    console.error("  - ai_collaboration: AI assistant coordination");
    console.error("  - docker_operations: Docker container management");
    console.error("  - git_operations: Git version control");

    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error("✅ Roo Automation MCP Server ready for requests");
  }

  async cleanup() {
    if (this.redis) {
      await this.redis.quit();
      console.error("🔌 Redis connection closed");
    }
  }
}

// Handle graceful shutdown
process.on('SIGINT', () => {
  console.error("\n🛑 Shutting down MCP server...");
  process.exit(0);
});

process.on('SIGTERM', () => {
  console.error("\n🛑 Shutting down MCP server...");
  process.exit(0);
});

// Start the server
const server = new RooAutomationServer();
server.start().catch((error) => {
  console.error("❌ Failed to start MCP server:", error);
  process.exit(1);
});