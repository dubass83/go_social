import { useEffect, useState, type FormEvent } from "react";
import { Link } from "react-router-dom";
import { getFeed } from "../api";
import { PostCard } from "../components/PostCard";
import type { PostWithMetadata } from "../types";

const PAGE_SIZE = 10;

export function FeedPage() {
  const [posts, setPosts] = useState<PostWithMetadata[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);

  const [search, setSearch] = useState("");
  const [tags, setTags] = useState("");
  const [searchInput, setSearchInput] = useState("");
  const [tagsInput, setTagsInput] = useState("");

  const loadFeed = async (
    newOffset: number,
    searchVal: string,
    tagsVal: string
  ) => {
    setLoading(true);
    setError(null);
    try {
      const data = await getFeed({
        limit: PAGE_SIZE,
        offset: newOffset,
        sort: "desc",
        search: searchVal || undefined,
        tags: tagsVal || undefined,
      });
      setPosts(data);
      setHasMore(data.length === PAGE_SIZE);
    } catch (err) {
      setError(err instanceof Error ? err.message : "Failed to load feed");
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    loadFeed(0, "", "");
  }, []);

  const handleSearch = (e: FormEvent) => {
    e.preventDefault();
    setSearch(searchInput);
    setTags(tagsInput);
    setOffset(0);
    loadFeed(0, searchInput, tagsInput);
  };

  const handleReset = () => {
    setSearchInput("");
    setTagsInput("");
    setSearch("");
    setTags("");
    setOffset(0);
    loadFeed(0, "", "");
  };

  const handlePrev = () => {
    const newOffset = Math.max(0, offset - PAGE_SIZE);
    setOffset(newOffset);
    loadFeed(newOffset, search, tags);
  };

  const handleNext = () => {
    const newOffset = offset + PAGE_SIZE;
    setOffset(newOffset);
    loadFeed(newOffset, search, tags);
  };

  return (
    <>
      <div className="feed-header">
        <h1>Your Feed</h1>
        <Link to="/posts/new" className="btn btn-primary">
          + New Post
        </Link>
      </div>

      <div className="search-bar">
        <form
          onSubmit={handleSearch}
          style={{ display: "contents" }}
        >
          <div className="form-group">
            <label htmlFor="search">Search</label>
            <input
              id="search"
              type="text"
              className="form-control"
              placeholder="Search posts…"
              value={searchInput}
              onChange={(e) => setSearchInput(e.target.value)}
            />
          </div>
          <div className="form-group">
            <label htmlFor="tags">Tags</label>
            <input
              id="tags"
              type="text"
              className="form-control"
              placeholder="e.g. go,react"
              value={tagsInput}
              onChange={(e) => setTagsInput(e.target.value)}
            />
          </div>
          <button type="submit" className="btn btn-primary">
            Search
          </button>
          {(search || tags) && (
            <button
              type="button"
              className="btn btn-secondary"
              onClick={handleReset}
            >
              Clear
            </button>
          )}
        </form>
      </div>

      {loading && <div className="loading-spinner">Loading…</div>}

      {error && <div className="error-message">{error}</div>}

      {!loading && !error && posts.length === 0 && (
        <div className="empty-state">
          <h3>No posts yet</h3>
          <p>
            {search || tags
              ? "No posts match your search."
              : "Follow some users to see their posts here, or create your own."}
          </p>
        </div>
      )}

      {posts.map((post) => (
        <PostCard key={post.id} post={post} />
      ))}

      {!loading && posts.length > 0 && (
        <div className="pagination">
          <button
            className="btn btn-secondary"
            onClick={handlePrev}
            disabled={offset === 0}
          >
            &larr; Previous
          </button>
          <span className="text-secondary">
            Posts {offset + 1}–{offset + posts.length}
          </span>
          <button
            className="btn btn-secondary"
            onClick={handleNext}
            disabled={!hasMore}
          >
            Next &rarr;
          </button>
        </div>
      )}
    </>
  );
}
