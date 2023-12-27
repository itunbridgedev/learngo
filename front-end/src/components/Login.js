import React, { useState } from 'react';
import { Link } from 'react-router-dom';
import '../css/Login.css'

const Login = ({ onLoginSuccess }) => {
    const [username, setUsername] = useState('');
    const [password, setPassword] = useState('');
    const [loginError, setLoginError] = useState('');

    const handleLogin = async (e) => {
        e.preventDefault();

        try {
            const response = await fetch(`/api/auth/login`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ username, password }),
            });

            if (!response.ok) {
                throw new Error(`Error: ${response.status}`);
            }

            const data = await response.json();
            console.log('Login successful:', data);
            localStorage.setItem('accessToken', data.token);
            localStorage.setItem('refreshToken', data.refreshToken);
            onLoginSuccess();
        } catch (error) {
            console.error('Login failed:', error);
            setLoginError('Login failed. Please try again.');
        }
    };

    return (
        <div className="login-container">
            <div className="login-header">
                <h1>Login to your account</h1>
            </div>
            <div className="login-form">
                <form onSubmit={handleLogin}>
                    <div className="form-group">
                        <label htmlFor="username">Username</label>
                        <input
                            type="text"
                            value={username}
                            onChange={(e) => setUsername(e.target.value)}
                        />
                    </div>
                    <div>
                        <label htmlFor="password">Password</label>
                        <input
                            type="password"
                            value={password}
                            onChange={(e) => setPassword(e.target.value)}
                        />
                    </div>
                    <div className="form-actions">
                        <button type="submit">Login</button>
                    </div>
                </form>
            </div>
            <div className="login-errors">
                {loginError && <p>{loginError}</p>}
            </div>
            <div className="login-footer">
                <p>
                    Don't have an account? <Link to="/register">Sign up</Link>
                </p>
            </div>
        </div>
    );
};

export default Login;