import React, { useState } from 'react';
import "../css/Register.css";

const Register = () => {
    const [formData, setFormData] = useState({
        username: '',
        password: '',
        email: '' // Add other fields as necessary
    });

    const [passwordMatch, setPasswordMatch] = useState(true);

    const handleChange = (e) => {
        setFormData({
            ...formData,
            [e.target.name]: e.target.value
        });
    };

    const handlePasswordBlur = () => {
        // Check if passwords match
        setPasswordMatch(formData.password === formData.confirmPassword);
      };

    const handleSubmit = async (e) => {
        e.preventDefault();
        if (passwordMatch) {
            try {
                const response = await fetch('/api/auth/register', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify(formData)
                });

                if (!response.ok) {
                    throw new Error('Registration failed');
                }

                const data = await response.json();
                console.log('Registration successful', data);
                // Handle successful registration (e.g., redirect to login page)
            } catch (error) {
                console.error('Error during registration:', error);
                // Handle errors in registration (e.g., display error message)
            }
        } else {
            alert("passwords don't match");
        }
    };

    return (
        <div className="register-container">
            <div className="register-header">
                <h1>Register for an Account</h1>
            </div>
            <div className="register-form">
                <form onSubmit={handleSubmit}>
                    <div class="form-group">
                        <label htmlFor="username">Username</label>
                        <input
                            type="text"
                            name="username"
                            value={formData.username}
                            onChange={handleChange}
                            required
                        />
                    </div>
                    <div className="form-group">
                        <label htmlFor="password">Password</label>
                            <input
                                type="password"
                                id="password"
                                name="password"
                                value={formData.password}
                                onChange={handleChange}
                                onBlur={handlePasswordBlur}
                            />
                    </div>
                    <div className="form-group">
                        <label htmlFor="confirmPassword">Confirm Password</label>
                            <input
                                type="password"
                                id="confirmPassword"
                                name="confirmPassword"
                                value={formData.confirmPassword}
                                onChange={handleChange}
                                onBlur={handlePasswordBlur}
                            />
                    </div>
                    {passwordMatch ? (
                        <p className="password-match">Passwords match</p>
                    ) : (
                        <p className="password-mismatch">Passwords do not match</p>
                    )}
                    <div className="form-group">
                        <label htmlFor="email">Email</label>
                        <input
                            type="email"
                            name="email"
                            value={formData.email}
                            onChange={handleChange}
                        />
                    </div>
                    <div className="register-actions">
                        <button type="submit">Register</button>
                    </div>
                </form>
            </div>
            
        </div>
        
    );
};

export default Register;
