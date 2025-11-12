/**
 * MAX WebApp Hook
 * Provides access to MAX Bridge API
 */

import { useEffect, useState } from 'react';
import type { MaxWebApp } from '@/shared/types/max-webapp';

export const useMaxWebApp = () => {
  const [webApp, setWebApp] = useState<MaxWebApp | null>(null);
  const [isReady, setIsReady] = useState(false);

  useEffect(() => {
    // Wait for MAX Bridge to load
    if (typeof window !== 'undefined' && window.WebApp) {
      const app = window.WebApp;
      setWebApp(app);
      
      // Notify MAX that mini-app is ready
      app.ready();
      setIsReady(true);

      console.log('[MAX WebApp] ✅ Initialized successfully', {
        version: app.version,
        platform: app.platform,
        hasInitData: !!app.initData,
        initDataLength: app.initData?.length || 0,
        initDataPreview: app.initData?.substring(0, 100) + '...',
        user: app.initDataUnsafe.user,
        viewportHeight: app.viewportHeight,
      });
    } else {
      // Running outside MAX (dev environment)
      console.warn('[MAX WebApp] ⚠️ Not available - running in dev mode');
      setIsReady(true);
    }
  }, []);

  return {
    webApp,
    isReady,
    isMaxEnvironment: !!webApp,
    initData: webApp?.initData || null,
    user: webApp?.initDataUnsafe.user || null,
  };
};
