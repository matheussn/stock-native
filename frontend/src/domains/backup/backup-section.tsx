import { useState } from 'react'

import { Message, Panel } from '@/components/panel'
import { useExportBackupMutation, useImportBackupMutation } from '@/domains/backup/api'

export const BackupSection = () => {
  const [exportPath, setExportPath] = useState('data/backup.db')
  const [importPath, setImportPath] = useState('')
  const [confirmReplace, setConfirmReplace] = useState(false)
  const [confirmUnderstand, setConfirmUnderstand] = useState(false)
  const [message, setMessage] = useState('')

  const exportBackup = useExportBackupMutation()
  const importBackup = useImportBackupMutation()

  return (
    <Panel title='Backup' subtitle='Exportar e importar banco SQLite com confirmação explícita'>
      <Message text={message} />
      <div className='grid two'>
        <div className='card'>
          <h3>Exportar</h3>
          <div className='form'>
            <input value={exportPath} onChange={(event) => setExportPath(event.target.value)} placeholder='Caminho do backup' />
            <button
              type='button'
              onClick={async () => {
                try {
                  await exportBackup.mutateAsync(exportPath)
                  setMessage('Backup exportado com sucesso')
                } catch {
                  setMessage('Erro ao exportar backup')
                }
              }}
            >
              Exportar
            </button>
          </div>
        </div>

        <div className='card'>
          <h3>Importar</h3>
          <div className='form'>
            <input value={importPath} onChange={(event) => setImportPath(event.target.value)} placeholder='Caminho do arquivo .db' />
            <label className='check'>
              <input type='checkbox' checked={confirmReplace} onChange={(event) => setConfirmReplace(event.target.checked)} />
              Confirmo substituição total do banco atual
            </label>
            <label className='check'>
              <input
                type='checkbox'
                checked={confirmUnderstand}
                onChange={(event) => setConfirmUnderstand(event.target.checked)}
              />
              Entendo que a operação é destrutiva
            </label>
            <button
              type='button'
              onClick={async () => {
                try {
                  await importBackup.mutateAsync({
                    filePath: importPath,
                    confirmReplace,
                    confirmUnderstand,
                  })
                  setMessage('Backup importado com sucesso')
                } catch {
                  setMessage('Erro ao importar backup')
                }
              }}
            >
              Importar
            </button>
          </div>
        </div>
      </div>
    </Panel>
  )
}
