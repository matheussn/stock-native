import {useEffect, useMemo, useState} from 'react';
import {
  CreateAssistentialWork,
  ListAssistentialWorks,
  SetAssistentialWorkActive,
  UpdateAssistentialWork
} from '../../../wailsjs/go/main/App';
import AssistentialWorkFilters from '../../components/AssistentialWorkFilters';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import AssistentialWorkForm from '../../components/AssistentialWorkForm';
import AssistentialWorkList from '../../components/AssistentialWorkList';
import './assistential-works-page.css';

const FILTERS = {
  active: 'Ativos',
  inactive: 'Inativos',
  all: 'Todos'
};

export default function AssistentialWorksPage() {
  const [works, setWorks] = useState([]);
  const [filter, setFilter] = useState('active');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [togglingActive, setTogglingActive] = useState(false);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [workToDeactivate, setWorkToDeactivate] = useState(null);
  const [form, setForm] = useState({
    id: null,
    name: '',
    description: ''
  });

  const isEditing = form.id !== null;

  const visibleWorks = useMemo(() => {
    if (filter === 'inactive') {
      return works.filter((item) => item.is_active === 0);
    }
    return works;
  }, [filter, works]);

  useEffect(() => {
    void loadWorks(filter);
  }, [filter]);

  async function loadWorks(currentFilter) {
    setLoading(true);
    setError('');
    try {
      const includeInactive = currentFilter !== 'active';
      const response = await ListAssistentialWorks(includeInactive);
      setWorks(response);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  async function handleSubmit(event) {
    event.preventDefault();
    const name = form.name.trim();

    if (!name) {
      setError('Informe o nome da obra assistencial.');
      return;
    }

    setSubmitting(true);
    setError('');
    try {
      if (isEditing) {
        await UpdateAssistentialWork(form.id, name, form.description);
      } else {
        await CreateAssistentialWork(name, form.description);
      }
      resetForm();
      setIsModalOpen(false);
      await loadWorks(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function handleToggleActive(item) {
    if (item.is_active === 1) {
      setWorkToDeactivate(item);
      return;
    }

    setError('');
    try {
      setTogglingActive(true);
      await SetAssistentialWorkActive(item.id, item.is_active === 0);
      await loadWorks(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setTogglingActive(false);
    }
  }

  async function confirmDeactivate() {
    if (!workToDeactivate) {
      return;
    }

    setError('');
    try {
      setTogglingActive(true);
      await SetAssistentialWorkActive(workToDeactivate.id, false);
      setWorkToDeactivate(null);
      await loadWorks(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setTogglingActive(false);
    }
  }

  function startEdit(item) {
    setForm({
      id: item.id,
      name: item.name,
      description: item.description ?? ''
    });
    setError('');
    setIsModalOpen(true);
  }

  function resetForm() {
    setForm({
      id: null,
      name: '',
      description: ''
    });
  }

  function startCreate() {
    resetForm();
    setError('');
    setIsModalOpen(true);
  }

  function closeModal() {
    setIsModalOpen(false);
    resetForm();
  }

  return (
    <main className='assistential-works-page'>
      <section className='panel panel-list'>
        <div className='list-header'>
          <div>
            <h1 className='title'>Obras Assistenciais</h1>
            <p className='subtitle'>Cadastro e manutenção</p>
          </div>
          <button type='button' className='btn btn-primary' onClick={startCreate}>
            + Novo
          </button>
        </div>

        <div className='list-tools'>
          <AssistentialWorkFilters filters={FILTERS} currentFilter={filter} onChange={setFilter} />
        </div>

        {error && <p className='feedback error'>{error}</p>}

        <AssistentialWorkList
          works={visibleWorks}
          loading={loading}
          onEdit={startEdit}
          onToggleActive={handleToggleActive}
        />
      </section>

      {isModalOpen && (
        <div className='modal-overlay' onClick={closeModal}>
          <section className='modal-content panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>{isEditing ? 'Editar obra assistencial' : 'Nova obra assistencial'}</h2>
              <button type='button' className='btn btn-ghost' onClick={closeModal} disabled={submitting}>
                Fechar
              </button>
            </div>

            <AssistentialWorkForm
              form={form}
              isEditing={isEditing}
              submitting={submitting}
              onNameChange={(value) => setForm((prev) => ({...prev, name: value}))}
              onDescriptionChange={(value) => setForm((prev) => ({...prev, description: value}))}
              onSubmit={handleSubmit}
              onCancel={closeModal}
            />
          </section>
        </div>
      )}

      {workToDeactivate && (
        <div className='modal-overlay' onClick={() => setWorkToDeactivate(null)}>
          <section className='modal-content modal-confirm panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>Desativar obra assistencial</h2>
            </div>
            <p className='confirm-text'>
              Deseja desativar <strong>{workToDeactivate.name}</strong>?
            </p>
            <div className='confirm-actions'>
              <button
                type='button'
                className='btn btn-ghost'
                onClick={() => setWorkToDeactivate(null)}
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
