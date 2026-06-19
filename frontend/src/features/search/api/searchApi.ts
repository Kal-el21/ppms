import axiosInstance from "../../../api/axiosInstance";
import type { ApiResponse } from "../../../types";
import type { SearchResponse } from "../types";

export const searchApi = {
  search: async (query: string) => {
    const res = await axiosInstance.get<ApiResponse<SearchResponse>>("/search", {
      params: { q: query },
    });
    return res.data.data;
  },
};