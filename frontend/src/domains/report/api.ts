import { useQuery } from '@tanstack/react-query'

import { callBackend } from '@/commons/utils/wails'
import type { Reports } from '@/types/inventory'

export const useReportsQuery = (dateFrom: string, dateTo: string) =>
  useQuery({
    queryKey: ['reports', dateFrom, dateTo],
    queryFn: () =>
      callBackend<Reports>('GetReports', {
        dateFrom: `${dateFrom}T00:00:00Z`,
        dateTo: `${dateTo}T23:59:59Z`,
      }),
    enabled: Boolean(dateFrom && dateTo),
  })