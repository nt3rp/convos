## Convos

A simple backend for an email-like application written in Go.

## Installation

Assuming that you have a basic Go setup ([as described here](https://golang.org/doc/code.html)):

- Get all dependencies: `go get -d -t -v ./...`

You will also need postgres v9.4, and a database called `convos` setup. You can create the database like so:

```
psql -U postgres -c 'CREATE DATABASE convos OWNER <your username>;'
```

You will also need to run database migrations. Migrations are based off of [`github.com/mattes/migrate`](https://github.com/mattes/migrate).

```
migrate -url "postgres://localhost/convos?sslmode=disable" -path ./migrations up
```

## Running the application

The application runs on port `8080` by default,

```bash
./server.sh
```

## Assumptions / Constraints

- A message only has one sender
- A message only has one recipient
- Authorization / Authentication are beyond the scope of this project
- Server / database configuration are beyond the scope of the project: the code only need work in simple local development environment
- Input sanitization is beyond the scope of this project

## API

All API methods can be called via HTTP requests, and try to follow RESTful conventions except as noted.

All responses are in JSON format, wrapped in a JSON envelope of the following form:

```
{
    "response": <response from method, typically a list or object>
    "meta": {
        "count": <number of results in `response`>
    }
    "error": { // Or null, if no errors
        "message": <Error message>
    }
}
```

Trailing slashes are required for all endpoints.

In order to 'authorize' the user, you will need to set the `X-USER-API-KEY` header in your request to the user of the
person using the system (this is by no means a good way to authorize a user, but setting up authentication and
authorization, as stated in the assumptions, is beyond the scope of this project).

```bash
curl -X GET \
-H 'X-USER-API-KEY: 1' \
...
```

All `convo` objects have a similar format:

```
{
    "id":12,                // integer; API / DB identifier for the `convo` object
    "sender":1,             // integer; user id of the person who sent the message
    "recipient":2,          // integer; user id of the person who will receive the message
    "parent":8,             // integer; API / DB identifier for the parent of this thread. If no parent, will match `id`
    "subject":"FIRST POST", // string (<= 140 characters); subject of the `convo`
    "body":"Woohoo",        // string (<= 64000 characters); body of the `convo`
    "read":true,            // boolean; if the user provided by `X-USER-API-KEY` has read this message
    "replies":null          // unused; planned feature to show replies of this convo, if any
}
```

All success status codes are simply **200 Success**.

### `GET` convos/

Retrieves all top-level conversations (i.e. conversations with no prior discussion)

#### Response

A list of of `convo` objects.

#### Errors

- **500 Server Error**: If there are problems connecting to the database or anything unexpected.

#### Example
```bash
curl -X GET \
-H 'X-USER-API-KEY: 1' \
http://localhost:8080/
```

### `POST` convos/

Create a new conversation.

#### Parameters
A JSON-encoded `convo` object. The following keys are required:

- **recipient**: *integer*, the user id of the recipient
- **subject**: *string (140 characters or less)*, the subject of the conversation
- **body**: *string (64k characters or less)*, the body of the conversation

The following keys are optional:
- **parent**: *integer*, the id of another conversation to reply to

#### Response

Returns the complete `convo` object created by this request.

#### Errors

- **404 Not Found**: The user is not a sender or reciever of the parent thread (if `parent` provided). See caveats.
- **500 Server Error**: If there are problems connecting to the database, there is a problem decoding the JSON, or anything unexpected.

#### Caveats

- The `sender` will always be set to the user provided by the `X-USER-API-KEY` header.
- Conversations are automatically marked as read for the current user.
- Normally, if a user tried to reply to a thread and they were neither a sender or receiver, we should return a
**403 Forbidden** or **401 Unauthorized**. Instead, we return a **404 Not Found** so that the user does not know about
other messages in the system.

#### Example

```bash
curl -X POST \
-H 'X-USER-API-KEY: 1' \
-d '{"recipient":2,"subject":"FIRST POST","body":"Woohoo"}' \
"http://localhost:8080/convos/"
```

### `GET` convos/:id/

Retrieves an individual conversation.

#### Response

A single `convo` objects.

#### Errors

- **404 Not Found**: The user is not a sender or reciever of the conversation. See caveats.
- **500 Server Error**: If there are problems connecting to the database, or anything unexpected.

#### Caveats

- Normally, if a user tried to reply to a thread and they were neither a sender or receiver, we should return a
**403 Forbidden** or **401 Unauthorized**. Instead, we return a **404 Not Found** so that the user does not know about
other messages in the system.

#### Example
```bash
curl -X GET \
-H 'X-USER-API-KEY: 1' \
http://localhost:8080/1/
```

### `PATCH` convos/:id/

Edits an existing conversation.

#### Parameters
A JSON-encoded patch object. It will only accept the following keys:

- **read**: *string*, whether the given conversation should be marked as read (*"true"*) or not (*"false"*)

#### Response

Returns the patched `convo` object.

#### Errors

- **404 Not Found**: The user is not a sender or reciever of the conversation. See caveats.
- **500 Server Error**: If there are problems connecting to the database, there is a problem decoding the JSON, or anything unexpected.

#### Caveats

- Normally, if a user tried to reply to a thread and they were neither a sender or receiver, we should return a
**403 Forbidden** or **401 Unauthorized**. Instead, we return a **404 Not Found** so that the user does not know about
other messages in the system.

#### Example
```bash
curl -X PATCH \
-H 'X-USER-API-KEY: 1' \
-d '{"read": "true"}' \
http://localhost:8080/convos/1/
```

### `DELETE` convos/:id/

Deletes an individual conversation and all conversations that have this conversation as a parent.

#### Response

A string, "success".

#### Errors

- **404 Not Found**: The user is not a sender or reciever of the thread. See caveats.
- **500 Server Error**: If there are problems connecting to the database, or anything unexpected.

#### Caveats

- Normally, if a user tried to reply to a thread and they were neither a sender or receiver, we should return a
**403 Forbidden** or **401 Unauthorized**. Instead, we return a **404 Not Found** so that the user does not know about
other messages in the system.

#### Example
```bash
curl -X DELETE \
-H 'X-USER-API-KEY: 1' \
http://localhost:8080/convos/1/
```

### `POST` convos/:id/reply/

Create a reply to an individual conversation

#### Parameters
See `POST convos/`

Unlike `POST convos/`, the **parent** key will be ignored (and replaced with `:id`)

#### Response

Returns the complete `convo` object created by this request.

#### Errors

- **404 Not Found**: The user is not a sender or reciever of the thread. See caveats.
- **500 Server Error**: If there are problems connecting to the database, there is a problem decoding the JSON, or anything unexpected.

#### Caveats

- The `sender` will always be set to the user provided by the `X-USER-API-KEY` header.
- The `parent` will automatically be set as `:id`
- Conversations are automatically marked as read for the current user.
- Normally, if a user tried to reply to a thread and they were neither a sender or receiver, we should return a
**403 Forbidden** or **401 Unauthorized**. Instead, we return a **404 Not Found** so that the user does not know about
other messages in the system.

#### Example

```bash
curl -X POST \
-H 'X-USER-API-KEY: 1' \
-d '{"recipient":2,"subject":"FIRST POST","body":"Woohoo"}' \
"http://localhost:8080/convos/5/"
```

## Database

A summary of the various tables and their structure, along with the design decisions made, are listed here.
The actual SQL migrations can be found in the `migrations` folder.

### `users`

This table (or a similar table with an `id` column) is assumed to be provided.

```
            Table "public.users"
  Column  |          Type          | Modifiers
----------+------------------------+-----------
 id       | integer                | not null
 fullname | character varying(255) |
Indexes:
    "users_pkey" PRIMARY KEY, btree (id)
Referenced by:
    TABLE "convos" CONSTRAINT "convos_recipient_id_fkey" FOREIGN KEY (recipient_id) REFERENCES users(id)
    TABLE "convos" CONSTRAINT "convos_sender_id_fkey" FOREIGN KEY (sender_id) REFERENCES users(id)
    TABLE "read_status" CONSTRAINT "read_status_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET DEFAULT
```

### `convos`

Stores information about conversations.

While it is not strictly necessary to have `parent_id` as part of this table (it could have been in its own table),
implementing them in the `convos` tables makes the implementation marginally simpler, as we need to obtain the
`parent_id` whenever we fetch a `convo` object.

`sender_id` and `recipient_id` could have been moved to a separate table, but since I assumed that there is only one
sender and one recipient, then it makes sense for these to be core components of a message.

`subject` and `body` are as they are because of the listed requirement (140 character subject, 64000 body).

As a nice side-benefit. Whenever a parent thread is deleted, so too are its children (as would be the case with
threaded email messages).

```
                                    Table "public.convos"
    Column    |           Type           |                      Modifiers
--------------+--------------------------+-----------------------------------------------------
 id           | integer                  | not null default nextval('convos_id_seq'::regclass)
 parent_id    | integer                  | not null
 sender_id    | integer                  | not null
 recipient_id | integer                  | not null
 subject      | character varying(140)   | not null
 body         | character varying(64000) | not null
Indexes:
    "convos_pkey" PRIMARY KEY, btree (id)
Foreign-key constraints:
    "convos_parent_id_fkey" FOREIGN KEY (parent_id) REFERENCES convos(id) ON DELETE CASCADE
    "convos_recipient_id_fkey" FOREIGN KEY (recipient_id) REFERENCES users(id)
    "convos_sender_id_fkey" FOREIGN KEY (sender_id) REFERENCES users(id)
Referenced by:
    TABLE "convos" CONSTRAINT "convos_parent_id_fkey" FOREIGN KEY (parent_id) REFERENCES convos(id) ON DELETE CASCADE
    TABLE "read_status" CONSTRAINT "read_status_thread_id_fkey" FOREIGN KEY (thread_id) REFERENCES convos(id) ON DELETE CASCADE
```

### `read_status`

A simple relationship table which stores information about which threads have been read by whom.

If we delete a thread, the read status will also be deleted.

If we delete a user, we change the `user_id` to 0 (to avoid any strangeness on things that may depend on `user_id`
being present)

```
        Table "public.read_status"
  Column   |  Type   |     Modifiers
-----------+---------+--------------------
 thread_id | integer | not null
 user_id   | integer | not null default 0
Foreign-key constraints:
    "read_status_thread_id_fkey" FOREIGN KEY (thread_id) REFERENCES convos(id) ON DELETE CASCADE
    "read_status_user_id_fkey" FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE SET DEFAULT
```