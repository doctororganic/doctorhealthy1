# Frontend Dockerfile for Nutrition Platform
FROM nginx:alpine

# Copy custom nginx configuration
COPY nginx.conf /etc/nginx/conf.d/default.conf

# Copy static files
COPY . /usr/share/nginx/html

# Create non-root user for security
RUN addgroup -g 1001 -S nginxuser && \
    adduser -S nginxuser -u 1001 -G nginxuser

# Create necessary directories
RUN mkdir -p /var/cache/nginx/client_temp \
             /var/cache/nginx/proxy_temp \
             /var/cache/nginx/fastcgi_temp \
             /var/cache/nginx/uwsgi_temp \
             /var/cache/nginx/scgi_temp && \
    chown -R nginxuser:nginxuser /var/cache/nginx && \
    chown -R nginxuser:nginxuser /usr/share/nginx/html && \
    chown -R nginxuser:nginxuser /etc/nginx/conf.d

# Switch to non-root user
USER nginxuser

# Expose port
EXPOSE 80

# Health check
HEALTHCHECK --interval=30s --timeout=10s --start-period=20s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost/ || exit 1

# Start nginx
CMD ["nginx", "-g", "daemon off;"]