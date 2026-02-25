import type { MeasureUnit } from '@/types/inventory'

export const formatNumber = (value: number): string => {
  return new Intl.NumberFormat('pt-BR', {
    minimumFractionDigits: 0,
    maximumFractionDigits: 3,
  }).format(value)
}

const measureUnitLabel: Record<MeasureUnit, string> = {
  kg: 'Kg',
  g: 'g',
  l: 'L',
  ml: 'ml',
  unidade: 'unidade',
}

export const formatMeasureUnit = (unit: MeasureUnit): string => measureUnitLabel[unit] ?? unit

export const formatProductMeasure = (amount: number, unit: MeasureUnit): string =>
  `${formatNumber(amount)}${formatMeasureUnit(unit)}`

export const toDateInputValue = (value: Date): string => {
  return value.toISOString().slice(0, 10)
}

export const startOfDayIso = (value: string): string => `${value}T00:00:00`
export const endOfDayIso = (value: string): string => `${value}T23:59:59`
