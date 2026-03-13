import {useEffect, useMemo, useState} from 'react';
import {
  CreateMovement,
  ListAssistentialWorks,
  ListMovementGroupItemResolutions,
  ListMovementGroupItems,
  ListMovementProductItems,
  ListMovements,
  ListProductGroupItems,
  ListProductGroups,
  ListProductVariationsByProduct,
  ListProducts
} from '../../../wailsjs/go/main/App';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import {getVariationLabel} from '../../utils/variation-label';
import './movements-page.css';

export default function MovementsPage() {
  const [assistentialWorks, setAssistentialWorks] = useState([]);
  const [products, setProducts] = useState([]);
  const [groups, setGroups] = useState([]);
  const [variationsByProduct, setVariationsByProduct] = useState({});
  const [groupItemsByGroup, setGroupItemsByGroup] = useState({});
  const [movements, setMovements] = useState([]);
  const [detailsByMovement, setDetailsByMovement] = useState({});

  const [historyTypeFilter, setHistoryTypeFilter] = useState('all');
  const [historyWorkFilter, setHistoryWorkFilter] = useState('all');
  const [expandedMovementId, setExpandedMovementId] = useState(null);

  const [loadingBaseData, setLoadingBaseData] = useState(true);
  const [loadingMovements, setLoadingMovements] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);

  const [form, setForm] = useState({
    assistentialWorkID: '',
    type: 'out',
    notes: '',
    productLines: [{variationId: '', quantity: '1'}],
    groupLines: [{groupId: '', quantity: '1'}]
  });

  const worksByID = useMemo(() => {
    const map = new Map();
    assistentialWorks.forEach((w) => map.set(w.id, w));
    return map;
  }, [assistentialWorks]);

  const productsByID = useMemo(() => {
    const map = new Map();
    products.forEach((p) => map.set(p.id, p));
    return map;
  }, [products]);

  const groupsByID = useMemo(() => {
    const map = new Map();
    groups.forEach((g) => map.set(g.id, g));
    return map;
  }, [groups]);

  const variationsByID = useMemo(() => {
    const map = new Map();
    Object.values(variationsByProduct).flat().forEach((v) => map.set(v.id, v));
    return map;
  }, [variationsByProduct]);

  const variationOptions = useMemo(() => {
    const options = [];
    products.forEach((p) => {
      const variations = variationsByProduct[p.id] ?? [];
      variations.forEach((v) => {
        options.push({
          id: v.id,
          label: `${p.name} - ${getVariationLabel({
            description: v.description,
            baseQuantity: v.base_quantity,
            baseUnit: p.base_unit,
            fallback: `Variação #${v.id}`
          })}`,
          stock: v.current_stock,
          productId: p.id,
          baseQuantity: v.base_quantity
        });
      });
    });
    return options;
  }, [products, variationsByProduct]);

  useEffect(() => {
    void loadBaseData();
  }, []);

  useEffect(() => {
    void loadMovements(historyTypeFilter, historyWorkFilter);
  }, [historyTypeFilter, historyWorkFilter]);

  async function loadBaseData() {
    setLoadingBaseData(true);
    setError('');
    try {
      const [worksRes, productsRes, groupsRes] = await Promise.all([
        ListAssistentialWorks(false),
        ListProducts(false),
        ListProductGroups(false)
      ]);

      setAssistentialWorks(worksRes);
      setProducts(productsRes);
      setGroups(groupsRes);

      if (worksRes.length > 0 && !form.assistentialWorkID) {
        setForm((prev) => ({...prev, assistentialWorkID: String(worksRes[0].id)}));
      }

      const variationEntries = await Promise.all(
        productsRes.map(async (p) => [p.id, await ListProductVariationsByProduct(p.id, false)])
      );
      const variationMap = {};
      variationEntries.forEach(([productId, variations]) => {
        variationMap[productId] = variations;
      });
      setVariationsByProduct(variationMap);

      const groupItemEntries = await Promise.all(
        groupsRes.map(async (g) => [g.id, await ListProductGroupItems(g.id)])
      );
      const groupItemMap = {};
      groupItemEntries.forEach(([groupId, items]) => {
        groupItemMap[groupId] = items;
      });
      setGroupItemsByGroup(groupItemMap);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoadingBaseData(false);
    }
  }

  async function loadMovements(typeFilter, workFilter) {
    setLoadingMovements(true);
    setError('');
    try {
      const movementType = typeFilter === 'all' ? '' : typeFilter;
      const assistentialWorkID = workFilter === 'all' ? 0 : Number(workFilter);
      const response = await ListMovements(movementType, assistentialWorkID);
      setMovements(response);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoadingMovements(false);
    }
  }

  function resetForm(nextAssistentialWorkID = assistentialWorks.length > 0 ? String(assistentialWorks[0].id) : '') {
    setForm({
      assistentialWorkID: nextAssistentialWorkID,
      type: 'out',
      notes: '',
      productLines: [{variationId: '', quantity: '1'}],
      groupLines: [{groupId: '', quantity: '1'}]
    });
  }

  function openCreateModal() {
    resetForm();
    setError('');
    setIsModalOpen(true);
  }

  function closeModal() {
    setIsModalOpen(false);
    resetForm(form.assistentialWorkID || (assistentialWorks.length > 0 ? String(assistentialWorks[0].id) : ''));
  }

  async function handleSubmit(event) {
    event.preventDefault();
    const assistentialWorkID = Number(form.assistentialWorkID);

    if (!Number.isInteger(assistentialWorkID) || assistentialWorkID <= 0) {
      setError('Selecione uma obra assistencial.');
      return;
    }

    const productItems = [];
    for (const line of form.productLines) {
      if (!line.variationId) {
        continue;
      }
      const quantity = Number(line.quantity);
      const variationID = Number(line.variationId);
      if (!Number.isInteger(quantity) || quantity <= 0) {
        setError('Quantidade de item direto deve ser inteiro maior que zero.');
        return;
      }
      productItems.push({
        product_variation_id: variationID,
        quantity
      });
    }

    const groupItems = [];
    for (const line of form.groupLines) {
      if (!line.groupId) {
        continue;
      }
      const quantity = Number(line.quantity);
      const groupID = Number(line.groupId);
      if (!Number.isInteger(quantity) || quantity <= 0) {
        setError('Quantidade de cesta deve ser inteiro maior que zero.');
        return;
      }
      groupItems.push({
        product_group_id: groupID,
        quantity
      });
    }

    if (productItems.length === 0 && groupItems.length === 0) {
      setError('Informe pelo menos um item direto ou uma cesta.');
      return;
    }

    if (form.type === 'out') {
      const validationError = validateOutStock(productItems, groupItems, variationOptions, groupItemsByGroup);
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
        type: form.type,
        notes: form.notes,
        product_items: productItems,
        group_items: groupItems
      });

      setIsModalOpen(false);
      resetForm(String(assistentialWorkID));

      await Promise.all([
        loadBaseData(),
        loadMovements(historyTypeFilter, historyWorkFilter)
      ]);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function toggleDetails(movementID) {
    if (expandedMovementId === movementID) {
      setExpandedMovementId(null);
      return;
    }

    setExpandedMovementId(movementID);
    if (detailsByMovement[movementID]) {
      return;
    }

    setDetailsByMovement((prev) => ({
      ...prev,
      [movementID]: {loading: true, productItems: [], groupItems: [], resolutionsByGroupItem: {}}
    }));

    try {
      const [productItems, groupItems] = await Promise.all([
        ListMovementProductItems(movementID),
        ListMovementGroupItems(movementID)
      ]);

      const resolutionEntries = await Promise.all(
        groupItems.map(async (g) => [g.id, await ListMovementGroupItemResolutions(g.id)])
      );
      const resolutionsByGroupItem = {};
      resolutionEntries.forEach(([groupItemID, list]) => {
        resolutionsByGroupItem[groupItemID] = list;
      });

      setDetailsByMovement((prev) => ({
        ...prev,
        [movementID]: {
          loading: false,
          productItems,
          groupItems,
          resolutionsByGroupItem
        }
      }));
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
      setDetailsByMovement((prev) => ({
        ...prev,
        [movementID]: {loading: false, productItems: [], groupItems: [], resolutionsByGroupItem: {}}
      }));
    }
  }

  function addProductLine() {
    setForm((prev) => ({
      ...prev,
      productLines: [...prev.productLines, {variationId: '', quantity: '1'}]
    }));
  }

  function removeProductLine(index) {
    setForm((prev) => ({
      ...prev,
      productLines: prev.productLines.filter((_, i) => i !== index)
    }));
  }

  function addGroupLine() {
    setForm((prev) => ({
      ...prev,
      groupLines: [...prev.groupLines, {groupId: '', quantity: '1'}]
    }));
  }

  function removeGroupLine(index) {
    setForm((prev) => ({
      ...prev,
      groupLines: prev.groupLines.filter((_, i) => i !== index)
    }));
  }

  return (
    <main className='movements-page'>
      <section className='panel panel-list'>
        <div className='page-header'>
          <div>
            <h1 className='title'>Movimentações</h1>
            <p className='subtitle'>Entrada e saída com itens diretos e cestas</p>
          </div>
          <button type='button' className='btn btn-primary' onClick={openCreateModal}>
            + Novo
          </button>
        </div>

        <div className='list-header'>
          <h2>Historico</h2>
          <div className='history-filters'>
            <select value={historyTypeFilter} onChange={(e) => setHistoryTypeFilter(e.target.value)}>
              <option value='all'>Todos os tipos</option>
              <option value='in'>Entradas</option>
              <option value='out'>Saidas</option>
            </select>
            <select value={historyWorkFilter} onChange={(e) => setHistoryWorkFilter(e.target.value)}>
              <option value='all'>Todas as obras</option>
              {assistentialWorks.map((work) => (
                <option key={work.id} value={work.id}>
                  {work.name}
                </option>
              ))}
            </select>
          </div>
        </div>

        {error && <p className='feedback error'>{error}</p>}

        {loadingMovements ? (
          <p className='feedback'>Carregando movimentações...</p>
        ) : movements.length === 0 ? (
          <p className='feedback'>Nenhuma movimentação encontrada.</p>
        ) : (
          <ul className='movement-list'>
            {movements.map((movement) => {
              const details = detailsByMovement[movement.id];
              return (
                <li key={movement.id} className='movement-item'>
                  <div className='movement-main'>
                    <div>
                      <strong>#{movement.id} - {movement.type === 'in' ? 'Entrada' : 'Saida'}</strong>
                      <p>Obra: {worksByID.get(movement.assistential_work_id)?.name ?? `#${movement.assistential_work_id}`}</p>
                      <p>Data: {movement.created_at}</p>
                      <p>Obs.: {movement.notes || 'Sem observações'}</p>
                    </div>
                    <div className='item-actions'>
                      <span className={`status ${movement.type === 'in' ? 'status-in' : 'status-out'}`}>
                        {movement.type === 'in' ? 'Entrada' : 'Saida'}
                      </span>
                      <button type='button' className='btn btn-ghost' onClick={() => toggleDetails(movement.id)}>
                        {expandedMovementId === movement.id ? 'Fechar detalhes' : 'Detalhes'}
                      </button>
                    </div>
                  </div>

                  {expandedMovementId === movement.id && (
                    <div className='movement-details'>
                      {details?.loading ? (
                        <p className='feedback'>Carregando detalhes...</p>
                      ) : (
                        <>
                          <h3>Itens diretos</h3>
                          {details && details.productItems.length > 0 ? (
                            <ul className='detail-list'>
                              {details.productItems.map((item) => {
                                const variation = variationsByID.get(item.product_variation_id);
                                const product = variation ? productsByID.get(variation.product_id) : null;
                                return (
                                  <li key={item.id}>
                                    {(product?.name ?? 'Produto')} - {getVariationLabel({
                                      description: variation?.description,
                                      baseQuantity: variation?.base_quantity,
                                      baseUnit: product?.base_unit,
                                      fallback: `Variação #${item.product_variation_id}`
                                    })} | Qtde: {item.quantity}
                                  </li>
                                );
                              })}
                            </ul>
                          ) : (
                            <p className='feedback'>Sem itens diretos.</p>
                          )}

                          <h3>Cestas e resolucoes</h3>
                          {details && details.groupItems.length > 0 ? (
                            <ul className='detail-list'>
                              {details.groupItems.map((groupItem) => (
                                <li key={groupItem.id}>
                                  <strong>{groupsByID.get(groupItem.product_group_id)?.name ?? `Cesta #${groupItem.product_group_id}`}</strong>
                                  <p>Quantidade de cestas: {groupItem.quantity}</p>
                                  <ul>
                                    {(details.resolutionsByGroupItem[groupItem.id] ?? []).map((res) => {
                                      const variation = variationsByID.get(res.product_variation_id);
                                      const product = variation ? productsByID.get(variation.product_id) : null;
                                      return (
                                        <li key={res.id}>
                                          {(product?.name ?? 'Produto')} - {getVariationLabel({
                                            description: variation?.description,
                                            baseQuantity: variation?.base_quantity,
                                            baseUnit: product?.base_unit,
                                            fallback: `Variação #${res.product_variation_id}`
                                          })} | Qtde: {res.quantity}
                                        </li>
                                      );
                                    })}
                                  </ul>
                                </li>
                              ))}
                            </ul>
                          ) : (
                            <p className='feedback'>Sem cestas nesta movimentação.</p>
                          )}
                        </>
                      )}
                    </div>
                  )}
                </li>
              );
            })}
          </ul>
        )}
      </section>

      {isModalOpen && (
        <div className='modal-overlay' onClick={closeModal}>
          <section className='modal-content panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>Nova movimentação</h2>
              <button type='button' className='btn btn-ghost' onClick={closeModal} disabled={submitting}>
                Fechar
              </button>
            </div>

            <form className='form' onSubmit={handleSubmit}>
              <label htmlFor='movement-work'>Obra assistencial</label>
              <select
                id='movement-work'
                value={form.assistentialWorkID}
                onChange={(event) => setForm((prev) => ({...prev, assistentialWorkID: event.target.value}))}
                disabled={submitting}
              >
                {assistentialWorks.length === 0 && <option value=''>Nenhuma obra ativa</option>}
                {assistentialWorks.map((work) => (
                  <option key={work.id} value={work.id}>
                    {work.name}
                  </option>
                ))}
              </select>

              <label htmlFor='movement-type'>Tipo</label>
              <select
                id='movement-type'
                value={form.type}
                onChange={(event) => setForm((prev) => ({...prev, type: event.target.value}))}
                disabled={submitting}
              >
                <option value='out'>Saida</option>
                <option value='in'>Entrada</option>
              </select>

              <label htmlFor='movement-notes'>Observacoes</label>
              <textarea
                id='movement-notes'
                value={form.notes}
                onChange={(event) => setForm((prev) => ({...prev, notes: event.target.value}))}
                placeholder='Opcional'
                rows={2}
                disabled={submitting}
              />

              <div className='section-header'>
                <h3>Itens diretos</h3>
                <button type='button' className='btn btn-ghost' onClick={addProductLine} disabled={submitting}>
                  + Item
                </button>
              </div>
              <div className='line-list'>
                {form.productLines.map((line, index) => (
                  <div key={`p-${index}`} className='line-item'>
                    <select
                      value={line.variationId}
                      onChange={(event) => setForm((prev) => ({
                        ...prev,
                        productLines: prev.productLines.map((item, i) => (i === index ? {...item, variationId: event.target.value} : item))
                      }))}
                      disabled={submitting}
                    >
                      <option value=''>Selecione a variação</option>
                      {variationOptions.map((option) => (
                        <option key={option.id} value={option.id}>
                          {option.label} (estoque: {option.stock})
                        </option>
                      ))}
                    </select>
                    <input
                      type='number'
                      min='1'
                      step='1'
                      value={line.quantity}
                      onChange={(event) => setForm((prev) => ({
                        ...prev,
                        productLines: prev.productLines.map((item, i) => (i === index ? {...item, quantity: event.target.value} : item))
                      }))}
                      disabled={submitting}
                    />
                    {form.productLines.length > 1 && (
                      <button type='button' className='btn btn-danger' onClick={() => removeProductLine(index)} disabled={submitting}>
                        Remover
                      </button>
                    )}
                  </div>
                ))}
              </div>

              <div className='section-header'>
                <h3>Cestas</h3>
                <button type='button' className='btn btn-ghost' onClick={addGroupLine} disabled={submitting}>
                  + Cesta
                </button>
              </div>
              <div className='line-list'>
                {form.groupLines.map((line, index) => (
                  <div key={`g-${index}`} className='line-item'>
                    <select
                      value={line.groupId}
                      onChange={(event) => setForm((prev) => ({
                        ...prev,
                        groupLines: prev.groupLines.map((item, i) => (i === index ? {...item, groupId: event.target.value} : item))
                      }))}
                      disabled={submitting}
                    >
                      <option value=''>Selecione a cesta</option>
                      {groups.map((group) => (
                        <option key={group.id} value={group.id}>
                          {group.name}
                        </option>
                      ))}
                    </select>
                    <input
                      type='number'
                      min='1'
                      step='1'
                      value={line.quantity}
                      onChange={(event) => setForm((prev) => ({
                        ...prev,
                        groupLines: prev.groupLines.map((item, i) => (i === index ? {...item, quantity: event.target.value} : item))
                      }))}
                      disabled={submitting}
                    />
                    {form.groupLines.length > 1 && (
                      <button type='button' className='btn btn-danger' onClick={() => removeGroupLine(index)} disabled={submitting}>
                        Remover
                      </button>
                    )}
                  </div>
                ))}
              </div>

              <button className='btn btn-primary' type='submit' disabled={submitting || loadingBaseData}>
                {submitting ? 'Salvando...' : 'Registrar movimentação'}
              </button>
            </form>

            {error && <p className='feedback error'>{error}</p>}
          </section>
        </div>
      )}
    </main>
  );
}

function validateOutStock(productItems, groupItems, variationOptions, groupItemsByGroup) {
  const stockByVariation = new Map();
  const variationsByProduct = new Map();

  variationOptions.forEach((v) => {
    stockByVariation.set(v.id, v.stock);
    if (!variationsByProduct.has(v.productId)) {
      variationsByProduct.set(v.productId, []);
    }
    variationsByProduct.get(v.productId).push(v);
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
      const variations = [...(variationsByProduct.get(comp.product_id) ?? [])].sort((a, b) => b.baseQuantity - a.baseQuantity);

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
        const alloc = Math.min(maxNeededUnits, availableUnits);
        stockByVariation.set(variation.id, availableUnits - alloc);
        remainingBase -= alloc * variation.baseQuantity;
      }

      if (remainingBase > 0) {
        return `Estoque insuficiente para resolver o produto #${comp.product_id} na cesta #${groupItem.product_group_id}.`;
      }
    }
  }

  return '';
}
