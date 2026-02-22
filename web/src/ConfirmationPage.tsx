import { useState } from "react";
import { useNavigate, useParams } from "react-router-dom";
import { API_URL } from "./config";

export const ConfirmationPage = () => {
  const { token } = useParams<{ token: string }>();
  const redirect = useNavigate();
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  const handleConfirm = async () => {
    setError(null);
    setLoading(true);
    try {
      const response = await fetch(`${API_URL}/users/activate/${token}`, {
        method: "PUT",
      });

      if (response.ok) {
        redirect("/");
      } else {
        setError("An error occurred while confirming your account.");
      }
    } catch {
      setError("An error occurred while confirming your account.");
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="auth-page">
      <div className="auth-card">
        <h1>Confirm your account</h1>
        <p>Thank you for your registration! Click below to activate your account.</p>

        {error && <div className="error-message">{error}</div>}

        <button
          className="btn btn-primary"
          onClick={handleConfirm}
          disabled={loading}
        >
          {loading ? "Confirming…" : "Confirm"}
        </button>
      </div>
    </div>
  );
};
