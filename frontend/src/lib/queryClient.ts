import { QueryClient } from "@tanstack/react-query";

export const queryClient = new QueryClient({
  defaultOptions: {
    queries: {
      retry: 1,
      refetchOnWindowFocus: false,
      staleTime: 1000 * 60,
    },
    mutations: {
      onError: (error: any) => {
        // Global fallback error — komponen individual bisa override ini
        // dengan onError di useMutation mereka sendiri
        console.error("Mutation error:", error?.response?.data?.message || error.message);
      },
    },
  },
});