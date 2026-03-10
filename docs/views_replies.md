# Replies and Views Flows

Based on the analysis of the codebase, here is the breakdown of the **Replies** and **Views** flows.

## Current Status Overview
The project has a solid foundation with the **Database Layer** (SQL schemas and queries) and **Data Access Layer** (generated `sqlc` code) already in place for both features. However, the **Application Layer** (Service logic) and **Transport Layer** (HTTP Handlers) are currently missing these implementations.

---

## 1. Replies Flow

### Concept
A user posts a text reply to a specific Locus (location-based message). This is a one-to-many relationship (One Locus -> Many Replies). When a reply is created, we also want to increment a counter on the main Locus to show activity.

### Architecture & Data Flow
1.  **Client:** Sends a `POST` request with `locus_id` and `content`.
2.  **HTTP Handler (`internal/transport/http/handler.go`):**
    - Receives the request.
    - Extracts `userID` from the context (via middleware).
    - Validates the input (e.g., content is not empty).
    - Calls the Service layer.
3.  **Service Layer (`internal/service/service.go`):**
    - Calls `CreateReply` in the database to save the new reply.
    - Calls `IncrementReplyCount` to update the `loci` table's `replies_count`.
    - *(Optional)* Broadcasts the new reply via the WebSocket `hub` so other users see it instantly.
4.  **Database:**
    - Inserts row into `replies` table.
    - Updates row in `loci` table.

### What is Needed
1.  **Service Layer:** Add a `CreateReply` method in `UserService`.
    - It needs to execute the `CreateReply` query.
    - It needs to execute the `IncrementReplyCount` query.
2.  **HTTP Layer:** Add a `ReplyToLocus` handler in `UserHandler`.
    - Define a request struct (e.g., `type ReplyReq struct { Content string }`).
    - Parse the `locus_id` from the URL parameters.
3.  **Routing:** Register the new endpoint (e.g., `POST /api/loci/{id}/reply`) in your router.

---

## 2. Views Flow

### Concept
We want to track when a user "views" a Locus. To prevent spamming the view count, we track unique views using a composite key of `(user_id, locus_id)`.

### Architecture & Data Flow
1.  **Client:** Sends a request indicating a Locus was viewed (or this happens automatically when fetching details).
2.  **HTTP Handler:**
    - Extracts `userID` and `locus_id`.
    - Calls the Service layer.
3.  **Service Layer:**
    - **Step 1 (Idempotency):** Try to insert a record into the `locus_views` table.
        - *Note:* Since the primary key is `(user_id, locus_id)`, the database will prevent duplicate inserts.
    - **Step 2 (Counter):** If the insert was successful (meaning this is a *new* view for this user), call `IncrementViewCount` to update the `loci` table.
4.  **Database:**
    - Inserts row into `locus_views`.
    - Updates `view_count` in `loci`.

### What is Needed
1.  **Fix SQL Query:** The file `internal/store/queries/create_locus_view.sql` is currently commented out. You need to uncomment the `INSERT` statement so `sqlc` can generate the Go code for it.
2.  **Service Layer:** Add a `RegisterView` method in `UserService`.
    - It should handle the logic: "Try to record view -> If successful, increment count".
3.  **HTTP Layer:** Add a `ViewLocus` handler in `UserHandler`.
4.  **Routing:** Register the endpoint (e.g., `POST /api/loci/{id}/view`).

---

## Summary of Work Required

You do not need to create new files, but you need to add code to these existing files:

1.  **`internal/store/queries/create_locus_view.sql`**: Uncomment the SQL query.
2.  **`internal/service/service.go`**: Add `CreateReply` and `RegisterView` methods.
3.  **`internal/transport/http/handler.go`**: Add `ReplyToLocus` and `ViewLocus` handler functions.
4.  **`cmd/main.go`** (or wherever your router is): Connect the new handlers to specific URL paths.
