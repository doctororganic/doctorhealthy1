#!/usr/bin/env bash
set -euo pipefail

# Deployment script for Nutrition Platform
# Usage: ./scripts/deploy.sh [environment] [version]

ENVIRONMENT="${1:-staging}"
VERSION="${2:-latest}"
REGION="${AWS_REGION:-us-east-1}"
BACKEND_IMAGE="nutrition-platform/backend"
FRONTEND_IMAGE="nutrition-platform/frontend"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Logging functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# Validate environment
validate_environment() {
    if [[ ! "$ENVIRONMENT" =~ ^(staging|production)$ ]]; then
        log_error "Invalid environment: $ENVIRONMENT. Must be 'staging' or 'production'"
        exit 1
    fi
    
    log_info "Deploying to $ENVIRONMENT (version: $VERSION)"
}

# Check prerequisites
check_prerequisites() {
    log_info "Checking prerequisites..."
    
    # Check if Docker is installed
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed"
        exit 1
    fi
    
    # Check if kubectl is installed (for Kubernetes deployment)
    if ! command -v kubectl &> /dev/null; then
        log_warning "kubectl is not installed. Kubernetes deployment will be skipped."
    fi
    
    # Check AWS CLI (for AWS deployment)
    if ! command -v aws &> /dev/null; then
        log_warning "AWS CLI is not installed. AWS deployment will be skipped."
    fi
    
    log_success "Prerequisites check completed"
}

# Build backend
build_backend() {
    log_info "Building backend..."
    
    cd backend
    
    # Build Go binary
    CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o nutrition-platform .
    
    # Build Docker image
    docker build -t "${BACKEND_IMAGE}:${VERSION}" .
    docker tag "${BACKEND_IMAGE}:${VERSION}" "${BACKEND_IMAGE}:latest"
    
    cd ..
    log_success "Backend built successfully"
}

# Build frontend
build_frontend() {
    log_info "Building frontend..."
    
    cd frontend-nextjs
    
    # Install dependencies
    npm ci --production=false
    
    # Build Next.js app
    npm run build
    
    # Build Docker image
    docker build -t "${FRONTEND_IMAGE}:${VERSION}" .
    docker tag "${FRONTEND_IMAGE}:${VERSION}" "${FRONTEND_IMAGE}:latest"
    
    cd ..
    log_success "Frontend built successfully"
}

# Run tests
run_tests() {
    log_info "Running tests..."
    
    # Backend tests
    cd backend
    if ! go test ./... -v; then
        log_error "Backend tests failed"
        exit 1
    fi
    cd ..
    
    # Frontend tests
    cd frontend-nextjs
    if ! npm test -- --watchAll=false; then
        log_error "Frontend tests failed"
        exit 1
    fi
    cd ..
    
    log_success "All tests passed"
}

# Deploy to staging
deploy_staging() {
    log_info "Deploying to staging environment..."
    
    # Push images to registry (Docker Hub or ECR)
    if command -v aws &> /dev/null; then
        # AWS ECR deployment
        log_info "Pushing to AWS ECR..."
        
        # Get ECR login token
        aws ecr get-login-password --region "$REGION" | docker login --username AWS --password-stdin "$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$REGION.amazonaws.com"
        
        # Tag and push backend
        BACKEND_ECR="$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$REGION.amazonaws.com/${BACKEND_IMAGE}"
        docker tag "${BACKEND_IMAGE}:${VERSION}" "${BACKEND_ECR}:${VERSION}"
        docker tag "${BACKEND_IMAGE}:${VERSION}" "${BACKEND_ECR}:latest"
        docker push "${BACKEND_ECR}:${VERSION}"
        docker push "${BACKEND_ECR}:latest"
        
        # Tag and push frontend
        FRONTEND_ECR="$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$REGION.amazonaws.com/${FRONTEND_IMAGE}"
        docker tag "${FRONTEND_IMAGE}:${VERSION}" "${FRONTEND_ECR}:${VERSION}"
        docker tag "${FRONTEND_IMAGE}:${VERSION}" "${FRONTEND_ECR}:latest"
        docker push "${FRONTEND_ECR}:${VERSION}"
        docker push "${FRONTEND_ECR}:latest"
        
        # Deploy to Kubernetes
        if command -v kubectl &> /dev/null; then
            log_info "Applying Kubernetes manifests..."
            kubectl apply -f k8s/staging/ --recursive
            kubectl set image deployment/backend backend="${BACKEND_ECR}:${VERSION}" -n staging
            kubectl set image deployment/frontend frontend="${FRONTEND_ECR}:${VERSION}" -n staging
            kubectl rollout status deployment/backend -n staging
            kubectl rollout status deployment/frontend -n staging
        fi
    else
        log_warning "AWS CLI not found. Skipping ECR deployment."
    fi
    
    log_success "Staging deployment completed"
}

# Deploy to production
deploy_production() {
    log_info "Deploying to production environment..."
    
    # Additional safety checks for production
    log_warning "Deploying to PRODUCTION environment. Press Ctrl+C to cancel within 10 seconds..."
    sleep 10
    
    # Push images to registry
    if command -v aws &> /dev/null; then
        # AWS ECR deployment
        log_info "Pushing to AWS ECR..."
        
        # Get ECR login token
        aws ecr get-login-password --region "$REGION" | docker login --username AWS --password-stdin "$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$REGION.amazonaws.com"
        
        # Tag and push backend
        BACKEND_ECR="$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$REGION.amazonaws.com/${BACKEND_IMAGE}"
        docker tag "${BACKEND_IMAGE}:${VERSION}" "${BACKEND_ECR}:${VERSION}"
        docker tag "${BACKEND_IMAGE}:${VERSION}" "${BACKEND_ECR}:latest"
        docker push "${BACKEND_ECR}:${VERSION}"
        docker push "${BACKEND_ECR}:latest"
        
        # Tag and push frontend
        FRONTEND_ECR="$(aws sts get-caller-identity --query Account --output text).dkr.ecr.$REGION.amazonaws.com/${FRONTEND_IMAGE}"
        docker tag "${FRONTEND_IMAGE}:${VERSION}" "${FRONTEND_ECR}:${VERSION}"
        docker tag "${FRONTEND_IMAGE}:${VERSION}" "${FRONTEND_ECR}:latest"
        docker push "${FRONTEND_ECR}:${VERSION}"
        docker push "${FRONTEND_ECR}:latest"
        
        # Deploy to Kubernetes
        if command -v kubectl &> /dev/null; then
            log_info "Applying Kubernetes manifests..."
            kubectl apply -f k8s/production/ --recursive
            kubectl set image deployment/backend backend="${BACKEND_ECR}:${VERSION}" -n production
            kubectl set image deployment/frontend frontend="${FRONTEND_ECR}:${VERSION}" -n production
            kubectl rollout status deployment/backend -n production
            kubectl rollout status deployment/frontend -n production
        fi
    else
        log_warning "AWS CLI not found. Skipping ECR deployment."
    fi
    
    log_success "Production deployment completed"
}

# Run health checks
run_health_checks() {
    log_info "Running health checks..."
    
    # Determine the base URL based on environment
    if [ "$ENVIRONMENT" = "staging" ]; then
        BASE_URL="https://staging-api.nutrition-platform.com"
    else
        BASE_URL="https://api.nutrition-platform.com"
    fi
    
    # Check backend health
    if command -v curl &> /dev/null; then
        if curl -f -s "${BASE_URL}/health" > /dev/null; then
            log_success "Backend health check passed"
        else
            log_error "Backend health check failed"
            exit 1
        fi
        
        # Check frontend health
        FRONTEND_URL="${BASE_URL/api./}"
        if curl -f -s "${FRONTEND_URL}" > /dev/null; then
            log_success "Frontend health check passed"
        else
            log_error "Frontend health check failed"
            exit 1
        fi
    else
        log_warning "curl not found. Skipping health checks."
    fi
}

# Rollback function
rollback() {
    log_info "Rolling back deployment..."
    
    if [ "$ENVIRONMENT" = "staging" ]; then
        NAMESPACE="staging"
    else
        NAMESPACE="production"
    fi
    
    if command -v kubectl &> /dev/null; then
        # Rollback backend
        kubectl rollout undo deployment/backend -n "$NAMESPACE"
        kubectl rollout status deployment/backend -n "$NAMESPACE"
        
        # Rollback frontend
        kubectl rollout undo deployment/frontend -n "$NAMESPACE"
        kubectl rollout status deployment/frontend -n "$NAMESPACE"
        
        log_success "Rollback completed"
    else
        log_error "kubectl not found. Cannot rollback."
        exit 1
    fi
}

# Cleanup function
cleanup() {
    log_info "Cleaning up..."
    
    # Remove Docker images
    docker rmi "${BACKEND_IMAGE}:${VERSION}" 2>/dev/null || true
    docker rmi "${FRONTEND_IMAGE}:${VERSION}" 2>/dev/null || true
    
    log_success "Cleanup completed"
}

# Main deployment flow
main() {
    log_info "Starting deployment process..."
    
    validate_environment
    check_prerequisites
    
    # Build applications
    build_backend
    build_frontend
    
    # Run tests
    run_tests
    
    # Deploy based on environment
    if [ "$ENVIRONMENT" = "staging" ]; then
        deploy_staging
    else
        deploy_production
    fi
    
    # Run health checks
    run_health_checks
    
    # Cleanup
    cleanup
    
    log_success "Deployment completed successfully!"
}

# Handle script interruption
trap 'log_error "Deployment interrupted"; exit 1' INT TERM

# Handle rollback if requested
if [ "${3:-}" = "rollback" ]; then
    rollback
    exit 0
fi

# Run main function
main "$@"
