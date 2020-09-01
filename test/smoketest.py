import json
from copy import deepcopy

import requests

URL = "http://localhost:8081/library-service/v1"
BOOKS_ENDPOINT = "/books"
BOOK_BY_ID_ENDPOINT = "/books/%s"

BOOKS_URL = URL + BOOKS_ENDPOINT
BOOKS_BY_ID_URL = URL + BOOK_BY_ID_ENDPOINT

BOOK_DATA = {
    "title": "The Great Gatsby",
    "authors": [
        "F. Scott Fitzgerald"
    ],
    "year": 1925,
    "comment": "The story of the mysteriously wealthy Jay Gatsby and his love for the beautiful Daisy Buchanan."
}


def check_status(res, expected_status, operation):
    if res.status_code != expected_status:
        print("There was an error with operation: {}!".format(operation))
        content, is_json = parse_response(res)
        print("\nResponse Status/Reason: {} / {}".format(res.status_code, res.reason))
        if is_json:
            print(json.dumps(content, indent=4))
        else:
            print(content)
        return False

    return True


def parse_response(response):
    is_json = False
    content = None
    try:
        content = response.json()
        is_json = True
    except ValueError:
        try:
            content = json.loads(response.content)
            is_json = True
        except Exception:
            pass
    if not content:
        if response.text:
            content = response.text
        else:
            content = response.content
    return content, is_json


def create_book():
    res = requests.post(BOOKS_URL, json=BOOK_DATA)

    if not check_status(res, 201, "createBook"):
        exit(-1)

    book = json.loads(res.content)
    print("Successfully created book with id: {}".format(book['id']))
    return book['id']


def get_book(id):
    res = requests.get(BOOKS_BY_ID_URL % id)

    if not check_status(res, 200, "getBookById"):
        exit(-1)

    book = json.loads(res.content)
    print("Retrieved book by id successfully: {}".format(book['title']))
    return book


def update_book(book):
    data = deepcopy(book)
    data['year'] = 1945
    headers = {'If-Match': "\"{}\"".format(data['version'])}
    res = requests.put(BOOKS_BY_ID_URL % book_id, json=data, headers=headers)

    if not check_status(res, 200, "getBookById"):
        exit(-1)

    book = json.loads(res.content)

    if book['year'] != 1945:
        print("Unable to update book data")
        exit(-2)

    print("Updated book by id successfully: {}".format(book['title']))
    return book


def list_books():
    res = requests.get(BOOKS_URL)

    if not check_status(res, 200, "listBooks"):
        exit(-1)

    books = json.loads(res.content)
    print("Successfully retrieved a list of books.")
    return books


def delete_book(id, book):
    headers = {'If-Match': "\"{}\"".format(book['version'])}
    res = requests.delete(BOOKS_BY_ID_URL % id, headers=headers)

    if not check_status(res, 204, "deleteBookById"):
        exit(-1)

    print("Successfully deleted book with id: {} with HTTP Status Code: {}".format(id, res.status_code))


def delete_all_books():
    books = list_books()

    if not books:
        return

    for x in range(0, len(books)):
        book = books[x]
        headers = {'If-Match': "\"{}\"".format(book['version'])}
        res = requests.delete(BOOKS_BY_ID_URL % id, headers=headers)

        if not check_status(res, 204, "deleteBookById"):
            exit(-1)

        print(res.content)


if __name__ == '__main__':
    # delete_all_books()

    book_id = create_book()

    created_book = get_book(book_id)

    updated_book = update_book(created_book)

    books = list_books()

    delete_book(book_id, updated_book)
