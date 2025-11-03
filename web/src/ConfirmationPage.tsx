import { useNavigate, useParams } from "react-router-dom";
import { API_URL } from "./config";

export const ConfirmationPage = () => {
  const { token } = useParams<{ token: string }>();
  const redirect = useNavigate();
  const handleConfirm = async () => {
    const response = await fetch(`${API_URL}/users/activate/${token}`, {
      method: "PUT",
    });

    if (response.ok) {
      // Handle successful activation
      redirect("/");
    } else {
      // Handle error
      alert("An error occurred while confirming your account.");
    }
  };

  return (
    <div>
      <h1>Confirmation</h1>
      <p>Thank you for your registration!</p>
      <button onClick={handleConfirm}>Confirm</button>
    </div>
  );
};
