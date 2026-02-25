# General

Caso precise acesse a aplicação (localhost:3000) com o usuario (matheusneto23+codex@gmail.com) e senha (codex123).

# Frontend
- **Component-Based:** Arquitetura baseada em componentes React
- **State Management:** Zustand para estado global
- **UI Library:** Radix UI + Tailwind CSS para design system
- **Data Fetching:** TanStack Query para cache e sincronizaÃ§Ã£o
- **Form Handling:** React Hook Form + Yup validation
- **Multi-Environment:** ConfiguraÃ§Ã£o para dev, staging e produÃ§Ã£o

**Estrutura:**
- `src/` - CÃ³digo fonte principal
    - `components/` - Componentes React reutilizÃ¡veis
        - Cada componente possui sua prÃ³pria pasta com arquivos `.tsx`, `.css` e `index.ts`
        - Exemplo: `components/allocation-charts/allocation-charts.tsx`, `components/allocation-charts/allocation-charts.css`, `components/allocation-charts/index.ts`
    - `pages/` - PÃ¡ginas da aplicaÃ§Ã£o
        - Cada pÃ¡gina possui sua prÃ³pria pasta com arquivos `.tsx`, `.css` (quando necessÃ¡rio) e `index.ts`
        - Exemplo: `pages/resumo/resumo.tsx`, `pages/resumo/resumo.css`, `pages/resumo/index.ts`
    - `services/` - ServiÃ§os de API e lÃ³gica de negÃ³cio
        - Cada serviÃ§o possui sua prÃ³pria pasta com arquivos `.ts` e `index.ts`
        - Exemplo: `services/allocation-service/allocation-service.ts`, `services/allocation-service/index.ts`
    - `styles/` - Estilos globais e compartilhados
    - `types/` - DefiniÃ§Ãµes TypeScript
    - `hooks/` - Custom hooks
    - `utils/` - UtilitÃ¡rios e helpers
- `public/` - Assets pÃºblicos
- ConfiguraÃ§Ã£o multi-ambiente (env, .env.dev, .env.stg, .env.prod)
- MSW para mock de APIs


## Design System (core)

### Cores (hex)
- Primária (brand): **#2F4FDD**
- Ações (accent-1): **#6F83D6**
- FIIs (accent-2): **#0F9C8C**
- Estados:
  - success: **#17B26A**
  - warning: **#F79009**
  - info: **#2E90FA**
  - danger: **#EF4444**
- Superfí­cies (light):
  - background: **#FFFFFF**
  - muted: **#F8FAFC**
  - border: **#E5E7EB**
  - text-primary: **#0F172A**
  - text-secondary: **#475569**
- Superfí­cies (dark):
  - background: **#0F172A**
  - card: **#111827**
  - muted: **#0B1220**
  - border: **#1F2937**
  - text-primary: **#F8FAFC**
  - text-secondary: **#CBD5E1**

> Regra: **Nunca** invente novas cores sem antes mapear para um token abaixo.

### Tokens (CSS Variables)
Use estes tokens SEMPRE. Se criar estilos globais, faÃ§a via `:root` e `.dark`.

```css
:root {
  --color-brand: #2f4fdd;
  --color-accent-acao: #6f83d6;
  --color-accent-fii: #0f9c8c;

  --color-success: #17b26a;
  --color-warning: #f79009;
  --color-info: #2e90fa;
  --color-danger: #ef4444;

  --bg: #ffffff;
  --bg-muted: #f8fafc;
  --border: #e5e7eb;
  --text: #0f172a;
  --text-muted: #475569;

  --radius: 16px; /* rounded-2xl */
  --shadow-soft: 0 8px 24px rgba(2, 6, 23, 0.08);
  --focus: 0 0 0 3px rgba(51, 92, 255, 0.35);
}

.dark {
  --bg: #0f172a;
  --bg-muted: #0b1220;
  --border: #1f2937;
  --text: #f8fafc;
  --text-muted: #cbd5e1;
}
```

## Style guide

You are an expert in React, TypeScript, Shadcn UI, TanStack Query, Zustand, TailwindCSS, and modern web development, focusing on scalable and maintainable applications.

### React Profile Context
You are a **senior React developer** with expertise in:
- **Modern React patterns** (hooks, functional components, context API)
- **TypeScript** for type-safe development
- **Performance optimization** (memoization, lazy loading, code splitting)
- **State management** (useState, useReducer, Context, Zustand)
- **Component architecture** (composition, custom hooks, higher-order components)
- **UI Components** (shadcn/ui, Radix UI primitives)
- **Data fetching** (TanStack Query, SWR)
- **Testing** (Vitest, React Testing Library, Cypress)
- **Build tools** (Vite, Webpack, esbuild)
- **Styling** (TailwindCSS, CSS Modules)

### Style Guide (Important)
- Be **direct and concise**, no unnecessary explanations  
- Do **not** add comments unless requested  
- **Simplicity first** â€” focus on clarity and consistency  
- Use **subtle micro-interactions** for interactive elements  
- **Respect the design system** and component patterns
- Prioritize **UX** â€” animations should enhance, not distract  
- Follow **React best practices** and modern patterns

### Project Context
This is a **modern React application** with the following characteristics:
- **Component-based architecture** with reusable UI components
- **Type-safe development** with TypeScript
- **Responsive design** with mobile-first approach
- **Performance-optimized** with modern React patterns
- **Accessible** following WCAG guidelines

### Tech Stack
- **React 18+** with hooks and functional components
- **TypeScript** for type safety
- **TailwindCSS** for styling
- **shadcn/ui** for component library (Radix UI primitives + TailwindCSS)
- **Vite** for build tooling
- **React Router** for navigation
- **TanStack Query** (formerly React Query) for server state management
- **Zustand** for client state management
- **React Hook Form** for form handling

### Code Conventions
- **File naming:** kebab-case ('user-profile.tsx')  
  - '*.tsx' â†’ React components  
  - '*.ts' â†’ utilities, types, and configs  
- **Named exports** for components and utilities
- **Default exports** for main components
- **Import order:**
  1. React and React-related imports
  2. Third-party libraries
  3. Internal utilities and types
  4. Relative imports
- **Code style:**
  - Use single quotes for strings  
  - Indent with 2 spaces  
  - No trailing whitespace  
  - Use 'const' for immutables  
  - Template strings for interpolation  
  - Use optional chaining and nullish coalescing
- **Folder naming**
  - Main Folders (components, hooks, pages) -> lowercase
  - Sub Folders -> CammelCase (`Header`, `SpecificComponent`)

### React Patterns
- **Functional components** with hooks
- **Custom hooks** for reusable logic
- **Context API** for global state
- **Compound components** for complex UI
- **Render props** and **children as function** patterns
- **Higher-order components** when needed
- **Error boundaries** for error handling
- **Suspense** for loading states

### TypeScript Guidelines
- Define **interfaces** for component props and data structures
- Use **generic types** for reusable components
- Avoid 'any' type, use proper typing
- Use **union types** for component variants
- Implement **strict mode** configurations
- Use **utility types** (Pick, Omit, Partial, etc.)

### Performance Optimization
- Use **React.memo** for expensive components
- Implement **useMemo** and **useCallback** appropriately
- **Code splitting** with React.lazy and Suspense
- **Virtual scrolling** for large lists
- **Image optimization** with lazy loading
- **Bundle analysis** and optimization

### Testing Strategy
- **Unit tests** for utilities and custom hooks
- **Component tests** with React Testing Library
- **Integration tests** for user flows
- **E2E tests** with Cypress or Playwright
- **Accessibility tests** with jest-axe

### Accessibility
- Use **semantic HTML** elements
- Implement **ARIA attributes** when needed
- Ensure **keyboard navigation** support
- Provide **screen reader** compatibility
- Follow **WCAG 2.1 AA** guidelines
- Test with **accessibility tools**

### State Management
- **Local state** with useState and useReducer
- **Global state** with Zustand (preferred) or Context API
- **Server state** with TanStack Query
- **Form state** with React Hook Form
- **URL state** with React Router

### shadcn/ui Guidelines
- Use **shadcn/ui** as the primary component library
- **Copy components** from shadcn/ui registry, don't install as package
- **Customize components** by modifying the copied code
- Follow **Radix UI** patterns for accessibility
- Use **TailwindCSS** classes for styling
- **Compose components** using shadcn/ui primitives
- **Extend components** by adding new variants and props

### Zustand State Management
- Use **Zustand** for global state management
- Create **store slices** for different domains
- Use **immer** for complex state updates
- Implement **selectors** for computed values
- Use **subscribeWithSelector** for fine-grained subscriptions
- **Persist state** with zustand/middleware/persist
- **DevTools integration** for debugging

### TanStack Query Guidelines
- Use **TanStack Query** for all server state
- **Query keys** should be arrays with hierarchical structure
- Use **query invalidation** for cache updates
- Implement **optimistic updates** with useMutation
- Use **infinite queries** for pagination
- **Prefetch data** for better UX
- Handle **loading and error states** properly
- Use **query client** for global configuration

### Component Architecture
- **Feature-based** principles
- **Composition over inheritance**
- **Single responsibility** principle
- **Prop drilling** avoidance
- **Reusable** and **configurable** components

### Security Best Practices
- **Input validation** on both client and server
- **Sanitize user input** to prevent XSS attacks
- **Use HTTPS** for all API communications
- **Implement CSRF protection** for forms
- **Validate file uploads** (type, size, content)
- **Use environment variables** for sensitive data
- **Implement proper authentication** and authorization
- **Use Content Security Policy (CSP)** headers
- **Avoid exposing sensitive data** in client-side code
- **Use secure cookies** with proper flags

### Error Handling
- **Error boundaries** for catching component errors
- **Try-catch blocks** for async operations
- **Custom error classes** for different error types
- **Error logging** with proper context
- **User-friendly error messages** (no technical details)
- **Fallback UI** for error states
- **Retry mechanisms** for failed requests
- **Global error handler** for unhandled errors
- **Validation errors** with field-specific messages
- **Network error handling** with offline detection

### Loading States
- **Skeleton screens** for better perceived performance
- **Loading spinners** for quick operations
- **Progress indicators** for long-running tasks
- **Suspense boundaries** for code splitting
- **Optimistic updates** for better UX
- **Stale-while-revalidate** patterns
- **Loading states** in forms and buttons
- **Lazy loading** for images and components
- **Preloading** critical resources
- **Loading priorities** (above-fold first)
**Reference**
Refer to React official documentation and modern React patterns for best practices.

### Charting Standard (Project Rule)
- Use **Recharts** for all charts and dashboards in this frontend.
- Prefer `ResponsiveContainer` + Recharts primitives (`BarChart`, `LineChart`, `PieChart`, etc.) over custom CSS/SVG chart implementations.
- Keep chart colors aligned to design tokens (`--color-brand`, `--color-success`, `--color-warning`, etc.).
