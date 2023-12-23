export const fetchWithToken = async (url, options = {}) => {
    let accessToken = localStorage.getItem('accessToken');

    const headers = {
        'Content-Type': 'application/json',
        Authorization: `Bearer ${accessToken}`,
        ...options.headers,
    };
    let response = await fetch(url, { ...options, headers });

    // If access token is expired and there is a refresh token
    if (response.status === 401) {
        const refreshToken = localStorage.getItem('refreshToken');
        if (refreshToken) {
            const refreshed = await refreshAccessToken(refreshToken);
            if (refreshed) {
                accessToken = localStorage.getItem('accessToken'); // Get the new access token
                headers.Authorization = `Bearer ${accessToken}`;
                response = await fetch(url, { ...options, headers }); // Retry the original request
            }
        }
    }

    return response;
};

const refreshAccessToken = async (refreshToken) => {
    try {
        const response = await fetch('/api/auth/refresh', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ refreshToken }),
        });

        if (response.ok) {
            const data = await response.json();
            localStorage.setItem('accessToken', data.accessToken); // Update the access token
            return true;
        } else {
            // Handle refresh token failure (e.g., remove tokens, redirect to login)
            localStorage.removeItem('accessToken');
            localStorage.removeItem('refreshToken');
            window.location.href = '/login';
            return false;
        }
    } catch (error) {
        console.error('Error refreshing access token:', error);
        window.location.href = '/login';
        return false;
    }
};