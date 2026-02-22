import { useEffect, useState, type FormEvent } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { createPost, getPost, updatePost } from "../api";

export function CreatePostPage() {
  const { postID } = useParams<{ postID: string }>();
  const isEdit = Boolean(postID);
  const navigate = useNavigate();

  const [title, setTitle] = useState("");
  const [content, setContent] = useState("");
  const [tagsInput, setTagsInput] = useState("");

  const [loading, setLoading] = useState(false);
  const [fetchLoading, setFetchLoading] = useState(isEdit);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    if (!isEdit || !postID) return;
    setFetchLoading(true);
    getPost(parseInt(postID, 10))
      .then((post) => {
        setTitle(post.title);
        setContent(post.content);
        setTagsInput((post.tags ?? []).join(", "));
      })
      .catch((err) =>
        setError(err instanceof Error ? err.message : "Failed to load post")
      )
      .finally(() => setFetchLoading(false));
  }, [isEdit, postID]);

  const parseTags = (raw: string): string[] =>
    raw
      .split(",")
      .map((t) => t.trim())
      .filter(Boolean);

  const handleSubmit = async (e: FormEvent) => {
    e.preventDefault();
    setError(null);
    setLoading(true);

    const tags = parseTags(tagsInput);

    try {
      if (isEdit && postID) {
        const updated = await updatePost(parseInt(postID, 10), {
          title,
          content,
          tags,
        });
        navigate(`/posts/${updated.id}`);
      } else {
        const created = await createPost(title, content, tags);
        navigate(`/posts/${created.id}`);
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to save post");
    } finally {
      setLoading(false);
    }
  };

  if (fetchLoading) return <div className="loading-spinner">Loading…</div>;

  return (
    <div className="create-post">
      <h1>{isEdit ? "Edit Post" : "New Post"}</h1>

      <form className="create-post-form" onSubmit={handleSubmit}>
        {error && <div className="error-message">{error}</div>}

        <div className="form-group">
          <label htmlFor="title">Title</label>
          <input
            id="title"
            type="text"
            className="form-control"
            value={title}
            onChange={(e) => setTitle(e.target.value)}
            required
            minLength={2}
            maxLength={100}
            placeholder="Post title"
          />
        </div>

        <div className="form-group">
          <label htmlFor="content">Content</label>
          <textarea
            id="content"
            className="form-control"
            rows={10}
            value={content}
            onChange={(e) => setContent(e.target.value)}
            required
            minLength={2}
            maxLength={1000}
            placeholder="What's on your mind?"
          />
          <span className="text-secondary">{content.length}/1000</span>
        </div>

        <div className="form-group">
          <label htmlFor="tags">Tags</label>
          <input
            id="tags"
            type="text"
            className="form-control"
            value={tagsInput}
            onChange={(e) => setTagsInput(e.target.value)}
            placeholder="go, programming, news"
          />
          <span className="text-secondary">Comma-separated</span>
        </div>

        <div className="create-post-actions">
          <button
            type="submit"
            className="btn btn-primary"
            disabled={loading}
          >
            {loading
              ? isEdit
                ? "Saving…"
                : "Publishing…"
              : isEdit
              ? "Save Changes"
              : "Publish Post"}
          </button>
          <button
            type="button"
            className="btn btn-secondary"
            onClick={() => navigate(-1)}
            disabled={loading}
          >
            Cancel
          </button>
        </div>
      </form>
    </div>
  );
}
