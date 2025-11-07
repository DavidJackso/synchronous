import { configureStore } from '@reduxjs/toolkit';
import sessionSetupReducer from '@/entities/session/model/sessionSetupSlice';

/**
 * Root Redux store configuration
 * Central state management for the application
 */
export const store = configureStore({
  reducer: {
    sessionSetup: sessionSetupReducer,
  },
  middleware: (getDefaultMiddleware) =>
    getDefaultMiddleware({
      serializableCheck: {
        ignoredActions: ['persist/PERSIST'],
      },
    }),
});

export type RootState = ReturnType<typeof store.getState>;
export type AppDispatch = typeof store.dispatch;
