#!/usr/bin/env python3
"""
Web Automation Commands for Nutrition Platform AI Assistant
Handles web automation for deployment verification and testing
"""

import os
import sys
import time
import json
import requests
from selenium import webdriver
from selenium.webdriver.common.by import By
from selenium.webdriver.support.ui import WebDriverWait
from selenium.webdriver.support import expected_conditions as EC
from selenium.webdriver.chrome.options import Options
from selenium.webdriver.firefox.options import FirefoxOptions
from selenium.webdriver.common.action_chains import ActionChains
from selenium.webdriver.common.keys import Keys
from typing import Dict, Any, List, Optional

# Add parent directory to path for imports
sys.path.append(os.path.dirname(os.path.dirname(os.path.abspath(__file__))))

from config.config_manager import ConfigManager
from utils.logger import Logger

class WebAutomation:
    """Web automation for deployment verification"""

    def __init__(self):
        """Initialize web automation"""
        self.config = ConfigManager()
        self.logger = Logger()
        self.driver = None
        self.browser_type = self.config.get_config("web_automation.browser", "chrome")

    def start_browser(self, headless: bool = True) -> bool:
        """Start web browser"""
        try:
            if self.browser_type.lower() == "chrome":
                options = Options()
                if headless:
                    options.add_argument('--headless')
                options.add_argument('--no-sandbox')
                options.add_argument('--disable-dev-shm-usage')
                options.add_argument('--disable-gpu')
                options.add_argument('--window-size=1920,1080')

                # Add user agent to avoid detection
                options.add_argument('--user-agent=Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36')

                self.driver = webdriver.Chrome(options=options)

            elif self.browser_type.lower() == "firefox":
                options = FirefoxOptions()
                if headless:
                    options.add_argument('--headless')

                self.driver = webdriver.Firefox(options=options)

            self.logger.info(f"Started {self.browser_type} browser")
            return True

        except Exception as e:
            self.logger.error(f"Failed to start browser: {str(e)}")
            return False

    def close_browser(self):
        """Close web browser"""
        try:
            if self.driver:
                self.driver.quit()
                self.driver = None
                self.logger.info("Browser closed")
        except Exception as e:
            self.logger.error(f"Error closing browser: {str(e)}")

    def navigate_to(self, url: str, wait_time: int = 10) -> bool:
        """Navigate to URL"""
        try:
            if not self.driver:
                if not self.start_browser():
                    return False

            self.logger.log_web_automation("navigation", url, True, f"Navigating to {url}")
            self.driver.get(url)

            # Wait for page to load
            WebDriverWait(self.driver, wait_time).until(
                lambda driver: driver.execute_script("return document.readyState") == "complete"
            )

            return True

        except Exception as e:
            self.logger.log_web_automation("navigation", url, False, str(e))
            return False

    def wait_for_element(self, selector: str, by: str = "css", timeout: int = 10) -> bool:
        """Wait for element to be present"""
        try:
            if not self.driver:
                return False

            if by.lower() == "css":
                locator = (By.CSS_SELECTOR, selector)
            elif by.lower() == "xpath":
                locator = (By.XPATH, selector)
            elif by.lower() == "id":
                locator = (By.ID, selector)
            elif by.lower() == "class":
                locator = (By.CLASS_NAME, selector)
            else:
                locator = (By.CSS_SELECTOR, selector)

            WebDriverWait(self.driver, timeout).until(
                EC.presence_of_element_located(locator)
            )

            return True

        except Exception as e:
            self.logger.error(f"Element not found: {selector} ({by}) - {str(e)}")
            return False

    def click_element(self, selector: str, by: str = "css", timeout: int = 10) -> bool:
        """Click on element"""
        try:
            if not self.driver:
                return False

            if by.lower() == "css":
                locator = (By.CSS_SELECTOR, selector)
            elif by.lower() == "xpath":
                locator = (By.XPATH, selector)
            elif by.lower() == "id":
                locator = (By.ID, selector)
            elif by.lower() == "class":
                locator = (By.CLASS_NAME, selector)
            else:
                locator = (By.CSS_SELECTOR, selector)

            element = WebDriverWait(self.driver, timeout).until(
                EC.element_to_be_clickable(locator)
            )

            element.click()
            self.logger.log_web_automation("click", selector, True, f"Clicked element: {selector}")
            return True

        except Exception as e:
            self.logger.log_web_automation("click", selector, False, str(e))
            return False

    def fill_form(self, selector: str, text: str, by: str = "css", timeout: int = 10) -> bool:
        """Fill form field"""
        try:
            if not self.driver:
                return False

            if by.lower() == "css":
                locator = (By.CSS_SELECTOR, selector)
            elif by.lower() == "xpath":
                locator = (By.XPATH, selector)
            elif by.lower() == "id":
                locator = (By.ID, selector)
            elif by.lower() == "name":
                locator = (By.NAME, selector)
            else:
                locator = (By.CSS_SELECTOR, selector)

            element = WebDriverWait(self.driver, timeout).until(
                EC.presence_of_element_located(locator)
            )

            element.clear()
            element.send_keys(text)
            self.logger.log_web_automation("form_fill", selector, True, f"Filled {selector} with: {text[:50]}...")
            return True

        except Exception as e:
            self.logger.log_web_automation("form_fill", selector, False, str(e))
            return False

    def get_text(self, selector: str, by: str = "css", timeout: int = 10) -> Optional[str]:
        """Get text from element"""
        try:
            if not self.driver:
                return None

            if by.lower() == "css":
                locator = (By.CSS_SELECTOR, selector)
            elif by.lower() == "xpath":
                locator = (By.XPATH, selector)
            elif by.lower() == "id":
                locator = (By.ID, selector)
            else:
                locator = (By.CSS_SELECTOR, selector)

            element = WebDriverWait(self.driver, timeout).until(
                EC.presence_of_element_located(locator)
            )

            text = element.text.strip()
            return text

        except Exception as e:
            self.logger.error(f"Error getting text from {selector}: {str(e)}")
            return None

    def take_screenshot(self, filename: str = None) -> Optional[str]:
        """Take screenshot"""
        try:
            if not self.driver:
                return None

            if not filename:
                timestamp = time.strftime("%Y%m%d_%H%M%S")
                filename = f"screenshot_{timestamp}.png"

            # Ensure logs directory exists
            logs_dir = os.path.join(os.path.dirname(os.path.dirname(os.path.abspath(__file__))), 'logs')
            os.makedirs(logs_dir, exist_ok=True)

            filepath = os.path.join(logs_dir, filename)
            self.driver.save_screenshot(filepath)

            self.logger.info(f"Screenshot saved: {filepath}")
            return filepath

        except Exception as e:
            self.logger.error(f"Error taking screenshot: {str(e)}")
            return None

    def test_deployment(self, url: str, health_check: bool = True, load_test: bool = False) -> Dict[str, Any]:
        """Test deployed application"""
        try:
            self.logger.log_web_automation("deployment_test", url, True, "Starting deployment tests")

            results = {
                "url": url,
                "tests": [],
                "overall_success": True,
                "timestamp": time.time()
            }

            # Basic connectivity test
            connectivity_result = self._test_connectivity(url)
            results["tests"].append(connectivity_result)
            if not connectivity_result["success"]:
                results["overall_success"] = False

            # Navigate to the application
            if connectivity_result["success"]:
                if not self.navigate_to(url):
                    results["tests"].append({
                        "name": "navigation",
                        "success": False,
                        "error": "Failed to navigate to application"
                    })
                    results["overall_success"] = False
                else:
                    results["tests"].append({
                        "name": "navigation",
                        "success": True,
                        "message": "Successfully navigated to application"
                    })

                    # Take screenshot
                    screenshot_path = self.take_screenshot()
                    if screenshot_path:
                        results["screenshot"] = screenshot_path

                    # Health check if requested
                    if health_check:
                        health_result = self._perform_web_health_check()
                        results["tests"].append(health_result)
                        if not health_result["success"]:
                            results["overall_success"] = False

                    # Load test if requested
                    if load_test:
                        load_result = self._perform_load_test()
                        results["tests"].append(load_result)
                        if not load_result["success"]:
                            results["overall_success"] = False

            return results

        except Exception as e:
            self.logger.log_web_automation("deployment_test", url, False, str(e))
            return {
                "url": url,
                "tests": [],
                "overall_success": False,
                "error": str(e),
                "timestamp": time.time()
            }

    def _test_connectivity(self, url: str) -> Dict[str, Any]:
        """Test basic connectivity"""
        try:
            start_time = time.time()
            response = requests.get(url, timeout=10)
            response_time = time.time() - start_time

            return {
                "name": "connectivity",
                "success": response.status_code < 400,
                "status_code": response.status_code,
                "response_time": response_time,
                "message": f"HTTP {response.status_code} in {response_time:.2f}s"
            }

        except requests.exceptions.RequestException as e:
            return {
                "name": "connectivity",
                "success": False,
                "error": str(e)
            }

    def _perform_web_health_check(self) -> Dict[str, Any]:
        """Perform web-based health check"""
        try:
            # Look for common health indicators
            health_indicators = [
                "api/health",
                "health",
                "status",
                "ping"
            ]

            health_results = []

            for indicator in health_indicators:
                try:
                    # Try to find health check endpoint
                    if indicator.startswith("api/"):
                        health_url = f"{self.driver.current_url.rstrip('/')}/{indicator}"
                    else:
                        health_url = f"{self.driver.current_url.rstrip('/')}/{indicator}"

                    # Navigate to health check URL
                    self.driver.get(health_url)

                    # Check if page loads successfully
                    WebDriverWait(self.driver, 5).until(
                        lambda driver: driver.execute_script("return document.readyState") == "complete"
                    )

                    health_results.append({
                        "endpoint": indicator,
                        "success": True,
                        "url": health_url
                    })

                except Exception:
                    health_results.append({
                        "endpoint": indicator,
                        "success": False,
                        "error": "Health check endpoint not accessible"
                    })

            # Return overall health check result
            success_count = sum(1 for result in health_results if result["success"])

            return {
                "name": "health_check",
                "success": success_count > 0,
                "details": health_results,
                "message": f"Health check: {success_count}/{len(health_indicators)} endpoints accessible"
            }

        except Exception as e:
            return {
                "name": "health_check",
                "success": False,
                "error": str(e)
            }

    def _perform_load_test(self, duration: int = 30, max_concurrent: int = 5) -> Dict[str, Any]:
        """Perform basic load test"""
        try:
            # This is a simplified load test
            # In a real scenario, you might use multiple threads or external tools

            start_time = time.time()
            request_count = 0
            errors = 0

            while time.time() - start_time < duration:
                try:
                    # Refresh the current page
                    self.driver.refresh()
                    WebDriverWait(self.driver, 5).until(
                        lambda driver: driver.execute_script("return document.readyState") == "complete"
                    )
                    request_count += 1

                    # Small delay between requests
                    time.sleep(1)

                except Exception as e:
                    errors += 1
                    time.sleep(0.5)

            total_time = time.time() - start_time

            return {
                "name": "load_test",
                "success": errors / request_count < 0.1 if request_count > 0 else False,
                "requests": request_count,
                "errors": errors,
                "duration": total_time,
                "rps": request_count / total_time if total_time > 0 else 0,
                "error_rate": errors / request_count if request_count > 0 else 0,
                "message": f"Load test: {request_count} requests in {total_time:.2f}s ({errors} errors)"
            }

        except Exception as e:
            return {
                "name": "load_test",
                "success": False,
                "error": str(e)
            }

    def scrape_page_data(self, selectors: List[str] = None, by: str = "css") -> Dict[str, Any]:
        """Scrape data from current page"""
        try:
            if not self.driver:
                return {"error": "Browser not started"}

            data = {}

            if not selectors:
                # Default selectors for common elements
                selectors = [
                    "title",
                    "h1",
                    "h2",
                    ".content",
                    ".main",
                    "#app"
                ]

            for selector in selectors:
                try:
                    if by.lower() == "css":
                        locator = (By.CSS_SELECTOR, selector)
                    elif by.lower() == "xpath":
                        locator = (By.XPATH, selector)
                    else:
                        locator = (By.CSS_SELECTOR, selector)

                    elements = self.driver.find_elements(*locator)

                    if elements:
                        if len(elements) == 1:
                            data[selector] = elements[0].text.strip()
                        else:
                            data[selector] = [elem.text.strip() for elem in elements if elem.text.strip()]

                except Exception as e:
                    data[selector] = f"Error: {str(e)}"

            return {
                "success": True,
                "data": data,
                "url": self.driver.current_url,
                "timestamp": time.time()
            }

        except Exception as e:
            return {
                "success": False,
                "error": str(e)
            }

    def check_page_performance(self) -> Dict[str, Any]:
        """Check page performance metrics"""
        try:
            if not self.driver:
                return {"error": "Browser not started"}

            # Get performance metrics using JavaScript
            performance_script = """
            const navigation = performance.getEntriesByType('navigation')[0];
            const paint = performance.getEntriesByType('paint');

            return {
                'dns_lookup': navigation.domainLookupEnd - navigation.domainLookupStart,
                'tcp_connect': navigation.connectEnd - navigation.connectStart,
                'server_response': navigation.responseEnd - navigation.requestStart,
                'dom_processing': navigation.domContentLoadedEventEnd - navigation.domContentLoadedEventStart,
                'total_load_time': navigation.loadEventEnd - navigation.navigationStart,
                'first_paint': paint.find(p => p.name === 'first-paint')?.startTime || 0,
                'first_contentful_paint': paint.find(p => p.name === 'first-contentful-paint')?.startTime || 0
            };
            """

            metrics = self.driver.execute_script(performance_script)

            return {
                "success": True,
                "metrics": metrics,
                "url": self.driver.current_url,
                "timestamp": time.time()
            }

        except Exception as e:
            return {
                "success": False,
                "error": str(e)
            }

    def check_broken_links(self, domain: str = None) -> Dict[str, Any]:
        """Check for broken links on the page"""
        try:
            if not self.driver:
                return {"error": "Browser not started"}

            # Find all links on the page
            link_script = """
            const links = Array.from(document.querySelectorAll('a[href]'));
            return links.map(link => ({
                'url': link.href,
                'text': link.textContent.trim(),
                'is_external': link.hostname !== window.location.hostname
            }));
            """

            links = self.driver.execute_script(link_script)

            broken_links = []
            working_links = []

            # Check each link (simplified - only checks if URL is reachable)
            for link in links[:10]:  # Limit to first 10 links for performance
                try:
                    if link["is_external"] and domain:
                        # For external links, just check if domain is reachable
                        response = requests.head(link["url"], timeout=5)
                        if response.status_code < 400:
                            working_links.append(link)
                        else:
                            broken_links.append({**link, "error": f"HTTP {response.status_code}"})
                    else:
                        working_links.append(link)

                except Exception as e:
                    broken_links.append({**link, "error": str(e)})

            return {
                "success": True,
                "total_links": len(links),
                "working_links": len(working_links),
                "broken_links": len(broken_links),
                "broken_links_details": broken_links[:5],  # Show first 5 broken links
                "url": self.driver.current_url,
                "timestamp": time.time()
            }

        except Exception as e:
            return {
                "success": False,
                "error": str(e)
            }

def handle(args) -> str:
    """Handle web automation commands"""
    automation = WebAutomation()

    try:
        if hasattr(args, 'web_action'):
            if args.web_action == 'test':
                url = getattr(args, 'url', '')
                health_check = getattr(args, 'health_check', True)
                load_test = getattr(args, 'load_test', False)

                if not url:
                    return "URL is required for web testing"

                result = automation.test_deployment(url, health_check, load_test)

                # Pretty print results
                return json.dumps(result, indent=2, default=str)

        return "Web automation action not recognized"

    finally:
        automation.close_browser()