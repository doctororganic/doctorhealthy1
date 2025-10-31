# Coolify MCP Server

A Model Context Protocol (MCP) server for interacting with Coolify API to manage deployments, projects, and servers.

## Features

- 🚀 **Project Management**: List, get details, and deploy projects
- 📦 **Deployment Control**: Initiate deployments with custom environments and branches
- 🖥️ **Server Management**: List and manage your Coolify servers
- 📊 **Deployment Monitoring**: Track deployment status and history
- 🔒 **Secure Authentication**: Bearer token authentication with your Coolify instance

## Quick Start

### 1. Configuration

Create a `.env` file in the project root:

```bash
COOLIFY_API_BASE_URL=https://api.doctorhealthy1.com/api/v1
COOLIFY_API_TOKEN=your_api_token_here
```

### 2. Installation & Build

```bash
npm install
npm run build
```

### 3. Usage

The server runs as an MCP server and can be integrated with MCP-compatible clients.

## Available Tools

### Project Management

#### `list_projects`
List all projects in your Coolify instance.

**Parameters:**
- `limit` (optional): Number of projects to return (default: 50, max: 100)
- `page` (optional): Page number for pagination (default: 1)

#### `get_project`
Get detailed information about a specific project.

**Parameters:**
- `project_id`: Project ID or UUID (required)

#### `deploy_project`
Initiate a deployment for a specific project.

**Parameters:**
- `project_id`: Project ID or UUID (required)
- `environment_name`: Environment name (optional, default: "production")
- `branch`: Branch to deploy (optional, default: "main")

### Deployment Management

#### `list_deployments`
List deployments, optionally filtered by project.

**Parameters:**
- `project_id` (optional): Filter by project ID
- `limit` (optional): Number of deployments to return (default: 50, max: 100)
- `page` (optional): Page number for pagination (default: 1)

#### `get_deployment`
Get detailed information about a specific deployment.

**Parameters:**
- `deployment_id`: Deployment ID or UUID (required)

### Server Management

#### `list_servers`
List all servers in your Coolify instance.

**Parameters:**
- `limit` (optional): Number of servers to return (default: 50, max: 100)
- `page` (optional): Page number for pagination (default: 1)

#### `get_server`
Get detailed information about a specific server.

**Parameters:**
- `server_id`: Server ID or UUID (required)

## Configuration

### Environment Variables

| Variable | Description | Required | Default |
|----------|-------------|----------|---------|
| `COOLIFY_API_BASE_URL` | Your Coolify API base URL | Yes | - |
| `COOLIFY_API_TOKEN` | Your Coolify API token | Yes | - |

### Getting Your API Token

1. Log into your Coolify dashboard
2. Navigate to Settings > API Keys
3. Create a new API key or use an existing one
4. Copy the token and add it to your `.env` file

## Development

### Building

```bash
npm run build
```

### Testing

```bash
node test-server.js
```

## Error Handling

The server includes comprehensive error handling for:
- Network connectivity issues
- Invalid API responses
- Authentication failures
- Missing or invalid parameters

All errors are returned with descriptive messages to help with troubleshooting.

## Security Best Practices

- Store API tokens securely (use `.env` files)
- Use read-only tokens when possible for querying operations
- Regularly rotate your API tokens
- Never commit API tokens to version control

## Troubleshooting

### Common Issues

1. **Authentication Failed**: Check your API token and base URL
2. **Connection Refused**: Verify your Coolify instance is accessible
3. **Invalid Project ID**: Ensure the project ID exists in your Coolify instance

### Getting Help

- Check the Coolify documentation for API reference
- Verify your API token has the necessary permissions
- Ensure your Coolify instance is running and accessible

## License

This project is part of the DoctorHealthy1 nutrition platform.