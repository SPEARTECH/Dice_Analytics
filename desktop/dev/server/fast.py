from fastapi import FastAPI, Request
from fastapi.responses import HTMLResponse, JSONResponse
from fastapi.templating import Jinja2Templates
import os
import random
import json
from starlette.concurrency import run_in_threadpool

app = FastAPI()
templates = Jinja2Templates(directory="templates")

# Since FastAPI is an ASGI application, it's inherently asynchronous.
# If you have blocking functions (like regular file I/O operations or CPU-bound tasks without async support),
# it's recommended to run them in a threadpool.

@app.get("/", response_class=HTMLResponse)
async def index(request: Request):
    # In FastAPI, request templates require a "Request" object to be passed explicitly.
    file_path = os.path.join(os.path.dirname(os.path.realpath(__file__)), "templates", "index.html")
    with open(file_path, 'r') as file:
        html_content = file.read()
        # For simplicity, rendering HTML directly from string here. 
        # In practice, consider using templates.render() for using Jinja2 templates.
    return HTMLResponse(content=html_content)

# Convert your helper functions here as needed. 
# Remember to update them to async functions if they perform I/O operations.

@app.get('/api/example_api_endpoint')
def example_api_endpoint():
    # Get the data from the request
    # data = request.json.get('data') # for POST requests with data
    data = {'welcome to':'Raptor!'}

    # Perform data processing

    # Return the modified data as JSON
    return json.dumps({'result': data})


@app.post("/api/simulate")
async def simulate(request: Request):
    json_data = await request.json()
    print(json_data)

    from cython_modules import run_python

    return run_python.main(    
        json_data['number'],
        json_data['winnings'],
        json_data['increase'],
        json_data['user_bet'],
        json_data['rolls'],
        json_data['recovery_rolls'],
        json_data['max_losing_streak'],
)


# Additional routes here...

if __name__ == "__main__":
    import uvicorn
    uvicorn.run(app, host="127.0.0.1", port=8000)