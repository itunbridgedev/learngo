import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './Login';
import ProductList from './ProductList';
import ShoppingCart from './ShoppinCart';
import { useCart } from '../contexts/CartContext';

const App = () => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);
    const { fetchCartData, cartLoading, cartError, cartData } = useCart();

    const handleLoginSuccess = () => {
        setIsLoggedIn(true);
    };   

    useEffect(() => {
        if (isLoggedIn && !cartData.length) { // Add a condition to fetch data only once when cartData is empty
            fetchCartData();
        }
    }, [isLoggedIn, fetchCartData, cartData]);

    if (cartLoading) return <div>Loading cart...</div>;
    if (cartError) return <div>Error: {cartError}</div>;
    return (
        <Router>
            <Routes>
                <Route path="/login" element={isLoggedIn ? <Navigate to="/products" /> : <Login onLoginSuccess={handleLoginSuccess} />} />
                <Route path="/products" element={<ProductList />} />
                <Route path="/cart" element={<ShoppingCart />} />
                <Route path="/" element={<Navigate to="/login" />} />
            </Routes>
        </Router>
    );
};

export default App;
