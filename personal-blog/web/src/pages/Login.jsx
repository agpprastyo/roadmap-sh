// src/pages/Login.jsx

import  { useState } from 'react';
import { useAuth } from '../contexts/AuthContext'; // Import Auth context
import { useNavigate } from 'react-router-dom'; // For navigation

const Login = () => {
    const { login } = useAuth(); // Use the login function from context
    const navigate = useNavigate();
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [error, setError] = useState(null);

    const handleSubmit = async (e) => {
        e.preventDefault();

        try {
            const response = await fetch('http://localhost:4444/api/v1/sign-in', {
                method: 'POST',
                headers: {
                    'Authorization': 'Basic ' + btoa(`${username}:${password}`),
                    'Content-Type': 'application/json'
                },
                credentials: 'include', // Include cookies in requests
            });

            if (response.ok) {
                const data = await response.json();
                if (data.message === "Sign in successful") {
                    // Call the login function to update the context
                    login();
                    navigate('/admin'); // Redirect to admin page
                }
            } else {
                throw new Error('Failed to log in');
            }
        } catch (err) {
            setError(err.message);
        }
    };

    return (
        <div className="flex items-center justify-center min-h-screen bg-gray-100">
            <form onSubmit={handleSubmit} className="bg-white p-6 rounded shadow-md w-96">
                <h2 className="text-lg font-bold mb-4">Login</h2>

                {error && <p className="text-red-500 text-xs mb-4">{error}</p>}

                <div className="mb-4">
                    <label className="block text-gray-700 mb-2">Username</label>
                    <input
                        type="text"
                        value={username}
                        onChange={(e) => setUsername(e.target.value)}
                        className="border border-gray-300 p-2 w-full rounded"
                        required
                    />
                </div>

                <div className="mb-4">
                    <label className="block text-gray-700 mb-2">Password</label>
                    <input
                        type="password"
                        value={password}
                        onChange={(e) => setPassword(e.target.value)}
                        className="border border-gray-300 p-2 w-full rounded"
                        required
                    />
                </div>

                <button type="submit" className="bg-blue-500 text-white p-2 rounded w-full">
                    Sign In
                </button>
            </form>
        </div>
    );
};

export default Login;
