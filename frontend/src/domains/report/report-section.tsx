import { yupResolver } from '@hookform/resolvers/yup'
import { useMemo } from 'react'
import { useForm } from 'react-hook-form'

import { toDateInputValue } from '@/commons/utils/format'
import { reportFilterSchema, type ReportFilterForm } from '@/commons/validators/schemas'
import { Panel } from '@/components/panel'
import { useReportsQuery } from '@/domains/report/api'

const today = toDateInputValue(new Date())
const firstDay = toDateInputValue(new Date(new Date().setDate(new Date().getDate() - 30)))

const exportCsv = (rows: Array<Record<string, string | number>>, filename: string) => {
  const headers = Object.keys(rows[0] ?? {})
  const content = [headers.join(','), ...rows.map((row) => headers.map((h) => JSON.stringify(row[h] ?? '')).join(','))].join('\n')
  const blob = new Blob([content], { type: 'text/csv;charset=utf-8;' })
  const link = document.createElement('a')
  link.href = URL.createObjectURL(blob)
  link.download = filename
  link.click()
}

export const ReportSection = () => {
  const form = useForm<ReportFilterForm>({
    resolver: yupResolver(reportFilterSchema),
    defaultValues: { dateFrom: firstDay, dateTo: today },
  })

  const dateFrom = form.watch('dateFrom')
  const dateTo = form.watch('dateTo')

  const reportQuery = useReportsQuery(dateFrom, dateTo)

  const summary = useMemo(() => {
    const data = reportQuery.data
    if (!data) {
      return { entries: 0, exits: 0 }
    }

    const entries = data.entriesByPeriod.reduce((acc, item) => acc + item.quantity, 0)
    const exits = data.exitsByPeriod.reduce((acc, item) => acc + item.quantity, 0)
    return { entries, exits }
  }, [reportQuery.data])

  return (
    <Panel title='Relatórios' subtitle='Entradas, saídas e agrupamentos por origem/destino'>
      <form className='filters'>
        <input type='date' {...form.register('dateFrom')} />
        <input type='date' {...form.register('dateTo')} />
      </form>

      <div className='grid three'>
        <article className='stat'>
          <h3>Total Entradas</h3>
          <strong>{summary.entries.toFixed(3)}</strong>
        </article>
        <article className='stat'>
          <h3>Total Saídas</h3>
          <strong>{summary.exits.toFixed(3)}</strong>
        </article>
        <article className='stat'>
          <h3>Movimentações</h3>
          <strong>{reportQuery.data?.movementHistory.length ?? 0}</strong>
        </article>
      </div>

      <div className='grid two'>
        <div className='card'>
          <header className='card-header'>
            <h3>Entradas por Origem</h3>
            <button
              type='button'
              onClick={() =>
                exportCsv(
                  (reportQuery.data?.entriesByOrigin ?? []).map((item) => ({ origem: item.name, quantidade: item.quantity })),
                  'entradas-por-origem.csv',
                )
              }
            >
              Exportar CSV
            </button>
          </header>
          <table className='table'>
            <thead>
              <tr>
                <th>Origem</th>
                <th>Quantidade</th>
              </tr>
            </thead>
            <tbody>
              {(reportQuery.data?.entriesByOrigin ?? []).map((item) => (
                <tr key={item.id}>
                  <td>{item.name}</td>
                  <td>{item.quantity}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>

        <div className='card'>
          <header className='card-header'>
            <h3>Saídas por Destino</h3>
            <button
              type='button'
              onClick={() =>
                exportCsv(
                  (reportQuery.data?.exitsByDestination ?? []).map((item) => ({ destino: item.name, quantidade: item.quantity })),
                  'saidas-por-destino.csv',
                )
              }
            >
              Exportar CSV
            </button>
          </header>
          <table className='table'>
            <thead>
              <tr>
                <th>Destino</th>
                <th>Quantidade</th>
              </tr>
            </thead>
            <tbody>
              {(reportQuery.data?.exitsByDestination ?? []).map((item) => (
                <tr key={item.id}>
                  <td>{item.name}</td>
                  <td>{item.quantity}</td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      </div>
    </Panel>
  )
}