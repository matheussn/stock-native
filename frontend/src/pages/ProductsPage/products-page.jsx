import {useEffect, useMemo, useState} from 'react';
import {
  CreateProduct,
  CreateProductVariation,
  ListProductVariationsByProduct,
  ListProducts,
  SetProductActive,
  SetProductVariationActive,
  UpdateProduct,
  UpdateProductVariation
} from '../../../wailsjs/go/main/App';
import AssistentialWorkFilters from '../../components/AssistentialWorkFilters';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import {getVariationLabel} from '../../utils/variation-label';
import './products-page.css';

const FILTERS = {
  active: 'Ativos',
  inactive: 'Inativos',
  all: 'Todos'
};

const BASE_UNITS = [
  {value: 'g', label: 'Grama (g)'},
  {value: 'kg', label: 'Quilograma (kg)'},
  {value: 'ml', label: 'Mililitro (ml)'},
  {value: 'L', label: 'Litro (L)'},
  {value: 'un', label: 'Unidade (un)'}
];

export default function ProductsPage() {
  const [products, setProducts] = useState([]);
  const [filter, setFilter] = useState('active');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [togglingActive, setTogglingActive] = useState(false);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [productToDeactivate, setProductToDeactivate] = useState(null);
  const [form, setForm] = useState({
    id: null,
    name: '',
    baseUnit: 'g',
    description: ''
  });
  const [draftVariations, setDraftVariations] = useState([]);
  const [variationDraft, setVariationDraft] = useState({
    baseQuantity: ''
  });
  const [productVariations, setProductVariations] = useState([]);
  const [loadingVariations, setLoadingVariations] = useState(false);
  const [variationSubmitting, setVariationSubmitting] = useState(false);
  const [variationToDeactivate, setVariationToDeactivate] = useState(null);

  const isEditing = form.id !== null;

  const visibleProducts = useMemo(() => {
    if (filter === 'inactive') {
      return products.filter((item) => item.is_active === 0);
    }
    return products;
  }, [filter, products]);

  useEffect(() => {
    void loadProducts(filter);
  }, [filter]);

  async function loadProducts(currentFilter) {
    setLoading(true);
    setError('');
    try {
      const includeInactive = currentFilter !== 'active';
      const response = await ListProducts(includeInactive);
      setProducts(response);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  function mapVariationForEdit(item) {
    return {
      ...item,
      draftBaseQuantity: String(item.base_quantity)
    };
  }

  async function loadVariations(productId) {
    setLoadingVariations(true);
    try {
      const response = await ListProductVariationsByProduct(productId, true);
      setProductVariations(response.map(mapVariationForEdit));
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoadingVariations(false);
    }
  }

  function resetForm() {
    setForm({
      id: null,
      name: '',
      baseUnit: 'g',
      description: ''
    });
  }

  function resetVariationDraft() {
    setVariationDraft({
      baseQuantity: ''
    });
  }

  function openCreateModal() {
    resetForm();
    resetVariationDraft();
    setDraftVariations([]);
    setProductVariations([]);
    setVariationToDeactivate(null);
    setError('');
    setIsModalOpen(true);
  }

  async function startEdit(item) {
    setForm({
      id: item.id,
      name: item.name,
      baseUnit: item.base_unit,
      description: item.description ?? ''
    });
    setDraftVariations([]);
    setProductVariations([]);
    resetVariationDraft();
    setVariationToDeactivate(null);
    setError('');
    setIsModalOpen(true);
    await loadVariations(item.id);
  }

  function closeModal() {
    setIsModalOpen(false);
    resetForm();
    setDraftVariations([]);
    setProductVariations([]);
    resetVariationDraft();
    setVariationToDeactivate(null);
  }

  async function handleAddVariation() {
    const baseQuantity = Number(variationDraft.baseQuantity);

    if (!Number.isInteger(baseQuantity) || baseQuantity <= 0) {
      setError('A quantidade base deve ser um inteiro maior que zero.');
      return;
    }

    if (isEditing) {
      setVariationSubmitting(true);
      setError('');
      try {
        await CreateProductVariation(form.id, '', baseQuantity);
        resetVariationDraft();
        await loadVariations(form.id);
      } catch (err) {
        setError(getFriendlyErrorMessage(err));
      } finally {
        setVariationSubmitting(false);
      }
      return;
    }

    setDraftVariations((prev) => [...prev, {baseQuantity}]);
    resetVariationDraft();
    setError('');
  }

  function removeDraftVariation(index) {
    setDraftVariations((prev) => prev.filter((_, currentIndex) => currentIndex !== index));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    const name = form.name.trim();

    if (!name) {
      setError('Informe o nome do produto.');
      return;
    }

    setSubmitting(true);
    setError('');
    try {
      if (isEditing) {
        await UpdateProduct(form.id, name, form.baseUnit, form.description);
      } else {
        const created = await CreateProduct(name, form.baseUnit, form.description);
        for (const variation of draftVariations) {
          await CreateProductVariation(
            created.id,
            '',
            variation.baseQuantity
          );
        }
      }

      closeModal();
      await loadProducts(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  function updateVariationDraft(id, field, value) {
    setProductVariations((prev) =>
      prev.map((item) => (item.id === id ? {...item, [field]: value} : item))
    );
  }

  async function handleSaveVariation(variation) {
    const baseQuantity = Number(variation.draftBaseQuantity);

    if (!Number.isInteger(baseQuantity) || baseQuantity <= 0) {
      setError('A quantidade base da variação deve ser um inteiro maior que zero.');
      return;
    }

    setVariationSubmitting(true);
    setError('');
    try {
      await UpdateProductVariation(variation.id, '', baseQuantity);
      await loadVariations(form.id);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setVariationSubmitting(false);
    }
  }

  async function handleToggleVariationActive(variation) {
    if (variation.is_active === 1) {
      setVariationToDeactivate(variation);
      return;
    }

    setVariationSubmitting(true);
    setError('');
    try {
      await SetProductVariationActive(variation.id, true);
      await loadVariations(form.id);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setVariationSubmitting(false);
    }
  }

  async function confirmDeactivateVariation() {
    if (!variationToDeactivate) {
      return;
    }

    setVariationSubmitting(true);
    setError('');
    try {
      await SetProductVariationActive(variationToDeactivate.id, false);
      setVariationToDeactivate(null);
      await loadVariations(form.id);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setVariationSubmitting(false);
    }
  }

  async function handleToggleActive(item) {
    if (item.is_active === 1) {
      setProductToDeactivate(item);
      return;
    }

    setError('');
    try {
      setTogglingActive(true);
      await SetProductActive(item.id, true);
      await loadProducts(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setTogglingActive(false);
    }
  }

  async function confirmDeactivate() {
    if (!productToDeactivate) {
      return;
    }

    setError('');
    try {
      setTogglingActive(true);
      await SetProductActive(productToDeactivate.id, false);
      setProductToDeactivate(null);
      await loadProducts(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setTogglingActive(false);
    }
  }

  return (
    <main className='products-page'>
      <section className='panel panel-list'>
        <div className='list-header'>
          <div>
            <h1 className='title'>Produtos</h1>
            <p className='subtitle'>Cadastro e manutenção</p>
          </div>
          <button type='button' className='btn btn-primary' onClick={openCreateModal}>
            + Novo
          </button>
        </div>

        <div className='list-tools'>
          <AssistentialWorkFilters filters={FILTERS} currentFilter={filter} onChange={setFilter} />
        </div>

        {error && <p className='feedback error'>{error}</p>}

        {loading ? (
          <p className='feedback'>Carregando produtos...</p>
        ) : visibleProducts.length === 0 ? (
          <p className='feedback'>Nenhum registro encontrado para este filtro.</p>
        ) : (
          <ul className='product-list'>
            {visibleProducts.map((item) => (
              <li key={item.id} className='product-item'>
                <div>
                  <strong>{item.name}</strong>
                  <p>Unidade base: {item.base_unit}</p>
                  <p>{item.description || 'Sem descrição'}</p>
                </div>

                <div className='item-actions'>
                  <span className={`status ${item.is_active === 1 ? 'status-on' : 'status-off'}`}>
                    {item.is_active === 1 ? 'Ativo' : 'Inativo'}
                  </span>
                  <button type='button' className='btn btn-ghost' onClick={() => void startEdit(item)}>
                    Editar
                  </button>
                  <button type='button' className='btn btn-danger' onClick={() => handleToggleActive(item)}>
                    {item.is_active === 1 ? 'Desativar' : 'Reativar'}
                  </button>
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>

      {isModalOpen && (
        <div className='modal-overlay' onClick={closeModal}>
          <section className='modal-content panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>{isEditing ? 'Editar produto' : 'Novo produto'}</h2>
              <button type='button' className='btn btn-ghost' onClick={closeModal} disabled={submitting}>
                Fechar
              </button>
            </div>

            <form className='form' onSubmit={handleSubmit}>
              <label htmlFor='product-name'>Nome</label>
              <input
                id='product-name'
                value={form.name}
                onChange={(event) => setForm((prev) => ({...prev, name: event.target.value}))}
                placeholder='Ex: Arroz'
                maxLength={120}
                disabled={submitting}
              />

              <label htmlFor='product-base-unit'>Unidade base</label>
              <select
                id='product-base-unit'
                value={form.baseUnit}
                onChange={(event) => setForm((prev) => ({...prev, baseUnit: event.target.value}))}
                disabled={submitting}
              >
                {BASE_UNITS.map((unit) => (
                  <option key={unit.value} value={unit.value}>
                    {unit.label}
                  </option>
                ))}
              </select>

              <label htmlFor='product-description'>Descrição</label>
              <textarea
                id='product-description'
                value={form.description}
                onChange={(event) => setForm((prev) => ({...prev, description: event.target.value}))}
                placeholder='Opcional'
                rows={3}
                maxLength={600}
                disabled={submitting}
              />

              <div className='variation-setup'>
                <h3>{isEditing ? 'Variações cadastradas' : 'Variações iniciais'}</h3>
                <p className='hint'>
                  {isEditing
                    ? 'Adicione novas variações e ajuste as existentes nesta mesma modal.'
                    : 'Adicione as variações antes de criar o produto.'}
                </p>
                <p className='stock-readonly-note'>
                  {isEditing
                    ? 'O estoque atual é somente leitura nesta tela. Use Movimentações para entradas e saídas.'
                    : 'Novas variações começam com estoque 0. O saldo é ajustado apenas em Movimentações.'}
                </p>

                <div className='variation-form'>
                  <div className='input-with-unit'>
                    <input
                      type='number'
                      min='1'
                      step='1'
                      value={variationDraft.baseQuantity}
                      onChange={(event) =>
                        setVariationDraft((prev) => ({...prev, baseQuantity: event.target.value}))
                      }
                      placeholder='Quantidade base'
                      disabled={submitting || variationSubmitting}
                    />
                    <span className='input-unit'>{form.baseUnit}</span>
                  </div>
                  <button
                    type='button'
                    className='btn btn-ghost'
                    onClick={() => void handleAddVariation()}
                    disabled={submitting || variationSubmitting}
                  >
                    + Variação
                  </button>
                </div>

                {!isEditing && (draftVariations.length === 0 ? (
                    <p className='feedback'>Nenhuma variação adicionada.</p>
                  ) : (
                    <ul className='variation-list'>
                      {draftVariations.map((variation, index) => (
                        <li key={`${variation.baseQuantity}-${index}`} className='variation-item'>
                          <div>
                            <strong>
                              {getVariationLabel({
                                description: '',
                                baseQuantity: variation.baseQuantity,
                                baseUnit: form.baseUnit,
                                fallback: `Variação ${index + 1}`
                              })}
                            </strong>
                            <p>Quantidade base: {variation.baseQuantity} {form.baseUnit}</p>
                          </div>
                          <button
                            type='button'
                            className='btn btn-danger'
                            onClick={() => removeDraftVariation(index)}
                            disabled={submitting || variationSubmitting}
                          >
                            Remover
                          </button>
                        </li>
                      ))}
                    </ul>
                  ))}

                {isEditing && (loadingVariations ? (
                  <p className='feedback'>Carregando variações...</p>
                ) : productVariations.length === 0 ? (
                  <p className='feedback'>Nenhuma variação cadastrada.</p>
                ) : (
                  <ul className='variation-list'>
                    {productVariations.map((variation) => (
                      <li key={variation.id} className='variation-item variation-item-editable'>
                        <div className='variation-item-body'>
                          <strong>
                            {getVariationLabel({
                              description: variation.description,
                              baseQuantity: variation.base_quantity,
                              baseUnit: form.baseUnit,
                              fallback: `Variação #${variation.id}`
                            })}
                          </strong>
                          <p>Estoque atual: {variation.current_stock}</p>
                          <p>Quantidade base: {variation.base_quantity} {form.baseUnit}</p>
                        </div>

                        <div className='variation-item-fields'>
                          <div className='input-with-unit'>
                            <input
                              type='number'
                              min='1'
                              step='1'
                              value={variation.draftBaseQuantity}
                              onChange={(event) => updateVariationDraft(variation.id, 'draftBaseQuantity', event.target.value)}
                              placeholder='Quantidade base'
                              disabled={submitting || variationSubmitting}
                            />
                            <span className='input-unit'>{form.baseUnit}</span>
                          </div>
                        </div>

                        <div className='variation-item-actions'>
                          <span className={`status ${variation.is_active === 1 ? 'status-on' : 'status-off'}`}>
                            {variation.is_active === 1 ? 'Ativa' : 'Inativa'}
                          </span>
                          <button
                            type='button'
                            className='btn btn-ghost'
                            onClick={() => void handleSaveVariation(variation)}
                            disabled={submitting || variationSubmitting}
                          >
                            Salvar
                          </button>
                          <button
                            type='button'
                            className={variation.is_active === 1 ? 'btn btn-danger' : 'btn btn-ghost'}
                            onClick={() => void handleToggleVariationActive(variation)}
                            disabled={submitting || variationSubmitting}
                          >
                            {variation.is_active === 1 ? 'Desativar' : 'Reativar'}
                          </button>
                        </div>
                      </li>
                    ))}
                  </ul>
                ))}
              </div>

              <div className='form-actions'>
                <button className='btn btn-primary' type='submit' disabled={submitting}>
                  {submitting ? 'Salvando...' : isEditing ? 'Salvar edição' : 'Criar produto'}
                </button>
                <button className='btn btn-ghost' type='button' onClick={closeModal} disabled={submitting}>
                  Cancelar
                </button>
              </div>
            </form>
          </section>
        </div>
      )}

      {variationToDeactivate && (
        <div className='modal-overlay' onClick={() => setVariationToDeactivate(null)}>
          <section className='modal-content modal-confirm panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>Desativar variação</h2>
            </div>
            <p className='confirm-text'>
              Deseja desativar{' '}
              <strong>
                {getVariationLabel({
                  description: variationToDeactivate.description,
                  baseQuantity: variationToDeactivate.base_quantity,
                  baseUnit: form.baseUnit,
                  fallback: `Variação #${variationToDeactivate.id}`
                })}
              </strong>
              ?
            </p>
            <div className='confirm-actions'>
              <button
                type='button'
                className='btn btn-ghost'
                onClick={() => setVariationToDeactivate(null)}
                disabled={variationSubmitting}
              >
                Cancelar
              </button>
              <button
                type='button'
                className='btn btn-danger'
                onClick={() => void confirmDeactivateVariation()}
                disabled={variationSubmitting}
              >
                {variationSubmitting ? 'Desativando...' : 'Desativar'}
              </button>
            </div>
          </section>
        </div>
      )}

      {productToDeactivate && (
        <div className='modal-overlay' onClick={() => setProductToDeactivate(null)}>
          <section className='modal-content modal-confirm panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>Desativar produto</h2>
            </div>
            <p className='confirm-text'>
              Deseja desativar <strong>{productToDeactivate.name}</strong>?
            </p>
            <div className='confirm-actions'>
              <button
                type='button'
                className='btn btn-ghost'
                onClick={() => setProductToDeactivate(null)}
                disabled={togglingActive}
              >
                Cancelar
              </button>
              <button
                type='button'
                className='btn btn-danger'
                onClick={confirmDeactivate}
                disabled={togglingActive}
              >
                {togglingActive ? 'Desativando...' : 'Desativar'}
              </button>
            </div>
          </section>
        </div>
      )}
    </main>
  );
}
