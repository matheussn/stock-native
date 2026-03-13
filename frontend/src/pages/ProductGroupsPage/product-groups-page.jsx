import {useEffect, useMemo, useState} from 'react';
import {
  CreateProductGroup,
  ListProductGroupItems,
  ListProductGroups,
  ListProducts,
  RemoveProductGroupItem,
  SetProductGroupActive,
  UpdateProductGroup,
  UpsertProductGroupItem
} from '../../../wailsjs/go/main/App';
import AssistentialWorkFilters from '../../components/AssistentialWorkFilters';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import './product-groups-page.css';

const FILTERS = {
  active: 'Ativos',
  inactive: 'Inativos',
  all: 'Todos'
};

export default function ProductGroupsPage() {
  const [groups, setGroups] = useState([]);
  const [products, setProducts] = useState([]);
  const [filter, setFilter] = useState('active');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [itemSubmitting, setItemSubmitting] = useState(false);
  const [togglingActive, setTogglingActive] = useState(false);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [loadingGroupItems, setLoadingGroupItems] = useState(false);
  const [groupItems, setGroupItems] = useState([]);
  const [draftItems, setDraftItems] = useState([]);
  const [groupToDeactivate, setGroupToDeactivate] = useState(null);
  const [itemToRemove, setItemToRemove] = useState(null);
  const [form, setForm] = useState({
    id: null,
    name: '',
    description: ''
  });
  const [itemForm, setItemForm] = useState({
    productId: '',
    baseQuantity: ''
  });

  const isEditing = form.id !== null;

  const visibleGroups = useMemo(() => {
    if (filter === 'inactive') {
      return groups.filter((item) => item.is_active === 0);
    }
    return groups;
  }, [filter, groups]);

  const productsByID = useMemo(() => {
    const map = new Map();
    products.forEach((p) => map.set(p.id, p));
    return map;
  }, [products]);

  const activeProducts = useMemo(() => products.filter((p) => p.is_active === 1), [products]);

  useEffect(() => {
    void loadGroups(filter);
  }, [filter]);

  useEffect(() => {
    void loadProducts();
  }, []);

  async function loadGroups(currentFilter) {
    setLoading(true);
    setError('');
    try {
      const includeInactive = currentFilter !== 'active';
      const response = await ListProductGroups(includeInactive);
      setGroups(response);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  async function loadProducts() {
    try {
      const response = await ListProducts(true);
      setProducts(response);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    }
  }

  function mapGroupItem(item) {
    return {
      ...item,
      draftBaseQuantity: String(item.base_quantity)
    };
  }

  async function loadGroupItems(groupId) {
    setLoadingGroupItems(true);
    try {
      const response = await ListProductGroupItems(groupId);
      setGroupItems(response.map(mapGroupItem));
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoadingGroupItems(false);
    }
  }

  function resetForm() {
    setForm({
      id: null,
      name: '',
      description: ''
    });
  }

  function resetItemForm() {
    setItemForm({
      productId: '',
      baseQuantity: ''
    });
  }

  function openCreateModal() {
    resetForm();
    resetItemForm();
    setDraftItems([]);
    setGroupItems([]);
    setGroupToDeactivate(null);
    setItemToRemove(null);
    setError('');
    setIsModalOpen(true);
  }

  async function startEdit(group) {
    setForm({
      id: group.id,
      name: group.name,
      description: group.description ?? ''
    });
    resetItemForm();
    setDraftItems([]);
    setGroupItems([]);
    setGroupToDeactivate(null);
    setItemToRemove(null);
    setError('');
    setIsModalOpen(true);
    await loadGroupItems(group.id);
  }

  function closeModal() {
    setIsModalOpen(false);
    resetForm();
    resetItemForm();
    setDraftItems([]);
    setGroupItems([]);
    setItemToRemove(null);
  }

  async function handleSubmit(event) {
    event.preventDefault();
    const name = form.name.trim();

    if (!name) {
      setError('Informe o nome da cesta.');
      return;
    }

    setSubmitting(true);
    setError('');
    try {
      if (isEditing) {
        await UpdateProductGroup(form.id, name, form.description);
      } else {
        const created = await CreateProductGroup(name, form.description);
        for (const item of draftItems) {
          await UpsertProductGroupItem(created.id, item.productId, item.baseQuantity);
        }
      }
      closeModal();
      await loadGroups(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function handleItemSubmit(event) {
    event?.preventDefault();
    const productId = Number(itemForm.productId);
    const baseQuantity = Number(itemForm.baseQuantity);

    if (!Number.isInteger(productId) || productId <= 0) {
      setError('Selecione um produto para adicionar na cesta.');
      return;
    }
    if (!Number.isInteger(baseQuantity) || baseQuantity <= 0) {
      setError('A quantidade base deve ser um inteiro maior que zero.');
      return;
    }

    if (!isEditing) {
      setDraftItems((prev) => {
        const index = prev.findIndex((item) => item.productId === productId);
        if (index === -1) {
          return [...prev, {productId, baseQuantity}];
        }

        const next = [...prev];
        next[index] = {...next[index], baseQuantity};
        return next;
      });
      resetItemForm();
      setError('');
      return;
    }

    setItemSubmitting(true);
    setError('');
    try {
      await UpsertProductGroupItem(form.id, productId, baseQuantity);
      resetItemForm();
      await loadGroupItems(form.id);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setItemSubmitting(false);
    }
  }

  function removeDraftItem(productId) {
    setDraftItems((prev) => prev.filter((item) => item.productId !== productId));
  }

  function updateGroupItemDraft(productId, value) {
    setGroupItems((prev) =>
      prev.map((item) => (item.product_id === productId ? {...item, draftBaseQuantity: value} : item))
    );
  }

  async function handleSaveGroupItem(item) {
    const baseQuantity = Number(item.draftBaseQuantity);

    if (!Number.isInteger(baseQuantity) || baseQuantity <= 0) {
      setError('A quantidade base deve ser um inteiro maior que zero.');
      return;
    }

    setItemSubmitting(true);
    setError('');
    try {
      await UpsertProductGroupItem(form.id, item.product_id, baseQuantity);
      await loadGroupItems(form.id);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setItemSubmitting(false);
    }
  }

  function requestRemoveGroupItem(item) {
    const productName = productsByID.get(item.product_id)?.name ?? `Produto #${item.product_id}`;
    setItemToRemove({
      groupId: form.id,
      productId: item.product_id,
      productName
    });
  }

  async function confirmRemoveGroupItem() {
    if (!itemToRemove) {
      return;
    }

    setItemSubmitting(true);
    setError('');
    try {
      await RemoveProductGroupItem(itemToRemove.groupId, itemToRemove.productId);
      setItemToRemove(null);
      await loadGroupItems(form.id);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setItemSubmitting(false);
    }
  }

  async function handleToggleGroupActive(group) {
    if (group.is_active === 1) {
      setGroupToDeactivate(group);
      return;
    }

    setError('');
    try {
      setTogglingActive(true);
      await SetProductGroupActive(group.id, true);
      await loadGroups(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setTogglingActive(false);
    }
  }

  async function confirmDeactivateGroup() {
    if (!groupToDeactivate) {
      return;
    }

    setError('');
    try {
      setTogglingActive(true);
      await SetProductGroupActive(groupToDeactivate.id, false);
      setGroupToDeactivate(null);
      await loadGroups(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setTogglingActive(false);
    }
  }

  return (
    <main className='product-groups-page'>
      <section className='panel panel-list'>
        <div className='list-header'>
          <div>
            <h1 className='title'>Cestas</h1>
            <p className='subtitle'>Cadastro de grupos de produtos</p>
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
          <p className='feedback'>Carregando cestas...</p>
        ) : visibleGroups.length === 0 ? (
          <p className='feedback'>Nenhum registro encontrado para este filtro.</p>
        ) : (
          <ul className='group-list'>
            {visibleGroups.map((group) => (
              <li key={group.id} className='group-item'>
                <div className='group-main'>
                  <div>
                    <strong>{group.name}</strong>
                    <p>{group.description || 'Sem descrição'}</p>
                  </div>

                  <div className='item-actions'>
                    <span className={`status ${group.is_active === 1 ? 'status-on' : 'status-off'}`}>
                      {group.is_active === 1 ? 'Ativa' : 'Inativa'}
                    </span>
                    <button type='button' className='btn btn-ghost' onClick={() => void startEdit(group)}>
                      Editar
                    </button>
                    <button type='button' className='btn btn-danger' onClick={() => void handleToggleGroupActive(group)}>
                      {group.is_active === 1 ? 'Desativar' : 'Reativar'}
                    </button>
                  </div>
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
              <h2>{isEditing ? 'Editar cesta' : 'Nova cesta'}</h2>
              <button type='button' className='btn btn-ghost' onClick={closeModal} disabled={submitting || itemSubmitting}>
                Fechar
              </button>
            </div>

            <form className='form' onSubmit={handleSubmit}>
              <label htmlFor='group-name'>Nome</label>
              <input
                id='group-name'
                value={form.name}
                onChange={(event) => setForm((prev) => ({...prev, name: event.target.value}))}
                placeholder='Ex: Cesta Básica'
                maxLength={120}
                disabled={submitting}
              />

              <label htmlFor='group-description'>Descrição</label>
              <textarea
                id='group-description'
                value={form.description}
                onChange={(event) => setForm((prev) => ({...prev, description: event.target.value}))}
                placeholder='Opcional'
                rows={3}
                maxLength={600}
                disabled={submitting}
              />

              <div className='composition-setup'>
                <h3>{isEditing ? 'Composição cadastrada' : 'Composição inicial'}</h3>
                <p className='hint'>
                  {isEditing
                    ? 'Adicione itens novos e ajuste as quantidades desta cesta nesta mesma modal.'
                    : 'Monte a composição antes de criar a cesta.'}
                </p>

                <div className='composition-form'>
                  <select
                    value={itemForm.productId}
                    onChange={(event) => setItemForm((prev) => ({...prev, productId: event.target.value}))}
                    disabled={itemSubmitting || submitting}
                  >
                    <option value=''>Selecione um produto</option>
                    {activeProducts.map((product) => (
                      <option key={product.id} value={product.id}>
                        {product.name} ({product.base_unit})
                      </option>
                    ))}
                  </select>
                  <input
                    type='number'
                    min='1'
                    step='1'
                    value={itemForm.baseQuantity}
                    onChange={(event) => setItemForm((prev) => ({...prev, baseQuantity: event.target.value}))}
                    placeholder='Quantidade base'
                    disabled={itemSubmitting || submitting}
                  />
                  <button
                    className='btn btn-ghost'
                    type='button'
                    onClick={() => void handleItemSubmit()}
                    disabled={itemSubmitting || submitting}
                  >
                    + Item
                  </button>
                </div>

                {!isEditing && (draftItems.length === 0 ? (
                  <p className='feedback'>Nenhum item adicionado nesta cesta.</p>
                ) : (
                  <ul className='composition-list'>
                    {draftItems.map((item) => {
                      const product = productsByID.get(item.productId);
                      const productName = product?.name ?? `Produto #${item.productId}`;
                      const unit = product?.base_unit ?? '-';

                      return (
                        <li key={item.productId} className='composition-item'>
                          <div className='composition-item-body'>
                            <strong>{productName}</strong>
                            <p>Quantidade base: {item.baseQuantity} {unit}</p>
                          </div>

                          <div className='composition-item-actions'>
                            <button
                              type='button'
                              className='btn btn-danger'
                              onClick={() => removeDraftItem(item.productId)}
                              disabled={submitting}
                            >
                              Remover
                            </button>
                          </div>
                        </li>
                      );
                    })}
                  </ul>
                ))}

                {isEditing && (loadingGroupItems ? (
                  <p className='feedback'>Carregando itens da cesta...</p>
                ) : groupItems.length === 0 ? (
                  <p className='feedback'>Nenhum item adicionado nesta cesta.</p>
                ) : (
                  <ul className='composition-list'>
                    {groupItems.map((item) => {
                      const product = productsByID.get(item.product_id);
                      const productName = product?.name ?? `Produto #${item.product_id}`;
                      const unit = product?.base_unit ?? '-';

                      return (
                        <li key={item.id} className='composition-item composition-item-editable'>
                          <div className='composition-item-body'>
                            <strong>{productName}</strong>
                            <p>Quantidade base atual: {item.base_quantity} {unit}</p>
                          </div>

                          <div className='composition-item-fields'>
                            <input
                              type='number'
                              min='1'
                              step='1'
                              value={item.draftBaseQuantity}
                              onChange={(event) => updateGroupItemDraft(item.product_id, event.target.value)}
                              placeholder='Quantidade base'
                              disabled={itemSubmitting || submitting}
                            />
                          </div>

                          <div className='composition-item-actions'>
                            <button
                              type='button'
                              className='btn btn-ghost'
                              onClick={() => void handleSaveGroupItem(item)}
                              disabled={itemSubmitting || submitting}
                            >
                              Salvar quantidade
                            </button>
                            <button
                              type='button'
                              className='btn btn-danger'
                              onClick={() => requestRemoveGroupItem(item)}
                              disabled={itemSubmitting || submitting}
                            >
                              Remover
                            </button>
                          </div>
                        </li>
                      );
                    })}
                  </ul>
                ))}
              </div>

              <div className='form-actions'>
                <button className='btn btn-primary' type='submit' disabled={submitting || itemSubmitting}>
                  {submitting ? 'Salvando...' : isEditing ? 'Salvar edição' : 'Criar cesta'}
                </button>
                <button className='btn btn-ghost' type='button' onClick={closeModal} disabled={submitting || itemSubmitting}>
                  Cancelar
                </button>
              </div>
            </form>
          </section>
        </div>
      )}

      {itemToRemove && (
        <div className='modal-overlay' onClick={() => setItemToRemove(null)}>
          <section className='modal-content modal-confirm panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>Remover item da cesta</h2>
            </div>
            <p className='confirm-text'>
              Deseja remover <strong>{itemToRemove.productName}</strong> da cesta?
            </p>
            <div className='confirm-actions'>
              <button
                type='button'
                className='btn btn-ghost'
                onClick={() => setItemToRemove(null)}
                disabled={itemSubmitting}
              >
                Cancelar
              </button>
              <button
                type='button'
                className='btn btn-danger'
                onClick={() => void confirmRemoveGroupItem()}
                disabled={itemSubmitting}
              >
                {itemSubmitting ? 'Removendo...' : 'Remover'}
              </button>
            </div>
          </section>
        </div>
      )}

      {groupToDeactivate && (
        <div className='modal-overlay' onClick={() => setGroupToDeactivate(null)}>
          <section className='modal-content modal-confirm panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>Desativar cesta</h2>
            </div>
            <p className='confirm-text'>
              Deseja desativar <strong>{groupToDeactivate.name}</strong>?
            </p>
            <div className='confirm-actions'>
              <button
                type='button'
                className='btn btn-ghost'
                onClick={() => setGroupToDeactivate(null)}
                disabled={togglingActive}
              >
                Cancelar
              </button>
              <button
                type='button'
                className='btn btn-danger'
                onClick={() => void confirmDeactivateGroup()}
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
