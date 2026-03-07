#!/usr/bin/env python3
"""Create demo SQLite databases for VHS recording.

Creates two databases:
  /tmp/asql-demo.db          — "production" data
  /tmp/asql-demo-staging.db  — "staging" data with intentional differences

The staging DB shares the same schema but has:
  - An extra user (Dave) with a different email domain
  - One email address changed (Bob → bob@staging.io)
  - Different published flags on posts
  - An additional post
"""
import os
import sqlite3
import sys

SCHEMA = """
CREATE TABLE users (
    id         INTEGER PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT UNIQUE,
    created_at TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE posts (
    id        INTEGER PRIMARY KEY,
    user_id   INTEGER REFERENCES users(id),
    title     TEXT NOT NULL,
    body      TEXT,
    published BOOLEAN DEFAULT 0
);
CREATE TABLE tags (
    id    INTEGER PRIMARY KEY,
    name  TEXT NOT NULL UNIQUE,
    color TEXT
);
"""

PROD_DATA = """
INSERT INTO users (name, email) VALUES
    ('Alice',   'alice@example.com'),
    ('Bob',     'bob@example.com'),
    ('Charlie', 'charlie@example.com');

INSERT INTO posts (user_id, title, body, published) VALUES
    (1, 'Getting Started with SQL', 'A beginner guide',             1),
    (1, 'Advanced Queries',         'Deep dive into joins',         1),
    (2, 'Database Design',          'Best practices for schemas',   0);

INSERT INTO tags (name, color) VALUES
    ('sql',      '#3B82F6'),
    ('tutorial', '#10B981'),
    ('database', '#F59E0B');
"""

STAGING_DATA = """
INSERT INTO users (name, email) VALUES
    ('Alice',   'alice@example.com'),
    ('Bob',     'bob@staging.io'),
    ('Charlie', 'charlie@example.com'),
    ('Dave',    'dave@staging.io');

INSERT INTO posts (user_id, title, body, published) VALUES
    (1, 'Getting Started with SQL', 'A beginner guide',             1),
    (1, 'Advanced Queries',         'Deep dive into joins',         0),
    (2, 'Database Design',          'Best practices for schemas',   1),
    (4, 'Staging Environment Tips', 'How to use staging effectively', 0);

INSERT INTO tags (name, color) VALUES
    ('sql',      '#3B82F6'),
    ('tutorial', '#10B981'),
    ('database', '#F59E0B');
"""


def create_db(path, data):
    if os.path.exists(path):
        os.remove(path)
    conn = sqlite3.connect(path)
    conn.executescript(SCHEMA + data)
    conn.close()


if __name__ == "__main__":
    prod_path = sys.argv[1] if len(sys.argv) > 1 else "/tmp/asql-demo.db"
    staging_path = sys.argv[2] if len(sys.argv) > 2 else "/tmp/asql-demo-staging.db"
    create_db(prod_path, PROD_DATA)
    create_db(staging_path, STAGING_DATA)
