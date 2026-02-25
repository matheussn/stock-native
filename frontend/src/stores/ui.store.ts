import { create } from 'zustand'

type Section = 'dashboard' | 'estoque' | 'cadastros' | 'movimentacoes' | 'relatorios' | 'backup'
export type CatalogScreen = 'categorias' | 'produtos' | 'origens' | 'destinos'

interface UiState {
  section: Section
  catalogScreen: CatalogScreen
  setSection: (section: Section) => void
  setCatalogScreen: (screen: CatalogScreen) => void
}

export const useUiStore = create<UiState>((set) => ({
  section: 'estoque',
  catalogScreen: 'categorias',
  setSection: (section) => set({ section }),
  setCatalogScreen: (catalogScreen) => set({ catalogScreen }),
}))
