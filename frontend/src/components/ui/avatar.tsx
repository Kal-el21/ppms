interface AvatarProps {
  name: string;
  size?: "sm" | "md";
  colorSeed?: number;
}

const colors = ["#2563EB", "#16A34A", "#D97706", "#DC2626", "#7C3AED", "#0EA5E9"];

function getInitials(name: string) {
  const parts = name.trim().split(" ");
  if (parts.length === 1) return parts[0].slice(0, 2).toUpperCase();
  return (parts[0][0] + parts[parts.length - 1][0]).toUpperCase();
}

function getColor(name: string) {
  const hash = name.split("").reduce((acc, ch) => acc + ch.charCodeAt(0), 0);
  return colors[hash % colors.length];
}

export function Avatar({ name, size = "md" }: AvatarProps) {
  const dim = size === "sm" ? "h-6 w-6 text-[10px]" : "h-9 w-9 text-[13px]";
  return (
    <div
      className={`${dim} rounded-full flex items-center justify-center font-semibold text-white flex-shrink-0 ring-2 ring-surface`}
      style={{ background: getColor(name) }}
      title={name}
    >
      {getInitials(name)}
    </div>
  );
}

export function AvatarStack({ names, max = 3 }: { names: string[]; max?: number }) {
  const visible = names.slice(0, max);
  const remaining = names.length - max;

  return (
    <div className="flex items-center">
      {visible.map((name, i) => (
        <div key={i} className={i > 0 ? "-ml-2" : ""}>
          <Avatar name={name} size="sm" />
        </div>
      ))}
      {remaining > 0 && (
        <div className="-ml-2 h-6 w-6 rounded-full bg-surface-tertiary text-ink-secondary text-[10px] font-semibold flex items-center justify-center ring-2 ring-surface">
          +{remaining}
        </div>
      )}
    </div>
  );
}