export interface User {
  id: number;
  username: string;
  email: string;
  created_at: string;
  active: boolean;
  role_id: number;
}

export interface Comment {
  id: number;
  post_id: number;
  user_id: number;
  content: string;
  created_at: string;
  user: User;
}

export interface Post {
  id: number;
  title: string;
  content: string;
  created_at: string;
  updated_at: string;
  user_id: number;
  version: number;
  tags: string[] | null;
  comments: Comment[] | null;
  user: User;
}

export interface PostWithMetadata extends Post {
  comments_count: number;
}

export interface FeedParams {
  limit?: number;
  offset?: number;
  sort?: "asc" | "desc";
  tags?: string;
  search?: string;
}
