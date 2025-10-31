# Nutrition Platform Deployment Information
# Generated: 2025-10-13

## Application Details
- Name: nutrition-platform-secure
- Description: AI-powered nutrition platform with enterprise security
- Version: 1.0.0

## Deployment Configuration
- Build Pack: Dockerfile
- Dockerfile Location: backend/Dockerfile
- Build Context: ./
- Start Command: (will be auto-detected)
- Port: 8080

## Environment Variables
All environment variables are in .env.production file

## Services Required
- PostgreSQL 15
- Redis 7-alpine

## Domain
- Primary: super.doctorhealthy1.com
- Secondary: my.doctorhealthy1.com

## Health Check
- Path: /health
- Interval: 30s
