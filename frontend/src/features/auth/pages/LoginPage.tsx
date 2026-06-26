import { useState } from "react";
import { useForm } from "react-hook-form";
import { zodResolver } from "@hookform/resolvers/zod";
import { z } from "zod";
import { useAuth } from "../context/AuthContext";
import { useNavigate } from "react-router-dom";
import { useTheme } from "../../../lib/theme";
import { useLogin, useVerifyOTP, useResendOTP } from "../hooks/useAuthActions";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { getErrorMessage } from "../../../lib/errorMessages";

const loginSchema = z.object({
  email: z.string().email("Email tidak valid"),
  password: z.string().min(1, "Password wajib diisi"),
});

const otpSchema = z.object({
  otp_code: z.string().length(6, "Kode OTP harus 6 digit"),
});

type LoginFormValues = z.infer<typeof loginSchema>;
type OTPFormValues = z.infer<typeof otpSchema>;

export default function LoginPage() {
  const { login } = useAuth();
  const navigate = useNavigate();
  const { theme, toggleTheme } = useTheme();

  const [step, setStep] = useState<"credentials" | "otp">("credentials");
  const [otpSessionToken, setOtpSessionToken] = useState("");
  const [isLoading, setIsLoading] = useState(false);
  const [errorMsg, setErrorMsg] = useState<string | null>(null);

  const loginForm = useForm<LoginFormValues>({ resolver: zodResolver(loginSchema) });
  const otpForm = useForm<OTPFormValues>({ resolver: zodResolver(otpSchema) });

  const { mutateAsync: doLogin } = useLogin();

  const handleLogin = async (values: LoginFormValues) => {
    setIsLoading(true);
    setErrorMsg(null);
    try {
      const data = await doLogin(values);

      if ("user" in data) {
        login(data.user);
        navigate("/dashboard");
      } else {
        setOtpSessionToken(data.otp_session_token);
        setStep("otp");
      }
    } catch (err: any) {
      setErrorMsg(getErrorMessage(err?.response?.data?.code, err?.response?.data?.message));
    } finally {
      setIsLoading(false);
    }
  };

  const { mutateAsync: doVerify } = useVerifyOTP();

  const handleVerifyOTP = async (values: OTPFormValues) => {
    setIsLoading(true);
    setErrorMsg(null);
    try {
      const data = await doVerify({ otp_session_token: otpSessionToken, otp_code: values.otp_code });
      login(data.user);
      navigate("/dashboard");
    } catch (err: any) {
      setErrorMsg(getErrorMessage(err?.response?.data?.code, err?.response?.data?.message));
    } finally {
      setIsLoading(false);
    }
  };

  const { mutate: doResend } = useResendOTP();

  const handleResendOTP = () => {
    try {
      doResend(otpSessionToken);
    } catch { /* ignore */ }
  };

  return (
    <div className="min-h-screen flex bg-surface-secondary">
      {/* Left panel */}
      <div className="hidden lg:flex lg:w-[44%] relative overflow-hidden bg-gradient-to-br from-primary-700 via-primary-600 to-danger-600 p-12 flex-col justify-between">
        <div className="absolute inset-0 opacity-[0.07]" style={{ backgroundImage: "radial-gradient(circle at 2px 2px, white 1px, transparent 0)", backgroundSize: "28px 28px" }} />
        <div className="relative flex items-center gap-2.5">
          <div className="h-8 w-8 rounded-lg bg-white/15 backdrop-blur flex items-center justify-center text-white font-bold text-[15px] border border-white/20">P</div>
          <span className="text-white font-semibold text-[15px] tracking-tight">PPMS</span>
        </div>
        <div className="relative">
          <h1 className="text-white text-[28px] font-semibold leading-tight tracking-tight mb-3 max-w-md">Kelola portofolio proyek perusahaan dalam satu tempat.</h1>
          <p className="text-white/75 text-[14px] leading-relaxed max-w-sm">Dari pengajuan, approval, hingga eksekusi — pantau setiap proyek, anggaran, dan tim Anda secara real-time.</p>
        </div>
        <div className="relative flex items-center gap-6 text-white/60 text-[12.5px]">
          <span>© 2026 PPMS Internal</span>
        </div>
      </div>

      {/* Right panel */}
      <div className="flex-1 flex items-center justify-center p-8 relative">
        <button onClick={toggleTheme} className="absolute top-6 right-6 h-9 w-9 rounded-md border border-border bg-surface flex items-center justify-center text-ink-secondary hover:bg-surface-secondary transition-colors">
          {theme === "light"
            ? <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M21 12.8A9 9 0 1111.2 3 7 7 0 0021 12.8z"/></svg>
            : <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><circle cx="12" cy="12" r="4"/><path d="M12 2v2M12 20v2M4.9 4.9l1.4 1.4M17.7 17.7l1.4 1.4M2 12h2M20 12h2M4.9 19.1l1.4-1.4M17.7 6.3l1.4-1.4"/></svg>
          }
        </button>

        <div className="w-full max-w-[360px]">
          <div className="lg:hidden flex items-center gap-2.5 mb-8">
            <div className="h-8 w-8 rounded-lg bg-gradient-to-br from-primary-600 to-danger-600 flex items-center justify-center text-white font-bold text-[15px]">P</div>
            <span className="font-semibold text-[15px] tracking-tight">PPMS</span>
          </div>

          {step === "credentials" ? (
            <>
              <h2 className="text-[20px] font-semibold tracking-tight mb-1.5">Masuk ke akun Anda</h2>
              <p className="text-[13.5px] text-ink-secondary mb-7">Gunakan kredensial perusahaan Anda untuk melanjutkan.</p>

              <form onSubmit={loginForm.handleSubmit(handleLogin)} className="space-y-4">
                <div>
                  <Label htmlFor="email">Email</Label>
                  <Input id="email" type="email" placeholder="nama@perusahaan.com" {...loginForm.register("email")} />
                  {loginForm.formState.errors.email && <p className="text-xs text-danger-600 mt-1.5">{loginForm.formState.errors.email.message}</p>}
                </div>
                <div>
                  <Label htmlFor="password">Password</Label>
                  <Input id="password" type="password" placeholder="••••••••" {...loginForm.register("password")} />
                  {loginForm.formState.errors.password && <p className="text-xs text-danger-600 mt-1.5">{loginForm.formState.errors.password.message}</p>}
                </div>
                {errorMsg && (
                  <div className="flex items-start gap-2 rounded-md bg-danger-50 dark:bg-danger-900/20 border border-danger-200 dark:border-danger-900/40 px-3 py-2.5">
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-danger-600 flex-shrink-0 mt-0.5"><circle cx="12" cy="12" r="9"/><path d="M12 8v5M12 16h.01"/></svg>
                    <p className="text-[12.5px] text-danger-700 dark:text-danger-400 m-0 leading-snug">{errorMsg}</p>
                  </div>
                )}
                <Button type="submit" variant="primary" className="w-full" size="lg" disabled={isLoading}>
                  {isLoading ? "Memproses..." : "Masuk"}
                </Button>
              </form>
            </>
          ) : (
            <>
              <button onClick={() => { setStep("credentials"); setErrorMsg(null); }} className="flex items-center gap-1.5 text-[12.5px] text-ink-secondary mb-6 hover:text-ink-primary transition-colors">
                <svg width="13" height="13" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2"><path d="M15 18l-6-6 6-6"/></svg>
                Kembali ke login
              </button>

              <div className="flex items-center justify-center w-12 h-12 rounded-full bg-primary-50 dark:bg-primary-900/30 text-primary-600 mb-5">
                <svg width="22" height="22" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2">
                  <path d="M22 16.92v3a2 2 0 01-2.18 2 19.79 19.79 0 01-8.63-3.07A19.5 19.5 0 013.07 9.81a19.79 19.79 0 01-3.07-8.68A2 2 0 012 0h3a2 2 0 012 1.72c.127.96.361 1.903.7 2.81a2 2 0 01-.45 2.11L6.09 7.91a16 16 0 006 6l1.27-1.27a2 2 0 012.11-.45c.907.339 1.85.573 2.81.7A2 2 0 0122 14.92z"/>
                </svg>
              </div>
              <h2 className="text-[20px] font-semibold tracking-tight mb-1.5">Verifikasi email Anda</h2>
              <p className="text-[13.5px] text-ink-secondary mb-7">Kode 6 digit telah dikirim ke email terdaftar Anda. Kode berlaku selama 10 menit.</p>

              <form onSubmit={otpForm.handleSubmit(handleVerifyOTP)} className="space-y-4">
                <div>
                  <Label htmlFor="otp_code">Kode OTP</Label>
                  <Input
                    id="otp_code"
                    placeholder="000000"
                    maxLength={6}
                    className="text-center text-[20px] tracking-[0.25em] font-semibold"
                    {...otpForm.register("otp_code")}
                  />
                  {otpForm.formState.errors.otp_code && <p className="text-xs text-danger-600 mt-1.5">{otpForm.formState.errors.otp_code.message}</p>}
                </div>
                {errorMsg && (
                  <div className="flex items-start gap-2 rounded-md bg-danger-50 dark:bg-danger-900/20 border border-danger-200 dark:border-danger-900/40 px-3 py-2.5">
                    <svg width="15" height="15" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-danger-600 flex-shrink-0 mt-0.5"><circle cx="12" cy="12" r="9"/><path d="M12 8v5M12 16h.01"/></svg>
                    <p className="text-[12.5px] text-danger-700 dark:text-danger-400 m-0 leading-snug">{errorMsg}</p>
                  </div>
                )}
                <Button type="submit" variant="primary" className="w-full" size="lg" disabled={isLoading}>
                  {isLoading ? "Memverifikasi..." : "Verifikasi"}
                </Button>
              </form>

              <p className="text-center text-[12.5px] text-ink-tertiary mt-5">
                Tidak menerima kode?{" "}
                <button onClick={handleResendOTP} className="text-primary-600 font-medium hover:text-primary-700">
                  Kirim ulang
                </button>
              </p>
            </>
          )}
        </div>
      </div>
    </div>
  );
}