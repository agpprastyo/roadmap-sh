import  { createContext, useContext, useState, useEffect } from 'react';

// Create a context to manage authentication state
const AuthContext = createContext(
    {}
);

// AuthProvider component to wrap around the application
export const AuthProvider = ({ children }) => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    useEffect(() => {
        // Check if there's a valid session in localStorage or cookies
        const userSession = localStorage.getItem('userSession');
        if (userSession) {
            setIsLoggedIn(true); // Set logged in state if user session exists
        }
    }, []);

    // Function to log in the user
    const login = () => {
        setIsLoggedIn(true);
        localStorage.setItem('userSession', 'true'); // Set a mock session in localStorage
    };

    // Function to log out the user
    const logout = () => {
        setIsLoggedIn(false);
        localStorage.removeItem('userSession'); // Clear the session on logout
    };

    // Exposing isLoggedIn, login, and logout
    return (
        <AuthContext.Provider value={{ isLoggedIn, login, logout }}>
            {children}
        </AuthContext.Provider>
    );
};

// Custom hook to use AuthContext
export const useAuth = () => useContext(AuthContext);
