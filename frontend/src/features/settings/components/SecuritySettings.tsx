import { useState } from "react";
import { useAuth } from "../../auth/context/AuthContext";
import { useNavigate } from "react-router-dom";
import { useToggle2FA, useChangePassword } from "../hooks/useSettings";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";
import { useToast } from "../../../components/ui/toast";

function Toggle({ checked, onChange, disabled }: { checked: boolean; onChange: (v: boolean) => void; disabled?: boolean }) {
  return (
    <button
      type="button"
      onClick={() => !disabled && onChange(!checked)}
      disabled={disabled}
      className={`relative h-6 w-11 rounded-full transition-colors flex-shrink-0 ${checked ? "bg-primary-600" : "bg-surface-tertiary"} ${disabled ? "opacity-50 cursor-not-allowed" : "cursor-pointer"}`}
    >
      <span className={`absolute top-0.5 h-5 w-5 rounded-full bg-white shadow-sm transition-transform ${checked ? "translate-x-[21px]" : "translate-x-[3px]"}`} />
    </button>
  );
}

export default function SecuritySettings() {
  const { user, logout } = useAuth();
  const navigate = useNavigate();
  const toast = useToast();
  const { mutate: toggle2FA, isPending: toggling2FA } = useToggle2FA();
  const { mutate: changePassword } = useChangePassword();

  const [changingPassword, setChangingPassword] = useState(false);
  const [showChangePassword, setShowChangePassword] = useState(false);
  const [passwords, setPasswords] = useState({ old: "", new: "", confirm: "" });

  const handleToggle2FA = async (enabled: boolean) => {
    toggle2FA(enabled);
  };

  const handleChangePassword = () => {
    if (passwords.new !== passwords.confirm) {
      toast.error("Konfirmasi password tidak cocok");
      return;
    }
    if (passwords.new.length < 8) {
      toast.error("Password baru minimal 8 karakter");
      return;
    }
    setChangingPassword(true);
    changePassword(
      { oldPassword: passwords.old, newPassword: passwords.new },
      {
        onSuccess: () => {
          setTimeout(() => {
            logout();
            navigate("/login");
          }, 1200);
        },
        onSettled: () => setChangingPassword(false),
      }
    );
  };

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader><CardTitle>Two-Factor Authentication</CardTitle></CardHeader>
        <CardContent>
          <div className="flex items-start justify-between gap-4">
            <div className="min-w-0">
              <p className="text-[13px] font-medium m-0">Login dengan verifikasi email</p>
              <p className="text-[11.5px] text-ink-secondary mt-1 m-0 leading-relaxed max-w-md">
                Saat diaktifkan, Anda akan menerima kode OTP 6 digit ke email{" "}
                <span className="font-medium">{user?.email}</span> setiap kali login.
                Kode berlaku selama 10 menit.
              </p>
            </div>
            <Toggle
              checked={user?.two_fa_enabled ?? false}
              onChange={handleToggle2FA}
              disabled={toggling2FA}
            />
          </div>

          {user?.two_fa_enabled && (
            <div className="mt-4 flex items-start gap-2 rounded-md bg-primary-50 dark:bg-primary-900/20 border border-primary-200 dark:border-primary-900/40 px-3 py-2.5">
              <svg width="14" height="14" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="2" className="text-primary-600 flex-shrink-0 mt-0.5">
                <circle cx="12" cy="12" r="9"/><path d="M12 8v5M12 16h.01"/>
              </svg>
              <p className="text-[12px] text-primary-700 dark:text-primary-400 m-0 leading-relaxed">
                2FA aktif. OTP akan dikirim ke email Anda saat login berikutnya.
                Pastikan email terdaftar aktif dan bisa diakses.
              </p>
            </div>
          )}
        </CardContent>
      </Card>

      <Card>
        <CardHeader>
          <CardTitle>Ubah Password</CardTitle>
        </CardHeader>
        <CardContent>
          {!showChangePassword ? (
            <div>
              <p className="text-[13px] text-ink-secondary mb-3">
                Setelah mengubah password, semua sesi aktif lain akan direvoke dan Anda perlu login ulang.
              </p>
              <Button variant="outline" onClick={() => setShowChangePassword(true)}>
                Ubah password
              </Button>
            </div>
          ) : (
            <div className="space-y-3">
              <div>
                <Label>Password saat ini</Label>
                <Input type="password" placeholder="••••••••" value={passwords.old} onChange={(e) => setPasswords({ ...passwords, old: e.target.value })} />
              </div>
              <div>
                <Label>Password baru</Label>
                <Input type="password" placeholder="Min. 8 karakter" value={passwords.new} onChange={(e) => setPasswords({ ...passwords, new: e.target.value })} />
              </div>
              <div>
                <Label>Konfirmasi password baru</Label>
                <Input type="password" placeholder="Ulangi password baru" value={passwords.confirm} onChange={(e) => setPasswords({ ...passwords, confirm: e.target.value })} />
              </div>
              <div className="flex gap-2 pt-1">
                <Button variant="danger" onClick={handleChangePassword} disabled={changingPassword}>
                  {changingPassword ? "Menyimpan..." : "Ubah password"}
                </Button>
                <Button variant="outline" onClick={() => { setShowChangePassword(false); setPasswords({ old: "", new: "", confirm: "" }); }}>
                  Batal
                </Button>
              </div>
            </div>
          )}
        </CardContent>
      </Card>
    </div>
  );
}
