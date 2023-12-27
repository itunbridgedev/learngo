// cartActions.js
import { fetchWithToken } from '../utils/authUtils'; // Assuming you have an 'api.js' file

// Define your Redux action creators
export const fetchCartData = () => async (dispatch) => {
    // Dispatch action to indicate that cart data fetch has started
    dispatch({ type: 'FETCH_CART_REQUEST' });

    try {
        const response = await fetchWithToken('/api/cart');
        if (!response.ok) {
            throw new Error('Failed to fetch cart');
        }
        const data = await response.json();
        if (data === null) {
            // Handle the scenario of an empty cart
            // For example, you might want to set an empty array
            dispatch(updateCartItems([])); // Dispatch your Redux action to update cart items
        } else {
            dispatch(updateCartItems(data)); // Dispatch your Redux action to update cart items
        }
    } catch (err) {
        // Dispatch an action to handle the error
        dispatch({ type: 'FETCH_CART_FAILURE', payload: err.message || 'A cart error occurred' }); 
    }
};

export const updateCartItems = (items) => ({
    type: 'UPDATE_CART_ITEMS',
    payload: items,
});
