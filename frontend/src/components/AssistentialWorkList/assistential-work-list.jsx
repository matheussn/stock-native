import {Pencil, RotateCcw, Trash2} from 'lucide-react';
import './assistential-work-list.css';

export default function AssistentialWorkList({works, loading, onEdit, onToggleActive}) {
  if (loading) {
    return <p className='feedback'>Carregando obras assistenciais...</p>;
  }

  if (works.length === 0) {
    return <p className='feedback'>Nenhum registro encontrado para este filtro.</p>;
  }

  return (
    <ul className='work-list'>
      {works.map((item) => (
        <li key={item.id} className='work-item'>
          <div>
            <strong>{item.name}</strong>
            <p>{item.description || 'Sem descrição'}</p>
          </div>

          <div className='item-actions'>
            <span className={`status ${item.is_active === 1 ? 'status-on' : 'status-off'}`}>
              {item.is_active === 1 ? 'Ativo' : 'Inativo'}
            </span>
            <button
              type='button'
              className='icon-action icon-action-edit'
              onClick={() => onEdit(item)}
              aria-label='Editar obra'
              title='Editar obra'
            >
              <Pencil aria-hidden='true' />
            </button>
            <button
              type='button'
              className={`icon-action ${item.is_active === 1 ? 'icon-action-delete' : 'icon-action-restore'}`}
              onClick={() => onToggleActive(item)}
              aria-label={item.is_active === 1 ? 'Desativar obra' : 'Reativar obra'}
              title={item.is_active === 1 ? 'Desativar obra' : 'Reativar obra'}
            >
              {item.is_active === 1 ? <Trash2 aria-hidden='true' /> : <RotateCcw aria-hidden='true' />}
            </button>
          </div>
        </li>
      ))}
    </ul>
  );
}
