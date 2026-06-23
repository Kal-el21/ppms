interface HealthBarProps {
  /** 0-100 */
  progress: number;
  /** true jika project/task berisiko (overdue, over-budget) — mengaktifkan gradient ke merah */
  atRisk?: boolean;
  className?: string;
}

// Komponen signature design system PPMS: progress bar yang gradiennya
// bergeser dari biru (primary brand) ke merah (danger brand) ketika
// project/task ditandai berisiko. Mengikat kedua warna brand jadi
// satu elemen yang fungsional, bukan dua warna terpisah yang dipasang
// sembarangan di tempat berbeda.
export function HealthBar({ progress, atRisk = false, className = "" }: HealthBarProps) {
  const clamped = Math.max(0, Math.min(100, progress));

  const fillStyle = atRisk
    ? { width: `${clamped}%`, background: "linear-gradient(90deg, #2563EB, #DC2626)" }
    : { width: `${clamped}%`, background: "#2563EB" };

  return (
    <div className={`h-1.5 w-full rounded-full bg-surface-tertiary overflow-hidden ${className}`}>
      <div className="h-full rounded-full transition-all duration-300" style={fillStyle} />
    </div>
  );
}