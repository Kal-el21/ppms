import { getErrorMessage } from "@/lib/errorMessages";
import axios from "axios";

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
});

axiosInstance.interceptors.request.use((config) => {
  const token = localStorage.getItem("access_token");
  if (token) {
    config.headers.Authorization = `Bearer ${token}`;
  }
  return config;
});

let isRefreshing = false;
let refreshSubscribers: ((token: string) => void)[] = [];

function onRefreshed(token: string) {
  refreshSubscribers.forEach((callback) => callback(token));
  refreshSubscribers = [];
}

axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.status === 401 && !originalRequest._retry) {
      if (isRefreshing) {
        return new Promise((resolve) => {
          refreshSubscribers.push((token: string) => {
            originalRequest.headers.Authorization = `Bearer ${token}`;
            resolve(axiosInstance(originalRequest));
          });
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      const refreshToken = localStorage.getItem("refresh_token");
      if (!refreshToken) {
        localStorage.clear();
        window.location.href = "/login";
        if (error.response?.data?.code) {
          error.friendlyMessage = getErrorMessage(error.response.data.code, error.response.data.message);
        }
        return Promise.reject(error);
      }

      try {
        const res = await axios.post(
          `${import.meta.env.VITE_API_BASE_URL || "http://localhost:8080/api/v1"}/auth/refresh`,
          { refresh_token: refreshToken }
        );

        const { access_token, refresh_token } = res.data.data;
        localStorage.setItem("access_token", access_token);
        localStorage.setItem("refresh_token", refresh_token);

        isRefreshing = false;
        onRefreshed(access_token);

        originalRequest.headers.Authorization = `Bearer ${access_token}`;
        return axiosInstance(originalRequest);
      } catch (refreshError) {
        isRefreshing = false;
        localStorage.clear();
        window.location.href = "/login";
        if (error.response?.data?.code) {
          error.friendlyMessage = getErrorMessage(error.response.data.code, error.response.data.message);
        }
        return Promise.reject(refreshError);
      }
    }
    if (error.response?.data?.code) {
      error.friendlyMessage = getErrorMessage(error.response.data.code, error.response.data.message);
    }
    return Promise.reject(error);
  }
);

export default axiosInstance;