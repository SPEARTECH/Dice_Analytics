import subprocess
import os
from flask import Flask
import time


def run_with_switches():
    # Check the default browser
    if os.path.exists("C:/Program Files/Google/Chrome/Application/chrome.exe"):
        command = [
            "C:/Program Files/Google/Chrome/Application/chrome.exe", 
            '--app=http://127.0.0.1:5000', 
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
            '--app=http://127.0.0.1:5000', 
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

if __name__ == '__main__':
    stop_previous_flask_server()
    
    subprocess.Popen(['python', f'{os.path.dirname(os.path.realpath(__file__))}/server/server.py'])
    # subprocess.Popen(['python', f'{os.path.dirname(os.path.realpath(__file__))}/server/fast.py'])
    # subprocess.Popen(['python', f'./server/server.py'])

    # ADD SPLASH SCREEN?

    # Run Apped Chrome Window
    run_with_switches()



