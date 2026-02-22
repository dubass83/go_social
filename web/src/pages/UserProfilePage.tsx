import { useEffect, useState } from "react";
import { useParams } from "react-router-dom";
import { followUser, getUser, unfollowUser } from "../api";
import { useAuth } from "../context/AuthContext";
import type { User } from "../types";

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
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);

  const [following, setFollowing] = useState(false);
  const [followLoading, setFollowLoading] = useState(false);
  const [followError, setFollowError] = useState<string | null>(null);

  useEffect(() => {
    if (!userID) return;
    setLoading(true);
    setFollowError(null);
    getUser(parseInt(userID, 10))
      .then(setProfile)
      .catch((err) =>
        setError(err instanceof Error ? err.message : "User not found")
      )
      .finally(() => setLoading(false));
  }, [userID]);

  const isOwnProfile = currentUser?.id === parseInt(userID ?? "0", 10);

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
      setFollowError(
        err instanceof Error ? err.message : "Action failed"
      );
    } finally {
      setFollowLoading(false);
    }
  };

  if (loading) return <div className="loading-spinner">Loading…</div>;
  if (error) return <div className="error-message">{error}</div>;
  if (!profile) return null;

  const initial = profile.username.charAt(0).toUpperCase();

  return (
    <div className="user-profile">
      <div className="user-profile-header">
        <div>
          <div className="user-avatar">{initial}</div>
          <h1>{profile.username}</h1>
          <div className="user-email">{profile.email}</div>
          <div
            className="text-secondary"
            style={{ marginTop: 6, fontSize: 13 }}
          >
            Joined {formatDate(profile.created_at)}
          </div>
          {!profile.active && (
            <div
              className="text-secondary"
              style={{ marginTop: 6, fontSize: 13, color: "#f59e0b" }}
            >
              Account not yet activated
            </div>
          )}
        </div>

        {!isOwnProfile && currentUser && (
          <div>
            <button
              className={following ? "btn btn-secondary" : "btn btn-primary"}
              onClick={handleFollowToggle}
              disabled={followLoading}
            >
              {followLoading
                ? "…"
                : following
                ? "Unfollow"
                : "Follow"}
            </button>
          </div>
        )}
      </div>

      {followError && (
        <div className="error-message" style={{ marginTop: 12 }}>
          {followError}
        </div>
      )}

      <div
        style={{
          marginTop: 24,
          padding: "20px",
          background: "#f8fafc",
          borderRadius: 8,
          border: "1px dashed #cbd5e1",
          color: "#64748b",
          fontSize: 14,
        }}
      >
        <strong>Posts by this user are not available yet.</strong>
        <br />
        The API endpoint <code>GET /v1/users/{"{userID}"}/posts</code> is not
        implemented. See <code>MISSING_APIS.md</code> for details.
      </div>
    </div>
  );
}
