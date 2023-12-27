import React, { useState, useEffect } from 'react';
import { fetchWithToken } from '../utils/authUtils';
import { connect } from 'react-redux';
import { fetchCartData } from '../actions/cartActions';

const ShoppingCart = ({ cartData, cartLoading, cartError, fetchCartData, hasFetched }) => {
    const [cartItemsWithDetails, setCartItemsWithDetails] = useState([]);
    
    // State to keep track of edited quantities
    const [editedQuantities, setEditedQuantities] = React.useState({});

    // Handle quantity change
    const handleQuantityChange = (productId, quantity) => {
        setEditedQuantities(prev => ({ ...prev, [productId]: quantity }));
    };

    const onUpdate = async (productId, newQuantity) => {
        // Logic to update the item on the server
        const response = await fetchWithToken(`/api/cart/items/${productId}`, {
            method: 'PUT',
            body: JSON.stringify({ quantity: newQuantity })
        });
        if (!response.ok) {
            throw new Error('Failed to update item');
        }
        // Re-fetch cart data
        fetchCartData();
    };

    const onDelete = async (productId) => {
        // Logic to delete the item from the server
        const response = await fetchWithToken(`/api/cart/items/${productId}`, {
            method: 'DELETE'
        });
        if (!response.ok) {
            throw new Error('Failed to delete item');
        }
        // Re-fetch cart data
        fetchCartData();
    }; 

    useEffect(() => {
        if(!hasFetched) {
            fetchCartData();
        }        
        // Fetch product details for each item in the cart
        const fetchProductDetails = async () => {
          const cartItemsWithDetails = await Promise.all(
            cartData.map(async (cartItem) => {
              const response = await fetchWithToken(`/api/products/${cartItem.product_id}`);
              if (response.ok) {
                const productData = await response.json();
                return { ...cartItem, productData };
              }
              return null; // Handle error or missing product
            })
          );
          setCartItemsWithDetails(cartItemsWithDetails.filter(Boolean));
        };
    
        if (cartData.length > 0) {
          fetchProductDetails();
        }
      }, [fetchCartData, cartData, hasFetched]);
    
    const calculateSubtotal = (item) => {
        return item.productData.price * (editedQuantities[item.product_id] || item.quantity);
    };

    const total = cartItemsWithDetails.reduce((acc, item) => {
        return acc + calculateSubtotal(item, "total");
    }, 0);
    // Render each cart item
    return (
        <div className="shopping-cart">
            {cartItemsWithDetails.map((item) => (
                <div key={item.productID} className="cart-item">
                    <div className="item-info">
                        <span>Product: {item.productData.name}</span>
                        <span>Price: {item.productData.price}</span>
                        <input
                            type="number"
                            value={editedQuantities[item.product_id] || item.quantity}
                            onChange={(e) => handleQuantityChange(item.product_id, parseInt(e.target.value, 10))}
                        />
                        <span>Subtotal: {calculateSubtotal(item, "sub")}</span>
                    </div>
                    <div className="item-actions">
                        <button onClick={() => onUpdate(item.product_id, editedQuantities[item.product_id] || item.quantity)}>
                            Update
                        </button>
                        <button onClick={() => onDelete(item.product_id)}>Delete</button>
                    </div>
                </div>
            ))}
            <div className="cart-total">
                <span>Total: {total}</span>
            </div>
        </div>
    );
};

const mapStateToProps = (state) => ({
    cartData: state.cart.cartData,
    cartLoading: state.cart.cartLoading,
    cartError: state.cart.cartError,
    hasFetched: state.cart.hasFetched
});

const mapDispatchToProps = {
    fetchCartData,
};

export default connect(mapStateToProps, mapDispatchToProps)(ShoppingCart);
