import * as React from "react";
import { cn } from "../../lib/utils";
import { Rows3, Rows2 } from "lucide-react";

export type TableDensity = "comfortable" | "compact";

const DensityContext = React.createContext<TableDensity>("comfortable");

const Table = React.forwardRef<
  HTMLTableElement,
  React.HTMLAttributes<HTMLTableElement> & { density?: TableDensity }
>(({ className, density, ...props }, ref) => (
  <div className="rounded-lg border border-border bg-surface overflow-hidden">
    <DensityContext.Provider value={density ?? "comfortable"}>
      <table ref={ref} className={cn("w-full border-collapse text-[13px]", className)} {...props} />
    </DensityContext.Provider>
  </div>
));
Table.displayName = "Table";

const TableHeader = React.forwardRef<HTMLTableSectionElement, React.HTMLAttributes<HTMLTableSectionElement>>(
  ({ className, ...props }, ref) => <thead ref={ref} className={cn("bg-surface-secondary", className)} {...props} />
);
TableHeader.displayName = "TableHeader";

const TableBody = React.forwardRef<HTMLTableSectionElement, React.HTMLAttributes<HTMLTableSectionElement>>(
  ({ className, ...props }, ref) => <tbody ref={ref} className={className} {...props} />
);
TableBody.displayName = "TableBody";

const TableRow = React.forwardRef<HTMLTableRowElement, React.HTMLAttributes<HTMLTableRowElement>>(
  ({ className, ...props }, ref) => (
    <tr
      ref={ref}
      className={cn("border-b border-border last:border-b-0 hover:bg-surface-secondary transition-colors", className)}
      {...props}
    />
  )
);
TableRow.displayName = "TableRow";

const TableHead = React.forwardRef<HTMLTableCellElement, React.ThHTMLAttributes<HTMLTableCellElement>>(
  ({ className, ...props }, ref) => {
    const density = React.useContext(DensityContext);
    return (
      <th
        ref={ref}
        className={cn(
          "text-left font-semibold text-[11.5px] uppercase tracking-wide text-ink-tertiary",
          density === "compact" ? "px-3 py-1.5" : "px-4 py-2.5",
          className
        )}
        {...props}
      />
    );
  }
);
TableHead.displayName = "TableHead";

const TableCell = React.forwardRef<HTMLTableCellElement, React.TdHTMLAttributes<HTMLTableCellElement>>(
  ({ className, ...props }, ref) => {
    const density = React.useContext(DensityContext);
    return (
      <td
        ref={ref}
        className={cn(density === "compact" ? "px-3 py-1.5" : "px-4 py-3", "align-middle", className)}
        {...props}
      />
    );
  }
);
TableCell.displayName = "TableCell";

interface DensityToggleProps {
  value: TableDensity;
  onChange: (density: TableDensity) => void;
}

export function TableDensityToggle({ value, onChange }: DensityToggleProps) {
  return (
    <div className="inline-flex items-center rounded-md border border-border bg-surface p-0.5">
      <button
        type="button"
        onClick={() => onChange("comfortable")}
        title="Tampilan nyaman"
        aria-label="Tampilan nyaman"
        className={cn(
          "h-7 w-7 flex items-center justify-center rounded transition-colors cursor-pointer",
          value === "comfortable" ? "bg-surface-secondary text-ink-primary" : "text-ink-tertiary hover:text-ink-primary"
        )}
      >
        <Rows3 className="h-3.5 w-3.5" />
      </button>
      <button
        type="button"
        onClick={() => onChange("compact")}
        title="Tampilan padat"
        aria-label="Tampilan padat"
        className={cn(
          "h-7 w-7 flex items-center justify-center rounded transition-colors cursor-pointer",
          value === "compact" ? "bg-surface-secondary text-ink-primary" : "text-ink-tertiary hover:text-ink-primary"
        )}
      >
        <Rows2 className="h-3.5 w-3.5" />
      </button>
    </div>
  );
}

export { Table, TableHeader, TableBody, TableRow, TableHead, TableCell };