import { STORAGE_KEYS } from "@/constants/storage";
import { isBrowser } from "@/lib/utils/window";

/**
 * Define the mapping between storage keys and their value types
 */
export interface StorageMap {
  [STORAGE_KEYS.ACCESS_TOKEN]: string;
  [STORAGE_KEYS.REFRESH_TOKEN]: string;
  [STORAGE_KEYS.USER_IDENTITY]: Record<string, unknown>;
}

class StorageService {
  /**
   * Generic setItem with type safety based on key
   */
  setItem<K extends keyof StorageMap>(key: K, value: StorageMap[K]): void {
    if (!isBrowser) return;
    
    const stringValue = typeof value === 'string' ? value : JSON.stringify(value);
    localStorage.setItem(key, stringValue);
  }

  /**
   * Generic getItem with automatic parsing and type safety
   */
  getItem<K extends keyof StorageMap>(key: K): StorageMap[K] | null {
    if (!isBrowser) return null;

    const value = localStorage.getItem(key);
    if (!value) return null;

    try {
      // If the expected type is string, return directly
      // Note: We check if it looks like JSON or if it's supposed to be an object
      if (key === STORAGE_KEYS.ACCESS_TOKEN || key === STORAGE_KEYS.REFRESH_TOKEN) {
        return value as unknown as StorageMap[K];
      }
      return JSON.parse(value) as StorageMap[K];
    } catch {
      return value as unknown as StorageMap[K];
    }
  }

  removeItem(key: keyof StorageMap): void {
    if (isBrowser) {
      localStorage.removeItem(key);
    }
  }

  // Sugar methods for common keys
  setAccessToken(token: string) {
    this.setItem(STORAGE_KEYS.ACCESS_TOKEN, token);
  }

  getAccessToken(): string | null {
    return this.getItem(STORAGE_KEYS.ACCESS_TOKEN);
  }

  setIdentity(identity: StorageMap[typeof STORAGE_KEYS.USER_IDENTITY]) {
    this.setItem(STORAGE_KEYS.USER_IDENTITY, identity);
  }

  getIdentity(): StorageMap[typeof STORAGE_KEYS.USER_IDENTITY] | null {
    return this.getItem(STORAGE_KEYS.USER_IDENTITY);
  }

  clearAll() {
    if (isBrowser) {
      Object.values(STORAGE_KEYS).forEach((key) => localStorage.removeItem(key));
    }
  }
}

export const storageService = new StorageService();
