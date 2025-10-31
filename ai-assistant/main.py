#!/usr/bin/env python3
"""
AI Assistant for Nutrition Platform - Deployment Automation
Main CLI framework for automated deployment management
"""

import argparse
import sys
import os
from datetime import datetime

# Add the ai-assistant directory to Python path
sys.path.append(os.path.dirname(os.path.abspath(__file__)))

from commands import deployment_commands, web_commands, desktop_commands, file_commands
from config.config_manager import ConfigManager
from utils.logger import Logger

class NutritionPlatformAssistant:
    def __init__(self):
        self.config = ConfigManager()
        self.logger = Logger()
        self.version = "1.0.0"

    def setup_cli(self):
        """Set up the command line interface"""
        parser = argparse.ArgumentParser(
            description='AI Assistant for Nutrition Platform - Deployment Automation',
            formatter_class=argparse.RawDescriptionHelpFormatter,
            epilog="""
Examples:
  # Deploy to Coolify
  python main.py deploy coolify --environment production

  # Deploy to VPS
  python main.py deploy vps --provider vultr --environment production

  # Deploy Docker containers
  python main.py deploy docker --compose-file docker-compose.production.yml

  # Monitor deployment status
  python main.py monitor status --deployment-id latest

  # Web automation for testing
  python main.py web test --url https://your-domain.com --health-check

  # File operations for deployment
  python main.py file backup --config --logs
            """
        )

        subparsers = parser.add_subparsers(dest='command', help='Available commands')

        # Global options
        parser.add_argument('--version', action='version', version=f'%(prog)s {self.version}')
        parser.add_argument('--verbose', '-v', action='store_true', help='Verbose output')
        parser.add_argument('--config', help='Configuration file path')

        # Deployment command group
        deploy_parser = subparsers.add_parser('deploy', help='Deployment operations')
        deploy_subparsers = deploy_parser.add_subparsers(dest='deploy_action', help='Deployment actions')

        # Coolify deployment
        coolify_parser = deploy_subparsers.add_parser('coolify', help='Deploy to Coolify')
        coolify_parser.add_argument('--environment', default='production',
                                  choices=['development', 'staging', 'production'],
                                  help='Deployment environment')
        coolify_parser.add_argument('--project', help='Coolify project name')
        coolify_parser.add_argument('--force', action='store_true', help='Force deployment')

        # VPS deployment
        vps_parser = deploy_subparsers.add_parser('vps', help='Deploy to VPS')
        vps_parser.add_argument('--provider', choices=['vultr', 'digitalocean', 'aws', 'custom'],
                              default='vultr', help='VPS provider')
        vps_parser.add_argument('--environment', default='production',
                              choices=['development', 'staging', 'production'])
        vps_parser.add_argument('--server-ip', help='Target server IP')
        vps_parser.add_argument('--ssh-key', help='SSH key path')

        # Docker deployment
        docker_parser = deploy_subparsers.add_parser('docker', help='Deploy Docker containers')
        docker_parser.add_argument('--compose-file', default='docker-compose.production.yml',
                                 help='Docker compose file')
        docker_parser.add_argument('--build', action='store_true', help='Build images before deployment')
        docker_parser.add_argument('--restart', action='store_true', help='Restart containers')

        # Monitor command group
        monitor_parser = subparsers.add_parser('monitor', help='Monitor deployments')
        monitor_subparsers = monitor_parser.add_subparsers(dest='monitor_action', help='Monitor actions')

        # Status monitoring
        status_parser = monitor_subparsers.add_parser('status', help='Check deployment status')
        status_parser.add_argument('--deployment-id', help='Specific deployment ID')
        status_parser.add_argument('--watch', action='store_true', help='Watch status in real-time')

        # Logs monitoring
        logs_parser = monitor_subparsers.add_parser('logs', help='View deployment logs')
        logs_parser.add_argument('--deployment-id', help='Specific deployment ID')
        logs_parser.add_argument('--tail', type=int, default=50, help='Number of lines to show')
        logs_parser.add_argument('--follow', '-f', action='store_true', help='Follow log output')

        # Health checks
        health_parser = monitor_subparsers.add_parser('health', help='Health check deployments')
        health_parser.add_argument('--url', help='Application URL for health check')
        health_parser.add_argument('--timeout', type=int, default=30, help='Timeout in seconds')

        # Web automation command group
        web_parser = subparsers.add_parser('web', help='Web automation for deployment verification')
        web_subparsers = web_parser.add_subparsers(dest='web_action', help='Web actions')

        # Web testing
        test_parser = web_subparsers.add_parser('test', help='Test deployed application')
        test_parser.add_argument('--url', required=True, help='Application URL')
        test_parser.add_argument('--health-check', action='store_true', help='Perform health check')
        test_parser.add_argument('--load-test', action='store_true', help='Perform load testing')

        # Desktop automation command group
        desktop_parser = subparsers.add_parser('desktop', help='Desktop automation for local deployment')
        desktop_subparsers = desktop_parser.add_subparsers(dest='desktop_action', help='Desktop actions')

        # Application management
        app_parser = desktop_subparsers.add_parser('app', help='Manage local applications')
        app_parser.add_argument('--action', choices=['start', 'stop', 'restart', 'status'],
                              help='Action to perform')
        app_parser.add_argument('--app-name', help='Application name')

        # File operations command group
        file_parser = subparsers.add_parser('file', help='File operations for deployment')
        file_subparsers = file_parser.add_subparsers(dest='file_action', help='File actions')

        # Backup operations
        backup_parser = file_subparsers.add_parser('backup', help='Backup deployment files')
        backup_parser.add_argument('--config', action='store_true', help='Backup configuration files')
        backup_parser.add_argument('--logs', action='store_true', help='Backup log files')
        backup_parser.add_argument('--database', action='store_true', help='Backup database')
        backup_parser.add_argument('--destination', help='Backup destination path')

        # Restore operations
        restore_parser = file_subparsers.add_parser('restore', help='Restore from backup')
        restore_parser.add_argument('--backup-path', required=True, help='Path to backup files')
        restore_parser.add_argument('--config', action='store_true', help='Restore configuration')
        restore_parser.add_argument('--logs', action='store_true', help='Restore logs')
        restore_parser.add_argument('--database', action='store_true', help='Restore database')

        # Configuration command group
        config_parser = subparsers.add_parser('config', help='Configuration management')
        config_subparsers = config_parser.add_subparsers(dest='config_action', help='Config actions')

        # Config management
        config_set_parser = config_subparsers.add_parser('set', help='Set configuration value')
        config_set_parser.add_argument('--key', required=True, help='Configuration key')
        config_set_parser.add_argument('--value', required=True, help='Configuration value')

        config_get_parser = config_subparsers.add_parser('get', help='Get configuration value')
        config_get_parser.add_argument('--key', help='Configuration key')

        return parser

    def main(self):
        """Main entry point"""
        parser = self.setup_cli()

        # Show help if no arguments provided
        if len(sys.argv) == 1:
            parser.print_help()
            return

        args = parser.parse_args()

        # Set up logging level based on verbose flag
        if args.verbose:
            self.logger.set_level('DEBUG')

        self.logger.info(f"Nutrition Platform AI Assistant v{self.version}")
        self.logger.info(f"Command: {' '.join(sys.argv[1:])}")

        try:
            # Route to appropriate command handler
            if args.command == 'deploy':
                result = self.handle_deployment(args)
            elif args.command == 'monitor':
                result = self.handle_monitoring(args)
            elif args.command == 'web':
                result = self.handle_web_automation(args)
            elif args.command == 'desktop':
                result = self.handle_desktop_automation(args)
            elif args.command == 'file':
                result = self.handle_file_operations(args)
            elif args.command == 'config':
                result = self.handle_configuration(args)
            else:
                parser.print_help()
                return

            # Log the result
            if result:
                self.logger.info(f"Command completed successfully: {result}")
            else:
                self.logger.warning("Command completed with no result")

        except KeyboardInterrupt:
            self.logger.info("Operation cancelled by user")
        except Exception as e:
            self.logger.error(f"Command failed: {str(e)}")
            if args.verbose:
                import traceback
                self.logger.error(traceback.format_exc())
            sys.exit(1)

    def handle_deployment(self, args):
        """Handle deployment commands"""
        return deployment_commands.handle(args)

    def handle_monitoring(self, args):
        """Handle monitoring commands"""
        if args.monitor_action == 'status':
            return deployment_commands.check_status(args)
        elif args.monitor_action == 'logs':
            return deployment_commands.view_logs(args)
        elif args.monitor_action == 'health':
            return deployment_commands.health_check(args)
        return "Monitoring action not implemented"

    def handle_web_automation(self, args):
        """Handle web automation commands"""
        return web_commands.handle(args)

    def handle_desktop_automation(self, args):
        """Handle desktop automation commands"""
        return desktop_commands.handle(args)

    def handle_file_operations(self, args):
        """Handle file operation commands"""
        return file_commands.handle(args)

    def handle_configuration(self, args):
        """Handle configuration commands"""
        if args.config_action == 'set':
            return self.config.set_config(args.key, args.value)
        elif args.config_action == 'get':
            value = self.config.get_config(args.key)
            if value:
                print(f"{args.key}: {value}")
                return f"Retrieved {args.key}"
            else:
                return f"Configuration key '{args.key}' not found"
        return "Configuration action not implemented"

def main():
    """Entry point function"""
    assistant = NutritionPlatformAssistant()
    assistant.main()

if __name__ == "__main__":
    main()