import type { RouteObject } from 'react-router';

/**
 * Application routes configuration
 * Defines all available routes in the application
 */
export const routes: RouteObject[] = [
  {
    path: '/',
    lazy: async () => {
      const { HomePage } = await import('@/pages/home');
      return { Component: HomePage };
    },
  },
  {
    path: '/onboarding',
    lazy: async () => {
      const { OnboardingPage } = await import('@/pages/onboarding');
      return { Component: OnboardingPage };
    },
  },
  {
    path: '/session-setup',
    lazy: async () => {
      const { SessionSetupPage } = await import('@/pages/session-setup');
      return { Component: SessionSetupPage };
    },
  },
  {
    path: '/lobby/:sessionId',
    lazy: async () => {
      const { LobbyPage } = await import('@/pages/lobby');
      return { Component: LobbyPage };
    },
  },
  {
    path: '/focus-session/:sessionId',
    lazy: async () => {
      const { FocusSessionPage } = await import('@/pages/focus-session');
      return { Component: FocusSessionPage };
    },
  },
  {
    path: '/session-report/:sessionId',
    lazy: async () => {
      const { SessionReportPage } = await import('@/pages/session-report');
      return { Component: SessionReportPage };
    },
  },
];
