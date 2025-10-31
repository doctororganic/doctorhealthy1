# Production-Ready Nutrition Platform - Optimized for Coolify
# Multi-stage build for security and performance

FROM node:18-alpine AS base
WORKDIR /app

# Install dumb-init for proper signal handling
RUN apk add --no-cache dumb-init curl

# Dependencies stage
FROM base AS dependencies
COPY production-nodejs/package*.json ./
RUN npm ci --only=production && \
    npm cache clean --force

# Production stage
FROM base AS production

# Set environment
ENV NODE_ENV=production \
    PORT=3000 \
    HOST=0.0.0.0

# Create non-root user for security
RUN addgroup -g 1001 -S nodejs && \
    adduser -S nodejs -u 1001

# Copy dependencies
COPY --from=dependencies --chown=nodejs:nodejs /app/node_modules ./node_modules

# Copy application code
COPY --chown=nodejs:nodejs production-nodejs/ ./

# Create necessary directories
RUN mkdir -p logs data && \
    chown -R nodejs:nodejs logs data

# Switch to non-root user
USER nodejs

# Expose port
EXPOSE 3000

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=40s --retries=3 \
    CMD node -e "require('http').get('http://localhost:3000/health', (r) => process.exit(r.statusCode === 200 ? 0 : 1))"

# Use dumb-init to handle signals properly
ENTRYPOINT ["dumb-init", "--"]

# Start application
CMD ["node", "server.js"]
