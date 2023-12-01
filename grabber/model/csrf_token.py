import hashlib
import os


def generate_csrf_token() -> str:
    return hashlib.md5(os.urandom(32)).hexdigest()
