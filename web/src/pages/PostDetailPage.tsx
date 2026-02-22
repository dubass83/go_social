import { useEffect, useState, type FormEvent } from "react";
import { Link, useNavigate, useParams } from "react-router-dom";
import { createComment, deletePost, getPost } from "../api";
import { useAuth } from "../context/AuthContext";
import type { Post } from "../types";

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
    day: "numeric",
    hour: "2-digit",
    minute: "2-digit",
  });
}

export function PostDetailPage() {
  const { postID } = useParams<{ postID: string }>();
  const { user } = useAuth();
  const navigate = useNavigate();

  const [post, setPost] = useState<Post | null>(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [commentText, setCommentText] = useState("");
  const [commentLoading, setCommentLoading] = useState(false);
  const [commentError, setCommentError] = useState<string | null>(null);

  const [deleting, setDeleting] = useState(false);

  useEffect(() => {
    if (!postID) return;
    setLoading(true);
    getPost(parseInt(postID, 10))
      .then(setPost)
      .catch((err) =>
        setError(err instanceof Error ? err.message : "Failed to load post")
      )
      .finally(() => setLoading(false));
  }, [postID]);

  const handleDelete = async () => {
    if (!post || !confirm("Delete this post?")) return;
    setDeleting(true);
    try {
      await deletePost(post.id);
      navigate("/");
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to delete post");
      setDeleting(false);
    }
  };

  const handleComment = async (e: FormEvent) => {
    e.preventDefault();
    if (!post || !user || !commentText.trim()) return;
    setCommentLoading(true);
    setCommentError(null);
    try {
      const comment = await createComment(post.id, user.id, commentText.trim());
      setPost((prev) =>
        prev
          ? { ...prev, comments: [...(prev.comments ?? []), comment] }
          : prev
      );
      setCommentText("");
    } catch (err) {
      setCommentError(
        err instanceof Error ? err.message : "Failed to post comment"
      );
    } finally {
      setCommentLoading(false);
    }
  };

  if (loading) return <div className="loading-spinner">Loading…</div>;
  if (error) return <div className="error-message">{error}</div>;
  if (!post) return null;

  const isOwner = user?.id === post.user_id;

  return (
    <>
      <div className="post-detail">
        <div className="post-detail-header">
          <h1>{post.title}</h1>
          <div className="post-detail-meta">
            <Link to={`/users/${post.user_id}`} style={{ fontWeight: 600 }}>
              @{post.user?.username ?? "unknown"}
            </Link>
            <span>{formatDate(post.created_at)}</span>
            {post.updated_at !== post.created_at && (
              <span>Edited {formatDate(post.updated_at)}</span>
            )}
          </div>
          {post.tags && post.tags.length > 0 && (
            <div className="tags-list" style={{ marginTop: 10 }}>
              {post.tags.map((tag) => (
                <span key={tag} className="tag">
                  {tag}
                </span>
              ))}
            </div>
          )}
        </div>

        <div className="post-detail-content">{post.content}</div>

        {isOwner && (
          <div className="post-detail-actions">
            <Link
              to={`/posts/${post.id}/edit`}
              className="btn btn-secondary"
            >
              Edit
            </Link>
            <button
              className="btn btn-danger"
              onClick={handleDelete}
              disabled={deleting}
            >
              {deleting ? "Deleting…" : "Delete"}
            </button>
          </div>
        )}
      </div>

      <div className="comments-section">
        <h3>
          {(post.comments?.length ?? 0) === 0
            ? "No comments yet"
            : `${post.comments!.length} Comment${post.comments!.length !== 1 ? "s" : ""}`}
        </h3>

        {(post.comments ?? []).map((c) => (
          <div key={c.id} className="comment">
            <div className="comment-header">
              <Link to={`/users/${c.user_id}`} className="comment-author">
                @{c.user?.username ?? "unknown"}
              </Link>
              <span className="comment-date">{formatDate(c.created_at)}</span>
            </div>
            <div className="comment-content">{c.content}</div>
          </div>
        ))}

        {user && (
          <form className="comment-form" onSubmit={handleComment}>
            {commentError && (
              <div className="error-message">{commentError}</div>
            )}
            <div className="form-group">
              <label htmlFor="comment">Add a comment</label>
              <textarea
                id="comment"
                className="form-control"
                rows={3}
                value={commentText}
                onChange={(e) => setCommentText(e.target.value)}
                placeholder="Write your comment…"
                required
              />
            </div>
            <button
              type="submit"
              className="btn btn-primary"
              disabled={commentLoading || !commentText.trim()}
            >
              {commentLoading ? "Posting…" : "Post Comment"}
            </button>
          </form>
        )}
      </div>
    </>
  );
}
