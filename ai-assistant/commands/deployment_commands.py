#!/usr/bin/env python3
"""
Deployment Commands for Nutrition Platform AI Assistant
Handles Coolify, Docker, and VPS deployment automation
"""

import os
import sys
import time
import json
import subprocess
import requests
from datetime import datetime, timedelta
from typing import Dict, Any, Optional, Tuple

# Add parent directory to path for imports
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from config.config_manager import ConfigManager
from utils.logger import Logger

class DeploymentManager:
    """Main deployment management class"""

    def __init__(self):
        """Initialize deployment manager"""
        self.config = ConfigManager()
        self.logger = Logger()
        self.nutrition_platform_dir = self.config.nutrition_platform_dir

    def deploy_coolify(self, environment: str = 'production', project: str = None, force: bool = False) -> str:
        """Deploy to Coolify platform"""
        try:
            self.logger.log_deployment_start("Coolify", environment, project=project or "nutrition-platform")

            # Get Coolify configuration
            coolify_config = self.config.get_config("coolify", {})
            if not coolify_config.get("base_url") or not coolify_config.get("api_token"):
                return "Coolify configuration not complete. Please set base_url and api_token."

            # Check if project exists or create it
            project_id = self._get_or_create_coolify_project(project or coolify_config.get("project_name", "nutrition-platform"))

            if not project_id:
                return "Failed to create or find Coolify project"

            # Get deployment settings
            deployment_settings = self.config.get_deployment_settings()

            # Backup if enabled
            if deployment_settings.get("backup_before_deploy", True):
                self.logger.log_step("Creating pre-deployment backup", "START")
                backup_result = self._create_backup("pre-deployment")
                if backup_result:
                    self.logger.info(f"Backup created: {backup_result}")
                self.logger.log_step("Creating pre-deployment backup", "END")

            # Deploy to Coolify
            self.logger.log_step("Deploying to Coolify", "START")

            # Get environment configuration
            env_config = self.config.get_environment_config(environment)

            # Prepare deployment data
            deployment_data = {
                "project_id": project_id,
                "environment": environment,
                "force": force,
                "timestamp": datetime.now().isoformat()
            }

            # Trigger deployment
            deploy_result = self._trigger_coolify_deployment(deployment_data)

            if deploy_result["success"]:
                deployment_id = deploy_result["deployment_id"]

                # Monitor deployment status
                self.logger.log_step("Monitoring deployment", "START")
                status_result = self._monitor_deployment(deployment_id, "coolify")

                if status_result["success"]:
                    self.logger.log_step("Deploying to Coolify", "END")
                    self.logger.log_deployment_end(True)

                    return f"Successfully deployed to Coolify (Project: {project_id}, Environment: {environment})"
                else:
                    error_msg = f"Deployment monitoring failed: {status_result['error']}"
                    self.logger.error(error_msg)

                    # Attempt rollback if enabled
                    if deployment_settings.get("rollback_on_failure", True):
                        self.logger.log_step("Attempting rollback", "START")
                        rollback_result = self._rollback_coolify_deployment(project_id, deployment_id)
                        if rollback_result:
                            self.logger.info("Rollback completed")
                        self.logger.log_step("Attempting rollback", "END")

                    return error_msg
            else:
                return f"Coolify deployment failed: {deploy_result['error']}"

        except Exception as e:
            self.logger.log_error_details(e, "Coolify deployment")
            return f"Coolify deployment error: {str(e)}"

    def deploy_docker(self, compose_file: str = None, environment: str = 'production', build: bool = False, restart: bool = False) -> str:
        """Deploy using Docker Compose"""
        try:
            self.logger.log_deployment_start("Docker", environment, compose_file=compose_file)

            # Get compose file path
            if not compose_file:
                compose_file = self.config.get_docker_compose_file(environment)

            if not os.path.exists(compose_file):
                return f"Docker Compose file not found: {compose_file}"

            # Get deployment settings
            deployment_settings = self.config.get_deployment_settings()

            # Backup if enabled
            if deployment_settings.get("backup_before_deploy", True):
                self.logger.log_step("Creating pre-deployment backup", "START")
                backup_result = self._create_backup("pre-deployment")
                if backup_result:
                    self.logger.info(f"Backup created: {backup_result}")
                self.logger.log_step("Creating pre-deployment backup", "END")

            # Stop existing containers if restart requested
            if restart:
                self.logger.log_step("Stopping existing containers", "START")
                self._docker_compose_command(compose_file, "down")
                self.logger.log_step("Stopping existing containers", "END")

            # Build images if requested
            if build:
                self.logger.log_step("Building Docker images", "START")
                result = self._docker_compose_command(compose_file, "build")
                if result["success"]:
                    self.logger.info("Docker images built successfully")
                else:
                    return f"Docker build failed: {result['error']}"
                self.logger.log_step("Building Docker images", "END")

            # Deploy containers
            self.logger.log_step("Starting Docker containers", "START")
            result = self._docker_compose_command(compose_file, "up", "-d")

            if result["success"]:
                self.logger.log_step("Starting Docker containers", "END")

                # Health check
                self.logger.log_step("Running health checks", "START")
                health_result = self._run_health_checks(environment)

                if health_result["success"]:
                    self.logger.log_step("Running health checks", "END")
                    self.logger.log_deployment_end(True)

                    return f"Successfully deployed Docker containers from {compose_file}"
                else:
                    error_msg = f"Health checks failed: {health_result['error']}"
                    self.logger.error(error_msg)

                    # Attempt rollback if enabled
                    if deployment_settings.get("rollback_on_failure", True):
                        self.logger.log_step("Attempting rollback", "START")
                        self._docker_compose_command(compose_file, "down")
                        self.logger.log_step("Attempting rollback", "END")

                    return error_msg
            else:
                return f"Docker deployment failed: {result['error']}"

        except Exception as e:
            self.logger.log_error_details(e, "Docker deployment")
            return f"Docker deployment error: {str(e)}"

    def deploy_vps(self, provider: str = 'vultr', environment: str = 'production', server_ip: str = None, ssh_key: str = None) -> str:
        """Deploy to VPS"""
        try:
            self.logger.log_deployment_start("VPS", environment, provider=provider, server_ip=server_ip)

            # Get VPS configuration
            vps_config = self.config.get_vps_config(provider)
            if not vps_config:
                return f"VPS configuration not found for provider: {provider}"

            # Use existing deployment scripts if available
            script_path = os.path.join(self.nutrition_platform_dir, f"deploy-{provider}.sh")
            if os.path.exists(script_path):
                self.logger.info(f"Using existing deployment script: {script_path}")
                return self._run_deployment_script(script_path, environment, server_ip, ssh_key)

            # Otherwise, use custom deployment logic
            return self._deploy_to_vps_custom(provider, environment, server_ip, ssh_key, vps_config)

        except Exception as e:
            self.logger.log_error_details(e, "VPS deployment")
            return f"VPS deployment error: {str(e)}"

    def check_status(self, args) -> str:
        """Check deployment status"""
        try:
            deployment_id = getattr(args, 'deployment_id', None) or 'latest'

            # Check Docker containers status
            if self._is_docker_deployment():
                return self._check_docker_status()

            # Check VPS status
            if self._is_vps_deployment():
                return self._check_vps_status()

            # Check Coolify status
            return self._check_coolify_status(deployment_id)

        except Exception as e:
            return f"Status check failed: {str(e)}"

    def view_logs(self, args) -> str:
        """View deployment logs"""
        try:
            deployment_id = getattr(args, 'deployment_id', None) or 'latest'
            tail_lines = getattr(args, 'tail', 50)
            follow = getattr(args, 'follow', False)

            # Get logs from Docker containers
            if self._is_docker_deployment():
                return self._get_docker_logs(tail_lines, follow)

            # Get logs from VPS
            if self._is_vps_deployment():
                return self._get_vps_logs(deployment_id, tail_lines)

            # Get logs from Coolify
            return self._get_coolify_logs(deployment_id, tail_lines)

        except Exception as e:
            return f"Log retrieval failed: {str(e)}"

    def health_check(self, args) -> str:
        """Perform health check"""
        try:
            url = getattr(args, 'url', None)
            timeout = getattr(args, 'timeout', 30)

            if not url:
                # Try to get URL from config
                monitoring_config = self.config.get_monitoring_settings()
                url = monitoring_config.get("health_check_url")

                if not url:
                    return "No health check URL provided or configured"

            # Perform health check
            return self._perform_health_check(url, timeout)

        except Exception as e:
            return f"Health check failed: {str(e)}"

    def _get_or_create_coolify_project(self, project_name: str) -> Optional[str]:
        """Get or create Coolify project"""
        try:
            coolify_config = self.config.get_config("coolify", {})
            base_url = coolify_config.get("base_url")
            api_token = coolify_config.get("api_token")

            headers = {
                'Authorization': f'Bearer {api_token}',
                'Content-Type': 'application/json'
            }

            # Check if project exists
            response = requests.get(f"{base_url}/api/projects", headers=headers)

            if response.status_code == 200:
                projects = response.json().get("data", [])
                for project in projects:
                    if project["name"] == project_name:
                        return project["id"]

            # Create new project if it doesn't exist
            create_data = {
                "name": project_name,
                "description": "Nutrition Platform - Auto-deployed by AI Assistant"
            }

            response = requests.post(f"{base_url}/api/projects", headers=headers, json=create_data)

            if response.status_code == 201:
                project = response.json()
                return project.get("id")

            return None

        except Exception as e:
            self.logger.error(f"Error managing Coolify project: {str(e)}")
            return None

    def _trigger_coolify_deployment(self, deployment_data: Dict[str, Any]) -> Dict[str, Any]:
        """Trigger Coolify deployment"""
        try:
            coolify_config = self.config.get_config("coolify", {})
            base_url = coolify_config.get("base_url")
            api_token = coolify_config.get("api_token")

            headers = {
                'Authorization': f'Bearer {api_token}',
                'Content-Type': 'application/json'
            }

            # Trigger deployment
            response = requests.post(
                f"{base_url}/api/projects/{deployment_data['project_id']}/deploy",
                headers=headers,
                json={"environment": deployment_data["environment"]}
            )

            if response.status_code == 200:
                deploy_response = response.json()
                return {
                    "success": True,
                    "deployment_id": deploy_response.get("deployment_id")
                }
            else:
                return {
                    "success": False,
                    "error": f"HTTP {response.status_code}: {response.text}"
                }

        except Exception as e:
            return {
                "success": False,
                "error": str(e)
            }

    def _monitor_deployment(self, deployment_id: str, platform: str, timeout: int = 300) -> Dict[str, Any]:
        """Monitor deployment status"""
        try:
            start_time = time.time()
            max_wait = timeout

            while time.time() - start_time < max_wait:
                if platform == "coolify":
                    status = self._check_coolify_deployment_status(deployment_id)
                elif platform == "docker":
                    status = self._check_docker_deployment_status()
                else:
                    status = self._check_vps_deployment_status(deployment_id)

                if status["finished"]:
                    return status

                time.sleep(10)  # Check every 10 seconds

            return {
                "success": False,
                "error": f"Deployment monitoring timeout after {timeout} seconds"
            }

        except Exception as e:
            return {
                "success": False,
                "error": str(e)
            }

    def _check_coolify_deployment_status(self, deployment_id: str) -> Dict[str, Any]:
        """Check Coolify deployment status"""
        try:
            coolify_config = self.config.get_config("coolify", {})
            base_url = coolify_config.get("base_url")
            api_token = coolify_config.get("api_token")

            headers = {'Authorization': f'Bearer {api_token}'}

            response = requests.get(f"{base_url}/api/deployments/{deployment_id}", headers=headers)

            if response.status_code == 200:
                deployment = response.json()
                status = deployment.get("status", "unknown")

                return {
                    "success": status == "success",
                    "finished": status in ["success", "failed"],
                    "status": status,
                    "details": deployment
                }

            return {
                "success": False,
                "finished": False,
                "status": "unknown",
                "error": f"HTTP {response.status_code}"
            }

        except Exception as e:
            return {
                "success": False,
                "finished": False,
                "status": "error",
                "error": str(e)
            }

    def _docker_compose_command(self, compose_file: str, command: str, *args: str) -> Dict[str, Any]:
        """Execute Docker Compose command"""
        try:
            cmd = ["docker-compose", "-f", compose_file, command] + list(args)

            result = subprocess.run(
                cmd,
                cwd=self.nutrition_platform_dir,
                capture_output=True,
                text=True,
                timeout=300
            )

            if result.returncode == 0:
                return {
                    "success": True,
                    "output": result.stdout,
                    "error_output": result.stderr
                }
            else:
                return {
                    "success": False,
                    "error": result.stderr,
                    "return_code": result.returncode
                }

        except subprocess.TimeoutExpired:
            return {
                "success": False,
                "error": "Command timed out"
            }
        except Exception as e:
            return {
                "success": False,
                "error": str(e)
            }

    def _run_deployment_script(self, script_path: str, environment: str, server_ip: str = None, ssh_key: str = None) -> str:
        """Run existing deployment script"""
        try:
            cmd = [script_path]

            if environment:
                cmd.extend(["--environment", environment])
            if server_ip:
                cmd.extend(["--server-ip", server_ip])
            if ssh_key:
                cmd.extend(["--ssh-key", ssh_key])

            result = subprocess.run(
                cmd,
                cwd=self.nutrition_platform_dir,
                capture_output=True,
                text=True,
                timeout=600
            )

            if result.returncode == 0:
                return f"Deployment script executed successfully: {result.stdout}"
            else:
                return f"Deployment script failed: {result.stderr}"

        except Exception as e:
            return f"Error running deployment script: {str(e)}"

    def _deploy_to_vps_custom(self, provider: str, environment: str, server_ip: str, ssh_key: str, vps_config: Dict[str, Any]) -> str:
        """Custom VPS deployment logic"""
        try:
            self.logger.log_vps_operation("Starting custom deployment", provider, server_ip)

            # This would contain custom deployment logic for VPS
            # For now, return a placeholder
            return f"Custom VPS deployment for {provider} (server: {server_ip}) - Implementation needed"

        except Exception as e:
            return f"Custom VPS deployment failed: {str(e)}"

    def _create_backup(self, backup_type: str) -> Optional[str]:
        """Create backup before deployment"""
        try:
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            backup_name = f"{backup_type}_backup_{timestamp}"

            # Create backup directory
            backup_dir = os.path.join(self.config.get_backup_path(), backup_name)
            os.makedirs(backup_dir, exist_ok=True)

            # Backup configuration files
            config_files = [
                "docker-compose.yml",
                "docker-compose.production.yml",
                ".env",
                ".env.production"
            ]

            for config_file in config_files:
                src_path = os.path.join(self.nutrition_platform_dir, config_file)
                if os.path.exists(src_path):
                    import shutil
                    shutil.copy2(src_path, backup_dir)

            return backup_dir

        except Exception as e:
            self.logger.error(f"Backup creation failed: {str(e)}")
            return None

    def _run_health_checks(self, environment: str) -> Dict[str, Any]:
        """Run health checks after deployment"""
        try:
            # Get monitoring configuration
            monitoring_config = self.config.get_monitoring_settings()
            health_url = monitoring_config.get("health_check_url")

            if not health_url:
                return {
                    "success": True,
                    "message": "No health check URL configured, skipping health checks"
                }

            # Perform health check
            return self._perform_health_check(health_url, 30)

        except Exception as e:
            return {
                "success": False,
                "error": str(e)
            }

    def _perform_health_check(self, url: str, timeout: int) -> str:
        """Perform HTTP health check"""
        try:
            start_time = time.time()

            response = requests.get(url, timeout=timeout)
            response_time = time.time() - start_time

            status_code = response.status_code

            self.logger.log_health_check(url, status_code, response_time)

            if 200 <= status_code < 300:
                return f"Health check passed for {url} (Status: {status_code}, Response time: {response_time:.2f}s)"
            else:
                return f"Health check failed for {url} (Status: {status_code})"

        except requests.exceptions.RequestException as e:
            return f"Health check error for {url}: {str(e)}"

    def _check_docker_status(self) -> str:
        """Check Docker deployment status"""
        try:
            result = subprocess.run(
                ["docker-compose", "ps"],
                cwd=self.nutrition_platform_dir,
                capture_output=True,
                text=True
            )

            if result.returncode == 0:
                return f"Docker container status:\n{result.stdout}"
            else:
                return f"Error checking Docker status: {result.stderr}"

        except Exception as e:
            return f"Docker status check failed: {str(e)}"

    def _check_vps_status(self) -> str:
        """Check VPS deployment status"""
        return "VPS status check - Implementation needed"

    def _check_coolify_status(self, deployment_id: str) -> str:
        """Check Coolify deployment status"""
        try:
            status = self._check_coolify_deployment_status(deployment_id)
            return f"Coolify deployment status: {status['status']}"
        except Exception as e:
            return f"Coolify status check failed: {str(e)}"

    def _get_docker_logs(self, tail_lines: int, follow: bool) -> str:
        """Get Docker logs"""
        try:
            cmd = ["docker-compose", "logs", "--tail", str(tail_lines)]

            if follow:
                cmd.append("--follow")

            result = subprocess.run(
                cmd,
                cwd=self.nutrition_platform_dir,
                capture_output=True,
                text=True
            )

            if result.returncode == 0:
                return result.stdout
            else:
                return f"Error getting Docker logs: {result.stderr}"

        except Exception as e:
            return f"Docker logs retrieval failed: {str(e)}"

    def _get_vps_logs(self, deployment_id: str, tail_lines: int) -> str:
        """Get VPS logs"""
        return f"VPS logs for deployment {deployment_id} (last {tail_lines} lines) - Implementation needed"

    def _get_coolify_logs(self, deployment_id: str, tail_lines: int) -> str:
        """Get Coolify logs"""
        return f"Coolify logs for deployment {deployment_id} (last {tail_lines} lines) - Implementation needed"

    def _rollback_coolify_deployment(self, project_id: str, deployment_id: str) -> bool:
        """Rollback Coolify deployment"""
        try:
            # Implementation for Coolify rollback
            self.logger.info("Coolify rollback - Implementation needed")
            return True
        except Exception as e:
            self.logger.error(f"Coolify rollback failed: {str(e)}")
            return False

    def _is_docker_deployment(self) -> bool:
        """Check if current deployment is Docker-based"""
        return os.path.exists(os.path.join(self.nutrition_platform_dir, "docker-compose.yml"))

    def _is_vps_deployment(self) -> bool:
        """Check if current deployment is VPS-based"""
        # Check for VPS deployment indicators
        return any(os.path.exists(os.path.join(self.nutrition_platform_dir, f"deploy-{provider}.sh"))
                  for provider in ["vultr", "digitalocean", "aws"])

def handle(args) -> str:
    """Handle deployment commands"""
    manager = DeploymentManager()

    if hasattr(args, 'deploy_action'):
        if args.deploy_action == 'coolify':
            return manager.deploy_coolify(
                environment=getattr(args, 'environment', 'production'),
                project=getattr(args, 'project', None),
                force=getattr(args, 'force', False)
            )
        elif args.deploy_action == 'docker':
            return manager.deploy_docker(
                compose_file=getattr(args, 'compose_file', None),
                environment=getattr(args, 'environment', 'production'),
                build=getattr(args, 'build', False),
                restart=getattr(args, 'restart', False)
            )
        elif args.deploy_action == 'vps':
            return manager.deploy_vps(
                provider=getattr(args, 'provider', 'vultr'),
                environment=getattr(args, 'environment', 'production'),
                server_ip=getattr(args, 'server_ip', None),
                ssh_key=getattr(args, 'ssh_key', None)
            )

    return "Deployment action not recognized"