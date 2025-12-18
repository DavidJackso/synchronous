/**
 * Telegram WebApp Hook
 * Provides access to Telegram WebApp API
 */

import { useEffect, useState } from 'react';

interface TelegramWebAppUser {
  id: number;
  first_name: string;
  last_name?: string;
  username?: string;
  language_code?: string;
  photo_url?: string;
}

interface TelegramWebAppChat {
  id: number;
  type: string;
  title?: string;
  username?: string;
  photo_url?: string;
}

interface TelegramWebAppData {
  query_id?: string;
  auth_date: number;
  hash: string;
  start_param?: string;
  user?: TelegramWebAppUser;
  chat?: TelegramWebAppChat;
}

interface TelegramWebApp {
  initData: string;
  initDataUnsafe: TelegramWebAppData;
  version: string;
  platform: 'ios' | 'android' | 'desktop' | 'web';
  ready: () => void;
  close: () => void;
  expand: () => void;
  enableClosingConfirmation: () => void;
  disableClosingConfirmation: () => void;
  openLink: (url: string) => void;
  openTelegramLink: (url: string) => void;
  isExpanded: boolean;
  viewportHeight: number;
  viewportStableHeight: number;
  headerColor: string;
  backgroundColor: string;
  isClosingConfirmationEnabled: boolean;
}

declare global {
  interface Window {
    Telegram?: {
      WebApp?: TelegramWebApp;
    };
  }
}

export const useTelegramWebApp = () => {
  const [webApp, setWebApp] = useState<TelegramWebApp | null>(null);
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    // Wait for Telegram WebApp to load
    if (typeof window !== 'undefined' && window.Telegram?.WebApp) {
      const app = window.Telegram.WebApp;
      setWebApp(app);
      
      // Notify Telegram that mini-app is ready
      app.ready();
      setIsReady(true);

      console.log('[Telegram WebApp] ✅ Initialized successfully', {
        version: app.version,
        platform: app.platform,
        hasInitData: !!app.initData,
        initDataLength: app.initData?.length || 0,
        initDataPreview: app.initData?.substring(0, 100) + '...',
        user: app.initDataUnsafe.user,
        viewportHeight: app.viewportHeight,
      });
    } else {
      // Running outside Telegram (dev environment)
      console.warn('[Telegram WebApp] ⚠️ Not available - running in dev mode');
      setIsReady(true);
    }
  }, []);

  return {
    webApp,
    isReady,
    isTelegramEnvironment: !!webApp,
    initData: webApp?.initData || null,
    user: webApp?.initDataUnsafe.user || null,
  };
};
