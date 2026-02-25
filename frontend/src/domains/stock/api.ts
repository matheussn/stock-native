import { useQuery } from '@tanstack/react-query'

import { callBackend } from '@/commons/utils/wails'
import type { StockItem } from '@/types/inventory'

interface StockFilter {
  categoryId?: string
  search?: string
}

export const useStockQuery = (filter: StockFilter) =>
  useQuery({
    queryKey: ['stock', filter],
    queryFn: () =>
      callBackend<StockItem[]>('GetCurrentStock', {
        categoryId: filter.categoryId ?? '',
        search: filter.search ?? '',
      }),
  })