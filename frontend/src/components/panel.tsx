import type { ReactNode } from 'react'

interface PanelProps {
  title: string
  subtitle?: string
  actions?: ReactNode
  children: ReactNode
}

export const Panel = ({ title, subtitle, actions, children }: PanelProps) => (
  <section className='panel'>
    <header className='panel-header'>
      <div>
        <h2>{title}</h2>
        {subtitle ? <p>{subtitle}</p> : null}
      </div>
      {actions ? <div>{actions}</div> : null}
    </header>
    <div>{children}</div>
  </section>
)

interface MessageProps {
  text?: string
  kind?: 'error' | 'success'
}

export const Message = ({ text, kind = 'success' }: MessageProps) =>
  text ? <p className={`message ${kind}`}>{text}</p> : null