#!/usr/bin/env python3
"""
File Operations Commands for Nutrition Platform AI Assistant
Handles file operations for deployment configuration management
"""

import os
import sys
import json
import shutil
import hashlib
import zipfile
import tarfile
from datetime import datetime
from pathlib import Path
from typing import Dict, Any, List, Optional

# Add parent directory to path for imports
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from config.config_manager import ConfigManager
from utils.logger import Logger

class FileOperations:
    """File operations for deployment management"""

    def __init__(self):
        """Initialize file operations"""
        self.config = ConfigManager()
        self.logger = Logger()
        self.nutrition_platform_dir = self.config.nutrition_platform_dir

    def create_file(self, file_path: str, content: str = "", overwrite: bool = False) -> str:
        """Create file with content"""
        try:
            full_path = os.path.join(self.nutrition_platform_dir, file_path)

            if os.path.exists(full_path) and not overwrite:
                return f"File already exists: {file_path}. Use --overwrite to replace."

            # Create directory if it doesn't exist
            os.makedirs(os.path.dirname(full_path), exist_ok=True)

            with open(full_path, 'w', encoding='utf-8') as f:
                f.write(content)

            self.logger.log_file_operation("create", file_path, True)
            return f"File created: {file_path}"

        except Exception as e:
            self.logger.log_file_operation("create", file_path, False, str(e))
            return f"Error creating file: {str(e)}"

    def edit_file(self, file_path: str, content: str, mode: str = 'replace') -> str:
        """Edit file content"""
        try:
            full_path = os.path.join(self.nutrition_platform_dir, file_path)

            if not os.path.exists(full_path):
                return f"File not found: {file_path}"

            if mode == 'append':
                with open(full_path, 'a', encoding='utf-8') as f:
                    f.write(content)
            elif mode == 'prepend':
                with open(full_path, 'r+', encoding='utf-8') as f:
                    original_content = f.read()
                    f.seek(0)
                    f.write(content + original_content)
            else:  # replace
                with open(full_path, 'w', encoding='utf-8') as f:
                    f.write(content)

            self.logger.log_file_operation("edit", file_path, True)
            return f"File edited: {file_path}"

        except Exception as e:
            self.logger.log_file_operation("edit", file_path, False, str(e))
            return f"Error editing file: {str(e)}"

    def delete_file(self, file_path: str) -> str:
        """Delete file"""
        try:
            full_path = os.path.join(self.nutrition_platform_dir, file_path)

            if not os.path.exists(full_path):
                return f"File not found: {file_path}"

            os.remove(full_path)
            self.logger.log_file_operation("delete", file_path, True)
            return f"File deleted: {file_path}"

        except Exception as e:
            self.logger.log_file_operation("delete", file_path, False, str(e))
            return f"Error deleting file: {str(e)}"

    def move_file(self, source_path: str, destination_path: str) -> str:
        """Move file"""
        try:
            src_full = os.path.join(self.nutrition_platform_dir, source_path)
            dst_full = os.path.join(self.nutrition_platform_dir, destination_path)

            if not os.path.exists(src_full):
                return f"Source file not found: {source_path}"

            # Create destination directory if it doesn't exist
            os.makedirs(os.path.dirname(dst_full), exist_ok=True)

            shutil.move(src_full, dst_full)
            self.logger.log_file_operation("move", f"{source_path} -> {destination_path}", True)
            return f"File moved: {source_path} -> {destination_path}"

        except Exception as e:
            self.logger.log_file_operation("move", f"{source_path} -> {destination_path}", False, str(e))
            return f"Error moving file: {str(e)}"

    def copy_file(self, source_path: str, destination_path: str) -> str:
        """Copy file"""
        try:
            src_full = os.path.join(self.nutrition_platform_dir, source_path)
            dst_full = os.path.join(self.nutrition_platform_dir, destination_path)

            if not os.path.exists(src_full):
                return f"Source file not found: {source_path}"

            # Create destination directory if it doesn't exist
            os.makedirs(os.path.dirname(dst_full), exist_ok=True)

            shutil.copy2(src_full, dst_full)
            self.logger.log_file_operation("copy", f"{source_path} -> {destination_path}", True)
            return f"File copied: {source_path} -> {destination_path}"

        except Exception as e:
            self.logger.log_file_operation("copy", f"{source_path} -> {destination_path}", False, str(e))
            return f"Error copying file: {str(e)}"

    def backup_files(self, config: bool = False, logs: bool = False, database: bool = False, destination: str = None) -> str:
        """Backup deployment files"""
        try:
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            backup_name = f"deployment_backup_{timestamp}"

            if not destination:
                backup_dir = os.path.join(self.config.get_backup_path(), backup_name)
            else:
                backup_dir = os.path.join(self.nutrition_platform_dir, destination, backup_name)

            os.makedirs(backup_dir, exist_ok=True)

            backed_up_files = []

            # Backup configuration files
            if config:
                config_files = [
                    "docker-compose.yml",
                    "docker-compose.production.yml",
                    "docker-compose.vps.yml",
                    ".env",
                    ".env.production",
                    ".env.vps",
                    "nginx.conf",
                    "package.json",
                    "package-lock.json"
                ]

                config_backup_dir = os.path.join(backup_dir, "config")
                os.makedirs(config_backup_dir, exist_ok=True)

                for config_file in config_files:
                    src_path = os.path.join(self.nutrition_platform_dir, config_file)
                    if os.path.exists(src_path):
                        shutil.copy2(src_path, config_backup_dir)
                        backed_up_files.append(config_file)

            # Backup log files
            if logs:
                logs_src_dir = os.path.join(self.nutrition_platform_dir, "logs")
                if os.path.exists(logs_src_dir):
                    logs_backup_dir = os.path.join(backup_dir, "logs")
                    shutil.copytree(logs_src_dir, logs_backup_dir)
                    backed_up_files.append("logs/")

            # Backup database (if applicable)
            if database:
                # This would depend on your database setup
                # For now, just note that database backup is requested
                db_backup_note = os.path.join(backup_dir, "database_backup_note.txt")
                with open(db_backup_note, 'w') as f:
                    f.write("Database backup requested but not implemented.\n")
                    f.write("Please implement database-specific backup logic.\n")
                backed_up_files.append("database_note")

            # Create backup manifest
            manifest = {
                "backup_name": backup_name,
                "timestamp": timestamp,
                "backed_up_files": backed_up_files,
                "backup_type": "deployment",
                "total_files": len(backed_up_files)
            }

            manifest_path = os.path.join(backup_dir, "backup_manifest.json")
            with open(manifest_path, 'w') as f:
                json.dump(manifest, f, indent=2)

            self.logger.log_backup_operation("create", backup_dir, True)
            return f"Backup created: {backup_dir} ({len(backed_up_files)} items)"

        except Exception as e:
            self.logger.log_backup_operation("create", destination or "auto", False, str(e))
            return f"Error creating backup: {str(e)}"

    def restore_files(self, backup_path: str, config: bool = False, logs: bool = False, database: bool = False) -> str:
        """Restore from backup"""
        try:
            full_backup_path = os.path.join(self.nutrition_platform_dir, backup_path)

            if not os.path.exists(full_backup_path):
                return f"Backup path not found: {backup_path}"

            # Read backup manifest
            manifest_path = os.path.join(full_backup_path, "backup_manifest.json")
            if not os.path.exists(manifest_path):
                return f"Backup manifest not found: {manifest_path}"

            with open(manifest_path, 'r') as f:
                manifest = json.load(f)

            restored_files = []

            # Restore configuration files
            if config:
                config_backup_dir = os.path.join(full_backup_path, "config")
                if os.path.exists(config_backup_dir):
                    for filename in os.listdir(config_backup_dir):
                        src_path = os.path.join(config_backup_dir, filename)
                        dst_path = os.path.join(self.nutrition_platform_dir, filename)

                        # Create backup of current file
                        if os.path.exists(dst_path):
                            backup_name = f"{filename}.backup.{manifest['timestamp']}"
                            backup_path_current = os.path.join(self.nutrition_platform_dir, backup_name)
                            shutil.copy2(dst_path, backup_path_current)

                        shutil.copy2(src_path, dst_path)
                        restored_files.append(filename)

            # Restore logs
            if logs:
                logs_backup_dir = os.path.join(full_backup_path, "logs")
                if os.path.exists(logs_backup_dir):
                    logs_dst_dir = os.path.join(self.nutrition_platform_dir, "logs")
                    if os.path.exists(logs_dst_dir):
                        shutil.rmtree(logs_dst_dir)
                    shutil.copytree(logs_backup_dir, logs_dst_dir)
                    restored_files.append("logs/")

            # Note: Database restore would need specific implementation

            self.logger.log_backup_operation("restore", backup_path, True)
            return f"Restored {len(restored_files)} items from backup: {backup_path}"

        except Exception as e:
            self.logger.log_backup_operation("restore", backup_path, False, str(e))
            return f"Error restoring backup: {str(e)}"

    def list_files(self, directory: str = ".", recursive: bool = False) -> str:
        """List files in directory"""
        try:
            full_path = os.path.join(self.nutrition_platform_dir, directory)

            if not os.path.exists(full_path):
                return f"Directory not found: {directory}"

            if recursive:
                result = []
                for root, dirs, files in os.walk(full_path):
                    level = root.replace(full_path, '').count(os.sep)
                    indent = ' ' * 2 * level
                    result.append(f"{indent}{os.path.basename(root)}/")
                    subindent = ' ' * 2 * (level + 1)
                    for file in files:
                        result.append(f"{subindent}{file}")
                return f"Directory listing for {directory} (recursive):\n" + "\n".join(result)
            else:
                files = os.listdir(full_path)
                return f"Directory listing for {directory}:\n" + "\n".join(f"  {f}" for f in files)

        except Exception as e:
            return f"Error listing files: {str(e)}"

    def find_files(self, pattern: str, directory: str = ".") -> str:
        """Find files matching pattern"""
        try:
            full_path = os.path.join(self.nutrition_platform_dir, directory)

            if not os.path.exists(full_path):
                return f"Directory not found: {directory}"

            matches = []
            for root, dirs, files in os.walk(full_path):
                for file in files:
                    if pattern.lower() in file.lower():
                        rel_path = os.path.relpath(os.path.join(root, file), self.nutrition_platform_dir)
                        matches.append(rel_path)

            if matches:
                return f"Found {len(matches)} files matching '{pattern}':\n" + "\n".join(f"  {m}" for m in matches)
            else:
                return f"No files found matching pattern: {pattern}"

        except Exception as e:
            return f"Error finding files: {str(e)}"

    def get_file_info(self, file_path: str) -> str:
        """Get file information"""
        try:
            full_path = os.path.join(self.nutrition_platform_dir, file_path)

            if not os.path.exists(full_path):
                return f"File not found: {file_path}"

            stat = os.stat(full_path)

            info = []
            info.append(f"Path: {file_path}")
            info.append(f"Full path: {full_path}")
            info.append(f"Size: {stat.st_size} bytes")
            info.append(f"Modified: {datetime.fromtimestamp(stat.st_mtime).strftime('%Y-%m-%d %H:%M:%S')}")
            info.append(f"Created: {datetime.fromtimestamp(stat.st_ctime).strftime('%Y-%m-%d %H:%M:%S')}")
            info.append(f"Type: {'Directory' if os.path.isdir(full_path) else 'File'}")

            # File hash
            if os.path.isfile(full_path):
                hash_md5 = hashlib.md5()
                with open(full_path, 'rb') as f:
                    for chunk in iter(lambda: f.read(4096), b""):
                        hash_md5.update(chunk)
                info.append(f"MD5 hash: {hash_md5.hexdigest()}")

            return "\n".join(info)

        except Exception as e:
            return f"Error getting file info: {str(e)}"

    def compress_files(self, files: List[str], archive_name: str, format: str = 'zip') -> str:
        """Compress files into archive"""
        try:
            archive_path = os.path.join(self.nutrition_platform_dir, archive_name)

            if format.lower() == 'zip':
                with zipfile.ZipFile(archive_path, 'w', zipfile.ZIP_DEFLATED) as archive:
                    for file_path in files:
                        full_path = os.path.join(self.nutrition_platform_dir, file_path)
                        if os.path.exists(full_path):
                            archive.write(full_path, file_path)

            elif format.lower() == 'tar':
                with tarfile.open(archive_path, 'w') as archive:
                    for file_path in files:
                        full_path = os.path.join(self.nutrition_platform_dir, file_path)
                        if os.path.exists(full_path):
                            archive.add(full_path, file_path)

            else:
                return f"Unsupported archive format: {format}"

            self.logger.log_file_operation("compress", f"{archive_name} ({format})", True)
            return f"Created archive: {archive_name} ({len(files)} files)"

        except Exception as e:
            self.logger.log_file_operation("compress", archive_name, False, str(e))
            return f"Error creating archive: {str(e)}"

    def extract_archive(self, archive_path: str, destination: str = ".") -> str:
        """Extract archive"""
        try:
            full_archive_path = os.path.join(self.nutrition_platform_dir, archive_path)
            full_destination = os.path.join(self.nutrition_platform_dir, destination)

            if not os.path.exists(full_archive_path):
                return f"Archive not found: {archive_path}"

            # Create destination directory
            os.makedirs(full_destination, exist_ok=True)

            if archive_path.endswith('.zip'):
                with zipfile.ZipFile(full_archive_path, 'r') as archive:
                    archive.extractall(full_destination)

            elif archive_path.endswith('.tar.gz') or archive_path.endswith('.tgz'):
                with tarfile.open(full_archive_path, 'r:gz') as archive:
                    archive.extractall(full_destination)

            elif archive_path.endswith('.tar'):
                with tarfile.open(full_archive_path, 'r') as archive:
                    archive.extractall(full_destination)

            else:
                return f"Unsupported archive format: {archive_path}"

            self.logger.log_file_operation("extract", archive_path, True)
            return f"Extracted archive: {archive_path} to {destination}"

        except Exception as e:
            self.logger.log_file_operation("extract", archive_path, False, str(e))
            return f"Error extracting archive: {str(e)}"

    def validate_deployment_files(self) -> str:
        """Validate deployment configuration files"""
        try:
            issues = []
            validated_files = []

            # Check essential deployment files
            essential_files = [
                "docker-compose.yml",
                "docker-compose.production.yml",
                "Dockerfile",
                ".env.production",
                "nginx.conf"
            ]

            for file_path in essential_files:
                full_path = os.path.join(self.nutrition_platform_dir, file_path)
                if os.path.exists(full_path):
                    validated_files.append(file_path)

                    # Basic validation for specific file types
                    if file_path.endswith('.yml') or file_path.endswith('.yaml'):
                        try:
                            with open(full_path, 'r') as f:
                                content = f.read()
                                # Basic YAML validation
                                if 'docker-compose' in file_path:
                                    if 'services:' not in content:
                                        issues.append(f"{file_path}: Missing 'services' section")
                                    if 'version:' not in content:
                                        issues.append(f"{file_path}: Missing 'version' field")
                        except Exception as e:
                            issues.append(f"{file_path}: Invalid YAML format - {str(e)}")

                    elif file_path.endswith('.env'):
                        try:
                            with open(full_path, 'r') as f:
                                lines = f.readlines()
                                for line in lines:
                                    if '=' not in line and line.strip() and not line.strip().startswith('#'):
                                        issues.append(f"{file_path}: Invalid environment variable format: {line.strip()}")
                        except Exception as e:
                            issues.append(f"{file_path}: Error reading file - {str(e)}")

                else:
                    issues.append(f"Missing essential file: {file_path}")

            result = []
            result.append(f"Validated {len(validated_files)} files")
            if issues:
                result.append(f"Found {len(issues)} issues:")
                result.extend(f"  - {issue}" for issue in issues)
            else:
                result.append("All deployment files are valid!")

            return "\n".join(result)

        except Exception as e:
            return f"Error validating deployment files: {str(e)}"

    def create_deployment_package(self, environment: str = 'production', include_git: bool = False) -> str:
        """Create deployment package"""
        try:
            timestamp = datetime.now().strftime("%Y%m%d_%H%M%S")
            package_name = f"nutrition-platform-{environment}-{timestamp}.tar.gz"

            package_path = os.path.join(self.nutrition_platform_dir, package_name)

            # Files to include in deployment package
            include_patterns = [
                "*.yml",
                "*.yaml",
                "Dockerfile*",
                ".env*",
                "nginx*",
                "package*.json",
                "src/**",
                "backend/**",
                "frontend/**",
                "deployment/**",
                "scripts/**"
            ]

            # Create tar.gz archive
            with tarfile.open(package_path, 'w:gz') as tar:
                for pattern in include_patterns:
                    if '**' in pattern:
                        # Handle recursive patterns
                        base_pattern = pattern.replace('**', '')
                        for root, dirs, files in os.walk(self.nutrition_platform_dir):
                            for file in files:
                                if base_pattern in file or pattern.endswith('**'):
                                    filepath = os.path.join(root, file)
                                    arcname = os.path.relpath(filepath, self.nutrition_platform_dir)
                                    tar.add(filepath, arcname=arcname)
                    else:
                        # Handle simple patterns
                        for file_path in Path(self.nutrition_platform_dir).glob(pattern):
                            if file_path.is_file():
                                tar.add(file_path, arcname=file_path.name)

            # Also include .gitignore patterns if requested
            if include_git:
                gitignore_path = os.path.join(self.nutrition_platform_dir, '.gitignore')
                if os.path.exists(gitignore_path):
                    tar.add(gitignore_path, arcname='.gitignore')

            self.logger.log_file_operation("create_package", package_name, True)
            return f"Deployment package created: {package_name}"

        except Exception as e:
            self.logger.log_file_operation("create_package", package_name, False, str(e))
            return f"Error creating deployment package: {str(e)}"

def handle(args) -> str:
    """Handle file operation commands"""
    file_ops = FileOperations()

    try:
        if hasattr(args, 'file_action'):
            if args.file_action == 'backup':
                return file_ops.backup_files(
                    config=getattr(args, 'config', False),
                    logs=getattr(args, 'logs', False),
                    database=getattr(args, 'database', False),
                    destination=getattr(args, 'destination', None)
                )
            elif args.file_action == 'restore':
                return file_ops.restore_files(
                    backup_path=getattr(args, 'backup_path', ''),
                    config=getattr(args, 'config', False),
                    logs=getattr(args, 'logs', False),
                    database=getattr(args, 'database', False)
                )

        return "File operation not recognized"

    except Exception as e:
        return f"File operation error: {str(e)}"