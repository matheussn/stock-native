import {lazy, Suspense, useCallback, useEffect, useMemo, useState} from 'react';
import {NavLink, Navigate, Route, Routes, useLocation, useNavigate} from 'react-router-dom';
import {
  ChooseRestoreSource,
  GetStartupStatus,
  RestoreDatabase,
  RetryDatabaseInitialization
} from '../wailsjs/go/main/App';
import AssistentialWorksPage from './pages/AssistentialWorksPage';
import FamiliesPage from './pages/FamiliesPage';
import InstitutionsPage from './pages/InstitutionsPage';
import MovementCreatePage from './pages/MovementsPage/movement-create-page';
import MovementsPage from './pages/MovementsPage';
import ProductGroupsPage from './pages/ProductGroupsPage';
import ProductsPage from './pages/ProductsPage';
import SettingsPage from './pages/SettingsPage';
import StockPage from './pages/StockPage';
import {AppNavigationGuardContext} from './utils/app-navigation-guard';
import {getFriendlyErrorMessage} from './utils/error-message';
import './App.css';

const DashboardPage = lazy(() => import('./pages/DashboardPage'));

const navigationItems = [
  {to: '/dashboard', label: 'Dashboard'},
  {to: '/obras-assistenciais', label: 'Obras Assistenciais'},
  {to: '/instituicoes', label: 'Instituições'},
  {to: '/produtos', label: 'Produtos'},
  {to: '/grupos', label: 'Grupos'},
  {to: '/familias', label: 'Famílias'},
  {to: '/movimentacoes', label: 'Movimentações'},
  {to: '/estoque', label: 'Estoque'},
  {to: '/configuracoes', label: 'Configurações'}
];

function App() {
  const navigate = useNavigate();
  const location = useLocation();
  const [startupStatus, setStartupStatus] = useState(null);
  const [startupLoading, setStartupLoading] = useState(true);
  const [startupBusy, setStartupBusy] = useState('');
  const [startupActionError, setStartupActionError] = useState('');
  const [restoreCandidatePath, setRestoreCandidatePath] = useState('');
  const [beforeLeaveHandler, setBeforeLeaveHandler] = useState(null);

  useEffect(() => {
    void loadStartupStatus();
  }, []);

  async function loadStartupStatus() {
    setStartupLoading(true);
    setStartupActionError('');

    try {
      const status = await GetStartupStatus();
      setStartupStatus(status);
    } catch (err) {
      setStartupStatus({
        ready: false,
        message: getFriendlyErrorMessage(err),
        technical_error: '',
        database_path: ''
      });
    } finally {
      setStartupLoading(false);
    }
  }

  async function handleRetryStartup() {
    setStartupBusy('retry');
    setStartupActionError('');

    try {
      const status = await RetryDatabaseInitialization();
      setStartupStatus(status);
    } catch (err) {
      setStartupActionError(getFriendlyErrorMessage(err));
    } finally {
      setStartupBusy('');
    }
  }

  async function handleSelectRestoreBackup() {
    setStartupActionError('');

    try {
      const sourcePath = await ChooseRestoreSource();
      if (!sourcePath) {
        return;
      }

      setRestoreCandidatePath(sourcePath);
    } catch (err) {
      setStartupActionError(getFriendlyErrorMessage(err));
    }
  }

  async function confirmStartupRestore() {
    if (!restoreCandidatePath) {
      return;
    }

    setStartupBusy('restore');
    setStartupActionError('');

    try {
      await RestoreDatabase(restoreCandidatePath);
      setRestoreCandidatePath('');
      await loadStartupStatus();
    } catch (err) {
      setStartupActionError(getFriendlyErrorMessage(err));
    } finally {
      setStartupBusy('');
    }
  }

  function cancelRestoreCandidate() {
    if (startupBusy === 'restore') {
      return;
    }

    setRestoreCandidatePath('');
  }

  const registerBeforeLeaveHandler = useCallback((handler) => {
    if (!handler) {
      setBeforeLeaveHandler(null);
      return;
    }

    setBeforeLeaveHandler(() => handler);
  }, []);

  const requestNavigation = useCallback((to) => {
    if (!to || to === location.pathname) {
      return;
    }

    if (beforeLeaveHandler) {
      const handled = beforeLeaveHandler(to);
      if (handled) {
        return;
      }
    }

    navigate(to);
  }, [beforeLeaveHandler, location.pathname, navigate]);

  const navigationGuardValue = useMemo(() => ({
    registerBeforeLeaveHandler,
    requestNavigation
  }), [registerBeforeLeaveHandler, requestNavigation]);

  if (startupLoading) {
    return (
      <div className='app-startup-state'>
        <p className='app-loading'>Preparando banco de dados...</p>
      </div>
    );
  }

  if (!startupStatus?.ready) {
    return (
      <>
        <div className='app-startup-state'>
          <section className='app-startup-panel'>
            <div className='app-startup-kicker'>Banco de dados</div>
            <h1 className='app-startup-title'>Não foi possível iniciar o sistema</h1>
            <p className='app-startup-copy'>
              {startupStatus?.message || 'Falha ao preparar o banco de dados local.'}
            </p>

            {startupStatus?.database_path && <div className='app-startup-path'>{startupStatus.database_path}</div>}

            {startupActionError && <p className='app-startup-feedback'>{startupActionError}</p>}

            <div className='app-startup-actions'>
              <button
                type='button'
                className='btn btn-primary'
                onClick={handleRetryStartup}
                disabled={startupBusy !== ''}
              >
                {startupBusy === 'retry' ? 'Tentando novamente...' : 'Tentar novamente'}
              </button>
              <button
                type='button'
                className='btn btn-ghost'
                onClick={handleSelectRestoreBackup}
                disabled={startupBusy !== ''}
              >
                Selecionar backup
              </button>
            </div>

            {startupStatus?.technical_error && (
              <details className='app-startup-details'>
                <summary>Detalhes técnicos</summary>
                <pre>{startupStatus.technical_error}</pre>
              </details>
            )}
          </section>
        </div>

        {restoreCandidatePath && (
          <div className='app-modal-overlay' onClick={cancelRestoreCandidate}>
            <section className='app-modal-content' onClick={(event) => event.stopPropagation()}>
              <h2>Restaurar banco de dados</h2>
              <p className='app-startup-copy'>
                O sistema vai substituir o banco atual pelo backup selecionado e preservar uma cópia de segurança do
                banco atual antes da troca.
              </p>
              <div className='app-modal-path'>{restoreCandidatePath}</div>
              <div className='app-modal-actions'>
                <button
                  type='button'
                  className='btn btn-ghost'
                  onClick={cancelRestoreCandidate}
                  disabled={startupBusy === 'restore'}
                >
                  Cancelar
                </button>
                <button
                  type='button'
                  className='btn btn-danger'
                  onClick={confirmStartupRestore}
                  disabled={startupBusy === 'restore'}
                >
                  {startupBusy === 'restore' ? 'Restaurando...' : 'Restaurar backup'}
                </button>
              </div>
            </section>
          </div>
        )}
      </>
    );
  }

  return (
    <AppNavigationGuardContext.Provider value={navigationGuardValue}>
      <div className='app-layout'>
      <aside className='app-sidebar'>
        <div className='brand'>
          <strong>Centro Espírita</strong>
          <span>Estoque</span>
        </div>

        <nav className='menu'>
          {navigationItems.map((item) => (
            <NavLink
              key={item.to}
              to={item.to}
              className={({isActive}) => `menu-link${isActive ? ' menu-link-active' : ''}`}
              onClick={(event) => {
                event.preventDefault();
                requestNavigation(item.to);
              }}
            >
              {item.label}
            </NavLink>
          ))}
        </nav>
      </aside>

      <div className='app-content'>
        <Suspense fallback={<p className='app-loading'>Carregando dashboard...</p>}>
          <Routes>
            <Route path='/' element={<Navigate to='/dashboard' replace />} />
            <Route path='/dashboard' element={<DashboardPage />} />
            <Route path='/obras-assistenciais' element={<AssistentialWorksPage />} />
            <Route path='/instituicoes' element={<InstitutionsPage />} />
            <Route path='/produtos' element={<ProductsPage />} />
            <Route path='/grupos' element={<ProductGroupsPage />} />
            <Route path='/familias' element={<FamiliesPage />} />
            <Route path='/movimentacoes' element={<MovementsPage />} />
            <Route path='/movimentacoes/nova' element={<MovementCreatePage />} />
            <Route path='/estoque' element={<StockPage />} />
            <Route path='/configuracoes' element={<SettingsPage />} />
            <Route path='*' element={<Navigate to='/dashboard' replace />} />
          </Routes>
        </Suspense>
      </div>
      </div>
    </AppNavigationGuardContext.Provider>
  );
}

export default App;
