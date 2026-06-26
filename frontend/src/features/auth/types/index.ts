import type { User } from "../../../types";

export interface LoginRequest {
  email: string;
  password: string;
}

export interface LoginInitResponse {
  two_fa_required: true;
  otp_session_token: string;
}

export interface LoginSuccessResponse {
  two_fa_required?: false;
  user: User;
  csrf_token: string;
}

export type LoginResponse = LoginInitResponse | LoginSuccessResponse;
