import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useLogin } from "../hooks/useLogin";
import { useTheme } from "../../../lib/theme";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { getErrorMessage } from "../../../lib/errorMessages";

const loginSchema = z.object({
  email: z.string().email("Email tidak valid"),
  password: z.string().min(1, "Password wajib diisi"),
});

type LoginFormValues = z.infer<typeof loginSchema>;

export default function LoginPage() {
  const { mutate, isPending, error } = useLogin();
  const { theme, toggleTheme } = useTheme();

  const {
    register,
    handleSubmit,
    formState: { errors },
  } = useForm<LoginFormValues>({
    resolver: zodResolver(loginSchema),
  });

  const onSubmit = (values: LoginFormValues) => {
    mutate(values);
  };

  const apiError = error as any;
  const errorMessage = apiError
    ? getErrorMessage(apiError?.response?.data?.code, apiError?.response?.data?.message)
    : null;

  return (
    <div className="min-h-screen flex bg-surface-secondary">
      {/* Left panel — brand showcase, hidden on small screens */}
      <div className="hidden lg:flex lg:w-[44%] relative overflow-hidden bg-gradient-to-br from-primary-700 via-primary-600 to-danger-600 p-12 flex-col justify-between">
        <div className="absolute inset-0 opacity-[0.07]" style={{
          backgroundImage: "radial-gradient(circle at 2px 2px, white 1px, transparent 0)",
          backgroundSize: "28px 28px",
        }} />

        <div className="relative flex items-center gap-2.5">
          <div className="h-8 w-8 rounded-lg bg-white/15 backdrop-blur flex items-center justify-center text-white font-bold text-[15px] border border-white/20">
            P
          </div>
          <span className="text-white font-semibold text-[15px] tracking-tight">PPMS</span>
        </div>

        <div className="relative">
          <h1 className="text-white text-[28px] font-semibold leading-tight tracking-tight mb-3 max-w-md">
            Kelola portofolio proyek perusahaan dalam satu tempat.
          </h1>
          <p className="text-white/75 text-[14px] leading-relaxed max-w-sm">
            Dari pengajuan, approval, hingga eksekusi — pantau setiap proyek,
            anggaran, dan tim Anda secara real-time.
          </p>
        </div>

        <div className="relative flex items-center gap-6 text-white/60 text-[12.5px]">
          <span>© 2026 PPMS Internal</span>
        </div>
      </div>

      {/* Right panel — form */}
      <div className="flex-1 flex items-center justify-center p-8 relative">
        <button
          onClick={toggleTheme}
          className="absolute top-6 right-6 h-9 w-9 rounded-md border border-border bg-surface flex items-center justify-center text-ink-secondary hover:bg-surface-secondary transition-colors"
          aria-label="Toggle theme"
        >
          {theme === "light" ? (
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <path d="M21 12.8A9 9 0 1111.2 3 7 7 0 0021 12.8z" />
            </svg>
          ) : (
            <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
              <circle cx="12" cy="12" r="4" />
              <path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4" />
            </svg>
          )}
        </button>

        <div className="w-full max-w-[360px]">
          <div className="lg:hidden flex items-center gap-2.5 mb-8">
            <div className="h-8 w-8 rounded-lg bg-gradient-to-br from-primary-600 to-danger-600 flex items-center justify-center text-white font-bold text-[15px]">
              P
            </div>
            <span className="font-semibold text-[15px] tracking-tight">PPMS</span>
          </div>

          <h2 className="text-[20px] font-semibold tracking-tight mb-1.5">Masuk ke akun Anda</h2>
          <p className="text-[13.5px] text-ink-secondary mb-7">
            Gunakan kredensial perusahaan Anda untuk melanjutkan.
          </p>

          <form onSubmit={handleSubmit(onSubmit)} className="space-y-4">
            <div>
              <Label htmlFor="email">Email</Label>
              <Input id="email" type="email" placeholder="nama@perusahaan.com" {...register("email")} />
              {errors.email && <p className="text-xs text-danger-600 mt-1.5">{errors.email.message}</p>}
            </div>

            <div>
              <div className="flex items-center justify-between mb-1.5">
                <Label htmlFor="password" className="!mb-0">
                  Password
                </Label>
                <a href="#" className="text-xs font-medium text-primary-600 hover:text-primary-700">
                  Lupa password?
                </a>
              </div>
              <Input id="password" type="password" placeholder="••••••••" {...register("password")} />
              {errors.password && <p className="text-xs text-danger-600 mt-1.5">{errors.password.message}</p>}
            </div>

            {errorMessage && (
              <div className="flex items-start gap-2 rounded-md bg-danger-50 dark:bg-danger-900/20 border border-danger-200 dark:border-danger-900/40 px-3 py-2.5">
                <svg
                  width="15"
                  height="15"
                  viewBox="0 0 24 24"
                  fill="none"
                  stroke="currentColor"
                  strokeWidth="2"
                  className="text-danger-600 flex-shrink-0 mt-0.5"
                >
                  <circle cx="12" cy="12" r="9" />
                  <path d="M12 8v5M12 16h.01" />
                </svg>
                <p className="text-[12.5px] text-danger-700 dark:text-danger-400 m-0 leading-snug">{errorMessage}</p>
              </div>
            )}

            <Button type="submit" variant="primary" className="w-full" size="lg" disabled={isPending}>
              {isPending ? "Memproses..." : "Masuk"}
            </Button>
          </form>

          <p className="text-center text-xs text-ink-tertiary mt-8">
            Butuh bantuan akses?{" "}
            <a href="#" className="text-primary-600 font-medium hover:text-primary-700">
              Hubungi administrator
            </a>
          </p>
        </div>
      </div>
    </div>
  );
}