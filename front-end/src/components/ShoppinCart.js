import React, { useState, useEffect } from 'react';
import { fetchWithToken } from '../utils/authUtils';
import { useCart } from '../contexts/CartContext';

const ShoppingCart = () => {
    const { cartData, fetchCartData } = useCart();
    const [cartItemsWithDetails, setCartItemsWithDetails] = useState([]);
    
    // State to keep track of edited quantities
    const [editedQuantities, setEditedQuantities] = React.useState({});

    // Handle quantity change
    const handleQuantityChange = (productId, quantity) => {
        setEditedQuantities(prev => ({ ...prev, [productId]: quantity }));
    };

    const onUpdate = async (productId, newQuantity) => {
        // Logic to update the item on the server
        const response = await fetchWithToken(`/api/cart/update/${productId}`, {
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
        const response = await fetchWithToken(`/api/cart/delete/${productId}`, {
            method: 'DELETE'
        });
        if (!response.ok) {
            throw new Error('Failed to delete item');
        }
        // Re-fetch cart data
        fetchCartData();
    }; 

    useEffect(() => {
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
      }, [cartData]);
    
    const calculateSubtotal = (item) => {
        return item.productData.price * (editedQuantities[item.productID] || item.productData.quantity);
    };

    const total = cartData.reduce((acc, item) => {
        return acc + calculateSubtotal(item);
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
                            value={editedQuantities[item.productID] || item.quantity}
                            onChange={(e) => handleQuantityChange(item.productID, parseInt(e.target.value, 10))}
                        />
                        <span>Subtotal: {calculateSubtotal(item)}</span>
                    </div>
                    <div className="item-actions">
                        <button onClick={() => onUpdate(item.productID, editedQuantities[item.productID] || item.quantity)}>
                            Update
                        </button>
                        <button onClick={() => onDelete(item.productID)}>Delete</button>
                    </div>
                </div>
            ))}
            <div className="cart-total">
                <span>Total: {total}</span>
            </div>
        </div>
    );
};


export default ShoppingCart;
