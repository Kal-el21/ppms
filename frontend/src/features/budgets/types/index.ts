export type TransactionType = "EXPENSE" | "REFUND" | "ADJUSTMENT";
export type AdjustmentType = "ERROR_CORRECTION" | "BUDGET_REALLOCATION" | "AUDIT_CORRECTION" | "MANUAL_OVERRIDE";

export interface Budget {
  id: number;
  project_id: number;
  budget_type: BudgetType | null;
  budget_name: string;
  allocated_budget: number;
  used_budget: number;
  remaining_budget: number;
  usage_percentage: number;
  version: number;
}

export interface PortfolioBudgetYear {
  id: number;
  year: number;
  capex_ceiling: number;
  opex_ceiling: number;
  version: number;
  created_at: string;
  updated_at: string;
}

export type BudgetType = "CAPEX" | "OPEX";

export interface Transaction {
  id: number;
  budget_id: number;
  type: TransactionType;
  adjustment_type: AdjustmentType | null;
  amount: number;
  reason: string;
  description: string;
  transaction_date: string;
  created_by: number;
  created_at: string;
}