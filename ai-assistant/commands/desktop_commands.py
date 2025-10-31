#!/usr/bin/env python3
"""
Desktop Automation Commands for Nutrition Platform AI Assistant
Handles desktop automation for local deployment tasks
"""

import os
import sys
import time
import subprocess
import pyautogui
import pygetwindow as gw
from typing import Dict, Any, Optional, Tuple

# Add parent directory to path for imports
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from config.config_manager import ConfigManager
from utils.logger import Logger

class DesktopAutomation:
    """Desktop automation for local deployment tasks"""

    def __init__(self):
        """Initialize desktop automation"""
        self.config = ConfigManager()
        self.logger = Logger()

        # Configure PyAutoGUI
        pyautogui.FAILSAFE = True
        pyautogui.PAUSE = 0.5

        # Get screen dimensions
        self.screen_width, self.screen_height = pyautogui.size()

        self.logger.info(f"Desktop automation initialized - Screen: {self.screen_width}x{self.screen_height}")

    def run_application(self, app_name: str, app_path: str = None) -> str:
        """Run application by name or path"""
        try:
            self.logger.log_desktop_automation("run_application", app_name, True)

            if app_path:
                # Run by path
                if os.path.exists(app_path):
                    subprocess.Popen([app_path])
                    return f"Started application: {app_path}"
                else:
                    return f"Application path not found: {app_path}"
            else:
                # Try to find and run by name
                if sys.platform == "win32":
                    # Windows
                    try:
                        subprocess.Popen([app_name])
                        return f"Started Windows application: {app_name}"
                    except Exception:
                        return f"Could not start Windows application: {app_name}"
                elif sys.platform == "darwin":
                    # macOS
                    try:
                        subprocess.Popen(["open", "-a", app_name])
                        return f"Started macOS application: {app_name}"
                    except Exception:
                        return f"Could not start macOS application: {app_name}"
                else:
                    # Linux
                    try:
                        subprocess.Popen([app_name])
                        return f"Started Linux application: {app_name}"
                    except Exception:
                        return f"Could not start Linux application: {app_name}"

        except Exception as e:
            self.logger.log_desktop_automation("run_application", app_name, False, str(e))
            return f"Error running application: {str(e)}"

    def click_coordinates(self, x: int, y: int) -> str:
        """Click at specific coordinates"""
        try:
            if 0 <= x <= self.screen_width and 0 <= y <= self.screen_height:
                pyautogui.click(x, y)
                self.logger.log_desktop_automation("click_coordinates", f"({x},{y})", True)
                return f"Clicked at coordinates ({x}, {y})"
            else:
                return f"Coordinates ({x}, {y}) are outside screen bounds (0,0)-({self.screen_width},{self.screen_height})"

        except Exception as e:
            self.logger.log_desktop_automation("click_coordinates", f"({x},{y})", False, str(e))
            return f"Error clicking coordinates: {str(e)}"

    def type_text(self, text: str, interval: float = 0.1) -> str:
        """Type text"""
        try:
            pyautogui.write(text, interval=interval)
            self.logger.log_desktop_automation("type_text", text[:50], True)
            return f"Typed text: {text[:50]}{'...' if len(text) > 50 else ''}"
        except Exception as e:
            self.logger.log_desktop_automation("type_text", text[:50], False, str(e))
            return f"Error typing text: {str(e)}"

    def press_key(self, key: str, modifiers: list = None) -> str:
        """Press a key or key combination"""
        try:
            if modifiers:
                # Press modifiers first
                for modifier in modifiers:
                    pyautogui.keyDown(modifier)

                # Press main key
                pyautogui.press(key)

                # Release modifiers
                for modifier in reversed(modifiers):
                    pyautogui.keyUp(modifier)
            else:
                pyautogui.press(key)

            key_combo = f"{'+'.join(modifiers)}+" if modifiers else "" + key
            self.logger.log_desktop_automation("press_key", key_combo, True)
            return f"Pressed key: {key_combo}"

        except Exception as e:
            self.logger.log_desktop_automation("press_key", key, False, str(e))
            return f"Error pressing key: {str(e)}"

    def take_screenshot(self, filename: str = None) -> str:
        """Take desktop screenshot"""
        try:
            if not filename:
                timestamp = time.strftime("%Y%m%d_%H%M%S")
                filename = f"desktop_screenshot_{timestamp}.png"

            # Ensure logs directory exists
            logs_dir = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'logs')
            os.makedirs(logs_dir, exist_ok=True)

            filepath = os.path.join(logs_dir, filename)
            screenshot = pyautogui.screenshot()
            screenshot.save(filepath)

            self.logger.log_desktop_automation("screenshot", filepath, True)
            return f"Screenshot saved: {filepath}"

        except Exception as e:
            self.logger.log_desktop_automation("screenshot", filename or "auto", False, str(e))
            return f"Error taking screenshot: {str(e)}"

    def find_and_click_image(self, image_path: str, confidence: float = 0.8, timeout: int = 10) -> str:
        """Find and click on image"""
        try:
            if not os.path.exists(image_path):
                return f"Image file not found: {image_path}"

            start_time = time.time()
            while time.time() - start_time < timeout:
                try:
                    location = pyautogui.locateOnScreen(image_path, confidence=confidence)
                    if location:
                        pyautogui.click(location)
                        self.logger.log_desktop_automation("click_image", image_path, True)
                        return f"Clicked on image: {image_path} at {location}"
                except Exception:
                    pass

                time.sleep(0.5)

            return f"Image not found on screen: {image_path}"

        except Exception as e:
            self.logger.log_desktop_automation("click_image", image_path, False, str(e))
            return f"Error finding/clicking image: {str(e)}"

    def get_window_info(self, window_title: str = None) -> str:
        """Get information about windows"""
        try:
            if window_title:
                # Get specific window
                windows = gw.getWindowsWithTitle(window_title)
                if windows:
                    window = windows[0]
                    return f"Window '{window_title}': Position({window.left},{window.top}), Size({window.width},{window.height}), Active: {window.isActive}"
                else:
                    return f"Window not found: {window_title}"
            else:
                # Get all windows
                windows = gw.getAllWindows()
                window_info = []

                for window in windows:
                    if window.title:  # Only include windows with titles
                        window_info.append({
                            'title': window.title,
                            'position': (window.left, window.top),
                            'size': (window.width, window.height),
                            'active': window.isActive
                        })

                return f"Found {len(window_info)} windows:\n" + "\n".join(
                    f"  '{w['title']}' at {w['position']} size {w['size']} active: {w['active']}"
                    for w in window_info[:10]  # Limit to first 10 windows
                )

        except Exception as e:
            return f"Error getting window info: {str(e)}"

    def activate_window(self, window_title: str) -> str:
        """Activate a window"""
        try:
            windows = gw.getWindowsWithTitle(window_title)
            if windows:
                window = windows[0]
                window.activate()
                self.logger.log_desktop_automation("activate_window", window_title, True)
                return f"Activated window: {window_title}"
            else:
                return f"Window not found: {window_title}"

        except Exception as e:
            self.logger.log_desktop_automation("activate_window", window_title, False, str(e))
            return f"Error activating window: {str(e)}"

    def scroll_screen(self, direction: str, amount: int = 5) -> str:
        """Scroll screen"""
        try:
            if direction.lower() == 'up':
                pyautogui.scroll(amount)
            elif direction.lower() == 'down':
                pyautogui.scroll(-amount)
            else:
                return f"Invalid scroll direction: {direction}. Use 'up' or 'down'"

            self.logger.log_desktop_automation("scroll", f"{direction} {amount}", True)
            return f"Scrolled {direction} by {amount} clicks"

        except Exception as e:
            self.logger.log_desktop_automation("scroll", f"{direction} {amount}", False, str(e))
            return f"Error scrolling: {str(e)}"

    def move_mouse(self, x: int, y: int, duration: float = 0.5) -> str:
        """Move mouse to coordinates"""
        try:
            if 0 <= x <= self.screen_width and 0 <= y <= self.screen_height:
                pyautogui.moveTo(x, y, duration=duration)
                self.logger.log_desktop_automation("move_mouse", f"({x},{y})", True)
                return f"Moved mouse to ({x}, {y})"
            else:
                return f"Coordinates ({x}, {y}) are outside screen bounds"

        except Exception as e:
            self.logger.log_desktop_automation("move_mouse", f"({x},{y})", False, str(e))
            return f"Error moving mouse: {str(e)}"

    def drag_mouse(self, from_x: int, from_y: int, to_x: int, to_y: int, duration: float = 0.5) -> str:
        """Drag mouse from one point to another"""
        try:
            if not (0 <= from_x <= self.screen_width and 0 <= from_y <= self.screen_height):
                return f"From coordinates ({from_x}, {from_y}) are outside screen bounds"

            if not (0 <= to_x <= self.screen_width and 0 <= to_y <= self.screen_height):
                return f"To coordinates ({to_x}, {to_y}) are outside screen bounds"

            pyautogui.moveTo(from_x, from_y)
            pyautogui.dragTo(to_x, to_y, duration=duration)

            self.logger.log_desktop_automation("drag_mouse", f"({from_x},{from_y})->({to_x},{to_y})", True)
            return f"Dragged from ({from_x}, {from_y}) to ({to_x}, {to_y})"

        except Exception as e:
            self.logger.log_desktop_automation("drag_mouse", f"({from_x},{from_y})->({to_x},{to_y})", False, str(e))
            return f"Error dragging mouse: {str(e)}"

    def wait_for_window(self, window_title: str, timeout: int = 10) -> str:
        """Wait for window to appear"""
        try:
            start_time = time.time()
            while time.time() - start_time < timeout:
                windows = gw.getWindowsWithTitle(window_title)
                if windows:
                    self.logger.log_desktop_automation("wait_for_window", window_title, True)
                    return f"Window appeared: {window_title}"
                time.sleep(0.5)

            return f"Timeout waiting for window: {window_title}"

        except Exception as e:
            self.logger.log_desktop_automation("wait_for_window", window_title, False, str(e))
            return f"Error waiting for window: {str(e)}"

    def close_window(self, window_title: str) -> str:
        """Close a window"""
        try:
            windows = gw.getWindowsWithTitle(window_title)
            if windows:
                window = windows[0]
                window.close()
                self.logger.log_desktop_automation("close_window", window_title, True)
                return f"Closed window: {window_title}"
            else:
                return f"Window not found: {window_title}"

        except Exception as e:
            self.logger.log_desktop_automation("close_window", window_title, False, str(e))
            return f"Error closing window: {str(e)}"

    def get_mouse_position(self) -> str:
        """Get current mouse position"""
        try:
            x, y = pyautogui.position()
            return f"Mouse position: ({x}, {y})"
        except Exception as e:
            return f"Error getting mouse position: {str(e)}"

    def open_terminal_and_run(self, command: str, terminal_app: str = None) -> str:
        """Open terminal and run command"""
        try:
            if sys.platform == "win32":
                # Windows
                if not terminal_app:
                    terminal_app = "cmd"
                subprocess.Popen([terminal_app, "/k", command])
                return f"Opened {terminal_app} and ran: {command}"

            elif sys.platform == "darwin":
                # macOS
                if not terminal_app:
                    terminal_app = "Terminal"
                # Use AppleScript to run command in terminal
                script = f'tell application "{terminal_app}" to do script "{command}"'
                subprocess.Popen(["osascript", "-e", script])
                return f"Opened {terminal_app} and ran: {command}"

            else:
                # Linux
                if not terminal_app:
                    terminal_app = "gnome-terminal"
                subprocess.Popen([terminal_app, "--", "bash", "-c", command])
                return f"Opened {terminal_app} and ran: {command}"

        except Exception as e:
            return f"Error opening terminal: {str(e)}"

    def simulate_user_interaction(self, actions: list) -> str:
        """Simulate a series of user interactions"""
        try:
            results = []

            for action in actions:
                action_type = action.get('type', '')
                params = action.get('params', {})

                if action_type == 'click':
                    result = self.click_coordinates(params.get('x', 0), params.get('y', 0))
                elif action_type == 'type':
                    result = self.type_text(params.get('text', ''))
                elif action_type == 'press':
                    result = self.press_key(params.get('key', ''), params.get('modifiers', []))
                elif action_type == 'wait':
                    time.sleep(params.get('seconds', 1))
                    result = f"Waited {params.get('seconds', 1)} seconds"
                elif action_type == 'screenshot':
                    result = self.take_screenshot(params.get('filename'))
                else:
                    result = f"Unknown action type: {action_type}"

                results.append(result)

                # Stop if any action fails
                if result.startswith("Error"):
                    break

            return f"User interaction simulation completed:\n" + "\n".join(f"  {r}" for r in results)

        except Exception as e:
            return f"Error in user interaction simulation: {str(e)}"

    def manage_deployment_app(self, action: str, app_name: str) -> str:
        """Manage deployment-related applications"""
        try:
            # Common deployment applications
            deployment_apps = {
                'docker': 'Docker Desktop',
                'vscode': 'Visual Studio Code',
                'terminal': 'Terminal',
                'browser': 'Google Chrome',
                'filemanager': 'Finder' if sys.platform == 'darwin' else 'Explorer'
            }

            app_display_name = deployment_apps.get(app_name.lower(), app_name)

            if action == 'start':
                return self.run_application(app_display_name)
            elif action == 'stop':
                return self.close_window(app_display_name)
            elif action == 'activate':
                return self.activate_window(app_display_name)
            elif action == 'status':
                return self.get_window_info(app_display_name)
            else:
                return f"Unknown action for deployment app: {action}"

        except Exception as e:
            return f"Error managing deployment app: {str(e)}"

def handle(args) -> str:
    """Handle desktop automation commands"""
    automation = DesktopAutomation()

    try:
        if hasattr(args, 'desktop_action'):
            if args.desktop_action == 'app':
                action = getattr(args, 'action', '')
                app_name = getattr(args, 'app_name', '')

                if not action or not app_name:
                    return "Both action and app_name are required for app management"

                return automation.manage_deployment_app(action, app_name)

        return "Desktop automation action not recognized"

    except Exception as e:
        return f"Desktop automation error: {str(e)}"