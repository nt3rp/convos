## Convos

TODO

## Installation

TODO

## Running the application

TODO

## API

All API methods can be called via HTTP requests, and try to follow RESTful conventions except as noted.

All responses are in JSON format, wrapped in a JSON envelope of the following form.

```
{
    "response": <response from method, typically a list or object>
    "meta": {
        "count": <number of results in `response`>
    }
    "error": { // Or null, if no errors
        "message": "Error message"
    }
}
```

### `GET` convos/

Retrieves all parent conversations

### `POST` convos/

Create a new conversation

### `GET` convos/:id/

Retrieves an individual conversation

### `PATCH` convos/:id/

Edits an existing conversation. Will only change the fields that are passed to it. Only certain fields can be edited:

Body?

### `DELETE` convos/:id/

Deletes an individual conversation and all conversations that have this conversation as a parent.

### `POST` convos/:id/reply/

Create a reply to an individual conversation

## Notes

TODO