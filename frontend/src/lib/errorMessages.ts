// Mapping kode error backend ke pesan ramah pengguna dalam Bahasa Indonesia.
// Backend tetap mengirim pesan Inggris (konsisten untuk API contract),
// frontend yang menerjemahkan untuk ditampilkan ke end-user.
export const errorMessageMap: Record<string, string> = {
  INSUFFICIENT_SYSTEM_ROLE: "Anda tidak memiliki izin sistem untuk melakukan aksi ini.",
  INSUFFICIENT_PROJECT_ROLE: "Role Anda di project ini tidak mengizinkan aksi tersebut.",
  NOT_PROJECT_MEMBER: "Anda bukan anggota aktif dari project ini.",
  RESOURCE_NOT_OWNED: "Anda bukan pemilik dari resource ini.",
  PROJECT_LOCKED: "Project ini sedang terkunci untuk perubahan.",
  LAST_PM_PROTECTION: "Tidak dapat menghapus Project Manager terakhir yang aktif.",
  INVALID_STATE_TRANSITION: "Perubahan status tidak valid dari kondisi saat ini.",
  NOT_FOUND: "Data yang dicari tidak ditemukan.",
  VALIDATION_ERROR: "Data yang dimasukkan tidak valid.",
  UNAUTHORIZED: "Sesi Anda telah berakhir, silakan login kembali.",
  CONFLICT: "Data telah diubah oleh proses lain, silakan refresh halaman.",
  RATE_LIMITED: "Terlalu banyak percobaan, silakan coba lagi sebentar.",
  FILE_TOO_LARGE: "Ukuran file melebihi batas maksimal 25MB.",
  UNSUPPORTED_FILE_TYPE: "Jenis file tidak didukung.",
  DUPLICATE_ENTRY: "Data ini sudah pernah dikirim sebelumnya.",
  INTERNAL_ERROR: "Terjadi kesalahan pada server, silakan coba lagi.",
};

export function getErrorMessage(code?: string, fallback?: string): string {
  if (code && errorMessageMap[code]) {
    return errorMessageMap[code];
  }
  return fallback || "Terjadi kesalahan yang tidak diketahui.";
}