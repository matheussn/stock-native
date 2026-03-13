import {useEffect, useState} from 'react';
import {
  BackupDatabase,
  ChooseBackupDestination,
  ChooseRestoreSource,
  GetDatabasePath,
  GetLogFilePath,
  RestoreDatabase
} from '../../../wailsjs/go/main/App';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import './settings-page.css';

export default function SettingsPage() {
  const [databasePath, setDatabasePath] = useState('');
  const [logFilePath, setLogFilePath] = useState('');
  const [loadingPath, setLoadingPath] = useState(true);
  const [busyAction, setBusyAction] = useState('');
  const [error, setError] = useState('');
  const [backupPath, setBackupPath] = useState('');
  const [restoreResult, setRestoreResult] = useState(null);
  const [restoreCandidatePath, setRestoreCandidatePath] = useState('');

  useEffect(() => {
    void loadDatabasePath();
  }, []);

  async function loadDatabasePath() {
    setLoadingPath(true);
    setError('');
    try {
      const [path, logPath] = await Promise.all([GetDatabasePath(), GetLogFilePath()]);
      setDatabasePath(path ?? '');
      setLogFilePath(logPath ?? '');
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoadingPath(false);
    }
  }

  async function handleBackup() {
    setError('');

    try {
      const destinationPath = await ChooseBackupDestination();
      if (!destinationPath) {
        return;
      }

      setBusyAction('backup');
      const savedPath = await BackupDatabase(destinationPath);
      setBackupPath(savedPath);
      setRestoreResult(null);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setBusyAction('');
    }
  }

  async function handleStartRestore() {
    setError('');

    try {
      const sourcePath = await ChooseRestoreSource();
      if (!sourcePath) {
        return;
      }

      setRestoreCandidatePath(sourcePath);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    }
  }

  async function confirmRestore() {
    if (!restoreCandidatePath) {
      return;
    }

    setError('');
    setBusyAction('restore');
    try {
      const result = await RestoreDatabase(restoreCandidatePath);
      setRestoreResult(result);
      setBackupPath('');
      setRestoreCandidatePath('');
      await loadDatabasePath();
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setBusyAction('');
    }
  }

  function cancelRestore() {
    if (busyAction === 'restore') {
      return;
    }
    setRestoreCandidatePath('');
  }

  return (
    <main className='settings-page'>
      <section className='panel settings-shell'>
        <header className='settings-header'>
          <div>
            <h1 className='title'>Configurações</h1>
            <p className='subtitle'>Backup e restauração do banco de dados local</p>
          </div>
        </header>

        <section className='settings-status'>
          <div className='status-block'>
            <div className='status-kicker'>Banco ativo</div>
            <div className='status-path'>{loadingPath ? 'Carregando caminho do banco...' : databasePath || 'Caminho indisponivel'}</div>
          </div>
          <div className='status-block'>
            <div className='status-kicker'>Log operacional</div>
            <div className='status-path'>{loadingPath ? 'Carregando caminho do log...' : logFilePath || 'Caminho indisponivel'}</div>
          </div>
          <p className='status-copy'>
            O sistema usa um banco SQLite local e registra eventos operacionais em arquivo. Mantenha backups periódicos
            fora da pasta padrão da aplicação.
          </p>
        </section>

        {error && <p className='feedback error'>{error}</p>}

        {(backupPath || restoreResult) && (
          <section className='settings-feedback panel'>
            <h2 className='settings-card-title'>Última operação</h2>
            {backupPath && (
              <div className='feedback-block'>
                <span className='feedback-label'>Backup criado em</span>
                <div className='feedback-path'>{backupPath}</div>
              </div>
            )}
            {restoreResult && (
              <>
                <div className='feedback-block'>
                  <span className='feedback-label'>Banco restaurado em</span>
                  <div className='feedback-path'>{restoreResult.restored_path}</div>
                </div>
                {restoreResult.previous_backup_path && (
                  <div className='feedback-block'>
                    <span className='feedback-label'>Cópia de segurança anterior</span>
                    <div className='feedback-path'>{restoreResult.previous_backup_path}</div>
                  </div>
                )}
              </>
            )}
          </section>
        )}

        <div className='settings-grid'>
          <section className='settings-card settings-card-backup'>
            <div className='settings-card-kicker'>Backup</div>
            <h2 className='settings-card-title'>Criar cópia do banco atual</h2>
            <p className='settings-card-copy'>
              Abre o seletor nativo para escolher onde salvar um arquivo `.db` com o estado atual do sistema.
            </p>
            <button type='button' className='btn btn-primary' onClick={handleBackup} disabled={busyAction !== ''}>
              {busyAction === 'backup' ? 'Criando backup...' : 'Criar backup'}
            </button>
          </section>

          <section className='settings-card settings-card-restore'>
            <div className='settings-card-kicker'>Restauração</div>
            <h2 className='settings-card-title'>Restaurar a partir de um backup</h2>
            <p className='settings-card-copy'>
              Substitui o banco ativo pelo arquivo selecionado. Antes da troca, o banco atual e salvo automaticamente
              como cópia de segurança.
            </p>
            <button type='button' className='btn btn-ghost' onClick={handleStartRestore} disabled={busyAction !== ''}>
              {busyAction === 'restore' ? 'Restaurando...' : 'Selecionar backup'}
            </button>
          </section>
        </div>
      </section>

      {restoreCandidatePath && (
        <div className='modal-overlay' onClick={cancelRestore}>
          <section className='modal-content modal-confirm panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>Restaurar banco de dados</h2>
            </div>

            <p className='confirm-text'>
              O banco atual sera substituido pelo arquivo selecionado abaixo.
            </p>
            <div className='confirm-path'>{restoreCandidatePath}</div>
            <p className='confirm-note'>
              Antes da restauração, o sistema salvará automaticamente uma cópia do banco atual na mesma pasta do
              `stock.db`.
            </p>

            <div className='confirm-actions'>
              <button type='button' className='btn btn-ghost' onClick={cancelRestore} disabled={busyAction === 'restore'}>
                Cancelar
              </button>
              <button
                type='button'
                className='btn btn-danger'
                onClick={confirmRestore}
                disabled={busyAction === 'restore'}
              >
                {busyAction === 'restore' ? 'Restaurando...' : 'Restaurar'}
              </button>
            </div>
          </section>
        </div>
      )}
    </main>
  );
}
