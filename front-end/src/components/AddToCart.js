// src/components/AddToCart.js
import React, { useState } from 'react';
import { fetchWithToken } from '../utils/authUtils'; // Adjust the import path as needed
import { connect } from 'react-redux';
import { updateCartItems } from '../actions/cartActions';

const AddToCart = ({ product_id, updateCartItems }) => {
    const [quantity, setQuantity] = useState(1);
    const addToCart = async () => {
        const response = await fetchWithToken('/api/cart/items', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ product_id, quantity }),
        });

        if (response.ok) {
            const cartData = await response.json();
            updateCartItems(cartData);
        } else {
            // Handle errors
        }
    };

    return (
        <div>
            <input
                type="number"
                value={quantity}
                onChange={(e) => setQuantity(e.target.value)}
                min="1"
                style={{ width: '50px', marginRight: '10px' }}
            />
            <button onClick={() => addToCart(product_id, quantity)}>Add to Cart</button>
        </div>
    );
};

const mapStateToProps = (state) => ({
    cartData: state.cart.cartData,
});

const mapDispatchToProps = {
    updateCartItems,
};

export default connect(mapStateToProps, mapDispatchToProps)(AddToCart);