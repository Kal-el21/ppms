import * as React from "react";

interface ChartDatum {
  label: string;
  value: number;
  color: string;
  value2?: number;
  color2?: string;
}

interface DoughnutChartProps {
  data: ChartDatum[];
  size?: number;
  thickness?: number;
}

const emptyDoughnut = (size: number, label: string) => (
  <div className="flex flex-col items-center justify-center" style={{ width: size, height: size }}>
    <span className="text-[12px] text-ink-tertiary">{label}</span>
  </div>
);

export function DoughnutChart({ data, size = 160, thickness = 22 }: DoughnutChartProps) {
  const total = data.reduce((sum, d) => sum + d.value, 0);
  const radius = (size - thickness) / 2;
  const circumference = 2 * Math.PI * radius;
  const center = size / 2;

  const segments = React.useMemo(() => {
    const offsets: number[] = [];
    let cumulative = 0;
    for (let i = 0; i < data.length; i++) {
      offsets.push(cumulative);
      cumulative += (data[i].value / total) * circumference;
    }
    return data.map((d, i) => ({
      label: d.label,
      dash: (d.value / total) * circumference,
      segOffset: -offsets[i],
      color: d.color,
      value: d.value,
    }));
  }, [data, total, circumference]);

  if (total === 0) {
    return emptyDoughnut(size, "Tidak ada data");
  }

  return (
    <div className="flex flex-col items-center gap-3">
      <svg width={size} height={size} viewBox={`0 0 ${size} ${size}`} className="-rotate-90">
        <circle cx={center} cy={center} r={radius} fill="none" stroke="var(--surface-tertiary)" strokeWidth={thickness} />
        {segments.map((s, i) => (
          <circle
            key={i}
            cx={center}
            cy={center}
            r={radius}
            fill="none"
            stroke={s.color}
            strokeWidth={thickness}
            strokeDasharray={`${s.dash} ${circumference - s.dash}`}
            strokeDashoffset={s.segOffset}
            strokeLinecap="butt"
          />
        ))}
      </svg>
      <div className="flex flex-wrap justify-center gap-x-3 gap-y-1">
        {data.map((d, i) => (
          <span key={i} className="flex items-center gap-1.5 text-[11.5px] text-ink-secondary">
            <span className="h-2 w-2 rounded-sm" style={{ background: d.color }} />
            {d.label} ({d.value})
          </span>
        ))}
      </div>
    </div>
  );
}

interface BarChartProps {
  data: ChartDatum[];
  height?: number;
}

export function BarChart({ data, height = 180 }: BarChartProps) {
  const max = Math.max(1, ...data.map((d) => d.value));

  if (data.length === 0) {
    return <p className="text-[12.5px] text-ink-tertiary">Tidak ada data.</p>;
  }

  return (
    <div style={{ height: height }} className="flex flex-col justify-center gap-2">
      {data.map((d, i) => (
        <div key={i} className="flex items-center gap-2">
          <span className="text-[11.5px] text-ink-secondary w-[96px] truncate flex-shrink-0">{d.label}</span>
          <div className="flex-1 h-3 rounded-full bg-surface-tertiary overflow-hidden">
            <div
              className="h-full rounded-full"
              style={{ width: `${(d.value / max) * 100}%`, background: d.color }}
            />
          </div>
          <span className="text-xs font-semibold text-ink-secondary w-8 text-right">{d.value}</span>
        </div>
      ))}
    </div>
  );
}

interface GroupedBarChartProps {
  data: ChartDatum[];
  height?: number;
}

export function GroupedBarChart({ data, height = 180 }: GroupedBarChartProps) {
  const max = Math.max(1, ...data.map((d) => Math.max(d.value, d.value2 ?? 0)));

  if (data.length === 0) {
    return <p className="text-[12.5px] text-ink-tertiary">Tidak ada data.</p>;
  }

  return (
    <div style={{ height: height }} className="flex flex-col justify-center gap-2">
      {data.map((d, i) => {
        const val1 = d.value;
        const val2 = d.value2 ?? 0;
        return (
          <div key={i} className="flex items-center gap-2">
            <span className="text-[11.5px] text-ink-secondary w-[96px] truncate flex-shrink-0">{d.label}</span>
            <div className="flex-1 flex gap-1">
              <div className="h-3 rounded-full bg-surface-tertiary overflow-hidden flex-1">
                <div
                  className="h-full rounded-full"
                  style={{ width: `${(val1 / max) * 100}%`, background: d.color }}
                />
              </div>
              {val2 > 0 && (
                <div className="h-3 rounded-full bg-surface-tertiary overflow-hidden flex-1">
                  <div
                    className="h-full rounded-full"
                    style={{ width: `${(val2 / max) * 100}%`, background: d.color2 ?? d.color }}
                  />
                </div>
              )}
            </div>
            <span className="text-[10px] font-semibold text-ink-secondary w-16 text-right">
              {val1 > 0 && `Rp ${Math.round(val1).toLocaleString("id-ID")}`}
              {val1 > 0 && val2 > 0 && " / "}
              {val2 > 0 && `Rp ${Math.round(val2).toLocaleString("id-ID")}`}
            </span>
          </div>
        );
      })}
    </div>
  );
}
