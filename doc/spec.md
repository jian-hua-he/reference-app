
# Notes App Spec

## 1) Goal

Create an simple app to create, list, and delete notes. It should have two differents interfaces

- HTTP API
- Command Line

## 2) Entity

Note contains the following attributes
- id: uuid
- text: string
- created_at: RFC3339 string

## 3) Interfaces

- Create note by provide text only
- List all notes
- Delete note by id

Example for Command line interfaces:

```bash
# Create note
$ ddd_note add "This is the note text"
$ ddd_note add "Second note!"
```

```bash
# List all notes
$ ddd_note list
+----+-----------------------+----------------------+
| id | text                  | created_at           |
+----+-----------------------+----------------------+
|  1 | This is the note text | 2025-12-31T00:00:00Z |
+----+-----------------------+----------------------+
|  2 | Second note!          | 2025-12-31T00:00:00Z |
+----+-----------------------+----------------------+
```

```bash
# Delete note
$ ddd_note delete 1
```

HTTP interfaces:

```
HTTP Method: POST
Path: notes
Request Body: {"text": "This is the note text"}
Response Code:
- 200 success
- 400 invalid input
- 500 server error

HTTP Method: GET
Path: notes
Response Body: [{"id":1, "text": "This is the note text"},...]
Response Code:
- 200 success
- 500 server error

HTTP Method: DELETE
Path: notes/{id}
Response Code:
- 204 success
- 400 invalid input
- 404 not found
- 500 server error
```