#!/usr/bin/env python3
"""
Logger utility for Nutrition Platform AI Assistant
Provides logging functionality for deployment operations
"""

import logging
import sys
import os
from datetime import datetime
from typing import Optional

class Logger:
    """Custom logger for deployment operations"""

    def __init__(self, log_file: str = None, level: str = 'INFO'):
        """Initialize logger"""
        self.logger = logging.getLogger('nutrition-platform-ai')
        self.logger.setLevel(getattr(logging, level.upper(), logging.INFO))

        # Remove existing handlers to avoid duplicates
        self.logger.handlers.clear()

        # Create formatters
        file_formatter = logging.Formatter(
            '%(asctime)s - %(name)s - %(levelname)s - %(funcName)s:%(lineno)d - %(message)s'
        )
        console_formatter = logging.Formatter(
            '%(levelname)s: %(message)s'
        )

        # Console handler
        console_handler = logging.StreamHandler(sys.stdout)
        console_handler.setFormatter(console_formatter)
        console_handler.setLevel(logging.INFO)
        self.logger.addHandler(console_handler)

        # File handler (if log file specified)
        if log_file:
            # Get the ai-assistant directory path
            ai_assistant_dir = os.path.dirname(os.path.dirname(os.path.abspath(__file__)))
            log_dir = os.path.join(ai_assistant_dir, 'logs')

            # Create logs directory if it doesn't exist
            os.makedirs(log_dir, exist_ok=True)

            # Use timestamped log file if not specified
            if log_file == 'auto':
                timestamp = datetime.now().strftime('%Y%m%d_%H%M%S')
                log_file = os.path.join(log_dir, f'deployment_{timestamp}.log')
            else:
                log_file = os.path.join(log_dir, log_file)

            file_handler = logging.FileHandler(log_file)
            file_handler.setFormatter(file_formatter)
            file_handler.setLevel(logging.DEBUG)
            self.logger.addHandler(file_handler)

        # Prevent duplicate messages
        self.logger.propagate = False

    def set_level(self, level: str):
        """Set logging level"""
        self.logger.setLevel(getattr(logging, level.upper(), logging.INFO))

        # Update console handler level
        for handler in self.logger.handlers:
            if isinstance(handler, logging.StreamHandler):
                handler.setLevel(getattr(logging, level.upper(), logging.INFO))

    def debug(self, message: str):
        """Log debug message"""
        self.logger.debug(message)

    def info(self, message: str):
        """Log info message"""
        self.logger.info(message)

    def warning(self, message: str):
        """Log warning message"""
        self.logger.warning(message)

    def error(self, message: str):
        """Log error message"""
        self.logger.error(message)

    def critical(self, message: str):
        """Log critical message"""
        self.logger.critical(message)

    def log_deployment_start(self, deployment_type: str, environment: str, **kwargs):
        """Log deployment start"""
        self.info(f"=== Starting {deployment_type} deployment to {environment} ===")
        for key, value in kwargs.items():
            self.info(f"{key}: {value}")

    def log_deployment_end(self, success: bool, duration: Optional[float] = None):
        """Log deployment end"""
        if success:
            self.info("=== Deployment completed successfully ===")
        else:
            self.error("=== Deployment failed ===")

        if duration:
            self.info(f"Duration: {duration:.2f} seconds")

    def log_step(self, step: str, status: str = "START"):
        """Log deployment step"""
        self.info(f"[{status}] {step}")

    def log_error_details(self, error: Exception, context: str = ""):
        """Log detailed error information"""
        self.error(f"Error in {context}: {str(error)}")
        self.debug(f"Error type: {type(error).__name__}")
        import traceback
        self.debug(f"Traceback: {traceback.format_exc()}")

    def log_health_check(self, url: str, status_code: int, response_time: float):
        """Log health check results"""
        if 200 <= status_code < 300:
            self.info(f"Health check PASSED for {url} - Status: {status_code} - Response time: {response_time:.2f}s")
        else:
            self.error(f"Health check FAILED for {url} - Status: {status_code} - Response time: {response_time:.2f}s")

    def log_backup_operation(self, operation: str, path: str, success: bool):
        """Log backup operation"""
        if success:
            self.info(f"Backup {operation} successful: {path}")
        else:
            self.error(f"Backup {operation} failed: {path}")

    def log_file_operation(self, operation: str, file_path: str, success: bool):
        """Log file operation"""
        if success:
            self.info(f"File {operation} successful: {file_path}")
        else:
            self.error(f"File {operation} failed: {file_path}")

    def log_vps_operation(self, operation: str, provider: str, server_ip: str = None):
        """Log VPS operation"""
        message = f"VPS {operation} on {provider}"
        if server_ip:
            message += f" (Server: {server_ip})"
        self.info(message)

    def log_docker_operation(self, operation: str, service: str = None, compose_file: str = None):
        """Log Docker operation"""
        message = f"Docker {operation}"
        if service:
            message += f" - Service: {service}"
        if compose_file:
            message += f" - Compose file: {compose_file}"
        self.info(message)

    def log_config_change(self, key: str, old_value: str, new_value: str):
        """Log configuration change"""
        self.info(f"Configuration changed - {key}: '{old_value}' -> '{new_value}'")

    def log_web_automation(self, action: str, url: str, success: bool, details: str = None):
        """Log web automation operation"""
        if success:
            message = f"Web automation {action} successful for {url}"
        else:
            message = f"Web automation {action} failed for {url}"

        if details:
            message += f" - {details}"

        if success:
            self.info(message)
        else:
            self.error(message)

    def log_desktop_automation(self, action: str, target: str, success: bool):
        """Log desktop automation operation"""
        if success:
            self.info(f"Desktop automation {action} successful: {target}")
        else:
            self.error(f"Desktop automation {action} failed: {target}")

    def get_recent_logs(self, lines: int = 50) -> str:
        """Get recent log entries"""
        # This is a simplified implementation
        # In a real scenario, you might want to read from the log file
        return f"Recent {lines} log entries (simplified view)"

    def cleanup_old_logs(self, retention_days: int = 30):
        """Clean up old log files"""
        try:
            log_dir = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'logs')
            if not os.path.exists(log_dir):
                return

            cutoff_date = datetime.now().timestamp() - (retention_days * 24 * 3600)

            for filename in os.listdir(log_dir):
                if filename.endswith('.log'):
                    filepath = os.path.join(log_dir, filename)
                    if os.path.getmtime(filepath) < cutoff_date:
                        os.remove(filepath)
                        self.info(f"Removed old log file: {filename}")

        except Exception as e:
            self.error(f"Error cleaning up old logs: {str(e)}")