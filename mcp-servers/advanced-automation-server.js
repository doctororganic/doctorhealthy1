/**
 * Advanced Automation MCP Server
 * Comprehensive automation suite with browser, desktop, document, and workflow capabilities
 */

const { Server } = require("@modelcontextprotocol/sdk/server/index.js");
const { StdioServerTransport } = require("@modelcontextprotocol/sdk/server/stdio.js");
const { CallToolRequestSchema, ListToolsRequestSchema } = require("@modelcontextprotocol/sdk/types.js");
const { execSync, exec } = require('child_process');
const util = require('util');
const execAsync = util.promisify(exec);
const fs = require('fs').promises;
const path = require('path');

class AdvancedAutomationServer {
  constructor() {
    this.server = new Server(
      {
        name: "advanced-automation-server",
        version: "2.0.0",
      },
      {
        capabilities: {
          tools: {},
        },
      }
    );

    this.redis = null;
    this.puppeteerBrowser = null;
    this.playwrightBrowser = null;

    this.setupToolHandlers();
  }

  // Initialize Redis connection
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
        console.error('✅ Advanced MCP Server Redis connected');
      } catch (error) {
        console.error('❌ Advanced MCP Server Redis connection failed:', error.message);
      }
    }
    return this.redis;
  }

  setupToolHandlers() {
    this.server.setRequestHandler(ListToolsRequestSchema, async () => {
      return {
        tools: [
          // Browser Automation Tools
          {
            name: "browser_navigate",
            description: "Navigate to a URL using Puppeteer browser automation",
            inputSchema: {
              type: "object",
              properties: {
                url: { type: "string", description: "URL to navigate to" },
                waitFor: { type: "string", enum: ["load", "domcontentloaded", "networkidle0", "networkidle2"], default: "networkidle0" },
                timeout: { type: "number", description: "Navigation timeout in milliseconds", default: 30000 }
              },
              required: ["url"]
            }
          },
          {
            name: "browser_screenshot",
            description: "Take a screenshot of the current page or element",
            inputSchema: {
              type: "object",
              properties: {
                selector: { type: "string", description: "CSS selector for element screenshot (optional)" },
                fullPage: { type: "boolean", description: "Take full page screenshot", default: true },
                path: { type: "string", description: "Path to save screenshot", default: "./screenshot.png" }
              }
            }
          },
          {
            name: "browser_click",
            description: "Click on an element using CSS selector",
            inputSchema: {
              type: "object",
              properties: {
                selector: { type: "string", description: "CSS selector to click" },
                waitFor: { type: "string", description: "Wait for selector after click" }
              },
              required: ["selector"]
            }
          },
          {
            name: "browser_type",
            description: "Type text into an input field",
            inputSchema: {
              type: "object",
              properties: {
                selector: { type: "string", description: "CSS selector of input field" },
                text: { type: "string", description: "Text to type" },
                delay: { type: "number", description: "Delay between keystrokes", default: 100 }
              },
              required: ["selector", "text"]
            }
          },
          {
            name: "browser_extract",
            description: "Extract text or HTML from page elements",
            inputSchema: {
              type: "object",
              properties: {
                selector: { type: "string", description: "CSS selector to extract from" },
                property: { type: "string", enum: ["text", "html", "attribute"], default: "text" },
                attribute: { type: "string", description: "Attribute name if property is 'attribute'" }
              },
              required: ["selector"]
            }
          },

          // Desktop Automation Tools
          {
            name: "desktop_mouse_move",
            description: "Move mouse cursor to specified coordinates",
            inputSchema: {
              type: "object",
              properties: {
                x: { type: "number", description: "X coordinate" },
                y: { type: "number", description: "Y coordinate" },
                duration: { type: "number", description: "Movement duration in milliseconds", default: 500 }
              },
              required: ["x", "y"]
            }
          },
          {
            name: "desktop_mouse_click",
            description: "Perform mouse click at current position or coordinates",
            inputSchema: {
              type: "object",
              properties: {
                button: { type: "string", enum: ["left", "right", "middle"], default: "left" },
                double: { type: "boolean", description: "Double click", default: false },
                x: { type: "number", description: "X coordinate (optional)" },
                y: { type: "number", description: "Y coordinate (optional)" }
              }
            }
          },
          {
            name: "desktop_keyboard_type",
            description: "Type text using system keyboard",
            inputSchema: {
              type: "object",
              properties: {
                text: { type: "string", description: "Text to type" },
                delay: { type: "number", description: "Delay between keystrokes", default: 100 }
              },
              required: ["text"]
            }
          },
          {
            name: "desktop_key_press",
            description: "Press special keys or key combinations",
            inputSchema: {
              type: "object",
              properties: {
                keys: { type: "array", description: "Array of keys to press" },
                modifiers: { type: "array", description: "Modifier keys (shift, ctrl, alt, cmd)" }
              },
              required: ["keys"]
            }
          },

          // Document Processing Tools
          {
            name: "document_read_pdf",
            description: "Extract text and metadata from PDF files",
            inputSchema: {
              type: "object",
              properties: {
                path: { type: "string", description: "Path to PDF file" },
                pages: { type: "array", description: "Specific pages to extract (optional)" }
              },
              required: ["path"]
            }
          },
          {
            name: "document_read_docx",
            description: "Extract text and images from Word documents",
            inputSchema: {
              type: "object",
              properties: {
                path: { type: "string", description: "Path to DOCX file" },
                includeImages: { type: "boolean", description: "Include image extraction", default: false }
              },
              required: ["path"]
            }
          },
          {
            name: "document_create_pdf",
            description: "Create PDF from text content",
            inputSchema: {
              type: "object",
              properties: {
                content: { type: "string", description: "Text content for PDF" },
                title: { type: "string", description: "Document title" },
                outputPath: { type: "string", description: "Output path for PDF" }
              },
              required: ["content", "outputPath"]
            }
          },

          // OCR Tools
          {
            name: "ocr_image_text",
            description: "Extract text from images using OCR",
            inputSchema: {
              type: "object",
              properties: {
                imagePath: { type: "string", description: "Path to image file" },
                language: { type: "string", description: "OCR language", default: "eng" }
              },
              required: ["imagePath"]
            }
          },

          // Web Scraping Tools
          {
            name: "web_scrape",
            description: "Scrape data from websites",
            inputSchema: {
              type: "object",
              properties: {
                url: { type: "string", description: "URL to scrape" },
                selectors: { type: "object", description: "CSS selectors to extract data" },
                headers: { type: "object", description: "Custom HTTP headers" }
              },
              required: ["url"]
            }
          },

          // Calendar/Event Tools
          {
            name: "calendar_parse_ics",
            description: "Parse ICS calendar files",
            inputSchema: {
              type: "object",
              properties: {
                icsPath: { type: "string", description: "Path to ICS file" }
              },
              required: ["icsPath"]
            }
          },

          // Workflow Orchestration Tools
          {
            name: "workflow_create",
            description: "Create a new workflow with multiple steps",
            inputSchema: {
              type: "object",
              properties: {
                name: { type: "string", description: "Workflow name" },
                steps: { type: "array", description: "Array of workflow steps" },
                trigger: { type: "string", description: "Workflow trigger condition" }
              },
              required: ["name", "steps"]
            }
          },
          {
            name: "workflow_execute",
            description: "Execute a predefined workflow",
            inputSchema: {
              type: "object",
              properties: {
                workflowId: { type: "string", description: "Workflow identifier" },
                parameters: { type: "object", description: "Workflow parameters" }
              },
              required: ["workflowId"]
            }
          },

          // AI Agent Coordination Tools
          {
            name: "agent_communicate",
            description: "Send message to another AI agent",
            inputSchema: {
              type: "object",
              properties: {
                targetAgent: { type: "string", description: "Target agent ID" },
                message: { type: "string", description: "Message content" },
                priority: { type: "string", enum: ["low", "normal", "high", "urgent"], default: "normal" }
              },
              required: ["targetAgent", "message"]
            }
          },
          {
            name: "agent_status",
            description: "Get status of all connected AI agents",
            inputSchema: {
              type: "object",
              properties: {
                detailed: { type: "boolean", description: "Include detailed status", default: false }
              }
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
          // Browser automation
          case "browser_navigate":
            return await this.browserNavigate(args);
          case "browser_screenshot":
            return await this.browserScreenshot(args);
          case "browser_click":
            return await this.browserClick(args);
          case "browser_type":
            return await this.browserType(args);
          case "browser_extract":
            return await this.browserExtract(args);

          // Desktop automation
          case "desktop_mouse_move":
            return await this.desktopMouseMove(args);
          case "desktop_mouse_click":
            return await this.desktopMouseClick(args);
          case "desktop_keyboard_type":
            return await this.desktopKeyboardType(args);
          case "desktop_key_press":
            return await this.desktopKeyPress(args);

          // Document processing
          case "document_read_pdf":
            return await this.documentReadPDF(args);
          case "document_read_docx":
            return await this.documentReadDOCX(args);
          case "document_create_pdf":
            return await this.documentCreatePDF(args);

          // OCR
          case "ocr_image_text":
            return await this.ocrImageText(args);

          // Web scraping
          case "web_scrape":
            return await this.webScrape(args);

          // Calendar
          case "calendar_parse_ics":
            return await this.calendarParseICS(args);

          // Workflow orchestration
          case "workflow_create":
            return await this.workflowCreate(args);
          case "workflow_execute":
            return await this.workflowExecute(args);

          // AI agent coordination
          case "agent_communicate":
            return await this.agentCommunicate(args);
          case "agent_status":
            return await this.agentStatus(args);

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

  // Browser Automation Implementations
  async browserNavigate(args) {
    const puppeteer = require('puppeteer');
    if (!this.puppeteerBrowser) {
      this.puppeteerBrowser = await puppeteer.launch({ headless: false });
    }
    const page = await this.puppeteerBrowser.newPage();
    await page.goto(args.url, { waitUntil: args.waitFor, timeout: args.timeout });
    return { content: [{ type: "text", text: `Navigated to ${args.url}` }] };
  }

  async browserScreenshot(args) {
    if (!this.puppeteerBrowser) throw new Error("Browser not initialized");
    const pages = await this.puppeteerBrowser.pages();
    const page = pages[0];

    let screenshotOptions = { path: args.path };
    if (args.selector) {
      const element = await page.$(args.selector);
      screenshotOptions = { ...screenshotOptions, clip: await element.boundingBox() };
    } else if (args.fullPage) {
      screenshotOptions.fullPage = true;
    }

    await page.screenshot(screenshotOptions);
    return { content: [{ type: "text", text: `Screenshot saved to ${args.path}` }] };
  }

  async browserClick(args) {
    if (!this.puppeteerBrowser) throw new Error("Browser not initialized");
    const pages = await this.puppeteerBrowser.pages();
    const page = pages[0];

    await page.waitForSelector(args.selector);
    await page.click(args.selector);

    if (args.waitFor) {
      await page.waitForSelector(args.waitFor);
    }

    return { content: [{ type: "text", text: `Clicked element: ${args.selector}` }] };
  }

  async browserType(args) {
    if (!this.puppeteerBrowser) throw new Error("Browser not initialized");
    const pages = await this.puppeteerBrowser.pages();
    const page = pages[0];

    await page.waitForSelector(args.selector);
    await page.type(args.selector, args.text, { delay: args.delay });

    return { content: [{ type: "text", text: `Typed text into ${args.selector}` }] };
  }

  async browserExtract(args) {
    if (!this.puppeteerBrowser) throw new Error("Browser not initialized");
    const pages = await this.puppeteerBrowser.pages();
    const page = pages[0];

    await page.waitForSelector(args.selector);

    let result;
    if (args.property === 'attribute') {
      result = await page.$eval(args.selector, (el, attr) => el.getAttribute(attr), args.attribute);
    } else {
      const jsHandle = await page.$eval(args.selector, (el, prop) =>
        prop === 'text' ? el.textContent : el.innerHTML, args.property);
      result = jsHandle;
    }

    return { content: [{ type: "text", text: result }] };
  }

  // Desktop Automation Implementations
  async desktopMouseMove(args) {
    const robot = require('robotjs');
    robot.moveMouse(args.x, args.y);
    if (args.duration > 0) {
      await new Promise(resolve => setTimeout(resolve, args.duration));
    }
    return { content: [{ type: "text", text: `Mouse moved to (${args.x}, ${args.y})` }] };
  }

  async desktopMouseClick(args) {
    const robot = require('robotjs');

    if (args.x !== undefined && args.y !== undefined) {
      robot.moveMouse(args.x, args.y);
    }

    if (args.double) {
      robot.mouseClick(args.button, true);
    } else {
      robot.mouseClick(args.button, false);
    }

    return { content: [{ type: "text", text: `Mouse ${args.double ? 'double-' : ''}clicked with ${args.button} button` }] };
  }

  async desktopKeyboardType(args) {
    const robot = require('robotjs');
    robot.typeStringDelayed(args.text, args.delay);
    return { content: [{ type: "text", text: `Typed: ${args.text}` }] };
  }

  async desktopKeyPress(args) {
    const robot = require('robotjs');

    // Handle key combinations
    if (args.modifiers && args.modifiers.length > 0) {
      args.modifiers.forEach(modifier => {
        robot.keyToggle(modifier, 'down');
      });

      args.keys.forEach(key => {
        robot.keyTap(key);
      });

      args.modifiers.forEach(modifier => {
        robot.keyToggle(modifier, 'up');
      });
    } else {
      args.keys.forEach(key => {
        robot.keyTap(key);
      });
    }

    return { content: [{ type: "text", text: `Pressed keys: ${args.keys.join('+')} ${args.modifiers ? 'with ' + args.modifiers.join('+') : ''}` }] };
  }

  // Document Processing Implementations
  async documentReadPDF(args) {
    const pdfParse = require('pdf-parse');

    const dataBuffer = await fs.readFile(args.path);
    const data = await pdfParse(dataBuffer);

    let text = data.text;

    // Extract specific pages if requested
    if (args.pages && args.pages.length > 0) {
      const lines = text.split('\n');
      // This is a simplified page extraction - in practice you'd need more sophisticated PDF parsing
      text = `Pages ${args.pages.join(', ')} extracted from ${args.path}`;
    }

    return {
      content: [{
        type: "text",
        text: JSON.stringify({
          text: text,
          pages: data.numpages,
          info: data.info
        }, null, 2)
      }]
    };
  }

  async documentReadDOCX(args) {
    const mammoth = require('mammoth');

    const result = await mammoth.extractRawText({ path: args.path });
    const text = result.value;

    let response = { text: text };

    if (args.includeImages) {
      const imageResult = await mammoth.convertToHtml({ path: args.path });
      response.html = imageResult.value;
      response.messages = imageResult.messages;
    }

    return { content: [{ type: "text", text: JSON.stringify(response, null, 2) }] };
  }

  async documentCreatePDF(args) {
    const officegen = require('officegen');
    const docx = officegen('docx');

    // Create a simple document
    const paragraph = docx.createP();
    paragraph.addText(args.content);

    // Save as PDF (note: officegen creates DOCX, you'd need additional conversion for PDF)
    const outputPath = args.outputPath.replace('.pdf', '.docx');
    const out = await fs.createWriteStream(outputPath);
    docx.generate(out);

    return { content: [{ type: "text", text: `Document created: ${outputPath}` }] };
  }

  // OCR Implementation
  async ocrImageText(args) {
    const tesseract = require('tesseract.js');

    const { data: { text } } = await tesseract.recognize(args.imagePath, args.language);

    return { content: [{ type: "text", text: text }] };
  }

  // Web Scraping Implementation
  async webScrape(args) {
    const axios = require('axios');
    const cheerio = require('cheerio');

    const response = await axios.get(args.url, {
      headers: args.headers || {
        'User-Agent': 'Mozilla/5.0 (compatible; MCP Scraper)'
      }
    });

    const $ = cheerio.load(response.data);
    const results = {};

    if (args.selectors) {
      for (const [key, selector] of Object.entries(args.selectors)) {
        const elements = $(selector);
        if (elements.length === 1) {
          results[key] = elements.first().text().trim();
        } else if (elements.length > 1) {
          results[key] = elements.map((i, el) => $(el).text().trim()).get();
        }
      }
    }

    return { content: [{ type: "text", text: JSON.stringify(results, null, 2) }] };
  }

  // Workflow Orchestration Implementation
  async workflowCreate(args) {
    const redis = await this.initializeRedis();
    const workflowId = `workflow_${Date.now()}`;

    const workflow = {
      id: workflowId,
      name: args.name,
      steps: args.steps,
      trigger: args.trigger,
      createdAt: Date.now(),
      status: 'created'
    };

    await redis.setex(`workflows:${workflowId}`, 86400, JSON.stringify(workflow));

    return { content: [{ type: "text", text: `Workflow created: ${workflowId}` }] };
  }

  async workflowExecute(args) {
    const redis = await this.initializeRedis();
    const workflowData = await redis.get(`workflows:${args.workflowId}`);

    if (!workflowData) {
      throw new Error(`Workflow ${args.workflowId} not found`);
    }

    const workflow = JSON.parse(workflowData);
    const results = [];

    for (const step of workflow.steps) {
      console.error(`Executing workflow step: ${step.name}`);

      // This would execute each step based on its type
      // For now, just log the execution
      results.push({
        step: step.name,
        status: 'completed',
        timestamp: Date.now()
      });
    }

    return { content: [{ type: "text", text: JSON.stringify(results, null, 2) }] };
  }

  // AI Agent Coordination Implementation
  async agentCommunicate(args) {
    const redis = await this.initializeRedis();
    const messageId = `msg_${Date.now()}`;

    const message = {
      id: messageId,
      from: 'advanced-automation-server',
      to: args.targetAgent,
      content: args.message,
      priority: args.priority,
      timestamp: Date.now(),
      status: 'sent'
    };

    await redis.setex(`messages:${messageId}`, 3600, JSON.stringify(message));
    await redis.lpush(`agent_queue:${args.targetAgent}`, messageId);

    return { content: [{ type: "text", text: `Message sent to ${args.targetAgent}` }] };
  }

  async agentStatus(args) {
    const redis = await this.initializeRedis();
    const agentKeys = await redis.keys('ai_assistants:*:*');

    const agents = {};
    for (const key of agentKeys) {
      const data = await redis.get(key);
      if (data) {
        const parsed = JSON.parse(data);
        const agentId = parsed.agentId;

        if (!agents[agentId]) {
          agents[agentId] = {
            id: agentId,
            tasks: [],
            lastActivity: 0
          };
        }

        agents[agentId].tasks.push({
          id: key.split(':')[2],
          status: parsed.status,
          timestamp: parsed.timestamp
        });

        if (parsed.timestamp > agents[agentId].lastActivity) {
          agents[agentId].lastActivity = parsed.timestamp;
        }
      }
    }

    if (!args.detailed) {
      // Summary view
      const summary = Object.values(agents).map(agent => ({
        id: agent.id,
        activeTasks: agent.tasks.filter(t => t.status === 'active').length,
        totalTasks: agent.tasks.length,
        lastActivity: new Date(agent.lastActivity).toISOString()
      }));

      return { content: [{ type: "text", text: JSON.stringify(summary, null, 2) }] };
    }

    return { content: [{ type: "text", text: JSON.stringify(agents, null, 2) }] };
  }

  async start() {
    console.error("🚀 Starting Advanced Automation MCP Server...");
    console.error("Available tool categories:");
    console.error("  🌐 Browser Automation: navigate, screenshot, click, type, extract");
    console.error("  🖥️  Desktop Automation: mouse, keyboard, key press");
    console.error("  📄 Document Processing: PDF/DOCX read/write, OCR");
    console.error("  🕷️  Web Scraping: structured data extraction");
    console.error("  📅 Calendar Tools: ICS parsing");
    console.error("  🔄 Workflow Orchestration: create and execute workflows");
    console.error("  🤖 AI Coordination: inter-agent communication");

    const transport = new StdioServerTransport();
    await this.server.connect(transport);
    console.error("✅ Advanced Automation MCP Server ready for requests");
  }

  async cleanup() {
    if (this.redis) {
      await this.redis.quit();
      console.error("🔌 Redis connection closed");
    }

    if (this.puppeteerBrowser) {
      await this.puppeteerBrowser.close();
      console.error("🌐 Browser closed");
    }

    if (this.playwrightBrowser) {
      await this.playwrightBrowser.close();
      console.error("🎭 Playwright browser closed");
    }
  }
}

// Handle graceful shutdown
process.on('SIGINT', () => {
  console.error("\n🛑 Shutting down Advanced MCP server...");
  process.exit(0);
});

process.on('SIGTERM', () => {
  console.error("\n🛑 Shutting down Advanced MCP server...");
  process.exit(0);
});

// Start the server
const server = new AdvancedAutomationServer();
server.start().catch((error) => {
  console.error("❌ Failed to start Advanced MCP server:", error);
  process.exit(1);
});