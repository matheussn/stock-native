import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { callBackend } from '@/commons/utils/wails'
import type { Movement } from '@/types/inventory'

const movementKey = ['movements'] as const

interface MovementFilter {
  type?: 'entrada' | 'saida'
  productId?: string
  dateFrom?: string
  dateTo?: string
}

export const useMovementsQuery = (filter: MovementFilter) =>
  useQuery({
    queryKey: [...movementKey, filter],
    queryFn: () =>
      callBackend<Movement[]>('ListMovements', {
        type: filter.type ?? '',
        productId: filter.productId ?? '',
        dateFrom: filter.dateFrom ? `${filter.dateFrom}T00:00:00Z` : null,
        dateTo: filter.dateTo ? `${filter.dateTo}T23:59:59Z` : null,
      }),
  })

export const useCreateEntryMovementMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      productId: string
      quantity: string
      movementDate: string
      sourceId: string
      note: string
    }) =>
      callBackend<Movement>('CreateEntryMovement', {
        productId: payload.productId,
        quantity: payload.quantity,
        movementDate: `${payload.movementDate}T00:00:00Z`,
        sourceId: payload.sourceId,
        note: payload.note,
      }),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: movementKey })
      void client.invalidateQueries({ queryKey: ['stock'] })
      void client.invalidateQueries({ queryKey: ['reports'] })
    },
  })
}

export const useCreateEntryBatchMovementMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      movementDate: string
      sourceId: string
      note: string
      items: Array<{
        productId: string
        quantity: string
      }>
    }) =>
      callBackend<Movement[]>('CreateEntryMovementsBatch', {
        movementDate: `${payload.movementDate}T00:00:00Z`,
        sourceId: payload.sourceId,
        note: payload.note,
        items: payload.items.map((item) => ({
          productId: item.productId,
          quantity: item.quantity,
        })),
      }),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: movementKey })
      void client.invalidateQueries({ queryKey: ['stock'] })
      void client.invalidateQueries({ queryKey: ['reports'] })
    },
  })
}

export const useCreateExitMovementMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      productId: string
      quantity: string
      movementDate: string
      sourceId: string
      note: string
    }) =>
      callBackend<Movement>('CreateExitMovement', {
        productId: payload.productId,
        quantity: payload.quantity,
        movementDate: `${payload.movementDate}T00:00:00Z`,
        sourceId: payload.sourceId,
        note: payload.note,
      }),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: movementKey })
      void client.invalidateQueries({ queryKey: ['stock'] })
      void client.invalidateQueries({ queryKey: ['reports'] })
    },
  })
}

export const useCreateExitBatchMovementMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      movementDate: string
      sourceId: string
      note: string
      items: Array<{
        productId: string
        quantity: string
      }>
    }) =>
      callBackend<Movement[]>('CreateExitMovementsBatch', {
        movementDate: `${payload.movementDate}T00:00:00Z`,
        sourceId: payload.sourceId,
        note: payload.note,
        items: payload.items.map((item) => ({
          productId: item.productId,
          quantity: item.quantity,
        })),
      }),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: movementKey })
      void client.invalidateQueries({ queryKey: ['stock'] })
      void client.invalidateQueries({ queryKey: ['reports'] })
    },
  })
}
