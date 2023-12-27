import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import { connect } from 'react-redux';
import { updateCartItems, fetchCartData } from '../actions/cartActions';

const ShoppingCartIcon = ({ cartData, itemCount, fetchCartData, hasFetched }) => {
    const navigate = useNavigate();
    const navigateToCart = () => {
        navigate('/cart');
    };
    useEffect(() => {
        if (!hasFetched) {
            fetchCartData();
        }
    }, [fetchCartData, hasFetched]);

    return (
        <div onClick={navigateToCart} style={{ display: 'flex', alignItems: 'center', cursor: 'pointer' }}>
            {/* Icon - you can replace this with an actual icon from a library like FontAwesome */}
            <svg
                style={{ width: '24px', height: '24px', marginRight: '10px' }}
                viewBox="0 0 24 24"
            >
                <path fill="currentColor" d="M7,18A2,2 0 0,1 5,20A2,2 0 0,1 3,18A2,2 0 0,1 5,16A2,2 0 0,1 7,18M1,2H3.27L6.89,13.59L5.25,16.04L5.2,16.18A1.73,1.73 0 0,0 5,17A2,2 0 0,0 7,19H19A2,2 0 0,0 21,17A2,2 0 0,0 19,15H7.42L8.27,12H17.55A1.73,1.73 0 0,0 17.7,11.8A1.73,1.73 0 0,0 17.55,11.6L15.67,7H9.59L8.44,4.56L7.25,2M19,18A2,2 0 0,1 17,20A2,2 0 0,1 15,18A2,2 0 0,1 17,16A2,2 0 0,1 19,18Z" />
            </svg>

            {/* Item Count */}
            <span>Cart: {itemCount}</span>
        </div>
    );
};
const mapStateToProps = (state) => ({
    cartData: state.cart.cartData || [],
    itemCount: (state.cart.cartData || []).reduce((total, item) => total + item.quantity, 0),
    hasFetched: state.cart.hasFetched
});

const mapDispatchToProps = {
    updateCartItems,
    fetchCartData,
};

export default connect(mapStateToProps, mapDispatchToProps)(ShoppingCartIcon);