import { useState, useCallback } from "react";

export function useTableSelection<T extends number | string>() {
  const [selectedIds, setSelectedIds] = useState<Set<T>>(new Set());

  const toggle = useCallback((id: T) => {
    setSelectedIds((prev) => {
      const next = new Set(prev);
      if (next.has(id)) {
        next.delete(id);
      } else {
        next.add(id);
      }
      return next;
    });
  }, []);

  const toggleAll = useCallback((ids: T[], isChecked: boolean) => {
    if (isChecked) {
      setSelectedIds(new Set(ids));
    } else {
      setSelectedIds(new Set());
    }
  }, []);

  const clear = useCallback(() => {
    setSelectedIds(new Set());
  }, []);

  const isSelected = useCallback(
    (id: T) => {
      return selectedIds.has(id);
    },
    [selectedIds]
  );

  const selectedCount = selectedIds.size;

  const isAllSelected = useCallback(
    (ids: T[]) => {
      return ids.length > 0 && ids.every((id) => selectedIds.has(id));
    },
    [selectedIds]
  );

  const isIndeterminate = useCallback(
    (ids: T[]) => {
      const count = ids.filter((id) => selectedIds.has(id)).length;
      return count > 0 && count < ids.length;
    },
    [selectedIds]
  );

  return {
    selectedIds,
    toggle,
    toggleAll,
    clear,
    isSelected,
    isAllSelected,
    isIndeterminate,
    selectedCount,
  };
}
