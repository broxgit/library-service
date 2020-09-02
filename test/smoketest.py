import json
import random
import re
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


def log_title(msg, decorator=None):
    d = decorator
    if not decorator:
        d = "-"
    message = ' {} '.format(msg)
    lines = d * 56
    print("\n\n{}".format(lines))
    print(message.center(56, d))
    print(lines)


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
        print("ERROR: Create book failed...")
        return None

    book = json.loads(res.content)
    return book


def delete_book(book):
    headers = {'If-Match': book['version']}
    return requests.delete(BOOKS_BY_ID_URL % book['id'], headers=headers)


def update_book(updated_book):
    headers = {'If-Match': updated_book['version']}
    return requests.put(BOOKS_BY_ID_URL % updated_book['id'], json=updated_book, headers=headers)


def create_book_test():
    log_title("Starting createBook Test")
    temp_book = create_book()
    if temp_book:
        print("SUCCESS: Created book with id: {}".format(temp_book['id']))
    return temp_book


def get_book(book):
    log_title("Starting getBookById Test")
    res = requests.get(BOOKS_BY_ID_URL % book['id'])

    if not check_status(res, 200, "getBookById"):
        print("ERROR: Get book by id failed...")
        return None

    book = json.loads(res.content)
    print("SUCCESS: Retrieved book by id successfully: {}".format(book['title']))
    return book


def update_book_test(book):
    log_title("Starting updateBookById Test")
    data = deepcopy(book)
    data['year'] = 1945

    res = update_book(data)

    if not check_status(res, 200, "updateBookById"):
        print("ERROR: Update book by id failed...")
        return None

    book = json.loads(res.content)

    if book['year'] != 1945:
        print("ERROR: Unable to update book data")
        return None

    print("SUCCESS: Updated book by id successfully: {}".format(book['title']))
    return book


def list_books():
    log_title("Starting listBooks Test")
    res = requests.get(BOOKS_URL)

    if not check_status(res, 200, "listBooks"):
        return None

    books = json.loads(res.content)
    print("SUCCESS: Retrieved a list of books.")
    return books


def delete_book_test(book):
    log_title("Starting deleteBookById Test")
    res = delete_book(book)

    if not check_status(res, 204, "deleteBookById"):
        return False

    print("SUCCESS: Deleted book with id: {}".format(book['id']))
    return True


def delete_all_books():
    res = requests.get(BOOKS_URL)
    if not check_status(res, 200, "listBooksForDeletion"):
        return None

    temp_books = json.loads(res.content)

    if not temp_books:
        return

    for x in range(0, len(temp_books)):
        book = temp_books[x]
        headers = {'If-Match': book['version']}
        res = requests.delete(BOOKS_BY_ID_URL % book['id'], headers=headers)

        if not check_status(res, 204, "deleteBookById"):
            return False

        print("Delete status code: {}".format(res.status_code))


# should fail
def try_create_duplicate_book():
    log_title("Starting Duplicated Book Test")
    temp_book = create_book()

    res = requests.post(BOOKS_URL, json=BOOK_DATA)

    if not check_status(res, 400, "createDuplicateBook"):
        return False

    print("SUCCESS: Duplicate book creation successfully blocked")

    delete_book(temp_book)
    return True


# should fail
def try_update_invalid_if_match():
    log_title("Starting Update with Invalid If-Match Tests")
    temp_book = create_book()
    ver = temp_book['version']

    # first try, completely different If-Match
    temp_book['version'] = re.sub(r'[a-zA-Z]', str(random.randint(0, 9)), ver)
    res = update_book(temp_book)

    if not check_status(res, 400, "updateInvalidIfMatch"):
        return False
    print("SUCCESS: Updating Book with Random If-Match Successfully Blocked")

    # second try, one character change in If-Match
    temp_book['version'] = ver[:-1] + "a"
    res_two = update_book(temp_book)

    if not check_status(res_two, 400, "updateInvalidIfMatch"):
        return False
    print("SUCCESS: Updating Book using If-Match with one character changed Successfully Blocked")

    temp_book['version'] = ver
    delete_book(temp_book)
    return True


# should fail
def try_delete_invalid_if_match():
    log_title("Starting Delete with Invalid If-Match Tests")
    temp_book = create_book()
    ver = temp_book['version']

    # first try, completely different If-Match
    temp_book['version'] = re.sub(r'[a-zA-Z]', str(random.randint(0, 9)), ver)
    res = delete_book(temp_book)

    if not check_status(res, 400, "deleteInvalidIfMatch"):
        return False
    print("SUCCESS: Deleting Book with Random If-Match Successfully Blocked")

    # second try, one character change in If-Match
    temp_book['version'] = ver[:-1] + "a"
    res_two = delete_book(temp_book)

    if not check_status(res_two, 400, "deleteInvalidIfMatch"):
        return False
    print("SUCCESS: Deleting Book using If-Match with one character changed Successfully Blocked")

    delete_book(temp_book)
    return True


def positive_tests():
    new_book = create_book_test()

    if not new_book:
        print("ERROR: Cannot continue positive tests. Skipping...")
        return

    created_book = get_book(new_book)

    if not created_book:
        print("ERROR: Cannot continue positive tests. Skipping...")
        return

    updated_book = update_book_test(created_book)

    if not updated_book:
        print("ERROR: Cannot continue positive tests. Skipping...")
        return

    list_books()

    delete_book_test(updated_book)


def negative_tests():
    try_create_duplicate_book()
    try_update_invalid_if_match()
    try_delete_invalid_if_match()


if __name__ == '__main__':
    delete_all_books()

    log_title("RUNNING POSITIVE SMOKE TESTS", '*')
    positive_tests()

    log_title("RUNNING NEGATIVE SMOKE TESTS", '*')
    negative_tests()
