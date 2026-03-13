import './assistential-work-filters.css';

export default function AssistentialWorkFilters({filters, currentFilter, onChange}) {
  return (
    <div className='filters' role='group' aria-label='Filtro de obras assistenciais'>
      {Object.entries(filters).map(([key, label]) => (
        <button
          key={key}
          type='button'
          className={`chip ${currentFilter === key ? 'chip-active' : ''}`}
          onClick={() => onChange(key)}
        >
          {label}
        </button>
      ))}
    </div>
  );
}
