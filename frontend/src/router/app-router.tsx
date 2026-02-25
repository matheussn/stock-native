import { BackupSection } from '@/domains/backup/backup-section'
import logo from '@/assets/images/logo-universal.png'
import { CatalogSection } from '@/domains/catalog/catalog-section'
import { DashboardSection } from '@/domains/dashboard/dashboard-section'
import { MovementSection } from '@/domains/movement/movement-section'
import { ReportSection } from '@/domains/report/report-section'
import { StockSection } from '@/domains/stock/stock-section'
import { useUiStore } from '@/stores/ui.store'

const sections = [
  { key: 'dashboard', label: 'Painel' },
  { key: 'estoque', label: 'Estoque' },
  { key: 'cadastros', label: 'Cadastros' },
  { key: 'movimentacoes', label: 'Movimentações' },
  { key: 'relatorios', label: 'Relatórios' },
  { key: 'backup', label: 'Backup' },
] as const

export const AppRouter = () => {
  const section = useUiStore((state) => state.section)
  const setSection = useUiStore((state) => state.setSection)

  return (
    <main className='layout'>
      <aside className='sidebar'>
        <div className='brand'>
          <img src={logo} alt='Logo Estoque Comunitário' className='brand-logo' />
          <div>
            <h1>Estoque Comunitário</h1>
            <p>Painel de controle</p>
          </div>
        </div>

        <nav className='sidebar-nav'>
          {sections.map((item) => (
            <button
              key={item.key}
              type='button'
              className={section === item.key ? 'active' : ''}
              onClick={() => setSection(item.key)}
            >
              {item.label}
            </button>
          ))}
        </nav>
      </aside>

      <section className='content'>
        {section === 'dashboard' ? <DashboardSection /> : null}
        {section === 'estoque' ? <StockSection /> : null}
        {section === 'cadastros' ? <CatalogSection /> : null}
        {section === 'movimentacoes' ? <MovementSection /> : null}
        {section === 'relatorios' ? <ReportSection /> : null}
        {section === 'backup' ? <BackupSection /> : null}
      </section>
    </main>
  )
}
