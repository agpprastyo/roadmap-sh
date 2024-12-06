// src/components/PrivateRoute.jsx

import { Navigate } from 'react-router-dom'; // Used to redirect
import { useAuth } from '../contexts/AuthContext'; // Import auth context

const PrivateRoute = ({ element }) => {
    const { isLoggedIn } = useAuth(); // Check login status from context

    if (!isLoggedIn) {
        // If not logged in, redirect to the login page
        return <Navigate to="/login" />;
    }

    return element; // If logged in, grant access to the requested element
};

export default PrivateRoute;
