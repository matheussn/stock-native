import { useMutation } from '@tanstack/react-query'

import { callBackend } from '@/commons/utils/wails'

export const useExportBackupMutation = () =>
  useMutation({
    mutationFn: (path: string) => callBackend<void>('ExportBackup', path),
  })

export const useImportBackupMutation = () =>
  useMutation({
    mutationFn: (payload: { filePath: string; confirmReplace: boolean; confirmUnderstand: boolean }) =>
      callBackend<void>('ImportBackup', payload),
  })