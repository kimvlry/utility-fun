# Simple Calendar HTTP Server

A minimal HTTP server for managing a small in-memory calendar of events.

## Features

### Supported Operations

| Method | Endpoint  | Description          | Parameters (JSON body for POST / Query string for GET)                                         |
| ------ | --------- | -------------------- | ---------------------------------------------------------------------------------------------- |
| POST   | `/events` | Create a new event   | `user_id` (int), `date` (YYYY-MM-DD), `event` (string)                                         |
| PUT    | `/events` | Update an event      | `user_id` (int), `date` (YYYY-MM-DD), `event` (string)                                         |
| DELETE | `/events` | Delete an event      | `user_id` (int), `date` (YYYY-MM-DD), `event_id` (int)                                         |
| GET    | `/events` | Get events by period | `user_id` (int), `date` (YYYY-MM-DD), `period` (string: `day`, `week`, `month`, default=`day`) |


### Request Format

- All `POST` endpoints expect data in the request body as **JSON**:
    - `Content-Type: application/json`
    - Example:

```json
{
  "user_id": 1,
  "date": "2025-07-16",
  "event": "Attend Go workshop"
}
```

### Response Format

- On successful execution, the server responds with JSON:
  ```json
  {"result": "..."}

- On business logic errors, the response is JSON:
  ```json
  {"error": "error description"}
  ```

| Status Code               | Description                                              |
|---------------------------|----------------------------------------------------------|
| 200 OK                    | Request was successful                                   |
| 400 Bad Request           | Input validation error (e.g., invalid date format)       |
| 503 Service Unavailable   | Business logic error (e.g., deleting non-existent event) |
| 500 Internal Server Error | Other unexpected errors                                  |

## Implementation

### Design

- Events are stored in memory using Go data structures.
- `user_id` represents the calendar user's identifier. Complex access control is not required for this project.
- An event is defined as a record containing an ID, date (YYYY-MM-DD), time (HH:MM) and a text description.
- A **middleware** component logs every HTTP request, including:
    - HTTP method
    - URL
    - Timestamp
- Logs are output to stdout or written to a file.
- The server listens on a port specified in configuration (via an environment variable).
- Business logic is separated from the HTTP layer. HTTP handlers only call methods from the business logic layer.

### Tests

- Passes `go vet` and `golint` checks, and is free from data races.
- Core business logic functions are covered by unit tests.
