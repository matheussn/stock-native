import { useMutation, useQuery, useQueryClient } from '@tanstack/react-query'

import { callBackend } from '@/commons/utils/wails'
import type { Category, Destination, Origin, Product } from '@/types/inventory'

const keys = {
  categories: ['categories'] as const,
  products: ['products'] as const,
  origins: ['origins'] as const,
  destinations: ['destinations'] as const,
}

export const useCategoriesQuery = () =>
  useQuery({
    queryKey: keys.categories,
    queryFn: () => callBackend<Category[]>('ListCategories'),
  })

export const useProductsQuery = () =>
  useQuery({
    queryKey: keys.products,
    queryFn: () => callBackend<Product[]>('ListProducts'),
  })

export const useOriginsQuery = () =>
  useQuery({
    queryKey: keys.origins,
    queryFn: () => callBackend<Origin[]>('ListOrigins'),
  })

export const useDestinationsQuery = () =>
  useQuery({
    queryKey: keys.destinations,
    queryFn: () => callBackend<Destination[]>('ListDestinations'),
  })

export const useCreateCategoryMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: { name: string }) => callBackend<Category>('CreateCategory', payload),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.categories })
    },
  })
}

export const useUpdateCategoryMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: { id: string; name: string }) => callBackend<Category>('UpdateCategory', payload),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.categories })
    },
  })
}

export const useDeactivateCategoryMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => callBackend<void>('DeactivateCategory', id),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.categories })
      void client.invalidateQueries({ queryKey: keys.products })
      void client.invalidateQueries({ queryKey: ['stock'] })
    },
  })
}

export const useCreateProductMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      name: string
      categoryId: string
      measureUnit: 'kg' | 'g' | 'l' | 'ml' | 'unidade'
      packageAmount: string
      lowStockThreshold: string
    }) =>
      callBackend<Product>('CreateProduct', payload),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.products })
      void client.invalidateQueries({ queryKey: ['stock'] })
    },
  })
}

export const useUpdateProductMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: {
      id: string
      name: string
      categoryId: string
      measureUnit: 'kg' | 'g' | 'l' | 'ml' | 'unidade'
      packageAmount: string
      lowStockThreshold: string
    }) => callBackend<Product>('UpdateProduct', payload),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.products })
      void client.invalidateQueries({ queryKey: ['stock'] })
    },
  })
}

export const useDeactivateProductMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => callBackend<void>('DeactivateProduct', id),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.products })
      void client.invalidateQueries({ queryKey: ['stock'] })
    },
  })
}

export const useCreateOriginMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: { name: string }) => callBackend<Origin>('CreateOrigin', payload),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.origins })
    },
  })
}

export const useUpdateOriginMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: { id: string; name: string }) =>
      callBackend<Origin>('UpdateOrigin', payload),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.origins })
    },
  })
}

export const useDeactivateOriginMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => callBackend<void>('DeactivateOrigin', id),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.origins })
    },
  })
}

export const useCreateDestinationMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: { name: string }) =>
      callBackend<Destination>('CreateDestination', payload),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.destinations })
    },
  })
}

export const useUpdateDestinationMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (payload: { id: string; name: string }) =>
      callBackend<Destination>('UpdateDestination', payload),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.destinations })
    },
  })
}

export const useDeactivateDestinationMutation = () => {
  const client = useQueryClient()
  return useMutation({
    mutationFn: (id: string) => callBackend<void>('DeactivateDestination', id),
    onSuccess: () => {
      void client.invalidateQueries({ queryKey: keys.destinations })
    },
  })
}
