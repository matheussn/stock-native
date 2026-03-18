import {ChevronLeft, Trash2} from 'lucide-react';
import {useContext, useEffect, useMemo, useState} from 'react';
import {Navigate, useNavigate} from 'react-router-dom';
import {
  CreateMovement,
  ListAssistentialWorks,
  ListInstitutions,
  ListProductGroupItems,
  ListProductGroups,
  ListProductVariationsByProduct,
  ListProducts
} from '../../../wailsjs/go/main/App';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import {
  clearMovementCreateAccess,
  hasMovementCreateAccess
} from '../../utils/movement-create-access';
import {AppNavigationGuardContext} from '../../utils/app-navigation-guard';
import {getVariationLabel} from '../../utils/variation-label';
import './movement-create-page.css';

function createInitialForm(assistentialWorkID = '') {
  return {
    assistentialWorkID,
    institutionID: '',
    type: 'out',
    notes: '',
    itemKey: '',
    itemQuantity: '1',
    items: []
  };
}

export default function MovementCreatePage() {
  const navigate = useNavigate();
  const {registerBeforeLeaveHandler} = useContext(AppNavigationGuardContext);
  const [assistentialWorks, setAssistentialWorks] = useState([]);
  const [institutions, setInstitutions] = useState([]);
  const [products, setProducts] = useState([]);
  const [groups, setGroups] = useState([]);
  const [variationsByProduct, setVariationsByProduct] = useState({});
  const [groupItemsByGroup, setGroupItemsByGroup] = useState({});

  const [loadingBaseData, setLoadingBaseData] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [showLeaveConfirm, setShowLeaveConfirm] = useState(false);
  const [isDirty, setIsDirty] = useState(false);
  const [pendingNavigationTarget, setPendingNavigationTarget] = useState('');
  const [form, setForm] = useState(createInitialForm());

  const canAccess = hasMovementCreateAccess();

  const activeInstitutions = useMemo(
    () => institutions.filter((institution) => institution.is_active === 1),
    [institutions]
  );

  const productOptions = useMemo(() => {
    const options = [];

    products.forEach((product) => {
      const variations = variationsByProduct[product.id] ?? [];

      variations.forEach((variation) => {
        options.push({
          id: variation.id,
          key: `product:${variation.id}`,
          kind: 'product',
            entityId: variation.id,
            label: `${product.name} - ${getVariationLabel({
              description: variation.description,
              baseQuantity: variation.base_quantity,
              baseUnit: product.base_unit,
              fallback: `Variação #${variation.id}`
            })}`,
          stock: variation.current_stock,
          productId: product.id,
          baseQuantity: variation.base_quantity
        });
      });
    });

    return options;
  }, [products, variationsByProduct]);

  const groupOptions = useMemo(
    () => groups.map((group) => ({
      key: `group:${group.id}`,
      kind: 'group',
      entityId: group.id,
      label: group.name
    })),
    [groups]
  );

  const itemOptionsByKey = useMemo(() => {
    const map = new Map();
    [...productOptions, ...groupOptions].forEach((option) => map.set(option.key, option));
    return map;
  }, [groupOptions, productOptions]);

  useEffect(() => {
    if (!canAccess) {
      return;
    }

    void loadBaseData();
  }, [canAccess]);

  useEffect(() => {
    if (!canAccess) {
      return undefined;
    }

    registerBeforeLeaveHandler((target) => {
      if (!isDirty) {
        clearMovementCreateAccess();
        navigate(target);
        return true;
      }

      setPendingNavigationTarget(target);
      setShowLeaveConfirm(true);
      return true;
    });

    return () => registerBeforeLeaveHandler(null);
  }, [canAccess, isDirty, navigate, registerBeforeLeaveHandler]);

  if (!canAccess) {
    return <Navigate to='/movimentacoes' replace />;
  }

  async function loadBaseData() {
    setLoadingBaseData(true);
    setError('');

    try {
      const [worksRes, institutionsRes, productsRes, groupsRes] = await Promise.all([
        ListAssistentialWorks(false),
        ListInstitutions(true),
        ListProducts(false),
        ListProductGroups(false)
      ]);

      setAssistentialWorks(worksRes);
      setInstitutions(institutionsRes);
      setProducts(productsRes);
      setGroups(groupsRes);
      setForm((prev) => ({
        ...prev,
        assistentialWorkID: prev.assistentialWorkID || (worksRes[0] ? String(worksRes[0].id) : '')
      }));

      const variationEntries = await Promise.all(
        productsRes.map(async (product) => [product.id, await ListProductVariationsByProduct(product.id, false)])
      );
      const variationMap = {};
      variationEntries.forEach(([productID, variations]) => {
        variationMap[productID] = variations;
      });
      setVariationsByProduct(variationMap);

      const groupItemEntries = await Promise.all(
        groupsRes.map(async (group) => [group.id, await ListProductGroupItems(group.id)])
      );
      const groupItemMap = {};
      groupItemEntries.forEach(([groupID, items]) => {
        groupItemMap[groupID] = items;
      });
      setGroupItemsByGroup(groupItemMap);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoadingBaseData(false);
    }
  }

  function requestBackToList() {
    if (submitting) {
      return;
    }

    if (!isDirty) {
      clearMovementCreateAccess();
      navigate('/movimentacoes');
      return;
    }

    setPendingNavigationTarget('/movimentacoes');
    setShowLeaveConfirm(true);
  }

  function cancelBackToList() {
    setShowLeaveConfirm(false);
    setPendingNavigationTarget('');
  }

  function confirmBackToList() {
    const target = pendingNavigationTarget || '/movimentacoes';
    clearMovementCreateAccess();
    setIsDirty(false);
    setShowLeaveConfirm(false);
    setPendingNavigationTarget('');
    navigate(target);
  }

  function addSelectedItem() {
    if (!form.itemKey) {
      setError('Selecione um item para adicionar.');
      return;
    }

    const quantity = Number(form.itemQuantity);
    if (!Number.isInteger(quantity) || quantity <= 0) {
      setError('Informe uma quantidade válida para o item.');
      return;
    }

    const selectedItem = itemOptionsByKey.get(form.itemKey);
    if (!selectedItem) {
      setError('O item selecionado não está mais disponível.');
      return;
    }

    setError('');
    setIsDirty(true);
    setForm((prev) => {
      const existingIndex = prev.items.findIndex((item) => item.key === selectedItem.key);

      if (existingIndex >= 0) {
        return {
          ...prev,
          itemKey: '',
          itemQuantity: '1',
          items: prev.items.map((item, index) => (
            index === existingIndex
              ? {...item, quantity: String(Number(item.quantity) + quantity)}
              : item
          ))
        };
      }

      return {
        ...prev,
        itemKey: '',
        itemQuantity: '1',
        items: [
          ...prev.items,
          {
            key: selectedItem.key,
            kind: selectedItem.kind,
            entityId: selectedItem.entityId,
            quantity: String(quantity)
          }
        ]
      };
    });
  }

  function updateItemQuantity(index, value) {
    setIsDirty(true);
    setForm((prev) => ({
      ...prev,
      items: prev.items.map((item, itemIndex) => (
        itemIndex === index ? {...item, quantity: value} : item
      ))
    }));
  }

  function removeItem(index) {
    setIsDirty(true);
    setForm((prev) => ({
      ...prev,
      items: prev.items.filter((_, itemIndex) => itemIndex !== index)
    }));
  }

  async function handleSubmit(event) {
    event.preventDefault();

    const assistentialWorkID = Number(form.assistentialWorkID);
    const institutionID = form.type === 'out' && form.institutionID ? Number(form.institutionID) : 0;

    if (!Number.isInteger(assistentialWorkID) || assistentialWorkID <= 0) {
      setError('Selecione uma obra assistencial.');
      return;
    }

    if (form.type === 'out' && form.institutionID && (!Number.isInteger(institutionID) || institutionID <= 0)) {
      setError('Selecione uma instituição válida.');
      return;
    }

    const productItems = [];
    const groupItems = [];

    for (const item of form.items) {
      const quantity = Number(item.quantity);
      if (!Number.isInteger(quantity) || quantity <= 0) {
        setError('Todas as quantidades devem ser inteiros maiores que zero.');
        return;
      }

      if (item.kind === 'product') {
        productItems.push({
          product_variation_id: item.entityId,
          quantity
        });
        continue;
      }

      groupItems.push({
        product_group_id: item.entityId,
        quantity
      });
    }

    if (productItems.length === 0 && groupItems.length === 0) {
      setError('Adicione pelo menos um item para registrar a movimentação.');
      return;
    }

    if (form.type === 'out') {
      const validationError = validateOutStock(productItems, groupItems, productOptions, groupItemsByGroup);
      if (validationError) {
        setError(validationError);
        return;
      }
    }

    setSubmitting(true);
    setError('');

    try {
      await CreateMovement({
        assistential_work_id: assistentialWorkID,
        institution_id: institutionID,
        type: form.type,
        notes: form.notes,
        product_items: productItems,
        group_items: groupItems
      });

      clearMovementCreateAccess();
      setIsDirty(false);
      navigate('/movimentacoes', {replace: true});
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  return (
    <main className='movement-create-page'>
      <header className='movement-create-header panel'>
        <div className='movement-create-header-main'>
          <button type='button' className='movement-create-back' onClick={requestBackToList} disabled={submitting} aria-label='Voltar'>
            <ChevronLeft aria-hidden='true' size={18} />
          </button>
          <h1>Nova movimentação</h1>
        </div>
      </header>

      {error && <p className='movement-create-feedback movement-create-feedback-error'>{error}</p>}

      <form className='movement-create-layout' onSubmit={handleSubmit}>
        <section className='movement-create-items panel'>
          <div className='movement-create-section-head'>
            <div>
              <h2>Itens da movimentação</h2>
            </div>
          </div>

          <div className='movement-create-picker'>
            <select
              value={form.itemKey}
              onChange={(event) => {
                setIsDirty(true);
                setForm((prev) => ({...prev, itemKey: event.target.value}));
              }}
              disabled={submitting || loadingBaseData}
            >
              <option value=''>Selecione um produto ou cesta</option>
              <optgroup label='Produtos'>
                {productOptions.map((option) => (
                  <option key={option.key} value={option.key}>
                    {option.label}
                  </option>
                ))}
              </optgroup>
              <optgroup label='Cestas'>
                {groupOptions.map((option) => (
                  <option key={option.key} value={option.key}>
                    {option.label}
                  </option>
                ))}
              </optgroup>
            </select>

            <input
              type='number'
              min='1'
              step='1'
              value={form.itemQuantity}
              onChange={(event) => {
                setIsDirty(true);
                setForm((prev) => ({...prev, itemQuantity: event.target.value}));
              }}
              disabled={submitting}
            />

            <button type='button' className='btn btn-primary' onClick={addSelectedItem} disabled={submitting || loadingBaseData}>
              Adicionar
            </button>
          </div>

          <div className='movement-create-list-wrap'>
            {form.items.length === 0 ? (
              <p className='movement-create-feedback'>Nenhum item adicionado ainda.</p>
            ) : (
              <ul className='movement-create-item-list'>
                {form.items.map((item, index) => {
                  const option = itemOptionsByKey.get(item.key);
                  return (
                    <li key={item.key} className='movement-create-item-row'>
                      <div className='movement-create-item-main'>
                        <div className='movement-create-item-title'>
                          <span className={`movement-create-item-tag movement-create-item-tag-${item.kind}`}>
                            {item.kind === 'product' ? 'Produto' : 'Cesta'}
                          </span>
                          <strong>{option?.label ?? 'Item indisponível'}</strong>
                        </div>
                      </div>

                      <div className='movement-create-item-side'>
                        <label>
                          Quantidade
                          <input
                            type='number'
                            min='1'
                            step='1'
                            value={item.quantity}
                            onChange={(event) => updateItemQuantity(index, event.target.value)}
                            disabled={submitting}
                          />
                        </label>
                      </div>

                      <button
                        type='button'
                        className='movement-create-item-remove'
                        onClick={() => removeItem(index)}
                        disabled={submitting}
                        aria-label='Remover item'
                        title='Remover item'
                      >
                        <Trash2 aria-hidden='true' size={16} />
                      </button>
                    </li>
                  );
                })}
              </ul>
            )}
          </div>
        </section>

        <aside className='movement-create-sidebar panel'>
          <div className='movement-create-sidebar-section'>
            <label htmlFor='movement-work'>Obra assistencial</label>
            <select
              id='movement-work'
              value={form.assistentialWorkID}
              onChange={(event) => {
                setIsDirty(true);
                setForm((prev) => ({...prev, assistentialWorkID: event.target.value}));
              }}
              disabled={submitting || loadingBaseData}
            >
              {assistentialWorks.length === 0 && <option value=''>Nenhuma obra ativa</option>}
              {assistentialWorks.map((work) => (
                <option key={work.id} value={work.id}>
                  {work.name}
                </option>
              ))}
            </select>
          </div>

          <div className='movement-create-sidebar-section'>
            <label htmlFor='movement-type'>Tipo</label>
            <select
              id='movement-type'
              value={form.type}
              onChange={(event) => {
                setIsDirty(true);
                setForm((prev) => ({
                  ...prev,
                  type: event.target.value,
                  institutionID: event.target.value === 'out' ? prev.institutionID : ''
                }));
              }}
              disabled={submitting || loadingBaseData}
            >
              <option value='out'>Saída</option>
              <option value='in'>Entrada</option>
            </select>
          </div>

          {form.type === 'out' && (
            <div className='movement-create-sidebar-section'>
              <label htmlFor='movement-institution'>Instituição</label>
              <select
                id='movement-institution'
                value={form.institutionID}
                onChange={(event) => {
                  setIsDirty(true);
                  setForm((prev) => ({...prev, institutionID: event.target.value}));
                }}
                disabled={submitting || loadingBaseData}
              >
                <option value=''>Sem instituição</option>
                {activeInstitutions.map((institution) => (
                  <option key={institution.id} value={institution.id}>
                    {institution.name}
                  </option>
                ))}
              </select>
            </div>
          )}

          <div className='movement-create-sidebar-section movement-create-sidebar-notes'>
            <label htmlFor='movement-notes'>Observações</label>
            <textarea
              id='movement-notes'
              value={form.notes}
              onChange={(event) => {
                setIsDirty(true);
                setForm((prev) => ({...prev, notes: event.target.value}));
              }}
              placeholder='Opcional'
              rows={8}
              disabled={submitting}
            />
          </div>

          <div className='movement-create-actions'>
            <button className='btn btn-primary' type='submit' disabled={submitting || loadingBaseData}>
              {submitting ? 'Salvando...' : 'Registrar movimentação'}
            </button>
          </div>
        </aside>
      </form>

      {showLeaveConfirm && (
        <div className='movement-create-confirm-overlay' onClick={cancelBackToList}>
          <section className='movement-create-confirm panel' onClick={(event) => event.stopPropagation()}>
            <h2>Voltar para movimentações</h2>
            <p>Você vai perder tudo o que fez nesta movimentação. Deseja voltar para a listagem?</p>
            <div className='movement-create-confirm-actions'>
              <button type='button' className='btn btn-ghost' onClick={cancelBackToList} disabled={submitting}>
                Cancelar
              </button>
              <button type='button' className='btn btn-danger' onClick={confirmBackToList} disabled={submitting}>
                Confirmar volta
              </button>
            </div>
          </section>
        </div>
      )}
    </main>
  );
}

function validateOutStock(productItems, groupItems, variationOptions, groupItemsByGroup) {
  const stockByVariation = new Map();
  const variationsByProduct = new Map();

  variationOptions.forEach((variation) => {
    const variationID = variation.id ?? variation.entityId;
    if (!variationID) {
      return;
    }

    stockByVariation.set(variationID, variation.stock);
    if (!variationsByProduct.has(variation.productId)) {
      variationsByProduct.set(variation.productId, []);
    }
    variationsByProduct.get(variation.productId).push({
      ...variation,
      id: variationID
    });
  });

  for (const item of productItems) {
    const current = stockByVariation.get(item.product_variation_id) ?? 0;
    if (current < item.quantity) {
      return `Estoque insuficiente na variação #${item.product_variation_id}.`;
    }
    stockByVariation.set(item.product_variation_id, current - item.quantity);
  }

  for (const groupItem of groupItems) {
    const composition = groupItemsByGroup[groupItem.product_group_id] ?? [];

    for (const comp of composition) {
      let remainingBase = comp.base_quantity * groupItem.quantity;
      const variations = [...(variationsByProduct.get(comp.product_id) ?? [])]
        .sort((left, right) => right.baseQuantity - left.baseQuantity);

      for (const variation of variations) {
        if (remainingBase <= 0) {
          break;
        }

        const availableUnits = stockByVariation.get(variation.id) ?? 0;
        if (availableUnits <= 0) {
          continue;
        }

        const maxNeededUnits = Math.floor(remainingBase / variation.baseQuantity);
        if (maxNeededUnits <= 0) {
          continue;
        }

        const allocation = Math.min(maxNeededUnits, availableUnits);
        stockByVariation.set(variation.id, availableUnits - allocation);
        remainingBase -= allocation * variation.baseQuantity;
      }

      if (remainingBase > 0) {
        return `Estoque insuficiente para resolver o produto #${comp.product_id} na cesta #${groupItem.product_group_id}.`;
      }
    }
  }

  return '';
}
