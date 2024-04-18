import subprocess
from http.server import HTTPServer, SimpleHTTPRequestHandler
import os
import sys
import random
from flask import Flask, render_template, render_template_string, request, jsonify, send_file, make_response
from werkzeug.utils import secure_filename
# import numpy as np
import random
import json
from gevent.pywsgi import WSGIServer

def run_with_switches():
    # Check the default browser
    if os.path.exists("C:/Program Files/Google/Chrome/Application/chrome.exe"):
        command = [
            "C:/Program Files/Google/Chrome/Application/chrome.exe", 
            '--app=http://127.0.0.1:8000', 
            '--disable-pinch', 
            '--disable-extensions', 
            '--guest'
        ]
        print("Running command:", command)
        subprocess.Popen(command)
        return
    elif os.path.exists("C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe"):
        command = [
            "C:/Program Files (x86)/Microsoft/Edge/Application/msedge.exe", 
            '--app=http://127.0.0.1:8000', 
            '--disable-pinch', 
            '--disable-extensions', 
            '--guest'
        ]
        print("Running command:", command)
        subprocess.Popen(command)
        return

    print("Chromium-based browser not found or default browser not set.")

def stop_previous_flask_server():
    try:
        # Read the PID from the file
        with open(f'{os.path.expanduser("~")}/flask_server.pid', 'r') as f:
            pid = int(f.read().strip())
        
        # # Check if the Flask server process is still running
        # while True:
        #     if not os.path.exists(f'/proc/{pid}'):
        #         break  # Exit the loop if the process has exited
        #     time.sleep(1)  # Sleep for a short duration before checking again

        # Terminate the Flask server process
        command = f'taskkill /F /PID {pid}'
        subprocess.run(command, shell=True, check=True)
        print("Previous Flask server process terminated.")
    except Exception as e:
        print(f"Error stopping previous Flask server: {e}")

app = Flask(__name__)
SPLASH_FOLDER = '/splash_screen/'

# getting the name of the directory
# where the this file is present.
path = os.path.dirname(os.path.realpath(__file__))

if os.path.isdir(path + SPLASH_FOLDER):
    SPLASH_FOLDER = path + SPLASH_FOLDER
else:
    print('Splash Screen folder not found at: ' + os.getcwd())

# Routes
@app.route('/')
def index():
    html = """
   
    """

    file_path = f'{os.path.dirname(os.path.realpath(__file__))}\\templates\\index.html'

    with open(file_path, 'r') as file:
        html = ''
        for line in file:
            html += line
            
        return render_template_string(html)
        # return render('index.html')

@app.route('/api/example_api_endpoint', methods=['GET'])
def example_api_endpoint():
    # Get the data from the request
    # data = request.json.get('data') # for POST requests with data
    data = {'welcome to':'Raptor!'}

    # Perform data processing

    # Return the modified data as JSON
    return jsonify({'result': data})

if __name__ == '__main__':
    stop_previous_flask_server()

    pid_file = f'{os.path.expanduser("~")}/flask_server.pid'
    with open(pid_file, 'w') as f:
        f.write(str(os.getpid()))  # Write the PID to the file
    
    # ADD SPLASH SCREEN?

    # Run Apped Chrome Window
    run_with_switches()

    http_server = WSGIServer(("127.0.0.1", 8000), app)
    http_server.serve_forever()
    # app.run(debug=True, threaded=True, port=8000)  





