import { API_URL } from "./App";
import { useNavigate, useParams } from "react-router-dom";

export default function ConfirmationPage() {
    const { token = '' }  = useParams();
    const redirect = useNavigate();

    const handleConfirm = async () => {
        const response = await fetch(`${API_URL}/users/activate/${token}`, {
            method: 'PUT'
        });

        if (response.ok) {
            // redirect to login page
            redirect('/');
        } else {
            // show error message
            alert('Failed to confirm email');
        }

    }

  return (
    <div>
      <h1>Confirmation Page</h1>
      <button onClick={handleConfirm}>Confirm</button>
    </div>
  );
}