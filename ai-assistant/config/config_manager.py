#!/usr/bin/env python3
"""
Configuration Manager for Nutrition Platform AI Assistant
Handles deployment settings and environment configurations
"""

import json
import os
import sys
from typing import Dict, Any, Optional

class ConfigManager:
    """Configuration manager for deployment settings"""

    def __init__(self, config_file: str = None):
        """Initialize configuration manager"""
        # Get the ai-assistant directory path
        ai_assistant_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
        nutrition_platform_dir = os.path.dirname(ai_assistant_dir)

        # Default configuration file path
        if config_file is None:
            config_file = os.path.join(ai_assistant_dir, 'config', 'deployment_config.json')

        self.config_file = config_file
        self.nutrition_platform_dir = nutrition_platform_dir
        self.ai_assistant_dir = ai_assistant_dir

        # Default configuration
        self.default_config = {
            # Deployment environments
            "environments": {
                "development": {
                    "docker_compose_file": "docker-compose.local.yml",
                    "build_args": {},
                    "environment_variables": {}
                },
                "staging": {
                    "docker_compose_file": "docker-compose.vps.yml",
                    "build_args": {},
                    "environment_variables": {}
                },
                "production": {
                    "docker_compose_file": "docker-compose.production.yml",
                    "build_args": {},
                    "environment_variables": {}
                }
            },

            # VPS providers configuration
            "vps_providers": {
                "vultr": {
                    "api_endpoint": "https://api.vultr.com/v2",
                    "default_region": "ewr",
                    "instance_type": "vc2-1c-1gb",
                    "os_id": 387,  # Ubuntu 20.04
                    "ssh_key_name": "nutrition-platform-key"
                },
                "digitalocean": {
                    "api_endpoint": "https://api.digitalocean.com/v2",
                    "default_region": "nyc1",
                    "instance_size": "s-1vcpu-1gb",
                    "image": "ubuntu-20-04-x64"
                }
            },

            # Coolify configuration
            "coolify": {
                "base_url": "https://your-coolify-instance.com",
                "api_token": "",
                "project_name": "nutrition-platform",
                "default_environment": "production"
            },

            # Deployment settings
            "deployment": {
                "max_retries": 3,
                "retry_delay": 30,
                "timeout": 300,
                "health_check_interval": 10,
                "backup_before_deploy": True,
                "rollback_on_failure": True
            },

            # Monitoring settings
            "monitoring": {
                "health_check_url": "/api/health",
                "log_retention_days": 30,
                "alert_webhook": "",
                "metrics_enabled": True
            },

            # Database settings
            "database": {
                "backup_path": "./backups",
                "retention_days": 7,
                "auto_backup": True,
                "backup_schedule": "0 2 * * *"  # Daily at 2 AM
            },

            # SSL/TLS settings
            "ssl": {
                "auto_renew": True,
                "provider": "letsencrypt",
                "email": "admin@your-domain.com",
                "staging": False
            },

            # Notification settings
            "notifications": {
                "deployment_success": True,
                "deployment_failure": True,
                "email_recipients": ["admin@your-domain.com"],
                "slack_webhook": ""
            }
        }

        self.config = self.load_config()

    def load_config(self) -> Dict[str, Any]:
        """Load configuration from file"""
        try:
            if os.path.exists(self.config_file):
                with open(self.config_file, 'r') as f:
                    file_config = json.load(f)
                    # Merge with default config
                    merged_config = self._deep_merge(self.default_config.copy(), file_config)
                    return merged_config
            else:
                # Create default config file
                self.save_config()
                return self.default_config.copy()
        except Exception as e:
            print(f"Warning: Could not load config file: {e}")
            return self.default_config.copy()

    def save_config(self) -> bool:
        """Save current configuration to file"""
        try:
            # Create config directory if it doesn't exist
            os.makedirs(os.path.dirname(self.config_file), exist_ok=True)

            with open(self.config_file, 'w') as f:
                json.dump(self.config, f, indent=2)
            return True
        except Exception as e:
            print(f"Error saving config: {e}")
            return False

    def set_config(self, key: str, value: Any) -> str:
        """Set a configuration value"""
        try:
            keys = key.split('.')
            current = self.config

            # Navigate to the nested location
            for k in keys[:-1]:
                if k not in current:
                    current[k] = {}
                current = current[k]

            # Set the value
            current[keys[-1]] = value
            self.save_config()
            return f"Configuration {key} set to {value}"
        except Exception as e:
            return f"Error setting configuration {key}: {str(e)}"

    def get_config(self, key: str, default: Any = None) -> Any:
        """Get a configuration value"""
        try:
            keys = key.split('.')
            current = self.config

            # Navigate to the nested location
            for k in keys:
                if isinstance(current, dict) and k in current:
                    current = current[k]
                else:
                    return default

            return current
        except Exception:
            return default

    def get_environment_config(self, environment: str) -> Dict[str, Any]:
        """Get configuration for specific environment"""
        return self.get_config(f"environments.{environment}", {})

    def get_vps_config(self, provider: str) -> Dict[str, Any]:
        """Get VPS provider configuration"""
        return self.get_config(f"vps_providers.{provider}", {})

    def get_deployment_settings(self) -> Dict[str, Any]:
        """Get deployment settings"""
        return self.get_config("deployment", {})

    def get_monitoring_settings(self) -> Dict[str, Any]:
        """Get monitoring settings"""
        return self.get_config("monitoring", {})

    def update_environment_config(self, environment: str, **kwargs) -> str:
        """Update environment-specific configuration"""
        try:
            if "environments" not in self.config:
                self.config["environments"] = {}

            if environment not in self.config["environments"]:
                self.config["environments"][environment] = {}

            self.config["environments"][environment].update(kwargs)
            self.save_config()
            return f"Updated {environment} environment configuration"
        except Exception as e:
            return f"Error updating environment config: {str(e)}"

    def add_notification_webhook(self, webhook_type: str, webhook_url: str) -> str:
        """Add notification webhook"""
        try:
            if "notifications" not in self.config:
                self.config["notifications"] = {}

            self.config["notifications"][webhook_type] = webhook_url
            self.save_config()
            return f"Added {webhook_type} webhook"
        except Exception as e:
            return f"Error adding webhook: {str(e)}"

    def get_docker_compose_file(self, environment: str) -> str:
        """Get Docker Compose file for environment"""
        env_config = self.get_environment_config(environment)
        compose_file = env_config.get("docker_compose_file", "docker-compose.yml")

        # Construct full path
        compose_path = os.path.join(self.nutrition_platform_dir, compose_file)

        # Check if file exists, fallback to default
        if not os.path.exists(compose_path):
            default_path = os.path.join(self.nutrition_platform_dir, "docker-compose.yml")
            if os.path.exists(default_path):
                return default_path
            return compose_path

        return compose_path

    def get_backup_path(self) -> str:
        """Get backup directory path"""
        backup_path = self.get_config("database.backup_path", "./backups")
        full_path = os.path.join(self.nutrition_platform_dir, backup_path)
        return full_path

    def _deep_merge(self, base: Dict[str, Any], update: Dict[str, Any]) -> Dict[str, Any]:
        """Deep merge two dictionaries"""
        result = base.copy()

        for key, value in update.items():
            if key in result and isinstance(result[key], dict) and isinstance(value, dict):
                result[key] = self._deep_merge(result[key], value)
            else:
                result[key] = value

        return result

    def validate_config(self) -> Dict[str, Any]:
        """Validate current configuration"""
        issues = []

        # Check required fields
        if not self.get_config("coolify.base_url"):
            issues.append("Coolify base URL not configured")

        if not self.get_config("coolify.api_token"):
            issues.append("Coolify API token not configured")

        # Check file paths
        for env in ["development", "staging", "production"]:
            compose_file = self.get_docker_compose_file(env)
            if not os.path.exists(compose_file):
                issues.append(f"Docker Compose file for {env} not found: {compose_file}")

        return {
            "valid": len(issues) == 0,
            "issues": issues
        }

    def export_config(self, export_path: str) -> str:
        """Export configuration to file"""
        try:
            with open(export_path, 'w') as f:
                json.dump(self.config, f, indent=2)
            return f"Configuration exported to {export_path}"
        except Exception as e:
            return f"Error exporting configuration: {str(e)}"

    def import_config(self, import_path: str) -> str:
        """Import configuration from file"""
        try:
            with open(import_path, 'r') as f:
                imported_config = json.load(f)

            # Merge with current config
            self.config = self._deep_merge(self.config, imported_config)
            self.save_config()
            return f"Configuration imported from {import_path}"
        except Exception as e:
            return f"Error importing configuration: {str(e)}"

    def reset_to_defaults(self) -> str:
        """Reset configuration to defaults"""
        try:
            self.config = self.default_config.copy()
            self.save_config()
            return "Configuration reset to defaults"
        except Exception as e:
            return f"Error resetting configuration: {str(e)}"