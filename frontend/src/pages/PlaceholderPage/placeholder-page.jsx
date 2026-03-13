import './placeholder-page.css';

export default function PlaceholderPage({title, phase, summary}) {
  return (
    <main className='panel placeholder-page'>
      <h1>{title}</h1>
      <p>{summary}</p>
      <span className='phase-tag'>{phase}</span>
    </main>
  );
}
