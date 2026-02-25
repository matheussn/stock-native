import { yupResolver } from '@hookform/resolvers/yup'
import { useState } from 'react'
import { useForm } from 'react-hook-form'

import { formatProductMeasure } from '@/commons/utils/format'
import { Message, Panel } from '@/components/panel'
import {
  useCategoriesQuery,
  useCreateCategoryMutation,
  useCreateDestinationMutation,
  useCreateOriginMutation,
  useCreateProductMutation,
  useDeactivateCategoryMutation,
  useDeactivateDestinationMutation,
  useDeactivateOriginMutation,
  useDeactivateProductMutation,
  useDestinationsQuery,
  useOriginsQuery,
  useProductsQuery,
} from '@/domains/catalog/api'
import { namedEntitySchema, productSchema, type NamedEntityForm, type ProductForm } from '@/commons/validators/schemas'
import { useUiStore } from '@/stores/ui.store'

const defaultNamedEntity: NamedEntityForm = {
  name: '',
}

const defaultProduct: ProductForm = {
  name: '',
  categoryId: '',
  measureUnit: 'kg',
  packageAmount: '',
  lowStockThreshold: '0',
}

export const CatalogSection = () => {
  const [message, setMessage] = useState('')
  const screen = useUiStore((state) => state.catalogScreen)
  const setCatalogScreen = useUiStore((state) => state.setCatalogScreen)

  const categoriesQuery = useCategoriesQuery()
  const productsQuery = useProductsQuery()
  const originsQuery = useOriginsQuery()
  const destinationsQuery = useDestinationsQuery()

  const createCategory = useCreateCategoryMutation()
  const deactivateCategory = useDeactivateCategoryMutation()
  const createProduct = useCreateProductMutation()
  const deactivateProduct = useDeactivateProductMutation()
  const createOrigin = useCreateOriginMutation()
  const deactivateOrigin = useDeactivateOriginMutation()
  const createDestination = useCreateDestinationMutation()
  const deactivateDestination = useDeactivateDestinationMutation()

  const categoryForm = useForm<NamedEntityForm>({
    resolver: yupResolver(namedEntitySchema),
    defaultValues: defaultNamedEntity,
  })

  const productForm = useForm<ProductForm>({
    resolver: yupResolver(productSchema),
    defaultValues: defaultProduct,
  })

  const originForm = useForm<NamedEntityForm>({
    resolver: yupResolver(namedEntitySchema),
    defaultValues: defaultNamedEntity,
  })

  const destinationForm = useForm<NamedEntityForm>({
    resolver: yupResolver(namedEntitySchema),
    defaultValues: defaultNamedEntity,
  })

  const subtitles = {
    categorias: 'Cadastro e gerenciamento de categorias',
    produtos: 'Cadastro e gerenciamento de produtos',
    origens: 'Cadastro e gerenciamento de origens',
    destinos: 'Cadastro e gerenciamento de destinos',
  }

  return (
    <Panel title='Cadastros' subtitle={subtitles[screen]}>
      <Message text={message} />
      <div className='tabs'>
        <button type='button' onClick={() => setCatalogScreen('categorias')} className={screen === 'categorias' ? 'active' : ''}>
          Categorias
        </button>
        <button type='button' onClick={() => setCatalogScreen('produtos')} className={screen === 'produtos' ? 'active' : ''}>
          Produtos
        </button>
        <button type='button' onClick={() => setCatalogScreen('origens')} className={screen === 'origens' ? 'active' : ''}>
          Origens
        </button>
        <button type='button' onClick={() => setCatalogScreen('destinos')} className={screen === 'destinos' ? 'active' : ''}>
          Destinos
        </button>
      </div>

      {screen === 'categorias' ? (
        <div className='card'>
          <h3>Categorias</h3>
          <form
            onSubmit={categoryForm.handleSubmit(async (values) => {
              try {
                await createCategory.mutateAsync({ name: values.name })
                categoryForm.reset(defaultNamedEntity)
                setMessage('Categoria criada')
              } catch {
                setMessage('Erro ao criar categoria')
              }
            })}
            className='form'
          >
            <input placeholder='Nome da categoria' {...categoryForm.register('name')} />
            <button type='submit'>Criar</button>
          </form>
          <ul className='list'>
            {(categoriesQuery.data ?? []).map((item) => (
              <li key={item.id}>
                <span>{item.name}</span>
                <button type='button' onClick={() => void deactivateCategory.mutateAsync(item.id)} disabled={!item.active}>
                  {item.active ? 'Desativar' : 'Inativo'}
                </button>
              </li>
            ))}
          </ul>
        </div>
      ) : null}

      {screen === 'produtos' ? (
        <div className='card'>
          <h3>Produtos</h3>
          <form
            onSubmit={productForm.handleSubmit(async (values) => {
              try {
                await createProduct.mutateAsync(values)
                productForm.reset(defaultProduct)
                setMessage('Produto criado')
              } catch {
                setMessage('Erro ao criar produto')
              }
            })}
            className='form'
          >
            <input placeholder='Nome do produto' {...productForm.register('name')} />
            <select {...productForm.register('categoryId')}>
              <option value=''>Categoria</option>
              {(categoriesQuery.data ?? []).filter((item) => item.active).map((item) => (
                <option key={item.id} value={item.id}>
                  {item.name}
                </option>
              ))}
            </select>
            <select {...productForm.register('measureUnit')}>
              <option value='kg'>Kg</option>
              <option value='g'>g</option>
              <option value='l'>L</option>
              <option value='ml'>ml</option>
              <option value='unidade'>unidade</option>
            </select>
            <input placeholder='Peso/volume por item' {...productForm.register('packageAmount')} />
            <input placeholder='Limite de estoque baixo' {...productForm.register('lowStockThreshold')} />
            <button type='submit'>Criar</button>
          </form>
          <ul className='list'>
            {(productsQuery.data ?? []).map((item) => (
              <li key={item.id}>
                <span>{`${item.name} - ${formatProductMeasure(item.packageAmount, item.measureUnit)}`}</span>
                <button type='button' onClick={() => void deactivateProduct.mutateAsync(item.id)} disabled={!item.active}>
                  {item.active ? 'Desativar' : 'Inativo'}
                </button>
              </li>
            ))}
          </ul>
        </div>
      ) : null}

      {screen === 'origens' ? (
        <div className='card'>
          <h3>Origens</h3>
          <form
            onSubmit={originForm.handleSubmit(async (values) => {
              try {
                await createOrigin.mutateAsync(values)
                originForm.reset(defaultNamedEntity)
                setMessage('Origem criada')
              } catch {
                setMessage('Erro ao criar origem')
              }
            })}
            className='form'
          >
            <input placeholder='Nome da origem' {...originForm.register('name')} />
            <button type='submit'>Criar</button>
          </form>
          <ul className='list'>
            {(originsQuery.data ?? []).map((item) => (
              <li key={item.id}>
                <span>{item.name}</span>
                <button type='button' onClick={() => void deactivateOrigin.mutateAsync(item.id)} disabled={!item.active}>
                  {item.active ? 'Desativar' : 'Inativo'}
                </button>
              </li>
            ))}
          </ul>
        </div>
      ) : null}

      {screen === 'destinos' ? (
        <div className='card'>
          <h3>Destinos</h3>
          <form
            onSubmit={destinationForm.handleSubmit(async (values) => {
              try {
                await createDestination.mutateAsync(values)
                destinationForm.reset(defaultNamedEntity)
                setMessage('Destino criado')
              } catch {
                setMessage('Erro ao criar destino')
              }
            })}
            className='form'
          >
            <input placeholder='Nome do destino' {...destinationForm.register('name')} />
            <button type='submit'>Criar</button>
          </form>
          <ul className='list'>
            {(destinationsQuery.data ?? []).map((item) => (
              <li key={item.id}>
                <span>{item.name}</span>
                <button type='button' onClick={() => void deactivateDestination.mutateAsync(item.id)} disabled={!item.active}>
                  {item.active ? 'Desativar' : 'Inativo'}
                </button>
              </li>
            ))}
          </ul>
        </div>
      ) : null}
    </Panel>
  )
}
