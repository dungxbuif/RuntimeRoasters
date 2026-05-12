import axios from 'axios';
import { ENV } from '@/constants/env';
import { storageService } from '@/services/storage.service';

const api = axios.create({
  baseURL: ENV.GATEWAY_URL,
  headers: {
    'Content-Type': 'application/json',
  },
  withCredentials: true,
});

api.interceptors.request.use((config) => {
  const token = storageService.getAccessToken();
  console.log(`[API Request] ${config.method?.toUpperCase()} ${config.url}`, {
    hasToken: !!token,
    tokenPrefix: token?.substring(0, 10),
  });
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

api.interceptors.response.use(
  (response) => {
    console.log(`[API Response OK] ${response.config.method?.toUpperCase()} ${response.config.url}`, response.status);
    return response;
  },
  (error) => {
    console.error(`[API Response ERROR] ${error.config?.method?.toUpperCase()} ${error.config?.url}`, {
      status: error.response?.status,
      data: error.response?.data,
      message: error.message,
    });
    return Promise.reject(error);
  }
);

export default api;
