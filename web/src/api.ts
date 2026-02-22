import { API_URL } from "./config";
import type { Comment, FeedParams, Post, PostWithMetadata, User } from "./types";

function requestHeaders(withBody = false): Record<string, string> {
  const headers: Record<string, string> = {};
  if (withBody) headers["Content-Type"] = "application/json";
  const token = localStorage.getItem("token");
  if (token) headers["Authorization"] = `Bearer ${token}`;
  return headers;
}

async function handleResponse<T>(res: Response): Promise<T> {
  const body = await res.json();
  if (!res.ok) {
    throw new Error(body.error || `HTTP ${res.status}`);
  }
  return body.data as T;
}

// --- Auth ---

export async function login(email: string, password: string): Promise<string> {
  const res = await fetch(`${API_URL}/authentication/token`, {
    method: "POST",
    headers: requestHeaders(true),
    body: JSON.stringify({ email, password }),
  });
  return handleResponse<string>(res);
}

export async function register(
  username: string,
  email: string,
  password: string
): Promise<User> {
  const res = await fetch(`${API_URL}/authentication/user`, {
    method: "POST",
    headers: requestHeaders(true),
    body: JSON.stringify({ username, email, password }),
  });
  return handleResponse<User>(res);
}

// --- Users ---

export async function getUser(userID: number): Promise<User> {
  const res = await fetch(`${API_URL}/users/${userID}`, {
    headers: requestHeaders(),
  });
  return handleResponse<User>(res);
}

export async function followUser(userID: number): Promise<void> {
  const res = await fetch(`${API_URL}/users/${userID}/follow`, {
    method: "PUT",
    headers: requestHeaders(),
  });
  return handleResponse<void>(res);
}

export async function unfollowUser(userID: number): Promise<void> {
  const res = await fetch(`${API_URL}/users/${userID}/unfollow`, {
    method: "PUT",
    headers: requestHeaders(),
  });
  return handleResponse<void>(res);
}

// --- Feed ---

export async function getFeed(
  params: FeedParams = {}
): Promise<PostWithMetadata[]> {
  const query = new URLSearchParams();
  if (params.limit != null) query.set("limit", String(params.limit));
  if (params.offset != null) query.set("offset", String(params.offset));
  if (params.sort) query.set("sort", params.sort);
  if (params.tags) query.set("tags", params.tags);
  if (params.search) query.set("search", params.search);

  const res = await fetch(`${API_URL}/users/feed?${query}`, {
    headers: requestHeaders(),
  });
  const data = await handleResponse<PostWithMetadata[] | null>(res);
  return data ?? [];
}

// --- Posts ---

export async function getPost(postID: number): Promise<Post> {
  const res = await fetch(`${API_URL}/posts/${postID}`, {
    headers: requestHeaders(),
  });
  return handleResponse<Post>(res);
}

export async function createPost(
  title: string,
  content: string,
  tags: string[]
): Promise<Post> {
  const res = await fetch(`${API_URL}/posts`, {
    method: "POST",
    headers: requestHeaders(true),
    body: JSON.stringify({ title, content, tags }),
  });
  return handleResponse<Post>(res);
}

export async function updatePost(
  postID: number,
  updates: { title?: string; content?: string; tags?: string[] }
): Promise<Post> {
  const res = await fetch(`${API_URL}/posts/${postID}`, {
    method: "PATCH",
    headers: requestHeaders(true),
    body: JSON.stringify(updates),
  });
  return handleResponse<Post>(res);
}

export async function deletePost(postID: number): Promise<void> {
  const res = await fetch(`${API_URL}/posts/${postID}`, {
    method: "DELETE",
    headers: requestHeaders(),
  });
  return handleResponse<void>(res);
}

// --- Comments ---

export async function createComment(
  postID: number,
  userID: number,
  content: string
): Promise<Comment> {
  const res = await fetch(`${API_URL}/posts/${postID}/comments`, {
    method: "POST",
    headers: requestHeaders(true),
    body: JSON.stringify({ user_id: userID, content }),
  });
  return handleResponse<Comment>(res);
}
