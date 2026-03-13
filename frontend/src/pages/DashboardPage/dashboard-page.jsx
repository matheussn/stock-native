import {useEffect, useMemo, useState} from 'react';
import {
  Bar,
  BarChart,
  CartesianGrid,
  Cell,
  Legend,
  Pie,
  PieChart,
  ResponsiveContainer,
  Tooltip,
  XAxis,
  YAxis
} from 'recharts';
import {
  GetAssistentialWorkOutflowKg,
  GetMonthlyMovementFlowKg,
  GetStockCoverageCheck
} from '../../../wailsjs/go/main/App';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import './dashboard-page.css';

const CHART_COLORS = [
  'var(--color-brand)',
  'var(--color-success)',
  'var(--color-warning)',
  'var(--color-info)',
  'var(--color-danger)',
  'var(--color-accent-acao)',
  'var(--color-accent-fii)'
];

const FLOW_COLORS = {
  input: 'var(--color-success)',
  output: 'var(--color-brand)'
};

const COVERAGE_COLORS = {
  demand: 'var(--color-warning)',
  available: 'var(--color-accent-fii)'
};

const kilogramFormatter = new Intl.NumberFormat('pt-BR', {
  minimumFractionDigits: 0,
  maximumFractionDigits: 1
});

const monthFormatter = new Intl.DateTimeFormat('pt-BR', {
  month: 'short',
  year: '2-digit'
});

const DONUT_TOOLTIP_SIZE = {
  width: 170,
  height: 54,
  padding: 8
};

function formatKg(value) {
  return `${kilogramFormatter.format(value)} kg`;
}

function formatMonthKey(monthKey) {
  const [year, month] = monthKey.split('-');
  const date = new Date(Number(year), Number(month) - 1, 1);
  return monthFormatter.format(date).replace('.', '');
}

function clamp(value, min, max) {
  return Math.min(Math.max(value, min), max);
}

function convertToKg(value, unit) {
  if (unit === 'kg') {
    return value;
  }
  if (unit === 'g') {
    return value / 1000;
  }
  return null;
}

function DonutTooltip({active, payload}) {
  if (!active || !payload || payload.length === 0) {
    return null;
  }

  const item = payload[0].payload;
  return (
    <div className='dashboard-tooltip'>
      <div className='dashboard-tooltip-head'>
        <span className='dashboard-tooltip-label'>Obra</span>
        <span className='dashboard-tooltip-label'>Kg</span>
      </div>
      <div className='dashboard-tooltip-row'>
        <div className='dashboard-tooltip-name'>
          <span className='dashboard-tooltip-swatch' style={{background: item.fill}} />
          <strong>{item.assistential_work_name}</strong>
        </div>
        <span className='dashboard-tooltip-value'>{formatKg(item.output_kg)}</span>
      </div>
    </div>
  );
}

function FlowTooltip({active, payload, label}) {
  if (!active || !payload || payload.length === 0) {
    return null;
  }

  return (
    <div className='dashboard-chart-tooltip'>
      <strong className='dashboard-chart-tooltip-title'>{label}</strong>
      <div className='dashboard-chart-tooltip-list'>
        {payload.map((entry) => (
          <div key={entry.dataKey} className='dashboard-chart-tooltip-row'>
            <div className='dashboard-chart-tooltip-name'>
              <span className='dashboard-chart-tooltip-swatch' style={{background: entry.color}} />
              <span>{entry.name}</span>
            </div>
            <span>{formatKg(entry.value)}</span>
          </div>
        ))}
      </div>
    </div>
  );
}

function CoverageTooltip({active, payload, label}) {
  if (!active || !payload || payload.length === 0) {
    return null;
  }

  const shortfall = Math.max((payload[0]?.payload?.requiredKg ?? 0) - (payload[0]?.payload?.availableKg ?? 0), 0);

  return (
    <div className='dashboard-chart-tooltip'>
      <strong className='dashboard-chart-tooltip-title'>{label}</strong>
      <div className='dashboard-chart-tooltip-list'>
        {payload.map((entry) => (
          <div key={entry.dataKey} className='dashboard-chart-tooltip-row'>
            <div className='dashboard-chart-tooltip-name'>
              <span className='dashboard-chart-tooltip-swatch' style={{background: entry.color}} />
              <span>{entry.name}</span>
            </div>
            <span>{formatKg(entry.value)}</span>
          </div>
        ))}
        <div className='dashboard-chart-tooltip-row'>
          <div className='dashboard-chart-tooltip-name'>
            <span className='dashboard-chart-tooltip-swatch dashboard-chart-tooltip-swatch-danger' />
            <span>Gap</span>
          </div>
          <span>{formatKg(shortfall)}</span>
        </div>
      </div>
    </div>
  );
}

export default function DashboardPage() {
  const [distribution, setDistribution] = useState([]);
  const [monthlyFlow, setMonthlyFlow] = useState([]);
  const [coverage, setCoverage] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');
  const [tooltipPosition, setTooltipPosition] = useState(null);

  useEffect(() => {
    void loadData();
  }, []);

  async function loadData() {
    setLoading(true);
    setError('');
    try {
      const [distributionRes, monthlyFlowRes, coverageRes] = await Promise.all([
        GetAssistentialWorkOutflowKg(),
        GetMonthlyMovementFlowKg(),
        GetStockCoverageCheck()
      ]);
      setDistribution(distributionRes);
      setMonthlyFlow(monthlyFlowRes);
      setCoverage(coverageRes);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  const outflowChartData = useMemo(
    () =>
      distribution.map((item, index) => ({
        ...item,
        fill: CHART_COLORS[index % CHART_COLORS.length]
      })),
    [distribution]
  );

  const totalKg = useMemo(
    () => outflowChartData.reduce((sum, item) => sum + item.output_kg, 0),
    [outflowChartData]
  );

  const monthlyFlowChartData = useMemo(
    () =>
      monthlyFlow.slice(-6).map((item) => ({
        ...item,
        label: formatMonthKey(item.month_key)
      })),
    [monthlyFlow]
  );

  const coverageChartData = useMemo(
    () =>
      coverage
        .map((item) => {
          const requiredKg = convertToKg(item.required_base_quantity, item.base_unit);
          const availableKg = convertToKg(item.available_base_quantity, item.base_unit);
          if (requiredKg === null || availableKg === null) {
            return null;
          }
          return {
            productId: item.product_id,
            productName: item.product_name,
            requiredKg,
            availableKg,
            shortfallKg: Math.max(requiredKg - availableKg, 0)
          };
        })
        .filter(Boolean)
        .sort((a, b) => b.requiredKg - a.requiredKg || b.shortfallKg - a.shortfallKg)
        .slice(0, 6),
    [coverage]
  );

  const coverageChartHeight = Math.max(280, coverageChartData.length * 52);

  function handleSliceEnter(slice) {
    if (!slice) {
      return;
    }

    const angle = (-slice.midAngle * Math.PI) / 180;
    const radialDistance = slice.outerRadius + 42;
    const centerX = slice.cx + Math.cos(angle) * radialDistance;
    const centerY = slice.cy + Math.sin(angle) * radialDistance;
    const chartWidth = slice.cx * 2;
    const chartHeight = slice.cy * 2;

    setTooltipPosition({
      x: clamp(centerX - DONUT_TOOLTIP_SIZE.width / 2, DONUT_TOOLTIP_SIZE.padding, chartWidth - DONUT_TOOLTIP_SIZE.width - DONUT_TOOLTIP_SIZE.padding),
      y: clamp(centerY - DONUT_TOOLTIP_SIZE.height / 2, DONUT_TOOLTIP_SIZE.padding, chartHeight - DONUT_TOOLTIP_SIZE.height - DONUT_TOOLTIP_SIZE.padding)
    });
  }

  return (
    <main className='dashboard-page'>
      <header className='dashboard-header'>
        <h1>Dashboard</h1>
        <p>Painel com distribuição das saídas, fluxo mensal e cobertura de estoque em kg.</p>
      </header>

      {error && <p className='dashboard-feedback dashboard-error'>{error}</p>}

      <section className='dashboard-card'>
        <div className='dashboard-card-copy'>
          <span className='dashboard-eyebrow'>Saídas por obra assistencial</span>
        </div>

        {loading ? (
          <p className='dashboard-feedback'>Carregando distribuicao...</p>
        ) : outflowChartData.length === 0 ? (
          <p className='dashboard-feedback'>Nenhuma saída em kg registrada até agora.</p>
        ) : (
          <div className='dashboard-chart-layout'>
            <div className='dashboard-chart-wrap'>
              <ResponsiveContainer width='100%' height='100%'>
                <PieChart>
                  <Pie
                    data={outflowChartData}
                    dataKey='output_kg'
                    nameKey='assistential_work_name'
                    cx='50%'
                    cy='50%'
                    innerRadius={92}
                    outerRadius={128}
                    paddingAngle={4}
                    cornerRadius={10}
                    startAngle={90}
                    endAngle={-270}
                    stroke='var(--bg)'
                    strokeWidth={4}
                    onMouseEnter={handleSliceEnter}
                    onMouseMove={handleSliceEnter}
                    onMouseLeave={() => setTooltipPosition(null)}
                  >
                    {outflowChartData.map((item) => (
                      <Cell key={item.assistential_work_id} fill={item.fill} />
                    ))}
                  </Pie>
                  <Tooltip content={<DonutTooltip />} offset={0} position={tooltipPosition ?? undefined} />
                </PieChart>
              </ResponsiveContainer>

              <div className='dashboard-chart-center'>
                <strong>{formatKg(totalKg)}</strong>
                <span>Total destinado</span>
              </div>
            </div>

            <div className='dashboard-legend'>
              <div className='dashboard-legend-head'>
                <span>Obra</span>
                <span>Kg</span>
              </div>

              <ul className='dashboard-legend-list'>
                {outflowChartData.map((item) => (
                  <li key={item.assistential_work_id} className='dashboard-legend-item'>
                    <div className='dashboard-legend-name'>
                      <span className='dashboard-legend-swatch' style={{background: item.fill}} />
                      <strong>{item.assistential_work_name}</strong>
                    </div>
                    <span className='dashboard-legend-value'>{formatKg(item.output_kg)}</span>
                  </li>
                ))}
              </ul>
            </div>
          </div>
        )}
      </section>

      <section className='dashboard-grid'>
        <article className='dashboard-panel'>
          <div className='dashboard-panel-header'>
            <div>
              <span className='dashboard-panel-kicker'>Fluxo mensal</span>
              <h2>Entradas x saídas por mês (kg)</h2>
            </div>
          </div>

          {loading ? (
            <p className='dashboard-feedback'>Carregando fluxo mensal...</p>
          ) : monthlyFlowChartData.length === 0 ? (
            <p className='dashboard-feedback'>Nenhuma movimentação em kg registrada.</p>
          ) : (
            <div className='dashboard-panel-chart'>
              <ResponsiveContainer width='100%' height='100%'>
                <BarChart data={monthlyFlowChartData} barGap={10}>
                  <CartesianGrid stroke='var(--border)' strokeDasharray='4 4' vertical={false} />
                  <XAxis dataKey='label' tickLine={false} axisLine={false} />
                  <YAxis tickLine={false} axisLine={false} width={46} />
                  <Tooltip content={<FlowTooltip />} cursor={{fill: 'rgba(47, 79, 221, 0.05)'}} />
                  <Legend />
                  <Bar dataKey='input_kg' name='Entradas' fill={FLOW_COLORS.input} radius={[8, 8, 0, 0]} />
                  <Bar dataKey='output_kg' name='Saidas' fill={FLOW_COLORS.output} radius={[8, 8, 0, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}
        </article>

        <article className='dashboard-panel'>
          <div className='dashboard-panel-header'>
            <div>
              <span className='dashboard-panel-kicker'>Cobertura</span>
              <h2>Demanda x disponivel por produto (kg)</h2>
            </div>
          </div>

          {loading ? (
            <p className='dashboard-feedback'>Carregando cobertura...</p>
          ) : coverageChartData.length === 0 ? (
            <p className='dashboard-feedback'>Sem produtos em kg para comparar cobertura.</p>
          ) : (
            <div className='dashboard-panel-chart' style={{height: `${coverageChartHeight}px`}}>
              <ResponsiveContainer width='100%' height='100%'>
                <BarChart data={coverageChartData} layout='vertical' barGap={8} margin={{left: 16}}>
                  <CartesianGrid stroke='var(--border)' strokeDasharray='4 4' horizontal={false} />
                  <XAxis type='number' tickLine={false} axisLine={false} />
                  <YAxis
                    type='category'
                    dataKey='productName'
                    tickLine={false}
                    axisLine={false}
                    width={92}
                  />
                  <Tooltip content={<CoverageTooltip />} cursor={{fill: 'rgba(47, 79, 221, 0.05)'}} />
                  <Legend />
                  <Bar dataKey='requiredKg' name='Demanda' fill={COVERAGE_COLORS.demand} radius={[0, 8, 8, 0]} />
                  <Bar dataKey='availableKg' name='Disponível' fill={COVERAGE_COLORS.available} radius={[0, 8, 8, 0]} />
                </BarChart>
              </ResponsiveContainer>
            </div>
          )}
        </article>
      </section>
    </main>
  );
}
