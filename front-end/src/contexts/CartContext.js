import React, { createContext, useState, useEffect, useContext } from 'react';
import { fetchWithToken } from '../utils/authUtils';

export const CartContext = createContext();

export const CartProvider = ({ children }) => {
    const [cartData, setCartData] = useState([]);
    const [itemCount, setItemCount] = useState(0);
    const [cartLoading, setCartLoading] = useState(false);
    const [cartError, setCartError] = useState(null);

    useEffect(() => {
        // Calculate the total item count
        const totalItemCount = cartData.reduce((total, item) => total + item.quantity, 0);
        setItemCount(totalItemCount);
    }, [cartData]);

    const updateCartItems = (newCartData) => {
        setCartData(newCartData);
    };

    const fetchCartData = async () => {
        setCartLoading(true);
        try {
            const response = await fetchWithToken('/api/cart');
            if(!response.ok) {
                throw new Error('Failed to fetch cart');
            }
            const data = await response.json();
            if (data === null) {
                // Handle the scenario of an empty cart
                // For example, you might want to set an empty array
                updateCartItems([]);
            } else {
                updateCartItems(data);
            }
        } catch (err) {
            setCartError(err.message || 'An cart error occurred');
        } finally {
            setCartLoading(false);
        }

    }

    return (
        <CartContext.Provider value={{ cartData, updateCartItems, itemCount, fetchCartData, cartLoading, cartError }}>
            {children}
        </CartContext.Provider>
    );
};

export const useCart = () => useContext(CartContext);
