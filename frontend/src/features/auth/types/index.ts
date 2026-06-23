import type { User } from "../../../types";

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginResponse {
  user: User;
  csrf_token: string;
}