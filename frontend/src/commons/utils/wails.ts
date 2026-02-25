import type { AppError } from '@/types/inventory'

type WailsApp = Record<string, (...args: any[]) => Promise<any>>

declare global {
  interface Window {
    go?: {
      app?: { App?: WailsApp }
      main?: { App?: WailsApp }
    }
  }
}

const getApp = (): WailsApp => {
  const app = window.go?.app?.App ?? window.go?.main?.App
  if (!app) {
    throw new Error('Wails bridge not available')
  }
  return app
}

export const callBackend = async <T>(method: string, ...args: unknown[]): Promise<T> => {
  const app = getApp()
  const fn = app[method]
  if (typeof fn !== 'function') {
    throw new Error(`Backend method not found: ${method}`)
  }

  try {
    return (await fn(...args)) as T
  } catch (error) {
    throw normalizeError(error)
  }
}

const normalizeError = (error: unknown): AppError => {
  if (typeof error === 'object' && error !== null && 'code' in error && 'message' in error) {
    return error as AppError
  }
  if (typeof error === 'string') {
    return { code: 'INTERNAL_ERROR', message: error }
  }
  return { code: 'INTERNAL_ERROR', message: 'Erro inesperado' }
}