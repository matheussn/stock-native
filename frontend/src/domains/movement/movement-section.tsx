import { useMemo, useState } from 'react'
import { useFieldArray, useForm } from 'react-hook-form'

import { formatProductMeasure, toDateInputValue } from '@/commons/utils/format'
import { Message, Panel } from '@/components/panel'
import { useDestinationsQuery, useOriginsQuery, useProductsQuery } from '@/domains/catalog/api'
import {
  useCreateEntryBatchMovementMutation,
  useCreateExitBatchMovementMutation,
  useMovementsQuery,
} from '@/domains/movement/api'

const initialDate = toDateInputValue(new Date())
const quantityPattern = /^\d+(\.\d{1,3})?$/

type MovementView = 'history' | 'create'

interface BatchMovementForm {
  sourceId: string
  movementDate: string
  note: string
  items: Array<{
    productId: string
    quantity: string
  }>
}

const defaultBatchValues: BatchMovementForm = {
  sourceId: '',
  movementDate: initialDate,
  note: '',
  items: [{ productId: '', quantity: '' }],
}

export const MovementSection = () => {
  const [message, setMessage] = useState('')
  const [messageKind, setMessageKind] = useState<'success' | 'error'>('success')
  const [view, setView] = useState<MovementView>('history')
  const [type, setType] = useState<'entrada' | 'saida'>('entrada')

  const productsQuery = useProductsQuery()
  const originsQuery = useOriginsQuery()
  const destinationsQuery = useDestinationsQuery()

  const movementQuery = useMovementsQuery({})
  const createEntryBatch = useCreateEntryBatchMovementMutation()
  const createExitBatch = useCreateExitBatchMovementMutation()

  const form = useForm<BatchMovementForm>({
    defaultValues: defaultBatchValues,
  })

  const { fields, append, remove } = useFieldArray({
    control: form.control,
    name: 'items',
  })

  const activeProducts = useMemo(
    () => (productsQuery.data ?? []).filter((item) => item.active),
    [productsQuery.data],
  )

  const sourceOptions = useMemo(() => {
    if (type === 'entrada') {
      return (originsQuery.data ?? []).filter((item) => item.active)
    }
    return (destinationsQuery.data ?? []).filter((item) => item.active)
  }, [destinationsQuery.data, originsQuery.data, type])

  const handleBackToHistory = () => {
    form.reset({ ...defaultBatchValues, movementDate: initialDate })
    setView('history')
  }

  const handleSubmit = form.handleSubmit(async (values) => {
    const validItems = values.items.filter((item) => item.productId.trim() && item.quantity.trim())

    if (!values.sourceId) {
      setMessageKind('error')
      setMessage(`Selecione ${type === 'entrada' ? 'uma origem' : 'um destino'}`)
      return
    }

    if (!values.movementDate) {
      setMessageKind('error')
      setMessage('Informe a data da movimentação')
      return
    }

    if (validItems.length === 0) {
      setMessageKind('error')
      setMessage(`Adicione ao menos um produto para ${type}`)
      return
    }

    if (!validItems.every((item) => quantityPattern.test(item.quantity))) {
      setMessageKind('error')
      setMessage('Use número positivo com até 3 casas decimais em todas as quantidades')
      return
    }

    try {
      if (type === 'entrada') {
        await createEntryBatch.mutateAsync({
          sourceId: values.sourceId,
          movementDate: values.movementDate,
          note: values.note,
          items: validItems,
        })
      } else {
        await createExitBatch.mutateAsync({
          sourceId: values.sourceId,
          movementDate: values.movementDate,
          note: values.note,
          items: validItems,
        })
      }

      setMessageKind('success')
      setMessage(`${type === 'entrada' ? 'Entrada' : 'Saída'} registrada para ${validItems.length} produto(s)`)
      handleBackToHistory()
    } catch (error) {
      setMessageKind('error')

      if (typeof error === 'object' && error !== null && 'message' in error) {
        const rawMessage = String((error as { message?: unknown }).message ?? '')
        if (rawMessage) {
          setMessage(rawMessage)
          return
        }
      }

      setMessage(`Erro ao registrar ${type === 'entrada' ? 'entrada' : 'saída'}`)
    }
  })

  if (view === 'create') {
    return (
      <Panel
        title='Nova movimentação'
        subtitle='Escolha entre entrada ou saída e preencha os dados'
        actions={
          <button type='button' className='movement-cancel-button' onClick={handleBackToHistory}>
            Voltar ao histórico
          </button>
        }
      >
        <Message text={message} kind={messageKind} />
        <div className='card'>
          <div className='tabs'>
            <button type='button' onClick={() => setType('entrada')} className={type === 'entrada' ? 'active' : ''}>
              Entrada
            </button>
            <button type='button' onClick={() => setType('saida')} className={type === 'saida' ? 'active' : ''}>
              Saída
            </button>
          </div>

          <form onSubmit={handleSubmit} className='form'>
            <input type='date' {...form.register('movementDate')} />
            <select {...form.register('sourceId')}>
              <option value=''>Selecione {type === 'entrada' ? 'origem' : 'destino'}</option>
              {sourceOptions.map((item) => (
                <option key={item.id} value={item.id}>
                  {item.name}
                </option>
              ))}
            </select>
            <textarea placeholder='Observação (opcional)' {...form.register('note')} rows={3} />

            <div className='entry-items'>
              {fields.map((field, index) => (
                <div key={field.id} className='entry-item-row'>
                  <select {...form.register(`items.${index}.productId` as const)}>
                    <option value=''>Produto</option>
                    {activeProducts.map((item) => (
                      <option key={item.id} value={item.id}>
                        {`${item.name} - ${formatProductMeasure(item.packageAmount, item.measureUnit)}`}
                      </option>
                    ))}
                  </select>
                  <input placeholder='Quantidade' {...form.register(`items.${index}.quantity` as const)} />
                  <button
                    type='button'
                    className='entry-remove-button'
                    onClick={() => remove(index)}
                    disabled={fields.length <= 1}
                  >
                    Remover
                  </button>
                </div>
              ))}
            </div>

            <button
              type='button'
              className='entry-add-button movement-cancel-button'
              onClick={() => append({ productId: '', quantity: '' })}
            >
              + Adicionar produto
            </button>

            <button type='submit'>Salvar {type}</button>
          </form>
        </div>
      </Panel>
    )
  }

  return (
    <Panel title='Movimentações' subtitle='Histórico geral de entradas e saídas'>
      <Message text={message} kind={messageKind} />
      <div className='card'>
        <div className='history-toolbar'>
          <h3>Histórico</h3>
          <button
            type='button'
            className='movement-add-button'
            onClick={() => {
              setType('entrada')
              setView('create')
            }}
          >
            +
          </button>
        </div>

        <table className='table'>
          <thead>
            <tr>
              <th>Data</th>
              <th>Tipo</th>
              <th>Produto</th>
              <th>Quantidade</th>
              <th>Origem/Destino</th>
            </tr>
          </thead>
          <tbody>
            {(movementQuery.data ?? []).map((item) => (
              <tr key={item.id}>
                <td>{new Date(item.movementDate).toLocaleDateString('pt-BR')}</td>
                <td>{item.type}</td>
                <td>{item.productName}</td>
                <td>{item.quantity}</td>
                <td>{item.originName ?? item.destinationName ?? '-'}</td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>
    </Panel>
  )
}
