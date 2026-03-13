import {useEffect, useMemo, useState} from 'react';
import {
  GetMonthlyDemandProjection,
  GetStockCoverageCheck,
  ListVariationStockStatus
} from '../../../wailsjs/go/main/App';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import './stock-page.css';

export default function StockPage() {
  const [variationStatus, setVariationStatus] = useState([]);
  const [monthlyDemand, setMonthlyDemand] = useState([]);
  const [coverage, setCoverage] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  const totals = useMemo(() => {
    const critical = coverage.filter((c) => !c.is_covered).length;
    return {
      totalVariations: variationStatus.length,
      uncoveredProducts: critical
    };
  }, [variationStatus, coverage]);

  useEffect(() => {
    void loadData();
  }, []);

  async function loadData() {
    setLoading(true);
    setError('');
    try {
      const [statusRes, demandRes, coverageRes] = await Promise.all([
        ListVariationStockStatus(false),
        GetMonthlyDemandProjection(),
        GetStockCoverageCheck()
      ]);
      setVariationStatus(statusRes);
      setMonthlyDemand(demandRes);
      setCoverage(coverageRes);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  return (
    <main className='stock-page'>
      <section className='cards'>
        <article className='panel card'>
          <h2>Variações monitoradas</h2>
          <strong>{totals.totalVariations}</strong>
        </article>
        <article className='panel card'>
          <h2>Produtos sem cobertura</h2>
          <strong className={totals.uncoveredProducts > 0 ? 'danger' : ''}>{totals.uncoveredProducts}</strong>
        </article>
      </section>

      {error && <p className='feedback error'>{error}</p>}

      <section className='panel table-panel'>
        <h2>Estoque por variação</h2>
        {loading ? (
          <p className='feedback'>Carregando estoque...</p>
        ) : variationStatus.length === 0 ? (
          <p className='feedback'>Nenhuma variação encontrada.</p>
        ) : (
          <div className='table-wrap'>
            <table>
              <thead>
                <tr>
                  <th>Produto</th>
                  <th>Variação</th>
                  <th>Atual</th>
                </tr>
              </thead>
              <tbody>
                {variationStatus.map((item) => (
                  <tr key={item.variation_id}>
                    <td>{item.product_name}</td>
                    <td>{item.base_quantity} {item.base_unit}</td>
                    <td>{item.current_stock}</td>
                  </tr>
                ))}
              </tbody>
            </table>
          </div>
        )}
      </section>

      <section className='grid-2'>
        <section className='panel table-panel'>
          <h2>Projeção mensal</h2>
          {loading ? (
            <p className='feedback'>Carregando projeção...</p>
          ) : monthlyDemand.length === 0 ? (
            <p className='feedback'>Sem atribuições ativas de cestas.</p>
          ) : (
            <div className='table-wrap'>
              <table>
                <thead>
                  <tr>
                    <th>Produto</th>
                    <th>Demanda</th>
                  </tr>
                </thead>
                <tbody>
                  {monthlyDemand.map((item) => (
                    <tr key={item.product_id}>
                      <td>{item.product_name}</td>
                      <td>{item.required_base_quantity} {item.base_unit}</td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>

        <section className='panel table-panel'>
          <h2>Cobertura</h2>
          {loading ? (
            <p className='feedback'>Carregando cobertura...</p>
          ) : coverage.length === 0 ? (
            <p className='feedback'>Sem dados para comparar cobertura.</p>
          ) : (
            <div className='table-wrap'>
              <table>
                <thead>
                  <tr>
                    <th>Produto</th>
                    <th>Demanda</th>
                    <th>Disponível</th>
                    <th>Gap</th>
                  </tr>
                </thead>
                <tbody>
                  {coverage.map((item) => (
                    <tr key={item.product_id}>
                      <td>{item.product_name}</td>
                      <td>{item.required_base_quantity} {item.base_unit}</td>
                      <td>{item.available_base_quantity} {item.base_unit}</td>
                      <td className={item.shortfall_base_quantity > 0 ? 'danger' : ''}>
                        {item.shortfall_base_quantity} {item.base_unit}
                      </td>
                    </tr>
                  ))}
                </tbody>
              </table>
            </div>
          )}
        </section>
      </section>
    </main>
  );
}
