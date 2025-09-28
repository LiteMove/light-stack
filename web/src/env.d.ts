// src/env.d.ts
interface ImportMetaEnv {
    readonly VITE_API_BASE_URL: string;
    readonly VITE_APP_TITLE: string;
    readonly VITE_APP_VERSION: string;
}

interface ImportMeta {
    readonly env: ImportMetaEnv;
}

declare module '@vue/runtime-core' {
  interface ComponentCustomProperties {
    $hasRole: (role: string) => boolean
    $hasAnyRole: (roles: string[]) => boolean
    $hasAllRole: (roles: string[]) => boolean
    $isAdmin: () => boolean
    $isSuperAdmin: () => boolean
    $hasAuth: (config: { roles?: string[], requireAll?: boolean }) => boolean
  }
}