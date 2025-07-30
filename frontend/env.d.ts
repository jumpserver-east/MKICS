/// <reference types="vite/client" />

interface ImportMetaEnv {
  readonly VITE_API_URL: string;
  // Add other VITE_ environment variables here
  readonly VITE_APP_TITLE?: string;
  // ... other env variables
}

interface ImportMeta {
  readonly env: ImportMetaEnv;
}
