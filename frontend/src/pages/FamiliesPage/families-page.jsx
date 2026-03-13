import {useEffect, useMemo, useState} from 'react';
import {
  AssignFamilyGroup,
  CreateFamily,
  ListAssistentialWorks,
  ListFamilies,
  ListFamilyGroupAssignments,
  ListProductGroups,
  SetFamilyActive,
  UpdateFamily
} from '../../../wailsjs/go/main/App';
import AssistentialWorkFilters from '../../components/AssistentialWorkFilters';
import {getFriendlyErrorMessage} from '../../utils/error-message';
import './families-page.css';

const STATUS_FILTERS = {
  active: 'Ativas',
  inactive: 'Inativas',
  all: 'Todas'
};

export default function FamiliesPage() {
  const [assistentialWorks, setAssistentialWorks] = useState([]);
  const [productGroups, setProductGroups] = useState([]);
  const [families, setFamilies] = useState([]);
  const [statusFilter, setStatusFilter] = useState('active');
  const [workFilter, setWorkFilter] = useState('all');
  const [loading, setLoading] = useState(true);
  const [submitting, setSubmitting] = useState(false);
  const [assignmentSubmitting, setAssignmentSubmitting] = useState(false);
  const [togglingActive, setTogglingActive] = useState(false);
  const [error, setError] = useState('');
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [familyToDeactivate, setFamilyToDeactivate] = useState(null);
  const [assignments, setAssignments] = useState([]);
  const [loadingAssignments, setLoadingAssignments] = useState(false);
  const [assignmentForm, setAssignmentForm] = useState({
    productGroupId: ''
  });
  const [form, setForm] = useState({
    id: null,
    assistentialWorkID: '',
    name: '',
    memberCount: '1',
    address: '',
    contact: ''
  });

  const isEditing = form.id !== null;

  const worksByID = useMemo(() => {
    const map = new Map();
    assistentialWorks.forEach((w) => map.set(w.id, w));
    return map;
  }, [assistentialWorks]);

  const groupsByID = useMemo(() => {
    const map = new Map();
    productGroups.forEach((g) => map.set(g.id, g));
    return map;
  }, [productGroups]);

  useEffect(() => {
    void loadAssistentialWorks();
    void loadProductGroups();
  }, []);

  useEffect(() => {
    void loadFamilies(statusFilter, workFilter);
  }, [statusFilter, workFilter]);

  async function loadAssistentialWorks() {
    try {
      const response = await ListAssistentialWorks(false);
      setAssistentialWorks(response);
      if (!form.assistentialWorkID && response.length > 0) {
        setForm((prev) => ({...prev, assistentialWorkID: String(response[0].id)}));
      }
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    }
  }

  async function loadProductGroups() {
    try {
      const response = await ListProductGroups(false);
      setProductGroups(response);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    }
  }

  async function loadFamilies(currentStatusFilter, currentWorkFilter) {
    setLoading(true);
    setError('');
    try {
      const includeInactive = currentStatusFilter !== 'active';
      const assistentialWorkID = currentWorkFilter === 'all' ? 0 : Number(currentWorkFilter);
      const response = await ListFamilies(includeInactive, assistentialWorkID);
      const filtered = currentStatusFilter === 'inactive'
        ? response.filter((item) => item.is_active === 0)
        : response;
      setFamilies(filtered);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoading(false);
    }
  }

  async function loadAssignments(familyID) {
    setLoadingAssignments(true);
    try {
      const response = await ListFamilyGroupAssignments(familyID);
      setAssignments(response);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setLoadingAssignments(false);
    }
  }

  function resetForm() {
    setForm({
      id: null,
      assistentialWorkID: assistentialWorks.length > 0 ? String(assistentialWorks[0].id) : '',
      name: '',
      memberCount: '1',
      address: '',
      contact: ''
    });
  }

  function resetAssignmentForm() {
    setAssignmentForm({
      productGroupId: ''
    });
  }

  function openCreateModal() {
    resetForm();
    resetAssignmentForm();
    setAssignments([]);
    setError('');
    setIsModalOpen(true);
  }

  async function startEdit(family) {
    setForm({
      id: family.id,
      assistentialWorkID: String(family.assistential_work_id),
      name: family.name,
      memberCount: String(family.member_count),
      address: family.address ?? '',
      contact: family.contact ?? ''
    });
    resetAssignmentForm();
    setAssignments([]);
    setError('');
    setIsModalOpen(true);
    await loadAssignments(family.id);
  }

  function closeModal() {
    setIsModalOpen(false);
    resetForm();
    resetAssignmentForm();
    setAssignments([]);
  }

  async function handleSubmit(event) {
    event.preventDefault();
    const assistentialWorkID = Number(form.assistentialWorkID);
    const memberCount = Number(form.memberCount);
    const name = form.name.trim();
    const initialGroupID = Number(assignmentForm.productGroupId);

    if (!Number.isInteger(assistentialWorkID) || assistentialWorkID <= 0) {
      setError('Selecione uma obra assistencial.');
      return;
    }
    if (!name) {
      setError('Informe o nome da família.');
      return;
    }
    if (!Number.isInteger(memberCount) || memberCount <= 0) {
      setError('Quantidade de membros deve ser inteiro maior que zero.');
      return;
    }

    setSubmitting(true);
    setError('');
    try {
      if (isEditing) {
        await UpdateFamily(form.id, name, memberCount, form.address, form.contact);
      } else {
        const created = await CreateFamily(assistentialWorkID, name, memberCount, form.address, form.contact);
        if (Number.isInteger(initialGroupID) && initialGroupID > 0) {
          await AssignFamilyGroup(created.id, initialGroupID);
        }
      }
      closeModal();
      await loadFamilies(statusFilter, workFilter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setSubmitting(false);
    }
  }

  async function handleAssignGroup() {
    const familyID = form.id;
    const productGroupID = Number(assignmentForm.productGroupId);

    if (!familyID) {
      setError('Salve a família antes de trocar a cesta.');
      return;
    }
    if (!Number.isInteger(productGroupID) || productGroupID <= 0) {
      setError('Selecione uma cesta para atribuir.');
      return;
    }

    setAssignmentSubmitting(true);
    setError('');
    try {
      await AssignFamilyGroup(familyID, productGroupID);
      resetAssignmentForm();
      await loadAssignments(familyID);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setAssignmentSubmitting(false);
    }
  }

  async function handleToggleActive(family) {
    if (family.is_active === 1) {
      setFamilyToDeactivate(family);
      return;
    }

    setError('');
    try {
      setTogglingActive(true);
      await SetFamilyActive(family.id, true);
      await loadFamilies(statusFilter, workFilter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setTogglingActive(false);
    }
  }

  async function confirmDeactivate() {
    if (!familyToDeactivate) {
      return;
    }

    setError('');
    try {
      setTogglingActive(true);
      await SetFamilyActive(familyToDeactivate.id, false);
      setFamilyToDeactivate(null);
      await loadFamilies(statusFilter, workFilter);
    } catch (err) {
      setError(getFriendlyErrorMessage(err));
    } finally {
      setTogglingActive(false);
    }
  }

  return (
    <main className='families-page'>
      <section className='panel panel-list'>
        <div className='list-header'>
          <div>
            <h1 className='title'>Famílias</h1>
            <p className='subtitle'>Cadastro de famílias por obra assistencial</p>
          </div>
          <button type='button' className='btn btn-primary' onClick={openCreateModal}>
            + Novo
          </button>
        </div>

        <div className='list-tools'>
          <AssistentialWorkFilters filters={STATUS_FILTERS} currentFilter={statusFilter} onChange={setStatusFilter} />
        </div>

        <div className='work-filter'>
          <label htmlFor='work-filter'>Filtrar por obra assistencial</label>
          <select id='work-filter' value={workFilter} onChange={(event) => setWorkFilter(event.target.value)}>
            <option value='all'>Todas as obras</option>
            {assistentialWorks.map((work) => (
              <option key={work.id} value={work.id}>
                {work.name}
              </option>
            ))}
          </select>
        </div>

        {error && <p className='feedback error'>{error}</p>}

        {loading ? (
          <p className='feedback'>Carregando famílias...</p>
        ) : families.length === 0 ? (
          <p className='feedback'>Nenhum registro encontrado para este filtro.</p>
        ) : (
          <ul className='family-list'>
            {families.map((family) => (
              <li key={family.id} className='family-item'>
                <div className='family-main'>
                  <div>
                    <strong>{family.name}</strong>
                    <p>Obra: {worksByID.get(family.assistential_work_id)?.name ?? `#${family.assistential_work_id}`}</p>
                    <p>Membros: {family.member_count}</p>
                    <p>Endereço: {family.address || 'Não informado'}</p>
                    <p>Contato: {family.contact || 'Não informado'}</p>
                  </div>

                  <div className='item-actions'>
                    <span className={`status ${family.is_active === 1 ? 'status-on' : 'status-off'}`}>
                      {family.is_active === 1 ? 'Ativa' : 'Inativa'}
                    </span>
                    <button type='button' className='btn btn-ghost' onClick={() => void startEdit(family)}>
                      Editar
                    </button>
                    <button type='button' className='btn btn-danger' onClick={() => void handleToggleActive(family)}>
                      {family.is_active === 1 ? 'Desativar' : 'Reativar'}
                    </button>
                  </div>
                </div>
              </li>
            ))}
          </ul>
        )}
      </section>

      {isModalOpen && (
        <div className='modal-overlay' onClick={closeModal}>
          <section className='modal-content panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>{isEditing ? 'Editar família' : 'Nova família'}</h2>
              <button type='button' className='btn btn-ghost' onClick={closeModal} disabled={submitting || assignmentSubmitting}>
                Fechar
              </button>
            </div>

            <form className='form' onSubmit={handleSubmit}>
              <label htmlFor='family-work'>Obra assistencial</label>
              <select
                id='family-work'
                value={form.assistentialWorkID}
                onChange={(event) => setForm((prev) => ({...prev, assistentialWorkID: event.target.value}))}
                disabled={submitting || isEditing}
              >
                {assistentialWorks.length === 0 && <option value=''>Nenhuma obra ativa</option>}
                {assistentialWorks.map((work) => (
                  <option key={work.id} value={work.id}>
                    {work.name}
                  </option>
                ))}
              </select>
              {isEditing && <p className='hint'>A obra assistencial não pode ser alterada após o cadastro.</p>}

              <label htmlFor='family-name'>Nome da família</label>
              <input
                id='family-name'
                value={form.name}
                onChange={(event) => setForm((prev) => ({...prev, name: event.target.value}))}
                placeholder='Ex: Família Santos'
                maxLength={120}
                disabled={submitting}
              />

              <label htmlFor='family-members'>Quantidade de membros</label>
              <input
                id='family-members'
                type='number'
                min='1'
                step='1'
                value={form.memberCount}
                onChange={(event) => setForm((prev) => ({...prev, memberCount: event.target.value}))}
                disabled={submitting}
              />

              <label htmlFor='family-address'>Endereço</label>
              <input
                id='family-address'
                value={form.address}
                onChange={(event) => setForm((prev) => ({...prev, address: event.target.value}))}
                placeholder='Opcional'
                maxLength={240}
                disabled={submitting}
              />

              <label htmlFor='family-contact'>Contato</label>
              <input
                id='family-contact'
                value={form.contact}
                onChange={(event) => setForm((prev) => ({...prev, contact: event.target.value}))}
                placeholder='Opcional'
                maxLength={120}
                disabled={submitting}
              />

              <div className='assignment-setup'>
                <h3>{isEditing ? 'Cestas da família' : 'Cesta inicial'}</h3>
                <p className='hint'>
                  {isEditing
                    ? 'Troque a cesta da família e acompanhe o histórico nesta mesma modal.'
                    : 'Opcional. Se informada, será vinculada logo após criar a família.'}
                </p>

                <div className='assignment-form'>
                  <select
                    value={assignmentForm.productGroupId}
                    onChange={(event) => setAssignmentForm((prev) => ({...prev, productGroupId: event.target.value}))}
                    disabled={assignmentSubmitting || submitting}
                  >
                    <option value=''>{isEditing ? 'Selecione a nova cesta' : 'Selecione a cesta inicial'}</option>
                    {productGroups.map((group) => (
                      <option key={group.id} value={group.id}>
                        {group.name}
                      </option>
                    ))}
                  </select>

                  {isEditing && (
                    <button
                      className='btn btn-ghost'
                      type='button'
                      onClick={() => void handleAssignGroup()}
                      disabled={assignmentSubmitting || submitting || productGroups.length === 0}
                    >
                      {assignmentSubmitting ? 'Trocando...' : 'Trocar cesta'}
                    </button>
                  )}
                </div>

                {isEditing && (loadingAssignments ? (
                  <p className='feedback'>Carregando histórico...</p>
                ) : assignments.length === 0 ? (
                  <p className='feedback'>Nenhuma cesta atribuída ainda.</p>
                ) : (
                  <ul className='assignment-list'>
                    {assignments.map((assignment) => (
                      <li key={assignment.id} className='assignment-item'>
                        <div>
                          <strong>{groupsByID.get(assignment.product_group_id)?.name ?? `Cesta #${assignment.product_group_id}`}</strong>
                          <p>Inicio: {assignment.started_at}</p>
                          <p>Fim: {assignment.ended_at || 'Ativa'}</p>
                        </div>
                        <span className={`status ${assignment.ended_at ? 'status-off' : 'status-on'}`}>
                          {assignment.ended_at ? 'Encerrada' : 'Ativa'}
                        </span>
                      </li>
                    ))}
                  </ul>
                ))}
              </div>

              <div className='form-actions'>
                <button className='btn btn-primary' type='submit' disabled={submitting || assistentialWorks.length === 0}>
                  {submitting ? 'Salvando...' : isEditing ? 'Salvar edição' : 'Criar família'}
                </button>
                <button className='btn btn-ghost' type='button' onClick={closeModal} disabled={submitting || assignmentSubmitting}>
                  Cancelar
                </button>
              </div>
            </form>
          </section>
        </div>
      )}

      {familyToDeactivate && (
        <div className='modal-overlay' onClick={() => setFamilyToDeactivate(null)}>
          <section className='modal-content modal-confirm panel' onClick={(event) => event.stopPropagation()}>
            <div className='modal-header'>
              <h2>Desativar família</h2>
            </div>
            <p className='confirm-text'>
              Deseja desativar <strong>{familyToDeactivate.name}</strong>?
            </p>
            <div className='confirm-actions'>
              <button
                type='button'
                className='btn btn-ghost'
                onClick={() => setFamilyToDeactivate(null)}
                disabled={togglingActive}
              >
                Cancelar
              </button>
              <button
                type='button'
                className='btn btn-danger'
                onClick={() => void confirmDeactivate()}
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
