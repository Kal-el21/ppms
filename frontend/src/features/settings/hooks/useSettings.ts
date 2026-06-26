import { useMutation } from "@tanstack/react-query";
import { settingsApi } from "../api/settingsApi";
import { authApi } from "../../auth/api/authApi";
import { useAuth } from "../../auth/context/AuthContext";
import { useToast } from "../../../components/ui/toast";

export function useUpdateProfile() {
  const { updateUser } = useAuth();
  const toast = useToast();

  return useMutation({
    mutationFn: settingsApi.updateProfile,
    onSuccess: (user) => {
      updateUser({ full_name: user.full_name });
      toast.success("Profile berhasil disimpan");
    },
    onError: (error: any) => {
      toast.error("Gagal menyimpan profile", error?.friendlyMessage);
    },
  });
}

export function useUploadProfilePhoto() {
  const { updateUser } = useAuth();
  const toast = useToast();

  return useMutation({
    mutationFn: settingsApi.uploadProfilePhoto,
    onSuccess: ({ photo_url }) => {
      updateUser({ profile_photo_url: photo_url });
      toast.success("Foto profile berhasil diupdate");
    },
    onError: (error: any) => {
      toast.error("Gagal upload foto", error?.friendlyMessage);
    },
  });
}

export function useToggle2FA() {
  const { updateUser } = useAuth();
  const toast = useToast();

  return useMutation({
    mutationFn: settingsApi.toggle2FA,
    onSuccess: (_, enabled) => {
      updateUser({ two_fa_enabled: enabled });
      toast.success(
        enabled
          ? "Two-factor authentication diaktifkan. Kode OTP akan dikirim ke email saat login."
          : "Two-factor authentication dinonaktifkan."
      );
    },
    onError: (error: any) => {
      toast.error("Gagal mengubah pengaturan 2FA", error?.friendlyMessage);
    },
  });
}

export function useToggleEmailNotification() {
  const { updateUser } = useAuth();
  const toast = useToast();

  return useMutation({
    mutationFn: settingsApi.toggleEmailNotification,
    onSuccess: (_, enabled) => {
      updateUser({ email_notification_enabled: enabled });
      toast.success(enabled ? "Notifikasi email diaktifkan" : "Notifikasi email dinonaktifkan");
    },
    onError: (error: any) => {
      toast.error("Gagal mengubah preferensi email", error?.friendlyMessage);
    },
  });
}

export function useChangePassword() {
  const toast = useToast();

  return useMutation({
    mutationFn: ({ oldPassword, newPassword }: { oldPassword: string; newPassword: string }) =>
      authApi.changePassword(oldPassword, newPassword),
    onSuccess: () => {
      toast.success("Password berhasil diubah");
    },
    onError: (error: any) => {
      toast.error("Gagal mengubah password", error?.friendlyMessage);
    },
  });
}
