import requests
import bs4
from urllib.parse import urlparse
from .utils import generate_headers, handle_x_migration, get_ondemand_file_url
from .generator import ClientTransaction


def main():
    # Initialize session
    session = requests.Session()
    session.headers = generate_headers()

    # Get home page response
    home_page_response = handle_x_migration(session=session)

    # Get ondemand.s file response
    ondemand_file_url = get_ondemand_file_url(response=home_page_response)
    if not ondemand_file_url:
        raise Exception("Couldn't get ondemand file URL")

    ondemand_file = session.get(url=ondemand_file_url)
    ondemand_file_response = bs4.BeautifulSoup(ondemand_file.content, "html.parser")

    # Example 1: Generate transaction ID for client_event.json endpoint
    url1 = "https://x.com/i/api/1.1/jot/client_event.json"
    method1 = "POST"
    path1 = urlparse(url=url1).path

    # Example 2: Generate transaction ID for UserByScreenName endpoint
    url2 = "https://x.com/i/api/graphql/1VOOyvKkiI3FMmkeDNxM9A/UserByScreenName"
    method2 = "GET"
    path2 = urlparse(url=url2).path

    # Create ClientTransaction instance
    ct = ClientTransaction(
        home_page_response=home_page_response,
        ondemand_file_response=ondemand_file_response,
    )

    # Generate transaction IDs
    transaction_id1 = ct.generate_transaction_id(method=method1, path=path1)
    transaction_id2 = ct.generate_transaction_id(method=method2, path=path2)

    print(f"Transaction ID for {url1}: {transaction_id1}")
    print(f"Transaction ID for {url2}: {transaction_id2}")

    return transaction_id1, transaction_id2


if __name__ == "__main__":
    main()
