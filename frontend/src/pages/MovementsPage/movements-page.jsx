import {useEffect, useMemo, useState} from 'react';
import {useNavigate} from 'react-router-dom';
import {
  ListAssistentialWorks,
  ListInstitutions,
  ListMovementGroupItemResolutions,
  ListMovementGroupItems,
  ListMovementProductItems,
  ListMovements,
  ListProductGroups,
  ListProductVariationsByProduct,
  ListProducts
} from '../../../wailsjs/go/main/App';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import {
  allowMovementCreateAccess,
  clearMovementCreateAccess
} from '../../utils/movement-create-access';
import {getVariationLabel} from '../../utils/variation-label';
import './movements-page.css';

export default function MovementsPage() {
  const navigate = useNavigate();

  const [assistentialWorks, setAssistentialWorks] = useState([]);
  const [institutions, setInstitutions] = useState([]);
  const [products, setProducts] = useState([]);
  const [groups, setGroups] = useState([]);
  const [variationsByProduct, setVariationsByProduct] = useState({});
  const [movements, setMovements] = useState([]);
  const [detailsByMovement, setDetailsByMovement] = useState({});

  const [historyTypeFilter, setHistoryTypeFilter] = useState('all');
  const [historyWorkFilter, setHistoryWorkFilter] = useState('all');
  const [expandedMovementId, setExpandedMovementId] = useState(null);

  const [loadingMovements, setLoadingMovements] = useState(true);
  const [error, setError] = useState('');

  const worksByID = useMemo(() => {
    const map = new Map();
    assistentialWorks.forEach((work) => map.set(work.id, work));
    return map;
  }, [assistentialWorks]);

  const institutionsByID = useMemo(() => {
    const map = new Map();
    institutions.forEach((institution) => map.set(institution.id, institution));
    return map;
  }, [institutions]);

  const productsByID = useMemo(() => {
    const map = new Map();
    products.forEach((product) => map.set(product.id, product));
    return map;
  }, [products]);

  const groupsByID = useMemo(() => {
    const map = new Map();
    groups.forEach((group) => map.set(group.id, group));
    return map;
  }, [groups]);

  const variationsByID = useMemo(() => {
    const map = new Map();
    Object.values(variationsByProduct).flat().forEach((variation) => map.set(variation.id, variation));
    return map;
  }, [variationsByProduct]);

  useEffect(() => {
    clearMovementCreateAccess();
    void loadBaseData();
  }, []);

  useEffect(() => {
    void loadMovements(historyTypeFilter, historyWorkFilter);
  }, [historyTypeFilter, historyWorkFilter]);

  async function loadBaseData() {
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

      const variationEntries = await Promise.all(
        productsRes.map(async (product) => [product.id, await ListProductVariationsByProduct(product.id, false)])
      );

      const variationMap = {};
      variationEntries.forEach(([productID, variations]) => {
        variationMap[productID] = variations;
      });

      setVariationsByProduct(variationMap);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
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
        groupItems.map(async (groupItem) => [
          groupItem.id,
          await ListMovementGroupItemResolutions(groupItem.id)
        ])
      );

      const resolutionsByGroupItem = {};
      resolutionEntries.forEach(([groupItemID, resolutions]) => {
        resolutionsByGroupItem[groupItemID] = resolutions;
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

  function openCreatePage() {
    allowMovementCreateAccess();
    navigate('/movimentacoes/nova');
  }

  return (
    <main className='movements-page'>
      <section className='panel panel-list'>
        <div className='page-header'>
          <div>
            <h1 className='title'>Movimentações</h1>
            <p className='subtitle'>Entrada e saída com itens diretos e cestas</p>
          </div>
          <button type='button' className='btn btn-primary' onClick={openCreatePage}>
            + Novo
          </button>
        </div>

        <div className='list-header'>
          <h2>Histórico</h2>
          <div className='history-filters'>
            <select value={historyTypeFilter} onChange={(event) => setHistoryTypeFilter(event.target.value)}>
              <option value='all'>Todos os tipos</option>
              <option value='in'>Entradas</option>
              <option value='out'>Saídas</option>
            </select>
            <select value={historyWorkFilter} onChange={(event) => setHistoryWorkFilter(event.target.value)}>
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
                      <strong>#{movement.id} - {movement.type === 'in' ? 'Entrada' : 'Saída'}</strong>
                      <p>Obra: {worksByID.get(movement.assistential_work_id)?.name ?? `#${movement.assistential_work_id}`}</p>
                      {movement.institution_id > 0 && (
                        <p>Instituição: {institutionsByID.get(movement.institution_id)?.name ?? `#${movement.institution_id}`}</p>
                      )}
                      <p>Data: {movement.created_at}</p>
                      <p>Obs.: {movement.notes || 'Sem observações'}</p>
                    </div>
                    <div className='item-actions'>
                      <span className={`status ${movement.type === 'in' ? 'status-in' : 'status-out'}`}>
                        {movement.type === 'in' ? 'Entrada' : 'Saída'}
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
                          <div className='movement-summary'>
                            <p>Obra: {worksByID.get(movement.assistential_work_id)?.name ?? `#${movement.assistential_work_id}`}</p>
                            {movement.institution_id > 0 && (
                              <p>Instituição: {institutionsByID.get(movement.institution_id)?.name ?? `#${movement.institution_id}`}</p>
                            )}
                          </div>

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

                          <h3>Cestas e resoluções</h3>
                          {details && details.groupItems.length > 0 ? (
                            <ul className='detail-list'>
                              {details.groupItems.map((groupItem) => (
                                <li key={groupItem.id}>
                                  <strong>{groupsByID.get(groupItem.product_group_id)?.name ?? `Cesta #${groupItem.product_group_id}`}</strong>
                                  <p>Quantidade de cestas: {groupItem.quantity}</p>
                                  <ul>
                                    {(details.resolutionsByGroupItem[groupItem.id] ?? []).map((resolution) => {
                                      const variation = variationsByID.get(resolution.product_variation_id);
                                      const product = variation ? productsByID.get(variation.product_id) : null;

                                      return (
                                        <li key={resolution.id}>
                                          {(product?.name ?? 'Produto')} - {getVariationLabel({
                                            description: variation?.description,
                                            baseQuantity: variation?.base_quantity,
                                            baseUnit: product?.base_unit,
                                            fallback: `Variação #${resolution.product_variation_id}`
                                          })} | Qtde: {resolution.quantity}
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
    </main>
  );
}
