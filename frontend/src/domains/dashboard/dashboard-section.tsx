import { useMemo, useState } from 'react'
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Legend,
  Line,
  LineChart,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis,
} from 'recharts'

import { formatNumber, toDateInputValue } from '@/commons/utils/format'
import { Panel } from '@/components/panel'
import { useCategoriesQuery, useDestinationsQuery, useOriginsQuery, useProductsQuery } from '@/domains/catalog/api'
import { useMovementsQuery } from '@/domains/movement/api'
import { useReportsQuery } from '@/domains/report/api'
import { useStockQuery } from '@/domains/stock/api'
import type { Movement } from '@/types/inventory'

type RangeKey = '12m' | 'ano' | '6m' | '3m'

type MonthlyChartData = {
  rows: Array<Record<string, string | number>>
  keys: string[]
}

const chartPalette = ['#2F4FDD', '#6F83D6', '#0F9C8C', '#2E90FA', '#17B26A', '#F79009', '#EF4444', '#475569']

const toIsoDate = (daysAgo: number): string => {
  const date = new Date()
  date.setDate(date.getDate() - daysAgo)
  return toDateInputValue(date)
}

const today = toIsoDate(0)
const last14Days = toIsoDate(13)

const getRangeStartDate = (range: RangeKey): string => {
  const now = new Date()

  if (range === 'ano') {
    return toDateInputValue(new Date(now.getFullYear(), 0, 1))
  }

  const monthsBack = range === '12m' ? 11 : range === '6m' ? 5 : 2
  return toDateInputValue(new Date(now.getFullYear(), now.getMonth() - monthsBack, 1))
}

const monthLabel = (monthKey: string): string => {
  const [year, month] = monthKey.split('-')
  return `${month}/${year.slice(-2)}`
}

const buildMonthKeys = (from: string, to: string): string[] => {
  const [fromYear, fromMonth] = from.slice(0, 7).split('-').map(Number)
  const [toYear, toMonth] = to.slice(0, 7).split('-').map(Number)

  const cursor = new Date(fromYear, fromMonth - 1, 1)
  const end = new Date(toYear, toMonth - 1, 1)
  const keys: string[] = []

  while (cursor <= end) {
    const y = cursor.getFullYear()
    const m = String(cursor.getMonth() + 1).padStart(2, '0')
    keys.push(`${y}-${m}`)
    cursor.setMonth(cursor.getMonth() + 1)
  }

  return keys
}

const buildMonthlySeriesByCategoryForSource = (
  movements: Movement[],
  fromDate: string,
  toDate: string,
  sourceId: string,
  getMovementSourceId: (item: Movement) => string | undefined,
  getCategoryNameByProductId: (productId: string) => string,
): MonthlyChartData => {
  const months = buildMonthKeys(fromDate, toDate)
  const monthCategoryMap = new Map<string, Record<string, number>>()
  const totals = new Map<string, number>()

  months.forEach((month) => monthCategoryMap.set(month, {}))

  movements.forEach((item) => {
    if (sourceId && getMovementSourceId(item) !== sourceId) {
      return
    }

    const monthKey = item.movementDate.slice(0, 7)
    const monthData = monthCategoryMap.get(monthKey)
    if (!monthData) {
      return
    }

    const category = getCategoryNameByProductId(item.productId)
    monthData[category] = (monthData[category] ?? 0) + item.quantity
    totals.set(category, (totals.get(category) ?? 0) + item.quantity)
  })

  const keys = [...totals.entries()].sort((a, b) => b[1] - a[1]).map(([name]) => name)
  const rows = months.map((month) => {
    const row: Record<string, string | number> = { mes: monthLabel(month) }
    const values = monthCategoryMap.get(month) ?? {}

    keys.forEach((key) => {
      row[key] = values[key] ?? 0
    })

    return row
  })

  return { rows, keys }
}

export const DashboardSection = () => {
  const [range, setRange] = useState<RangeKey>('12m')
  const [selectedOriginId, setSelectedOriginId] = useState('')
  const [selectedDestinationId, setSelectedDestinationId] = useState('')

  const rangeStartDate = useMemo(() => getRangeStartDate(range), [range])

  const categoriesQuery = useCategoriesQuery()
  const originsQuery = useOriginsQuery()
  const destinationsQuery = useDestinationsQuery()
  const productsQuery = useProductsQuery()
  const stockQuery = useStockQuery({})
  const reportsQuery = useReportsQuery(rangeStartDate, today)
  const movementQuery = useMovementsQuery({ dateFrom: last14Days, dateTo: today })

  const stock = stockQuery.data ?? []
  const reports = reportsQuery.data
  const movements = movementQuery.data ?? []
  const categories = categoriesQuery.data ?? []
  const origins = (originsQuery.data ?? []).filter((item) => item.active)
  const destinations = (destinationsQuery.data ?? []).filter((item) => item.active)
  const products = productsQuery.data ?? []

  const productCategoryMap = useMemo(() => {
    const categoryById = new Map(categories.map((item) => [item.id, item.name]))
    const map = new Map<string, string>()

    products.forEach((item) => {
      map.set(item.id, categoryById.get(item.categoryId) ?? 'Sem categoria')
    })

    return map
  }, [categories, products])

  const totalStock = useMemo(() => stock.reduce((acc, item) => acc + item.currentStock, 0), [stock])
  const lowStockCount = useMemo(() => stock.filter((item) => item.isLowStock).length, [stock])

  const entriesTotal = useMemo(
    () => (reports?.entriesByPeriod ?? []).reduce((acc, item) => acc + item.quantity, 0),
    [reports?.entriesByPeriod],
  )
  const exitsTotal = useMemo(
    () => (reports?.exitsByPeriod ?? []).reduce((acc, item) => acc + item.quantity, 0),
    [reports?.exitsByPeriod],
  )

  const topProducts = useMemo(
    () =>
      [...stock]
        .sort((a, b) => b.currentStock - a.currentStock)
        .slice(0, 6),
    [stock],
  )

  const entriesByCategoryMonthly = useMemo(
    () =>
      buildMonthlySeriesByCategoryForSource(
        reports?.entriesByPeriod ?? [],
        rangeStartDate,
        today,
        selectedOriginId,
        (item) => item.originId,
        (productId) => productCategoryMap.get(productId) ?? 'Sem categoria',
      ),
    [productCategoryMap, rangeStartDate, reports?.entriesByPeriod, selectedOriginId],
  )

  const exitsByCategoryMonthly = useMemo(
    () =>
      buildMonthlySeriesByCategoryForSource(
        reports?.exitsByPeriod ?? [],
        rangeStartDate,
        today,
        selectedDestinationId,
        (item) => item.destinationId,
        (productId) => productCategoryMap.get(productId) ?? 'Sem categoria',
      ),
    [productCategoryMap, rangeStartDate, reports?.exitsByPeriod, selectedDestinationId],
  )

  const last14DaysSeries = useMemo(() => {
    const map = new Map<string, number>()
    for (let offset = 13; offset >= 0; offset -= 1) {
      const dateKey = toIsoDate(offset)
      map.set(dateKey, 0)
    }

    movements.forEach((item) => {
      const key = item.movementDate.slice(0, 10)
      map.set(key, (map.get(key) ?? 0) + item.quantity)
    })

    return Array.from(map.entries()).map(([date, total]) => ({ date, total }))
  }, [movements])

  const entryExitData = [
    { name: 'Entradas', value: entriesTotal },
    { name: 'Saídas', value: exitsTotal },
  ]

  return (
    <Panel title='Painel' subtitle='Visão geral de estoque e movimentações dos últimos dias'>
      <div className='dashboard-range tabs'>
        <button type='button' className={range === '12m' ? 'active' : ''} onClick={() => setRange('12m')}>
          12 meses
        </button>
        <button type='button' className={range === 'ano' ? 'active' : ''} onClick={() => setRange('ano')}>
          No ano
        </button>
        <button type='button' className={range === '6m' ? 'active' : ''} onClick={() => setRange('6m')}>
          6 meses
        </button>
        <button type='button' className={range === '3m' ? 'active' : ''} onClick={() => setRange('3m')}>
          3 meses
        </button>
      </div>

      <div className='grid three'>
        <article className='stat'>
          <h3>Saldo em estoque</h3>
          <strong>{formatNumber(totalStock)}</strong>
        </article>
        <article className='stat'>
          <h3>Produtos com alerta</h3>
          <strong>{lowStockCount}</strong>
        </article>
        <article className='stat'>
          <h3>Movimentações (período)</h3>
          <strong>{reports?.movementHistory.length ?? 0}</strong>
        </article>
      </div>

      <div className='dashboard-grid'>
        <div className='card'>
          <h3>Top produtos por saldo</h3>
          <div className='dashboard-chart'>
            <ResponsiveContainer width='100%' height='100%'>
              <BarChart
                data={topProducts.map((item) => ({ name: item.productName, saldo: item.currentStock }))}
                margin={{ top: 12, right: 12, left: 0, bottom: 0 }}
              >
                <CartesianGrid strokeDasharray='3 3' stroke='#E5E7EB' />
                <XAxis dataKey='name' tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} />
                <Tooltip formatter={(value) => formatNumber(Number(value ?? 0))} />
                <Bar dataKey='saldo' fill='var(--color-brand)' radius={[6, 6, 0, 0]} />
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className='card'>
          <h3>Entradas vs Saídas (período)</h3>
          <div className='dashboard-donut-wrap'>
            <div className='dashboard-pie'>
              <ResponsiveContainer width='100%' height='100%'>
                <PieChart>
                  <Pie data={entryExitData} dataKey='value' nameKey='name' innerRadius={45} outerRadius={72}>
                    <Cell fill='var(--color-success)' />
                    <Cell fill='var(--color-warning)' />
                  </Pie>
                  <Tooltip formatter={(value) => formatNumber(Number(value ?? 0))} />
                  <Legend />
                </PieChart>
              </ResponsiveContainer>
            </div>
            <div className='dashboard-legend'>
              <p>
                <span className='dot success' /> Entradas: <strong>{formatNumber(entriesTotal)}</strong>
              </p>
              <p>
                <span className='dot warning' /> Saídas: <strong>{formatNumber(exitsTotal)}</strong>
              </p>
            </div>
          </div>
        </div>
      </div>

      <div className='dashboard-grid'>
        <div className='card'>
          <h3>Entradas por categoria, por mês</h3>
          <div className='dashboard-filter'>
            <label htmlFor='entry-origin-select'>Origem</label>
            <select id='entry-origin-select' value={selectedOriginId} onChange={(event) => setSelectedOriginId(event.target.value)}>
              <option value=''>Todas as origens</option>
              {origins.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.name}
                </option>
              ))}
            </select>
          </div>
          <div className='dashboard-chart dashboard-chart-tall'>
            <ResponsiveContainer width='100%' height='100%'>
              <BarChart data={entriesByCategoryMonthly.rows} margin={{ top: 12, right: 12, left: 0, bottom: 0 }}>
                <CartesianGrid strokeDasharray='3 3' stroke='#E5E7EB' />
                <XAxis dataKey='mes' tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} label={{ value: 'quantidade', angle: -90, position: 'insideLeft' }} />
                <Tooltip formatter={(value) => formatNumber(Number(value ?? 0))} />
                <Legend />
                {entriesByCategoryMonthly.keys.map((key, index) => (
                  <Bar key={key} dataKey={key} stackId='entrada' fill={chartPalette[index % chartPalette.length]} />
                ))}
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>

        <div className='card'>
          <h3>Saídas por categoria, por mês</h3>
          <div className='dashboard-filter'>
            <label htmlFor='exit-destination-select'>Destino</label>
            <select id='exit-destination-select' value={selectedDestinationId} onChange={(event) => setSelectedDestinationId(event.target.value)}>
              <option value=''>Todos os destinos</option>
              {destinations.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.name}
                </option>
              ))}
            </select>
          </div>
          <div className='dashboard-chart dashboard-chart-tall'>
            <ResponsiveContainer width='100%' height='100%'>
              <BarChart data={exitsByCategoryMonthly.rows} margin={{ top: 12, right: 12, left: 0, bottom: 0 }}>
                <CartesianGrid strokeDasharray='3 3' stroke='#E5E7EB' />
                <XAxis dataKey='mes' tick={{ fontSize: 11 }} />
                <YAxis tick={{ fontSize: 11 }} label={{ value: 'quantidade', angle: -90, position: 'insideLeft' }} />
                <Tooltip formatter={(value) => formatNumber(Number(value ?? 0))} />
                <Legend />
                {exitsByCategoryMonthly.keys.map((key, index) => (
                  <Bar key={key} dataKey={key} stackId='saida' fill={chartPalette[index % chartPalette.length]} />
                ))}
              </BarChart>
            </ResponsiveContainer>
          </div>
        </div>
      </div>

      <div className='card'>
        <h3>Volume diário de movimentações (14 dias)</h3>
        <div className='dashboard-chart'>
          <ResponsiveContainer width='100%' height='100%'>
            <LineChart
              data={last14DaysSeries.map((item) => ({ dia: item.date.slice(5), volume: item.total }))}
              margin={{ top: 12, right: 12, left: 0, bottom: 0 }}
            >
              <CartesianGrid strokeDasharray='3 3' stroke='#E5E7EB' />
              <XAxis dataKey='dia' tick={{ fontSize: 11 }} />
              <YAxis tick={{ fontSize: 11 }} />
              <Tooltip formatter={(value) => formatNumber(Number(value ?? 0))} />
              <Line type='monotone' dataKey='volume' stroke='var(--color-accent-acao)' strokeWidth={2} dot={false} />
            </LineChart>
          </ResponsiveContainer>
        </div>
        <div className='dashboard-line-footer'>
          <span>{last14DaysSeries[0]?.date ?? '-'}</span>
          <span>{last14DaysSeries[last14DaysSeries.length - 1]?.date ?? '-'}</span>
        </div>
      </div>
    </Panel>
  )
}
