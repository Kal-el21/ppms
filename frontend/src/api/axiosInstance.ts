import axios from "axios";
import { getErrorMessage } from "../lib/errorMessages";

const axiosInstance = axios.create({
  baseURL: import.meta.env.VITE_API_BASE_URL || "http://localhost:8081/api/v1",
  headers: {
    "Content-Type": "application/json",
  },
  // withCredentials WAJIB true agar browser menyertakan httpOnly cookie
  // (access token, refresh token, csrf token) di setiap request — tanpa ini,
  // cookie tidak pernah terkirim sama sekali meskipun ada di browser.
  withCredentials: true,
});

function getCSRFTokenFromCookie(): string | null {
  // CSRF cookie SENGAJA tidak httpOnly (lihat backend cookie.go), jadi bisa
  // dibaca di sini untuk disertakan sebagai header di setiap request mutasi.
  const match = document.cookie.match(/(?:^|;\s*)ppms_csrf_token=([^;]*)/);
  return match ? decodeURIComponent(match[1]) : null;
}

axiosInstance.interceptors.request.use((config) => {
  const method = config.method?.toUpperCase();
  if (method && !["GET", "HEAD", "OPTIONS"].includes(method)) {
    const csrfToken = getCSRFTokenFromCookie();
    if (csrfToken) {
      config.headers["X-CSRF-Token"] = csrfToken;
    }
  }
  return config;
});

let isRefreshing = false;
let refreshSubscribers: (() => void)[] = [];

function onRefreshed() {
  refreshSubscribers.forEach((callback) => callback());
  refreshSubscribers = [];
}

axiosInstance.interceptors.response.use(
  (response) => response,
  async (error) => {
    const originalRequest = error.config;

    if (error.response?.data?.code) {
      error.friendlyMessage = getErrorMessage(error.response.data.code, error.response.data.message);
    }

    if (error.response?.status === 401 && !originalRequest._retry && !originalRequest.url?.includes("/auth/")) {
      if (isRefreshing) {
        return new Promise((resolve) => {
          refreshSubscribers.push(() => {
            resolve(axiosInstance(originalRequest));
          });
        });
      }

      originalRequest._retry = true;
      isRefreshing = true;

      try {
        // Refresh token dikirim otomatis via cookie (path-scoped ke /api/v1/auth),
        // tidak perlu dikirim manual di body seperti versi localStorage sebelumnya.
        await axios.post(
          `${import.meta.env.VITE_API_BASE_URL || "http://localhost:8081/api/v1"}/auth/refresh`,
          {},
          { withCredentials: true }
        );

        isRefreshing = false;
        onRefreshed();

        return axiosInstance(originalRequest);
      } catch (refreshError) {
        isRefreshing = false;
        window.location.href = "/login";
        return Promise.reject(refreshError);
      }
    }

    return Promise.reject(error);
  }
);

export default axiosInstance;
