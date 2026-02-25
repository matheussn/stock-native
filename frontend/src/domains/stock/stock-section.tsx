import { useMemo, useState } from 'react'

import { Panel } from '@/components/panel'
import { useCategoriesQuery } from '@/domains/catalog/api'
import { useStockQuery } from '@/domains/stock/api'
import { formatNumber, formatProductMeasure } from '@/commons/utils/format'

export const StockSection = () => {
  const [categoryId, setCategoryId] = useState('')
  const [search, setSearch] = useState('')

  const categoriesQuery = useCategoriesQuery()
  const stockQuery = useStockQuery({ categoryId, search })

  const items = useMemo(() => stockQuery.data ?? [], [stockQuery.data])

  return (
    <Panel title='Estoque Atual' subtitle='Saldo por produto com alerta de baixo estoque'>
      <div className='filters'>
        <input
          value={search}
          onChange={(event) => setSearch(event.target.value)}
          placeholder='Buscar produto'
        />
        <select value={categoryId} onChange={(event) => setCategoryId(event.target.value)}>
          <option value=''>Todas as categorias</option>
          {(categoriesQuery.data ?? []).map((category) => (
            <option key={category.id} value={category.id}>
              {category.name}
            </option>
          ))}
        </select>
      </div>

      <table className='table'>
        <thead>
          <tr>
            <th>Produto</th>
            <th>Categoria</th>
            <th>Embalagem</th>
            <th>Saldo</th>
            <th>Status</th>
          </tr>
        </thead>
        <tbody>
          {items.map((item) => (
            <tr key={item.productId}>
              <td>{item.productName}</td>
              <td>{item.categoryName}</td>
              <td>{formatProductMeasure(item.packageAmount, item.measureUnit)}</td>
              <td>{formatNumber(item.currentStock)}</td>
              <td>
                {item.isLowStock ? (
                  <span className='badge warning'>Baixo</span>
                ) : (
                  <span className='badge ok'>OK</span>
                )}
              </td>
            </tr>
          ))}
        </tbody>
      </table>
    </Panel>
  )
}
