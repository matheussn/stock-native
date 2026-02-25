export interface AppError {
  code: 'VALIDATION_ERROR' | 'NOT_FOUND' | 'INSUFFICIENT_STOCK' | 'CONFLICT' | 'INTERNAL_ERROR'
  message: string
  details?: unknown
}

export interface Category {
  id: string
  name: string
  active: boolean
  createdAt: string
  updatedAt: string
}

export interface Product {
  id: string
  name: string
  categoryId: string
  measureUnit: MeasureUnit
  packageAmount: number
  lowStockThreshold: number
  active: boolean
  createdAt: string
  updatedAt: string
}

export type MeasureUnit = 'kg' | 'g' | 'l' | 'ml' | 'unidade'

export interface Origin {
  id: string
  name: string
  active: boolean
  createdAt: string
  updatedAt: string
}

export interface Destination {
  id: string
  name: string
  active: boolean
  createdAt: string
  updatedAt: string
}

export interface Movement {
  id: string
  type: 'entrada' | 'saida'
  productId: string
  productName: string
  quantity: number
  movementDate: string
  originId?: string
  originName?: string
  destinationId?: string
  destinationName?: string
  note: string
  createdAt: string
}

export interface StockItem {
  productId: string
  productName: string
  categoryId: string
  categoryName: string
  measureUnit: MeasureUnit
  packageAmount: number
  lowStockThreshold: number
  currentStock: number
  isLowStock: boolean
}

export interface GroupedTotal {
  id: string
  name: string
  quantity: number
}

export interface Reports {
  entriesByPeriod: Movement[]
  entriesByOrigin: GroupedTotal[]
  exitsByPeriod: Movement[]
  exitsByDestination: GroupedTotal[]
  movementHistory: Movement[]
}
