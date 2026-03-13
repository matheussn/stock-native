import './assistential-work-form.css';

export default function AssistentialWorkForm({
  form,
  isEditing,
  submitting,
  onNameChange,
  onDescriptionChange,
  onSubmit,
  onCancel
}) {
  return (
    <form className='form' onSubmit={onSubmit}>
      <label htmlFor='work-name'>Nome</label>
      <input
        id='work-name'
        value={form.name}
        onChange={(event) => onNameChange(event.target.value)}
        placeholder='Ex: Fraternidade'
        maxLength={120}
        disabled={submitting}
      />

      <label htmlFor='work-description'>Descrição</label>
      <textarea
        id='work-description'
        value={form.description}
        onChange={(event) => onDescriptionChange(event.target.value)}
        placeholder='Opcional'
        rows={3}
        maxLength={600}
        disabled={submitting}
      />

      <div className='form-actions'>
        <button className='btn btn-primary' type='submit' disabled={submitting}>
          {isEditing ? 'Salvar edição' : 'Criar obra'}
        </button>
        {isEditing && (
          <button className='btn btn-ghost' type='button' onClick={onCancel} disabled={submitting}>
            Cancelar
          </button>
        )}
      </div>
    </form>
  );
}
