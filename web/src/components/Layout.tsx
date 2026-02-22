import { Link, Outlet, useNavigate } from "react-router-dom";
import { useAuth } from "../context/AuthContext";

export function Layout() {
  const { user, signOut } = useAuth();
  const navigate = useNavigate();

  const handleSignOut = () => {
    signOut();
    navigate("/login");
  };

  return (
    <div className="layout">
      <nav className="navbar">
        <div className="navbar-inner">
          <Link to="/" className="navbar-brand">
            GoSocial
          </Link>
          <div className="navbar-nav">
            {user ? (
              <>
                <Link to="/posts/new" className="btn btn-primary" style={{ fontSize: 13 }}>
                  New Post
                </Link>
                <Link
                  to={`/users/${user.id}`}
                  className="navbar-user"
                >
                  @{user.username}
                </Link>
                <button className="btn btn-secondary" onClick={handleSignOut}>
                  Sign Out
                </button>
              </>
            ) : (
              <>
                <Link to="/login" className="btn btn-secondary">
                  Sign In
                </Link>
                <Link to="/register" className="btn btn-primary">
                  Register
                </Link>
              </>
            )}
          </div>
        </div>
      </nav>
      <main className="main-content">
        <Outlet />
      </main>
    </div>
  );
}
