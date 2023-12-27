// cartReducer.js
const initialState = {
    cartData: [],
    cartLoading: false,
    cartError: null,
    hasFetched: false,
};

const cartReducer = (state = initialState, action) => {
    switch (action.type) {
        case 'FETCH_CART_REQUEST':
            return {
                ...state,
                cartLoading: true,
            };
        case 'FETCH_CART_SUCCESS':
            return {
                ...state,
                cartData: action.payload,
                cartLoading: false,
                hasFetched: true,
                cartError: null,
            };
        case 'FETCH_CART_FAILURE':
            return {
                ...state,
                cartData: [],
                cartLoading: false,
                hasFetched: true,
                cartError: action.payload,
            };
        case 'UPDATE_CART_ITEMS':
            return {
                ...state,
                cartData: action.payload,
                cartLoading: false,
                hasFetched: true,
                cartError: null,
            };
        default:
            return state;
    }
};

export default cartReducer;
