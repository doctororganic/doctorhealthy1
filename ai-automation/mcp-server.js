/**
 * AI Automation MCP Server
 * Provides comprehensive tools for AI agent coordination and automation
 */

const { McpServer } = require('@modelcontextprotocol/sdk/server/mcp.js');
const { StdioServerTransport } = require('@modelcontextprotocol/sdk/server/stdio.js');
const { z } = require('zod');
const puppeteer = require('puppeteer');
const fs = require('fs').promises;
const path = require('path');
const { exec } = require('child_process');
const util = require('util');
const execAsync = util.promisify(exec);

// Import our AI collaboration modules
const sharedMemory = require('../ai-collaboration/file-shared-memory');
const workflowCoordinator = require('../ai-collaboration/workflow-coordinator');
const roleAssignmentManager = require('../ai-collaboration/role-assignments');
const kiloTestingFramework = require('../ai-collaboration/kilo-testing-framework');
const rooValidationSystem = require('../ai-collaboration/roo-validation-system');
const orchestrationController = require('../ai-collaboration/orchestration-controller');

class AIAutomationMCPServer {
  constructor() {
    this.server = new McpServer({
      name: "ai-automation-server",
      version: "1.0.0"
    });

    this.browser = null;
    this.setupTools();
  }

  /**
   * Set up all MCP tools for AI automation
   */
  setupTools() {
    // AI Agent Coordination Tools
    this.setupAICoordinationTools();

    // Browser Automation Tools
    this.setupBrowserAutomationTools();

    // Document Processing Tools
    this.setupDocumentProcessingTools();

    // System Automation Tools
    this.setupSystemAutomationTools();

    // Development Workflow Tools
    this.setupDevelopmentWorkflowTools();
  }

  /**
   * AI Agent Coordination Tools
   */
  setupAICoordinationTools() {
    // Get AI collaboration status
    this.server.tool(
      "get_ai_collaboration_status",
      "Get current status of AI agent collaboration and workflows",
      {},
      async () => {
        try {
          const status = await sharedMemory.getCollaborationStatus();
          const workflowStatus = await workflowCoordinator.getWorkflowStatus();
          const orchestrationStatus = await orchestrationController.getOrchestrationStatus();

          return {
            content: [{
              type: "text",
              text: JSON.stringify({
                collaboration: status,
                workflow: workflowStatus,
                orchestration: orchestrationStatus,
                timestamp: new Date().toISOString()
              }, null, 2)
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );

    // Assign task to AI agent
    this.server.tool(
      "assign_ai_task",
      "Assign a specific task to an AI agent",
      {
        agentId: z.string().describe("AI agent ID (kilo, roo, codesupernova)"),
        taskId: z.string().describe("Unique task identifier"),
        taskData: z.object({
          type: z.string(),
          description: z.string(),
          priority: z.enum(["low", "medium", "high", "critical"]).optional(),
          estimatedHours: z.number().optional()
        }).describe("Task details and metadata")
      },
      async ({ agentId, taskId, taskData }) => {
        try {
          const result = await workflowCoordinator.assignTask(agentId, taskId, taskData);
          return {
            content: [{
              type: "text",
              text: `Task assigned successfully: ${JSON.stringify(result, null, 2)}`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );

    // Start workflow stage
    this.server.tool(
      "start_workflow_stage",
      "Start a specific workflow stage",
      {
        stageName: z.string().describe("Name of the workflow stage to start")
      },
      async ({ stageName }) => {
        try {
          const result = await workflowCoordinator.startWorkflowStage(stageName);
          return {
            content: [{
              type: "text",
              text: `Workflow stage started: ${JSON.stringify(result, null, 2)}`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );

    // Execute test suite
    this.server.tool(
      "execute_test_suite",
      "Execute comprehensive test suite for a component",
      {
        suiteName: z.enum(["frontend", "backend", "integration", "performance"]).describe("Test suite to execute"),
        agentId: z.string().optional().describe("AI agent executing the tests")
      },
      async ({ suiteName, agentId = "kilo" }) => {
        try {
          const result = await kiloTestingFramework.executeTestSuite(suiteName, agentId);
          return {
            content: [{
              type: "text",
              text: `Test suite executed: ${JSON.stringify(result, null, 2)}`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );

    // Perform validation
    this.server.tool(
      "perform_validation",
      "Perform comprehensive validation on a component",
      {
        validationType: z.enum(["security", "performance", "codeQuality", "architecture"]).describe("Type of validation to perform"),
        component: z.string().describe("Component to validate"),
        agentId: z.string().optional().describe("AI agent performing validation")
      },
      async ({ validationType, component, agentId = "roo" }) => {
        try {
          const result = await rooValidationSystem.performValidation(validationType, component, agentId);
          return {
            content: [{
              type: "text",
              text: `Validation completed: ${JSON.stringify(result, null, 2)}`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );
  }

  /**
   * Browser Automation Tools
   */
  setupBrowserAutomationTools() {
    // Launch browser
    this.server.tool(
      "launch_browser",
      "Launch a new browser instance for automation",
      {
        headless: z.boolean().optional().describe("Run in headless mode"),
        url: z.string().optional().describe("Initial URL to navigate to")
      },
      async ({ headless = true, url }) => {
        try {
          this.browser = await puppeteer.launch({
            headless,
            args: ['--no-sandbox', '--disable-setuid-sandbox']
          });

          if (url) {
            const page = await this.browser.newPage();
            await page.goto(url);
          }

          return {
            content: [{
              type: "text",
              text: `Browser launched successfully. Headless: ${headless}, URL: ${url || 'none'}`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );

    // Navigate to URL
    this.server.tool(
      "navigate_browser",
      "Navigate browser to a specific URL",
      {
        url: z.string().describe("URL to navigate to")
      },
      async ({ url }) => {
        try {
          if (!this.browser) {
            throw new Error("Browser not launched. Use launch_browser first.");
          }

          const page = await this.browser.newPage();
          await page.goto(url);

          return {
            content: [{
              type: "text",
              text: `Navigated to: ${url}`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );

    // Take screenshot
    this.server.tool(
      "take_screenshot",
      "Take a screenshot of the current browser page",
      {
        filename: z.string().optional().describe("Filename for the screenshot")
      },
      async ({ filename = "screenshot.png" }) => {
        try {
          if (!this.browser) {
            throw new Error("Browser not launched. Use launch_browser first.");
          }

          const pages = await this.browser.pages();
          const page = pages[pages.length - 1];

          await page.screenshot({ path: filename, fullPage: true });

          return {
            content: [{
              type: "text",
              text: `Screenshot saved: ${filename}`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );

    // Close browser
    this.server.tool(
      "close_browser",
      "Close the browser instance",
      {},
      async () => {
        try {
          if (this.browser) {
            await this.browser.close();
            this.browser = null;
          }

          return {
            content: [{
              type: "text",
              text: "Browser closed successfully"
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );
  }

  /**
   * Document Processing Tools
   */
  setupDocumentProcessingTools() {
    // Read PDF file
    this.server.tool(
      "read_pdf",
      "Extract text content from a PDF file",
      {
        filePath: z.string().describe("Path to the PDF file")
      },
      async ({ filePath }) => {
        try {
          const pdf = require('pdf-parse');
          const dataBuffer = await fs.readFile(filePath);
          const data = await pdf(dataBuffer);

          return {
            content: [{
              type: "text",
              text: `PDF content extracted (${data.text.length} characters): ${data.text.substring(0, 1000)}...`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );

    // Read Word document
    this.server.tool(
      "read_word_doc",
      "Extract text content from a Word document",
      {
        filePath: z.string().describe("Path to the Word document")
      },
      async ({ filePath }) => {
        try {
          const mammoth = require('mammoth');
          const result = await mammoth.extractRawText({ path: filePath });

          return {
            content: [{
              type: "text",
              text: `Word document content extracted (${result.value.length} characters): ${result.value.substring(0, 1000)}...`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );

    // Read any text file
    this.server.tool(
      "read_text_file",
      "Read and return content of any text file",
      {
        filePath: z.string().describe("Path to the text file")
      },
      async ({ filePath }) => {
        try {
          const content = await fs.readFile(filePath, 'utf8');
          return {
            content: [{
              type: "text",
              text: content
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );
  }

  /**
   * System Automation Tools
   */
  setupSystemAutomationTools() {
    // Execute shell command
    this.server.tool(
      "execute_command",
      "Execute a shell command on the system",
      {
        command: z.string().describe("Command to execute"),
        cwd: z.string().optional().describe("Working directory for command execution")
      },
      async ({ command, cwd }) => {
        try {
          const options = cwd ? { cwd } : {};
          const { stdout, stderr } = await execAsync(command, options);

          return {
            content: [{
              type: "text",
              text: `Command executed successfully:\nSTDOUT:\n${stdout}\nSTDERR:\n${stderr}`
            }]
          };
        } catch (error) {
          return {
            content: [{
              type: "text",
              text: `Command failed:\nSTDOUT: ${error.stdout}\nSTDERR: ${error.stderr}\nError: ${error.message}`
            }],
            isError: true
          };
        }
      }
    );

    // List directory contents
    this.server.tool(
      "list_directory",
      "List contents of a directory",
      {
        dirPath: z.string().optional().describe("Directory path to list"),
        recursive: z.boolean().optional().describe("List recursively")
      },
      async ({ dirPath = ".", recursive = false }) => {
        try {
          const command = recursive ? `find "${dirPath}" -type f` : `ls -la "${dirPath}"`;
          const { stdout } = await execAsync(command);

          return {
            content: [{
              type: "text",
              text: `Directory listing:\n${stdout}`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );

    // Create file
    this.server.tool(
      "create_file",
      "Create a new file with specified content",
      {
        filePath: z.string().describe("Path for the new file"),
        content: z.string().describe("Content to write to the file")
      },
      async ({ filePath, content }) => {
        try {
          await fs.writeFile(filePath, content);
          return {
            content: [{
              type: "text",
              text: `File created successfully: ${filePath}`
            }]
          };
        } catch (error) {
          return {
            content: [{ type: "text", text: `Error: ${error.message}` }],
            isError: true
          };
        }
      }
    );
  }

  /**
   * Development Workflow Tools
   */
  setupDevelopmentWorkflowTools() {
    // Deploy application
    this.server.tool(
      "deploy_application",
      "Deploy the nutrition platform application",
      {
        environment: z.enum(["development", "staging", "production"]).optional().describe("Deployment environment"),
        components: z.array(z.string()).optional().describe("Specific components to deploy")
      },
      async ({ environment = "development", components }) => {
        try {
          const deployScript = path.join(__dirname, '../deploy-nodejs.sh');
          const command = `bash ${deployScript} deploy`;

          const { stdout } = await execAsync(command);

          return {
            content: [{
              type: "text",
              text: `Deployment completed:\n${stdout}`
            }]
          };
        } catch (error) {
          return {
            content: [{
              type: "text",
              text: `Deployment failed:\n${error.stdout}\n${error.stderr}`
            }],
            isError: true
          };
        }
      }
    );

    // Run tests
    this.server.tool(
      "run_tests",
      "Execute test suite for the application",
      {
        testType: z.enum(["unit", "integration", "e2e", "all"]).optional().describe("Type of tests to run")
      },
      async ({ testType = "all" }) => {
        try {
          const testScript = path.join(__dirname, '../test-ai-automation.js');
          const { stdout } = await execAsync(`node ${testScript}`);

          return {
            content: [{
              type: "text",
              text: `Tests completed:\n${stdout}`
            }]
          };
        } catch (error) {
          return {
            content: [{
              type: "text",
              text: `Tests failed:\n${error.stdout}\n${error.stderr}`
            }],
            isError: true
          };
        }
      }
    );

    // Start monitoring
    this.server.tool(
      "start_monitoring",
      "Start the complete monitoring system",
      {},
      async () => {
        try {
          const monitorScript = path.join(__dirname, '../deploy-monitoring.sh');
          const { stdout } = await execAsync(`bash ${monitorScript} deploy`);

          return {
            content: [{
              type: "text",
              text: `Monitoring system started:\n${stdout}`
            }]
          };
        } catch (error) {
          return {
            content: [{
              type: "text",
              text: `Failed to start monitoring:\n${error.stdout}\n${error.stderr}`
            }],
            isError: true
          };
        }
      }
    );
  }

  /**
   * Start the MCP server
   */
  async start() {
    try {
      const transport = new StdioServerTransport();
      await this.server.connect(transport);
      console.error('🤖 AI Automation MCP Server running on stdio');
      console.error('🔧 Available tools:');
      console.error('  - AI Agent Coordination (collaboration, workflow, tasks)');
      console.error('  - Browser Automation (launch, navigate, screenshot)');
      console.error('  - Document Processing (PDF, Word, text files)');
      console.error('  - System Automation (commands, file operations)');
      console.error('  - Development Workflow (deploy, test, monitor)');
    } catch (error) {
      console.error('Failed to start MCP server:', error);
      process.exit(1);
    }
  }

  /**
   * Stop the MCP server
   */
  async stop() {
    try {
      if (this.browser) {
        await this.browser.close();
      }
      process.exit(0);
    } catch (error) {
      console.error('Error during shutdown:', error);
      process.exit(1);
    }
  }
}

// Handle graceful shutdown
process.on('SIGINT', async () => {
  console.error('Received SIGINT, shutting down gracefully...');
  const server = new AIAutomationMCPServer();
  await server.stop();
});

process.on('SIGTERM', async () => {
  console.error('Received SIGTERM, shutting down gracefully...');
  const server = new AIAutomationMCPServer();
  await server.stop();
});

// Start the server if this file is run directly
if (require.main === module) {
  const server = new AIAutomationMCPServer();
  server.start().catch(error => {
    console.error('Failed to start AI automation MCP server:', error);
    process.exit(1);
  });
}

module.exports = AIAutomationMCPServer;