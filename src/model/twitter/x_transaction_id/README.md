# X-Client-Transaction-Id Generator

This module generates the X-Client-Transaction-Id header needed for Twitter/X API requests.

## What is X-Client-Transaction-Id?

X-Client-Transaction-Id is a required header for many Twitter/X API endpoints. It's a complex hash generated based on:
- The Twitter site verification key from the homepage
- Animation data from SVG elements on the Twitter homepage
- The HTTP method and path of the request
- Current timestamp

## Usage

### Basic Usage

```python
import requests
import bs4
from urllib.parse import urlparse
from src.model.twitter.x_transaction_id.utils import generate_headers, handle_x_migration, get_ondemand_file_url
from src.model.twitter.x_transaction_id import ClientTransaction

# Initialize session
session = requests.Session()
session.headers = generate_headers()

# Get home page response
home_page_response = handle_x_migration(session=session)

# Get ondemand.s file response
ondemand_file_url = get_ondemand_file_url(response=home_page_response)
ondemand_file = session.get(url=ondemand_file_url)
ondemand_file_response = bs4.BeautifulSoup(ondemand_file.content, 'html.parser')

# Prepare request details
url = "https://x.com/i/api/1.1/jot/client_event.json"
method = "POST"
path = urlparse(url=url).path

# Generate transaction ID
ct = ClientTransaction(home_page_response=home_page_response, ondemand_file_response=ondemand_file_response)
transaction_id = ct.generate_transaction_id(method=method, path=path)

# Use the transaction ID in your request
session.headers["X-Client-Transaction-Id"] = transaction_id
```

### Async Usage

```python
import httpx
import bs4
from urllib.parse import urlparse
from src.model.twitter.x_transaction_id.utils import generate_headers, handle_x_migration_async, get_ondemand_file_url
from src.model.twitter.x_transaction_id import ClientTransaction

async def get_transaction_id():
    # Initialize session
    session = httpx.AsyncClient(headers=generate_headers())
    
    # Get home page response
    home_page_response = await handle_x_migration_async(session=session)
    
    # Get ondemand.s file response
    ondemand_file_url = get_ondemand_file_url(response=home_page_response)
    ondemand_file = await session.get(url=ondemand_file_url)
    ondemand_file_response = bs4.BeautifulSoup(ondemand_file.content, 'html.parser')
    
    # Prepare request details
    url = "https://x.com/i/api/1.1/jot/client_event.json"
    method = "POST"
    path = urlparse(url=url).path
    
    # Generate transaction ID
    ct = ClientTransaction(home_page_response=home_page_response, ondemand_file_response=ondemand_file_response)
    transaction_id = ct.generate_transaction_id(method=method, path=path)
    
    # Use the transaction ID in your request
    session.headers["X-Client-Transaction-Id"] = transaction_id
    
    return transaction_id, session
```

## Requirements

- requests (or httpx for async)
- beautifulsoup4 