#!/usr/bin/env python3
import http.server
import socketserver
import os
import sys

# Change to the current directory (where this script is located)
os.chdir(os.path.dirname(os.path.abspath(__file__)))

PORT = 3000

class MyHTTPRequestHandler(http.server.SimpleHTTPRequestHandler):
    def end_headers(self):
        self.send_header('Cache-Control', 'no-cache, no-store, must-revalidate')
        self.send_header('Pragma', 'no-cache')
        self.send_header('Expires', '0')
        super().end_headers()

if __name__ == "__main__":
    try:
        httpd = socketserver.TCPServer(("", PORT), MyHTTPRequestHandler)
        print(f"Server started at http://localhost:{PORT}", flush=True)
        print("Press Ctrl+C to stop the server", flush=True)
        
        httpd.serve_forever()
    except KeyboardInterrupt:
        print("\nServer stopped.", flush=True)
    except Exception as e:
        print(f"Error: {e}", flush=True)
        sys.exit(1)
