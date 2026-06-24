import { useState, useRef } from "react";
import { useAuth } from "../../auth/context/AuthContext";
import axiosInstance from "../../../api/axiosInstance";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Label } from "../../../components/ui/label";
import { Card, CardContent, CardHeader, CardTitle } from "../../../components/ui/card";
import { Avatar } from "../../../components/ui/avatar";
import { useToast } from "../../../components/ui/toast";

export default function ProfileSettings() {
  const { user, updateUser } = useAuth(); // bukan login
  const toast = useToast();
  const fileInputRef = useRef<HTMLInputElement>(null);

  const [fullName, setFullName] = useState(user?.full_name || "");
  const [isSaving, setIsSaving] = useState(false);
  const [isUploading, setIsUploading] = useState(false);

  const handleSave = async () => {
    setIsSaving(true);
    try {
      const res = await axiosInstance.put("/me", { full_name: fullName });
      updateUser({ full_name: res.data.data.full_name });
      toast.success("Profile berhasil disimpan");
    } catch (err: any) {
      toast.error("Gagal menyimpan profile", err?.friendlyMessage);
    } finally {
      setIsSaving(false);
    }
  };

  const handlePhotoUpload = async (e: React.ChangeEvent<HTMLInputElement>) => {
    const file = e.target.files?.[0];
    if (!file) return;
    setIsUploading(true);
    try {
      const formData = new FormData();
      formData.append("photo", file);
      const res = await axiosInstance.post("/me/photo", formData, {
        headers: { "Content-Type": "multipart/form-data" },
      });
      updateUser({ profile_photo_url: res.data.data.photo_url });
      toast.success("Foto profile berhasil diupdate");
    } catch (err: any) {
      toast.error("Gagal upload foto", err?.friendlyMessage);
    } finally {
      setIsUploading(false);
      e.target.value = "";
    }
  };

  return (
    <div className="space-y-4">
      <Card>
        <CardHeader><CardTitle>Foto Profile</CardTitle></CardHeader>
        <CardContent>
          <div className="flex items-center gap-4">
            <div className="relative">
              {user?.profile_photo_url ? (
                <img
                  src={user.profile_photo_url}
                  alt={user.full_name}
                  className="h-16 w-16 rounded-full object-cover ring-2 ring-border"
                />
              ) : (
                <div className="h-16 w-16">
                  <Avatar name={user?.full_name || "?"} size="md" />
                </div>
              )}
            </div>
            <div>
              <input
                ref={fileInputRef}
                type="file"
                accept="image/jpeg,image/png,image/webp"
                className="hidden"
                onChange={handlePhotoUpload}
              />
              <Button
                variant="outline"
                size="sm"
                onClick={() => fileInputRef.current?.click()}
                disabled={isUploading}
              >
                {isUploading ? "Mengupload..." : "Ganti foto"}
              </Button>
              <p className="text-[11.5px] text-ink-tertiary mt-1.5">JPG, PNG, WebP. Maks 5MB.</p>
            </div>
          </div>
        </CardContent>
      </Card>

      <Card>
        <CardHeader><CardTitle>Informasi Dasar</CardTitle></CardHeader>
        <CardContent className="space-y-4">
          <div>
            <Label>Nama lengkap</Label>
            <Input value={fullName} onChange={(e) => setFullName(e.target.value)} />
          </div>
          <div>
            <Label>Email</Label>
            <Input value={user?.email} disabled className="opacity-60 cursor-not-allowed" />
            <p className="text-[11.5px] text-ink-tertiary mt-1.5">Email tidak bisa diubah sendiri. Hubungi Admin.</p>
          </div>
          <Button variant="primary" onClick={handleSave} disabled={isSaving}>
            {isSaving ? "Menyimpan..." : "Simpan perubahan"}
          </Button>
        </CardContent>
      </Card>
    </div>
  );
}