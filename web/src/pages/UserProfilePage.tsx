import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { followUser, getUser, getUserPosts, unfollowUser } from "../api";
import { PostCard } from "../components/PostCard";
import { useAuth } from "../context/AuthContext";
import type { PostWithMetadata, User } from "../types";

const PAGE_SIZE = 10;

function formatDate(iso: string) {
  return new Date(iso).toLocaleDateString("en-US", {
    year: "numeric",
    month: "long",
  });
}

export function UserProfilePage() {
  const { userID } = useParams<{ userID: string }>();
  const { user: currentUser } = useAuth();

  const [profile, setProfile] = useState<User | null>(null);
  const [profileLoading, setProfileLoading] = useState(true);
  const [profileError, setProfileError] = useState<string | null>(null);

  const [posts, setPosts] = useState<PostWithMetadata[]>([]);
  const [postsLoading, setPostsLoading] = useState(true);
  const [postsError, setPostsError] = useState<string | null>(null);
  const [offset, setOffset] = useState(0);
  const [hasMore, setHasMore] = useState(true);

  const [following, setFollowing] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);
  const [followError, setFollowError] = useState<string | null>(null);

  const parsedID = parseInt(userID ?? "0", 10);

  useEffect(() => {
    if (!userID) return;
    setProfileLoading(true);
    setProfileError(null);
    getUser(parsedID)
      .then(setProfile)
      .catch((err) =>
        setProfileError(err instanceof Error ? err.message : "User not found")
      )
      .finally(() => setProfileLoading(false));
  }, [userID]);

  const loadPosts = (newOffset: number) => {
    if (!userID) return;
    setPostsLoading(true);
    setPostsError(null);
    getUserPosts(parsedID, { limit: PAGE_SIZE, offset: newOffset, sort: "desc" })
      .then((data) => {
        setPosts(data);
        setHasMore(data.length === PAGE_SIZE);
      })
      .catch((err) =>
        setPostsError(err instanceof Error ? err.message : "Failed to load posts")
      )
      .finally(() => setPostsLoading(false));
  };

  useEffect(() => {
    loadPosts(0);
  }, [userID]);

  const isOwnProfile = currentUser?.id === parsedID;

  const handleFollowToggle = async () => {
    if (!profile) return;
    setFollowLoading(true);
    setFollowError(null);
    try {
      if (following) {
        await unfollowUser(profile.id);
        setFollowing(false);
      } else {
        await followUser(profile.id);
        setFollowing(true);
      }
    } catch (err) {
      setFollowError(err instanceof Error ? err.message : "Action failed");
    } finally {
      setFollowLoading(false);
    }
  };

  const handlePrev = () => {
    const newOffset = Math.max(0, offset - PAGE_SIZE);
    setOffset(newOffset);
    loadPosts(newOffset);
  };

  const handleNext = () => {
    const newOffset = offset + PAGE_SIZE;
    setOffset(newOffset);
    loadPosts(newOffset);
  };

  if (profileLoading) return <div className="loading-spinner">Loading…</div>;
  if (profileError) return <div className="error-message">{profileError}</div>;
  if (!profile) return null;

  const initial = profile.username.charAt(0).toUpperCase();

  return (
    <>
      <div className="user-profile">
        <div className="user-profile-header">
          <div>
            <div className="user-avatar">{initial}</div>
            <h1>{profile.username}</h1>
            <div className="user-email">{profile.email}</div>
            <div className="text-secondary" style={{ marginTop: 6 }}>
              Joined {formatDate(profile.created_at)}
            </div>
            {!profile.active && (
              <div style={{ marginTop: 6, fontSize: 13, color: "#f59e0b" }}>
                Account not yet activated
              </div>
            )}
          </div>

          {!isOwnProfile && currentUser && (
            <button
              className={following ? "btn btn-secondary" : "btn btn-primary"}
              onClick={handleFollowToggle}
              disabled={followLoading}
            >
              {followLoading ? "…" : following ? "Unfollow" : "Follow"}
            </button>
          )}
        </div>

        {followError && (
          <div className="error-message" style={{ marginTop: 12 }}>
            {followError}
          </div>
        )}
      </div>

      <div style={{ marginTop: 16 }}>
        <h2 style={{ fontSize: 17, fontWeight: 700, marginBottom: 12 }}>
          Posts by {profile.username}
        </h2>

        {postsLoading && <div className="loading-spinner">Loading posts…</div>}
        {postsError && <div className="error-message">{postsError}</div>}

        {!postsLoading && !postsError && posts.length === 0 && (
          <div className="empty-state">
            <h3>No posts yet</h3>
            <p>{profile.username} hasn't published anything.</p>
          </div>
        )}

        {posts.map((post) => (
          <PostCard key={post.id} post={post} />
        ))}

        {!postsLoading && posts.length > 0 && (
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
      </div>
    </>
  );
}
