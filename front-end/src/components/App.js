import React, { useState, useEffect } from 'react';
import { BrowserRouter as Router, Route, Routes, Navigate } from 'react-router-dom';
import Login from './Login';
import Register from './Register';
import ProductList from './ProductList';
import ShoppingCart from './ShoppinCart';
import { connect } from 'react-redux';
import { fetchCartData } from '../actions/cartActions';

const App = ({ cartLoading, cartError, fetchCartData }) => {
    const [isLoggedIn, setIsLoggedIn] = useState(false);

    const handleLoginSuccess = () => {
        setIsLoggedIn(true);
    };   

    useEffect(() => {
        if (isLoggedIn) {
            fetchCartData();
        }
    }, [isLoggedIn, fetchCartData]); 

    if (cartLoading) return <div>Loading cart...</div>;
    if (cartError) return <div>Error: {cartError}</div>;
    return (
        <Router>
            <Routes>
                <Route path="/login" element={isLoggedIn ? <Navigate to="/products" /> : <Login onLoginSuccess={handleLoginSuccess} />} />
                <Route path="/register" element={<Register />} />
                <Route path="/products" element={<ProductList />} />
                <Route path="/cart" element={<ShoppingCart />} />
                <Route path="/" element={<Navigate to="/login" />} />
            </Routes>
        </Router>
    );
};

const mapStateToProps = (state) => ({
    cartData: state.cart.cartData,
    cartLoading: state.cart.cartLoading,
    cartError: state.cart.cartError,
});

const mapDispatchToProps = {
    fetchCartData,
};

export default connect(mapStateToProps, mapDispatchToProps)(App);