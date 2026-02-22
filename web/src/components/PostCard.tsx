import { Link } from "react-router-dom";
import type { PostWithMetadata } from "../types";

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString("en-US", {
    year: "numeric",
    month: "short",
    day: "numeric",
  });
}

interface PostCardProps {
  post: PostWithMetadata;
}

export function PostCard({ post }: PostCardProps) {
  const preview =
    post.content.length > 200
      ? post.content.slice(0, 200) + "…"
      : post.content;

  return (
    <article className="post-card">
      <div className="post-card-header">
        <span className="post-card-meta">
          <Link
            to={`/users/${post.user_id}`}
            className="post-card-author"
          >
            @{post.user?.username ?? "unknown"}
          </Link>
          <span style={{ margin: "0 6px" }}>·</span>
          {formatDate(post.created_at)}
        </span>
        {post.comments_count > 0 && (
          <span className="comment-count">
            {post.comments_count} comment{post.comments_count !== 1 ? "s" : ""}
          </span>
        )}
      </div>

      <h2>
        <Link to={`/posts/${post.id}`}>{post.title}</Link>
      </h2>

      <p className="post-card-content">{preview}</p>

      {post.tags && post.tags.length > 0 && (
        <div className="tags-list post-card-footer">
          {post.tags.map((tag) => (
            <span key={tag} className="tag">
              {tag}
            </span>
          ))}
        </div>
      )}
    </article>
  );
}
