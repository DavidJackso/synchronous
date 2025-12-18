/**
 * Authentication API
 * Endpoints for login, logout, and token refresh
 * 
 * SECURITY: Tokens are stored in http-only cookies by backend
 * Frontend does not handle tokens directly
 */

import { apiClient, isAxiosError } from './client';
import type {
  LoginRequest,
  LoginResponse,
  RefreshTokenResponse,
} from './types';

// ============================================================================
// Auth Endpoints
// ============================================================================

/**
 * Login with Telegram initData
 * Backend will set http-only cookies for access and refresh tokens
 * 
 * @param initData - Telegram WebApp initData string (validated by backend)
 * @param deviceId - Unique device identifier
 * @returns User profile data (tokens are in cookies)
 */
export const login = async (
  initData: string,
  deviceId: string
): Promise<LoginResponse> => {
  console.log('[Auth API] 📤 Sending login request', {
    initDataLength: initData.length,
    initDataFull: initData, // Full initData for debugging
    deviceIdPreview: deviceId.substring(0, 50),
  });

  const payload: LoginRequest = {
    initData,
    deviceId,
  };

  try {
    const response = await apiClient.post<LoginResponse>('/auth/login', payload);

    const allCookies = document.cookie;
    const setCookieHeader = response.headers?.['set-cookie'];
    
    console.log('[Auth API] ✅ Login successful', {
      user: response.data.user,
      cookiesReceived: allCookies.includes('access_token'),
      allCookies: allCookies || '(no cookies - http-only cookies are not visible in JS)',
      setCookieHeader: setCookieHeader || '(no Set-Cookie header)',
      responseHeaders: response.headers,
    });
    
    // Note: http-only cookies are NOT visible in document.cookie
    // This is expected behavior - they're stored by browser but not accessible to JS
    if (!allCookies.includes('access_token')) {
      console.warn('[Auth API] ⚠️ access_token not visible in document.cookie - this is NORMAL for http-only cookies');
    }

    // Tokens are automatically stored in http-only cookies by backend
    // No need to manually store them
    return response.data;
  } catch (error: unknown) {
    console.error('[Auth API] ❌ Login failed', error);
    
    // Log detailed error information
    if (isAxiosError(error)) {
      // TypeScript now knows error is AxiosError
      const axiosError = error;
      console.error('[Auth API] Error details:', {
        status: axiosError.response?.status,
        statusText: axiosError.response?.statusText,
        data: axiosError.response?.data,
        headers: axiosError.response?.headers,
      });
    }
    
    throw error;
  }
};

/**
 * Refresh access token using refresh token from http-only cookie
 * Backend reads refresh token from cookie and sets new access token cookie
 * 
 * @returns New access token (in cookie)
 */
export const refreshAccessToken = async (): Promise<RefreshTokenResponse> => {
  const response = await apiClient.post<RefreshTokenResponse>('/auth/refresh');

  // New access token is automatically stored in http-only cookie by backend
  return response.data;
};

/**
 * Logout current user and clear cookies
 * Backend will clear http-only cookies
 */
export const logout = async (): Promise<void> => {
  await apiClient.post('/auth/logout');
  // Cookies are cleared by backend
};

/**
 * Check if user is authenticated
 * We can't access http-only cookies from JS, so we need to check via API call
 * or maintain auth state in Redux after login
 */
export const checkAuth = async (): Promise<boolean> => {
  try {
    // Try to get user profile - if succeeds, user is authenticated
    await apiClient.get('/users/me');
    return true;
  } catch {
    return false;
  }
};
