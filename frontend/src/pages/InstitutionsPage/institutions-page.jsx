import {Pencil, RotateCcw, Trash2} from 'lucide-react';
import {useEffect, useMemo, useState} from 'react';
import {
  CreateInstitution,
  ListInstitutions,
  SetInstitutionActive,
  UpdateInstitution
} from '../../../wailsjs/go/main/App';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import './institutions-page.css';

const FILTERS = {
  active: 'Ativas',
  inactive: 'Inativas',
  all: 'Todas'
};

const INITIAL_FORM = {
  id: null,
  name: '',
  address: '',
  cnpj: '',
  responsibleName: '',
  phone: ''
};

export default function InstitutionsPage() {
  const [institutions, setInstitutions] = useState([]);
  const [filter, setFilter] = useState('active');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [togglingActive, setTogglingActive] = useState(false);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [institutionToDeactivate, setInstitutionToDeactivate] = useState(null);
  const [form, setForm] = useState(INITIAL_FORM);

  const isEditing = form.id !== null;

  const visibleInstitutions = useMemo(() => {
    if (filter === 'inactive') {
      return institutions.filter((item) => item.is_active === 0);
    }
    return institutions;
  }, [filter, institutions]);

  useEffect(() => {
    void loadInstitutions(filter);
  }, [filter]);

  async function loadInstitutions(currentFilter) {
    setLoading(true);
    setError('');
    try {
      const includeInactive = currentFilter !== 'active';
      const response = await ListInstitutions(includeInactive);
      setInstitutions(response);
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
      setError('Informe o nome da instituição.');
      return;
    }

    setSubmitting(true);
    setError('');
    try {
      if (isEditing) {
        await UpdateInstitution(form.id, name, form.address, form.cnpj, form.responsibleName, form.phone);
      } else {
        await CreateInstitution(name, form.address, form.cnpj, form.responsibleName, form.phone);
      }
      resetForm();
      setIsModalOpen(false);
      await loadInstitutions(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function handleToggleActive(item) {
    if (item.is_active === 1) {
      setInstitutionToDeactivate(item);
      return;
    }

    setError('');
    try {
      setTogglingActive(true);
      await SetInstitutionActive(item.id, true);
      await loadInstitutions(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setTogglingActive(false);
    }
  }

  async function confirmDeactivate() {
    if (!institutionToDeactivate) {
      return;
    }

    setError('');
    try {
      setTogglingActive(true);
      await SetInstitutionActive(institutionToDeactivate.id, false);
      setInstitutionToDeactivate(null);
      await loadInstitutions(filter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setTogglingActive(false);
    }
  }

  function startCreate() {
    resetForm();
    setError('');
    setIsModalOpen(true);
  }

  function startEdit(item) {
    setForm({
      id: item.id,
      name: item.name ?? '',
      address: item.address ?? '',
      cnpj: item.cnpj ?? '',
      responsibleName: item.responsible_name ?? '',
      phone: item.phone ?? ''
    });
    setError('');
    setIsModalOpen(true);
  }

  function closeModal() {
    setIsModalOpen(false);
    resetForm();
  }

  function resetForm() {
    setForm(INITIAL_FORM);
  }

  return (
    <main className='institutions-page'>
      <section className='panel panel-list'>
        <div className='list-header'>
          <div>
            <h1 className='title'>Instituições</h1>
            <p className='subtitle'>Cadastro e manutenção</p>
          </div>
          <button type='button' className='btn btn-primary' onClick={startCreate}>
            + Novo
          </button>
        </div>

        <div className='list-tools'>
          <div className='filters' role='group' aria-label='Filtro de instituições'>
            {Object.entries(FILTERS).map(([key, label]) => (
              <button
                key={key}
                type='button'
                className={`chip ${filter === key ? 'chip-active' : ''}`}
                onClick={() => setFilter(key)}
              >
                {label}
              </button>
            ))}
          </div>
        </div>

        {error && <p className='feedback error'>{error}</p>}

        {loading ? (
          <p className='feedback'>Carregando instituições...</p>
        ) : visibleInstitutions.length === 0 ? (
          <p className='feedback'>Nenhum registro encontrado para este filtro.</p>
        ) : (
          <ul className='institution-list'>
            {visibleInstitutions.map((item) => {
              const metadata = [
                item.cnpj ? `CNPJ: ${item.cnpj}` : '',
                item.responsible_name ? `Responsável: ${item.responsible_name}` : '',
                item.phone ? `Telefone: ${item.phone}` : '',
                item.address ? `Endereço: ${item.address}` : ''
              ].filter(Boolean);

              return (
                <li key={item.id} className='institution-item'>
                  <div>
                    <strong>{item.name}</strong>
                    {metadata.length > 0 ? (
                      <div className='institution-meta'>
                        {metadata.map((line) => (
                          <p key={line}>{line}</p>
                        ))}
                      </div>
                    ) : (
                      <p className='institution-empty'>Sem dados complementares.</p>
                    )}
                  </div>

                  <div className='item-actions'>
                    <span className={`status ${item.is_active === 1 ? 'status-on' : 'status-off'}`}>
                      {item.is_active === 1 ? 'Ativa' : 'Inativa'}
                    </span>
                    <button
                      type='button'
                      className='icon-action icon-action-edit'
                      onClick={() => startEdit(item)}
                      aria-label='Editar instituição'
                      title='Editar instituição'
                    >
                      <Pencil aria-hidden='true' />
                    </button>
                    <button
                      type='button'
                      className={`icon-action ${item.is_active === 1 ? 'icon-action-delete' : 'icon-action-restore'}`}
                      onClick={() => handleToggleActive(item)}
                      aria-label={item.is_active === 1 ? 'Desativar instituição' : 'Reativar instituição'}
                      title={item.is_active === 1 ? 'Desativar instituição' : 'Reativar instituição'}
                    >
                      {item.is_active === 1 ? <Trash2 aria-hidden='true' /> : <RotateCcw aria-hidden='true' />}
                    </button>
                  </div>
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
              <h2>{isEditing ? 'Editar instituição' : 'Nova instituição'}</h2>
              <button type='button' className='btn btn-ghost' onClick={closeModal} disabled={submitting}>
                Fechar
              </button>
            </div>

            <form className='form' onSubmit={handleSubmit}>
              <label htmlFor='institution-name'>Nome</label>
              <input
                id='institution-name'
                value={form.name}
                onChange={(event) => setForm((prev) => ({...prev, name: event.target.value}))}
                placeholder='Ex: Casa Esperança'
                maxLength={120}
                disabled={submitting}
              />

              <label htmlFor='institution-cnpj'>CNPJ</label>
              <input
                id='institution-cnpj'
                value={form.cnpj}
                onChange={(event) => setForm((prev) => ({...prev, cnpj: event.target.value}))}
                placeholder='Opcional'
                maxLength={32}
                disabled={submitting}
              />

              <label htmlFor='institution-responsible'>Nome do responsável</label>
              <input
                id='institution-responsible'
                value={form.responsibleName}
                onChange={(event) => setForm((prev) => ({...prev, responsibleName: event.target.value}))}
                placeholder='Opcional'
                maxLength={120}
                disabled={submitting}
              />

              <label htmlFor='institution-phone'>Telefone</label>
              <input
                id='institution-phone'
                value={form.phone}
                onChange={(event) => setForm((prev) => ({...prev, phone: event.target.value}))}
                placeholder='Opcional'
                maxLength={40}
                disabled={submitting}
              />

              <label htmlFor='institution-address'>Endereço</label>
              <textarea
                id='institution-address'
                value={form.address}
                onChange={(event) => setForm((prev) => ({...prev, address: event.target.value}))}
                placeholder='Opcional'
                rows={3}
                maxLength={400}
                disabled={submitting}
              />

              <div className='form-actions'>
                <button className='btn btn-primary' type='submit' disabled={submitting}>
                  {isEditing ? 'Salvar edição' : 'Criar instituição'}
                </button>
                {isEditing && (
                  <button className='btn btn-ghost' type='button' onClick={closeModal} disabled={submitting}>
                    Cancelar
                  </button>
                )}
              </div>
            </form>
          </section>
        </div>
      )}

      {institutionToDeactivate && (
        <div className='modal-overlay' onClick={() => setInstitutionToDeactivate(null)}>
          <section className='modal-content modal-confirm panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>Desativar instituição</h2>
            </div>
            <p className='confirm-text'>
              Deseja desativar <strong>{institutionToDeactivate.name}</strong>?
            </p>
            <div className='confirm-actions'>
              <button
                type='button'
                className='btn btn-ghost'
                onClick={() => setInstitutionToDeactivate(null)}
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
