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
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

export default api;
