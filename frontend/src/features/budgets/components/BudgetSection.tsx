import { useState } from "react";
import { useBudget, useCreateBudget, useTransactions, useCreateTransaction, useUpdateBudget } from "../hooks/useBudget";
import { Card, CardHeader, CardTitle, CardContent } from "../../../components/ui/card";
import { Button } from "../../../components/ui/button";
import { Input } from "../../../components/ui/input";
import { Select } from "@/components/ui/select";
import { StatusBadge } from "../../../components/ui/status-badge";
import { EmptyState } from "../../../components/shared/EmptyState";

interface BudgetSectionProps {
  projectId: number;
}

const formatCurrency = (value: number) =>
  `Rp ${Math.round(value || 0).toLocaleString("id-ID")}`;

export default function BudgetSection({ projectId }: BudgetSectionProps) {
  const { data: budget, isError } = useBudget(projectId);
  const { mutate: createBudget } = useCreateBudget(projectId);
  const { mutate: updateBudget } = useUpdateBudget(projectId, budget?.id ?? 0);
  const { data: transactions } = useTransactions(projectId, budget?.id ?? 0);
  const { mutate: createTransaction } = useCreateTransaction(projectId, budget?.id ?? 0);

  const [allocatedInput, setAllocatedInput] = useState("");
  const [budgetType, setBudgetType] = useState("CAPEX");
  const [budgetName, setBudgetName] = useState("");
  const [editType, setEditType] = useState("CAPEX");
  const [editName, setEditName] = useState("");

  const [txType, setTxType] = useState("EXPENSE");
  const [txAmount, setTxAmount] = useState("");
  const [txReason, setTxReason] = useState("");
  const [adjustmentType, setAdjustmentType] = useState("ERROR_CORRECTION");

  if (isError || !budget) {
    return (
      <Card>
        <CardHeader>
          <CardTitle>Budget</CardTitle>
        </CardHeader>
        <CardContent>
          <EmptyState
            icon={
              <svg width="20" height="20" viewBox="0 0 24 24" fill="none" stroke="currentColor" strokeWidth="1.8">
                <line x1="12" y1="1" x2="12" y2="23" />
                <path d="M17 5H9.5a3.5 3.5 0 000 7h5a3.5 3.5 0 010 7H6" />
              </svg>
            }
            title="Belum ada budget"
            description="Tetapkan alokasi anggaran untuk mulai melacak pengeluaran project ini."
            action={
              <div className="flex flex-col gap-2 w-full max-w-xs">
                <Select
                  value={budgetType}
                  onChange={(e) => setBudgetType(e.target.value)}
                  options={[
                    { value: "CAPEX", label: "CAPEX" },
                    { value: "OPEX", label: "OPEX" },
                  ]}
                />
                <Input
                  placeholder="Nama budget (cth. Server Infrastruktur)"
                  value={budgetName}
                  onChange={(e) => setBudgetName(e.target.value)}
                />
                <Input
                  type="number"
                  placeholder="Alokasi (Rp)"
                  value={allocatedInput}
                  onChange={(e) => setAllocatedInput(e.target.value)}
                />
                <Button
                  variant="primary"
                  onClick={() =>
                    createBudget({
                      allocated_budget: Number(allocatedInput),
                      budget_type: budgetType as any,
                      budget_name: budgetName,
                    })
                  }
                >
                  Set
                </Button>
              </div>
            }
          />
        </CardContent>
      </Card>
    );
  }

  const handleAddTransaction = () => {
    createTransaction({
      type: txType,
      adjustment_type: txType === "ADJUSTMENT" ? adjustmentType : undefined,
      amount: Number(txAmount),
      reason: txReason,
      idempotency_key: crypto.randomUUID(),
    });
    setTxAmount("");
    setTxReason("");
  };

  const saveBudgetMeta = () => {
    updateBudget({
      allocated_budget: budget.allocated_budget,
      budget_type: editType as any,
      budget_name: editName,
      version: budget.version,
    });
  };

  const usageColor =
    budget.usage_percentage >= 100 ? "#DC2626" : budget.usage_percentage >= 80 ? "#D97706" : "#2563EB";

  return (
    <Card>
      <CardHeader>
        <CardTitle>
          Budget
          {budget.budget_type && (
            <StatusBadge color={budget.budget_type === "CAPEX" ? "indigo" : "teal"}>
              {budget.budget_type}
            </StatusBadge>
          )}
        </CardTitle>
      </CardHeader>
      <CardContent>
        <div className="grid grid-cols-3 gap-4 text-sm mb-4">
          <div>
            <p className="text-[11.5px] text-ink-tertiary mb-1">Allocated</p>
            <p className="text-[15px] font-semibold m-0">{formatCurrency(budget.allocated_budget)}</p>
          </div>
          <div>
            <p className="text-[11.5px] text-ink-tertiary mb-1">Used</p>
            <p className="text-[15px] font-semibold m-0">{formatCurrency(budget.used_budget)}</p>
          </div>
          <div>
            <p className="text-[11.5px] text-ink-tertiary mb-1">Remaining</p>
            <p className="text-[15px] font-semibold m-0">{formatCurrency(budget.remaining_budget)}</p>
          </div>
        </div>

        <div className="h-2 w-full rounded-full bg-surface-tertiary overflow-hidden mb-1.5">
          <div
            className="h-full rounded-full transition-all"
            style={{ width: `${Math.min(budget.usage_percentage, 100)}%`, background: usageColor }}
          />
        </div>
        <p className="text-xs text-ink-tertiary mb-2">{budget.usage_percentage.toFixed(1)}% digunakan</p>

        {budget.budget_name && (
          <p className="text-[12.5px] text-ink-secondary mb-4">Nama budget: {budget.budget_name}</p>
        )}

        <div className="flex flex-wrap items-end gap-2 border-t border-border pt-4 mb-4">
          <div className="flex-1 min-w-[140px]">
            <p className="text-[11.5px] text-ink-tertiary mb-1">Tipe</p>
            <Select
              value={editType}
              onChange={(e) => setEditType(e.target.value)}
              options={[
                { value: "CAPEX", label: "CAPEX" },
                { value: "OPEX", label: "OPEX" },
              ]}
            />
          </div>
          <div className="flex-1 min-w-[180px]">
            <p className="text-[11.5px] text-ink-tertiary mb-1">Nama budget</p>
            <Input
              placeholder={budget.budget_name || "Nama budget"}
              value={editName}
              onChange={(e) => setEditName(e.target.value)}
            />
          </div>
          <Button
            variant="outline"
            size="sm"
            onClick={() => {
              if (!editName) setEditName(budget.budget_name);
              saveBudgetMeta();
            }}
          >
            Simpan
          </Button>
        </div>

        <div className="space-y-2.5 border-t border-border pt-5 mb-6">
          <p className="text-[13px] font-semibold mb-2">Tambah transaksi</p>
          <div className="flex gap-2">
            <select
              value={txType}
              onChange={(e) => setTxType(e.target.value)}
              className="h-9 px-2.5 rounded-md border border-border-strong bg-surface text-[13px] font-medium cursor-pointer"
            >
              <option value="EXPENSE">Expense</option>
              <option value="REFUND">Refund</option>
              <option value="ADJUSTMENT">Adjustment</option>
            </select>
            <Input type="number" placeholder="Jumlah" value={txAmount} onChange={(e) => setTxAmount(e.target.value)} className="w-36" />
          </div>

          {txType === "ADJUSTMENT" && (
            <select
              value={adjustmentType}
              onChange={(e) => setAdjustmentType(e.target.value)}
              className="w-full h-9 px-2.5 rounded-md border border-border-strong bg-surface text-[13px] font-medium cursor-pointer"
            >
              <option value="ERROR_CORRECTION">Error correction</option>
              <option value="BUDGET_REALLOCATION">Budget reallocation</option>
              <option value="AUDIT_CORRECTION">Audit correction</option>
              <option value="MANUAL_OVERRIDE">Manual override</option>
            </select>
          )}

          <Input placeholder="Alasan (wajib untuk adjustment)" value={txReason} onChange={(e) => setTxReason(e.target.value)} />

          <Button variant="primary" size="sm" onClick={handleAddTransaction}>
            Catat transaksi
          </Button>
        </div>

        <div className="space-y-2 border-t border-border pt-5">
          <p className="text-[13px] font-semibold mb-2">Riwayat transaksi</p>
          {transactions?.length === 0 && <p className="text-[12.5px] text-ink-tertiary">Belum ada transaksi.</p>}
          {transactions?.map((t) => (
            <div key={t.id} className="flex items-center justify-between py-1.5 border-b border-border last:border-b-0">
              <div className="flex items-center gap-2">
                <StatusBadge color={t.type === "EXPENSE" ? "red" : t.type === "REFUND" ? "green" : "amber"}>
                  {t.type}
                </StatusBadge>
                {t.adjustment_type && <span className="text-[11px] text-ink-tertiary">{t.adjustment_type}</span>}
              </div>
              <span className={`text-[13px] font-semibold ${t.amount < 0 ? "text-success-600" : "text-danger-600"}`}>
                Rp {Math.abs(t.amount).toLocaleString("id-ID")}
              </span>
            </div>
          ))}
        </div>
      </CardContent>
    </Card>
  );
}
