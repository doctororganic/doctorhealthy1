# Security Policy

## Supported Versions

We actively support the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |

## Reporting a Vulnerability

If you discover a security vulnerability, please do NOT create a public GitHub issue. Instead, please email security@yourdomain.com with:

- A description of the vulnerability
- Steps to reproduce
- Potential impact
- Suggested fix (if any)

We will respond within 48 hours.

## Security Best Practices

### Environment Variables
- Never commit `.env` files
- Use `.env.example` files for documentation
- Rotate secrets regularly
- Use strong, unique passwords

### API Keys
- Store API keys in environment variables
- Never hardcode credentials
- Use secret management services in production

### Database
- Use connection strings from environment variables
- Never commit database credentials
- Use read-only database users when possible

### Dependencies
- Keep dependencies up to date
- Review security advisories regularly
- Use `npm audit` and `go list -json -m all | nvd-cli` to check for vulnerabilities


