#!/usr/bin/env python3
import http.server
import socketserver
import os
import webbrowser
import time
import threading

# Change to the current directory (where this script is located)
os.chdir(os.path.dirname(os.path.abspath(__file__)))

PORT = 3000

class MyHTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
        self.send_header('Pragma', 'no-cache')
        self.send_header('Expires', '0')
        super().end_headers()

def open_browser():
    time.sleep(1)
    webbrowser.open(f'http://localhost:{PORT}')

if __name__ == "__main__":
    try:
        httpd = socketserver.TCPServer(("", PORT), MyHTTPRequestHandler)
        print(f"Server started at http://localhost:{PORT}")
        print("Press Ctrl+C to stop the server")
        
        # Open browser in a separate thread
        browser_thread = threading.Thread(target=open_browser)
        browser_thread.daemon = True
        browser_thread.start()
        
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nServer stopped.")
    except Exception as e:
        print(f"Error: {e}")
