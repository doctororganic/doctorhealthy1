# Pre-Push Security Checklist

Before pushing to GitHub, verify:

## ✅ Security Checks

- [ ] No `.env` files committed (only `.env.example` files)
- [ ] No API keys or secrets in code
- [ ] No database credentials hardcoded
- [ ] No JWT secrets in code
- [ ] No AWS/cloud credentials
- [ ] No SSL certificates or private keys
- [ ] No database files (`.db`, `.sqlite`)
- [ ] No log files with sensitive data
- [ ] No build artifacts with secrets

## ✅ Files to Verify

Run these commands before pushing:

```bash
# Check for sensitive files
git ls-files | grep -E "\.(env|key|pem|secret|db|sqlite)"

# Should only show .env.example files, nothing else

# Check staged files
git diff --cached --name-only | grep -E "\.(env|key|pem|secret)"

# Should be empty

# Verify .gitignore is working
git check-ignore -v .env backend/.env frontend-nextjs/.env.local

# Should show all .env files are ignored
```

## ✅ Safe to Push

If all checks pass, you're safe to push!

