export type SystemRole = "ADMIN" | "USER" | "VIEWER";

export interface User {
  id: number;
  full_name: string;
  email: string;
  system_role: SystemRole;
  division_id: number | null;
  is_active?: boolean;
}

export interface Division {
  id: number;
  name: string;
  description: string;
}

export interface ApiResponse<T> {
  success: boolean;
  data: T;
  message?: string;
}

export interface ApiErrorResponse {
  success: false;
  code: string;
  message: string;
}

export interface PaginatedResponse<T> {
  success: boolean;
  data: T[];
  meta: {
    page: number;
    limit: number;
    total: number;
  };
}