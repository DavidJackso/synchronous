import { useEffect, useState } from 'react';
import { AppRouter } from '@/app/providers/RouterProvider';
import { useAuth } from '@/app/store';
import { useMaxWebApp } from '@/shared/hooks/useMaxWebApp';
import { Spin } from 'antd';
import './App.css';

/**
 * Root App component
 * 
 * Handles automatic authentication via MAX initData
 * User doesn't need to interact - login happens automatically
 */
function App() {
  const { login, isAuthenticated, isLoading: authLoading } = useAuth();
  const { initData, user, isReady } = useMaxWebApp();
  const [isInitializing, setIsInitializing] = useState(true);

  // Automatic login when MAX initData is available
  useEffect(() => {
    const performAutoLogin = async () => {
      // Wait for MAX WebApp to be ready
      if (!isReady) {
        return;
      }

      // If already authenticated, stop initialization
      if (isAuthenticated) {
        console.log('[App] Already authenticated');
        setIsInitializing(false);
        return;
      }

      // If no initData available (dev mode), skip auth
      if (!initData || !user) {
        console.warn('[App] No initData - running in dev mode without auth');
        setIsInitializing(false);
        return;
      }

      // Perform automatic login
      try {
        console.log('[App] Auto-login with MAX initData', { user });
        const deviceId = navigator.userAgent;
        await login(initData, deviceId);
        console.log('[App] Auto-login successful');
      } catch (error) {
        console.error('[App] Auto-login failed:', error);
      } finally {
        setIsInitializing(false);
      }
    };

    performAutoLogin();
  }, [isReady, initData, user, isAuthenticated, login]);

  // Show loading spinner during initialization
  if (isInitializing || authLoading) {
    return (
      <div style={{
        display: 'flex',
        justifyContent: 'center',
        alignItems: 'center',
        minHeight: '100vh',
        background: 'linear-gradient(135deg, #1e293b 0%, #334155 100%)',
      }}>
        <Spin size="large" tip="Загрузка приложения..." />
      </div>
    );
  }

  return <AppRouter />;
}

export default App;
