# Missing API Endpoints

These endpoints are needed by the Web UI but not yet implemented in the Go backend.
Each description includes the expected request/response shape so you can implement them as a study exercise.

---

## 1. `GET /v1/users/{userID}/posts`

**Purpose:** Fetch all posts authored by a specific user, for display on their profile page.

**Auth:** Bearer token (required)

**Path parameters:**
- `userID` – integer user ID

**Query parameters:**
| Parameter | Type   | Default | Description                              |
|-----------|--------|---------|------------------------------------------|
| `limit`   | int    | 10      | Number of posts to return (1–100)        |
| `offset`  | int    | 0       | Pagination offset                        |
| `sort`    | string | `desc`  | Sort by `created_at`: `asc` or `desc`    |

**Success response `200 OK`:**
```json
{
  "data": [
    {
      "id": 1,
      "title": "Hello World",
      "content": "Post body...",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z",
      "user_id": 42,
      "version": 1,
      "tags": ["go", "programming"],
      "comments_count": 3,
      "user": { "id": 42, "username": "alice", ... }
    }
  ]
}
```

**Where it is used:** `UserProfilePage` — to list a user's posts beneath their profile info.

**Implementation hints:**
- Reuse the `PostStore.GetPostsByUserID` (or similar) query.
- Return `PostWithMetadata` (include `comments_count`).
- Apply the same pagination validation used in `GetUserFeed`.
- The user must exist; return `404` if not.

---

## 2. `GET /v1/users/me`

**Purpose:** Return the currently authenticated user's profile without requiring the client to decode the JWT and make a second request.

**Auth:** Bearer token (required)

**No path or query parameters.**

**Success response `200 OK`:**
```json
{
  "data": {
    "id": 42,
    "username": "alice",
    "email": "alice@example.com",
    "created_at": "2024-01-01T00:00:00Z",
    "active": true,
    "role_id": 1
  }
}
```

**Where it is used:** Navbar / AuthContext — to display the signed-in user's username and link to their profile.

**Implementation hints:**
- Extract the user from the request context (already done by `AuthTokenMiddleware`).
- This is a one-liner handler; no database query needed if the user object is already in the context.
- Register the route **before** `/v1/users/{userID}` to avoid the `me` string being treated as a userID.

---

## 3. `GET /v1/posts` — Public / discovery feed

**Purpose:** Return a paginated list of all posts (not filtered by follows), useful as a public discovery feed for new users who do not follow anyone yet.

**Auth:** Optional (show public posts without auth; show personalized hints when authenticated)

**Query parameters:** Same as `GET /v1/users/feed` (`limit`, `offset`, `sort`, `tags`, `search`).

**Success response `200 OK`:**
```json
{
  "data": [
    { /* PostWithMetadata */ }
  ]
}
```

**Where it is used:** `FeedPage` — show as a fallback when the user's personal feed is empty (no follows yet).

**Implementation hints:**
- Add a `GetAllPosts(ctx, PaginatedFeedQuery) ([]PostWithMetadata, error)` method to `PostStore`.
- The SQL query is similar to `GetUserFeed` but without the `followers` JOIN filter.
- This endpoint can be public (`/v1/posts`) or authenticated — your choice.
