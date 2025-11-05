import { createBrowserRouter, RouterProvider } from 'react-router';
import { routes } from '@/app/routes';

/**
 * Router configuration provider
 * Wraps the application with React Router v7
 */
const router = createBrowserRouter(routes);

export const AppRouter = () => {
  return <RouterProvider router={router} />;
};
