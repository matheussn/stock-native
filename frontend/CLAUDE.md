# General

Caso precise acesse a aplicação (localhost:3000) com o usuario (matheusneto23+codex@gmail.com) e senha (codex123).

# Architecture

## Tech Stack
- **Component-Based:** React component architecture
- **State Management:** Zustand for global state
- **UI Library:** Radix UI + Tailwind CSS design system
- **Data Fetching:** TanStack Query for cache and synchronization
- **Form Handling:** React Hook Form + Yup validation
- **Charts:** Recharts for all data visualizations (PieChart, BarChart, etc.) — never use manual SVG or CSS charts
- **Multi-Environment:** Configuration for dev, staging, and production

## Domain-Driven Architecture

The application follows a **domain-driven structure** where related functionality is grouped together by business domain. Each domain is self-contained with its own types, API hooks, components, pages, and stores.

### What is a Domain?

A **domain** represents a business capability or entity within the application. Domains encapsulate related business logic and provide clear boundaries for organizing code.

**Domain Types:**
- **Entity-only domains**: Contain types and API hooks, but no UI components or pages. These are used across multiple domains and represent core business entities (e.g., `user`, `product`).
- **UI-focused domains**: Contain pages and components, but may have minimal or no entity logic. These represent user-facing features.
- **Mixed domains**: Contain both entity logic (types, API hooks) and UI (pages, components). Most domains fall into this category.

**Domain Characteristics:**
- Self-contained with clear boundaries
- Encapsulates related business logic
- Exposes a public API via barrel file (`index.ts`)
- Can depend on other domains, but only through their barrel files
- Owns its data fetching, state management, and UI concerns

### Domain Structure

```
src/
├── domains/                   # Domain-organized code
│   ├── [domain-name]/         # e.g., user, product, order
│   │   ├── index.ts           # Barrel file - public API exports
│   │   ├── types/             # Domain-specific types
│   │   │   └── [Entity].ts
│   │   ├── api/               # TanStack Query hooks
│   │   │   ├── query-keys.ts
│   │   │   ├── use-*-query.ts
│   │   │   └── use-*-mutation.ts
│   │   ├── components/        # Domain-specific components
│   │   │   └── [ComponentName]/
│   │   ├── pages/             # Domain pages
│   │   │   └── [PageName]/
│   │   ├── stores/            # Domain-specific Zustand stores (if needed)
│   │   └── utils/             # Domain-specific utilities (if needed)
│   └── ...
├── components/                # Shared UI primitives (Button, Dialog, etc.)
├── commons/                   # Technical utilities (no business logic, pure functions/helpers)
│   ├── hooks/                 # Shared hooks
│   ├── utils/                 # Shared utilities
│   ├── constants/             # Shared constants
│   └── validators/            # Shared validators
├── providers/                 # React context providers
├── router/                    # React Router configuration
├── locales/                   # i18n translation files
└── assets/                    # Images, fonts, styles
```

### Domain Principles

1. **Self-Contained**: All domain-related code (types, API, components, pages) lives within the domain folder
2. **Public API via Barrel**: Export reusable components/hooks/types through `index.ts` barrel file
3. **Domain Dependencies**: Domains can import from other domains, but only through barrel files
4. **Component Organization**: 
   - `components/` = UI primitives (Button, Dialog, etc.) - domain-agnostic, truly reusable
   - `domains/[domain]/components/` = **All domain-specific components** - even if used by other domains, they stay in their domain and are exported via barrel file
5. **Domain vs Commons Distinction**:
   - `domains/` = All business domains (with or without UI, contains business logic)
   - `commons/` = Technical utilities (no business logic, pure functions/helpers)

### Domain Examples

**Product Domain:**
```
domains/product/
├── index.ts                    # Exports: ProductsPage, ProductDetailsPage, ProductSelector, useGetAllProductsQuery, Product type
├── types/
│   └── Product.ts
├── api/
│   ├── query-keys.ts
│   ├── use-get-all-products-query.ts
│   └── use-post-product-mutation.ts
├── components/
│   ├── ProductSelector/
│   └── ProductStatusTag/
└── pages/
    ├── Products/
    └── ProductDetails/
```

**User Domain:**
```
domains/user/
├── index.ts                    # Exports: UsersPage, UserSelector, usePostUserMutation
├── types/
│   └── User.ts
├── api/
│   ├── use-post-user-mutation.ts
│   └── use-get-all-users-query.ts
├── components/
│   ├── UserSelector/           # Used by order, payment domains
│   └── UserCard/
├── pages/
│   └── Users/
└── constants.ts
```

### Cross-Domain Dependencies

**Shared Entities** (e.g., User, Product):
- Create a dedicated domain: `domains/user/` or `domains/product/`
- Contains base API hooks and types
- Other domains can create domain-specific wrappers if needed

**Shared Components** (e.g., UserSelector, ProductCard):
- Keep components in their domain
- Export via domain barrel file when needed by other domains
- Other domains import from barrel: `import { UserSelector } from '@/domains/user'`

**Import Rules:**
- ✅ `import { X } from '@/domains/[domain]'` - Import from barrel
- ❌ `import { X } from '@/domains/[domain]/components/X'` - Deep imports (not allowed)
- ✅ `import { X } from '@/components'` - Shared UI primitives
- ✅ `import { X } from '@/commons/utils'` - Shared utilities

### Migration Strategy

- New code: Use domain structure immediately
- Existing code: Migrate incrementally, domain by domain
- During migration: Both structures may coexist temporarily

## Legacy Folder Structure (During Migration)

```
src/
├── components/           # Reusable React components
│   └── [ComponentName]/  # PascalCase folders with index.tsx
│       ├── index.tsx
│       └── [sub-components]
├── pages/               # Application pages
│   └── [Feature]/       # Feature-based organization
│       ├── index.tsx    # Main page component
│       ├── hooks/       # Page-specific hooks (if only used in this page)
│       └── components/  # Page-specific components
├── commons/            # Shared utilities and helpers
│   ├── hooks/          # Shared hooks (used across multiple pages/components)
│   ├── utils/          # Utility functions
│   ├── constants/      # Constants
│   └── validators/     # Validation utilities
├── resources/          # API resources and configurations
│   └── apis/           # Server state hooks (TanStack Query)
│       └── [resource]/
│           ├── query-keys.ts
│           ├── use-*-query.ts
│           └── use-*-mutation.ts
├── stores/             # Zustand global state stores
├── types/              # TypeScript type definitions
├── providers/          # React context providers
├── router/             # React Router configuration
│   ├── index.tsx       # Main router configuration
│   └── routes/         # Route definitions by feature
├── locales/            # i18n translation files
└── assets/             # Images, fonts, styles
```

## Key Patterns

### Component Organization

**Domain Components:**
- **All domain-specific components** stay in `domains/[domain]/components/`
- Use PascalCase folders with `index.tsx`
- Even if a component is used by multiple domains, it stays in its domain
- Export reusable components via domain barrel file (`domains/[domain]/index.ts`) when needed by other domains

**Shared Components:**
- UI primitives (Button, Dialog, etc.) in `components/` - domain-agnostic, truly reusable

**Principle**: Domain-specific components always stay in their domain. Export via barrel file when needed by other domains.

### Custom Hooks Organization

**Domain Hooks:**
- **domains/[domain]/api/**: Domain API hooks (TanStack Query) - domain server state
- **domains/[domain]/pages/*/hooks/**: Page-scoped hooks within domain
- **domains/[domain]/components/*/hooks/**: Component-scoped hooks within domain

**Shared Hooks:**
- **commons/hooks/**: Shared hooks (used across multiple domains)
- **Component-level**: Hooks used only within a single component

**Principle**: Keep hooks in their domain, move to commons only when used across multiple domains

### Domain Barrel Files

Each domain should have an `index.ts` barrel file that exports:
- **Pages** (for router configuration)
- Reusable components (used by other domains)
- Public API hooks (if needed by other domains)
- Public types (if needed by other domains)
- Constants (if needed by other domains)

**Example:**
```typescript
// domains/user/index.ts
export { UsersPage } from './pages/Users'
export { UserSelector } from './components/UserSelector'
export { usePostUserMutation } from './api/use-post-user-mutation'
export { MAX_USERS_PER_PAGE } from './constants'
export type { User, UserStatus } from './types/User'
```

**Do NOT export:**
- Internal implementation details
- Components only used within the domain
- Private utilities or helpers

### File Naming
- **PascalCase**: Components (folders with `index.tsx`), Pages, Types
  - Components: `Button/`, `Dialog/`
  - Pages: `Detail.tsx`, `Content.tsx`
  - Types: `User.ts`, `Product.ts`
- **camelCase**: All hooks (API and non-API)
  - API hooks: `useGetAllUsersQuery.ts`
  - Shared hooks: `useAuth.ts`, `useLocalStorage.ts`
  - Page hooks: `useMetricsDashboardData.ts`
- **kebab-case**: Utils, stores, route files, query keys, constants
  - Utils: `format-number.ts`
  - Stores: `auth.store.ts`
  - Route files: `admin-route.tsx`
  - Query keys: `query-keys.ts`
- **Index files**: `index.ts` or `index.tsx` for cleaner imports

### Cross-Domain Dependencies

**Handling Shared Entities:**
- If an entity (e.g., User, Product) is used by multiple domains, create a dedicated domain: `domains/user/` or `domains/product/`
- Contains base API hooks, types, and utilities
- Other domains can create domain-specific wrappers if needed:
  ```typescript
  // domains/order/api/use-get-order-users-query.ts
  import { useGetAllUsersQuery } from '@/domains/user'
  
  export const useGetOrderUsersQuery = (orderId: string) => {
    return useGetAllUsersQuery({ order: orderId })
  }
  ```

**Handling Shared Components:**
- If a component is used by 2+ domains but is domain-specific (e.g., UserSelector uses user domain logic):
  - **Keep it in the domain** where it belongs
  - Export via domain barrel file: `domains/user/index.ts`
  - Other domains import from barrel: `import { UserSelector } from '@/domains/user'`
  - **Never move domain-specific components to a shared folder** - they stay in their domain

**Dependency Rules:**
- ✅ Domain → Domain: Allowed (import from barrel file)
- ✅ Domain → Shared/Commons: Allowed
- ❌ Shared/Commons → Domain: Avoid (creates coupling)
- ✅ Domain → Components: Allowed (UI primitives)

### Dependency Visualization

Visualizing domain dependencies helps identify coupling, circular dependencies, and architectural issues.

**Purpose:**
- Identify circular dependencies between domains
- Measure coupling between domains
- Discover domains that import from too many other domains
- Guide refactoring decisions

**Tools:**
- **dependency-cruiser**: Analyze and visualize dependencies with rules enforcement
- **madge**: Generate dependency graphs from code
- **Custom scripts**: Build custom visualization using AST parsing or import analysis

**What to Visualize:**
- Domain → Domain imports (via barrel files)
- Dependency depth (how many levels of domain dependencies exist)
- Circular dependencies (domains that depend on each other)
- Import frequency (which domains are most depended upon)

**When to Review:**
- During refactoring efforts
- When adding new cross-domain dependencies
- Periodic architectural reviews (e.g., quarterly)
- Before major feature additions

**What to Look For:**
- **High coupling**: Domains that import from many other domains (consider splitting or restructuring)
- **Circular dependencies**: Domains that depend on each other (refactor to break cycles)
- **Hub domains**: Domains that are imported by many others (may indicate shared entity domains)
- **Isolated domains**: Domains with no cross-domain dependencies (may be good or indicate missing integration)

### Router Organization
- Routes are defined in `src/router/routes/` with one file per feature
- Main router configuration is in `src/router/index.tsx`
- Routes import pages from domain barrel files: `import { ProductsPage } from '@/domains/product'`
- Pages are exported via domain barrel file, not imported directly from `pages/` folder

### Multi-Environment
- Configuration files: `.env`, `.env.dev`, `.env.stg`, `.env.prod`
- Storybook for component development
- MSW for API mocking

# Component Patterns

## Folder Structure

Components use PascalCase folders with `index.tsx`:

```
components/
└── Button/
    ├── index.tsx          # Main component
    ├── CloseButton.tsx    # Sub-components (if needed)
    └── KebabButton.tsx
```

## Index Files for Cleaner Imports

Use `index.ts` files to enable cleaner imports:

**components/index.ts:**
```typescript
export { Button, buttonVariants } from './Button'
export { Dialog, DialogContent } from './Dialog'
```

**Usage:**
```typescript
import { Button, Dialog } from '@/components'  // Clean import
```

Apply this pattern everywhere (components, pages, hooks, etc.) for better import readability.

## Component Composition

- **Composition over inheritance** - Build complex UIs by composing smaller components
- **Single responsibility** - Each component should have one clear purpose
- **Compound components** - Use for complex UI patterns (e.g., Dialog with DialogHeader, DialogContent)
- **Prop interfaces** - Always define TypeScript interfaces for component props

## When to Split Components

Split a component when:
- It exceeds ~200-300 lines
- It has multiple distinct responsibilities
- Parts are reused independently
- It improves readability and testability

## Export Patterns

- **Named exports** preferred for most components
- **Default exports** required only for lazy-loaded components (React.lazy())

## Component Structure

```typescript
import { ... } from '...'

interface ComponentProps {
  // Props definition
}

export const Component = ({ prop1, prop2 }: ComponentProps) => {
  // Hooks
  // State
  // Effects
  // Handlers
  // Render
  return <div>...</div>
}

Component.displayName = 'Component'
```

# Error Handling

## Error Boundaries

Use Error Boundaries to catch component errors:

```typescript
import { ErrorBoundary } from '@/components/ErrorBoundary'

<ErrorBoundary>
  <YourComponent />
</ErrorBoundary>
```

## Error Logging

- **Acceptable**: `console.error` for error logging (user reporting, monitoring)
- **Remove**: Development debug logs (e.g., `console.log('API response:', data)`)
- **Context**: Include relevant context in error logs for debugging

## API Error Handling

TanStack Query handles API errors automatically. Access error states in components:

```typescript
const { data, error, isLoading } = useGetResourceQuery(id)

if (error) {
  // Handle error state
  return <ErrorState error={error} />
}
```

## Error State UI

Provide user-friendly error messages:
- Show clear, actionable error messages
- Avoid technical details in user-facing errors
- Provide fallback UI for error states
- Include retry mechanisms when appropriate

## Error Types

- **Network errors**: Handle offline scenarios, connection issues
- **Validation errors**: Field-specific messages from forms
- **API errors**: Server error responses with user-friendly messages
- **Component errors**: Caught by Error Boundaries with fallback UI

## Best Practices

- **User-friendly messages** - No technical details in user-facing errors
- **Error context** - Log technical details for debugging, show simple messages to users
- **Retry mechanisms** - Provide retry options for transient errors
- **Fallback UI** - Always provide fallback UI for error states
- **Error boundaries** - Use at appropriate levels to catch component errors


# Forms and Validation

## Schema Definition

Define form schemas using Yup:

```typescript
import * as yup from 'yup'

const formSchema = yup.object({
  name: yup
    .string()
    .required('Name is required')
    .max(100, 'Name must be less than 100 characters'),
  description: yup
    .string()
    .optional()
    .default('')
    .max(200, 'Description must be less than 200 characters'),
  type: yup
    .string()
    .oneOf(['OPTION1', 'OPTION2'])
    .required('Type is required'),
})

type FormType = yup.InferType<typeof formSchema>
```

## Form Setup

Use React Hook Form with Yup resolver:

```typescript
import { useForm } from 'react-hook-form'
import { yupResolver } from '@hookform/resolvers/yup'

const form = useForm<FormType>({
  resolver: yupResolver(formSchema),
  defaultValues: {
    // Optional default values
  },
})
```

## Form Component Structure

```typescript
import { Form, FormField, FormItem, FormLabel, FormMessage } from '@/components'

export const MyForm = () => {
  const form = useForm<FormType>({
    resolver: yupResolver(formSchema),
  })

  const onSubmit = async (data: FormType) => {
    // Handle submission
  }

  return (
    <Form {...form}>
      <form onSubmit={form.handleSubmit(onSubmit)}>
        <FormField
          control={form.control}
          name="name"
          render={({ field }) => (
            <FormItem>
              <FormLabel>Name</FormLabel>
              <Input {...field} />
              <FormMessage />
            </FormItem>
          )}
        />
      </form>
    </Form>
  )
}
```

## Best Practices

- **Type inference** - Use `yup.InferType<typeof schema>` for form types
- **Validation messages** - Provide clear, user-friendly error messages
- **Field-specific errors** - Use `FormMessage` component for field-level errors
- **Default values** - Set appropriate defaults in schema or form config
- **Optional fields** - Use `.optional()` for non-required fields
- **Clear empty attributes** - Use `clearEmptyObjectAttributes` before API submission


# Custom Hooks Organization

## Placement Rules

Organize hooks based on scope and reusability:

### Server State Hooks
**Location**: `resources/apis/[resource]/`
- TanStack Query hooks for server state
- Query hooks: `use-get-*-query.ts`
- Mutation hooks: `use-*-mutation.ts`
- Control server state and API interactions

### Page-Scoped Hooks
**Location**: `pages/[Feature]/hooks/`
- Hooks used only within a specific page
- Keep close to where they're used
- Move to `commons/hooks/` if reused across multiple pages

### Shared Hooks
**Location**: `commons/hooks/`
- Hooks used across multiple pages/components
- General-purpose utilities (e.g., `useLocalStorage`, `useAuth`)
- Move up from page-level when reused

### Component-Level Hooks
**Location**: Same file as component or `[Component]/hooks/`
- Hooks used only within a single component
- Keep in component file if simple, separate if complex

## Naming Conventions

**File naming:** All hooks use camelCase (matching React's built-in hooks convention)

- **Server state hooks**: `useGetTrainingByIdQuery.ts`, `usePostTrainingMutation.ts`
- **Shared hooks**: `useAuth.ts`, `useLocalStorage.ts`, `useDebouncedValue.ts`
- **Page-specific hooks**: `useMetricsDashboardData.ts`, `useBenchmarkPreferences.ts`
- **Component-level hooks**: `useButtonState.ts`, `useSelectHandlers.ts`

**Function naming:** Exported hook functions use camelCase (e.g., `useGetTrainingByIdQuery`, `useAuth`)

## Reusability Principle

**Keep close to usage, move up layers only when used in multiple scopes:**

1. Start at component/page level
2. Move to `commons/hooks/` when reused across pages
3. Keep server state hooks in `resources/apis/`

## Examples

```typescript
// Server state - resources/apis/training/
export const useGetTrainingByIdQuery = (id: string) => { ... }

// Shared - commons/hooks/
export const useAuth = () => { ... }
export const useLocalStorage = <T>(key: string) => { ... }

// Page-specific - pages/NeoData/Training/hooks/
export const useTrainingTime = (training: Training) => { ... }

// Component-level - components/Button/hooks/ or in component file
const useButtonState = () => { ... }
```


# MSW Mocking Patterns

## Directory Structure

```
src/mocks/
├── browser.ts          # MSW worker setup and configuration
├── config.ts           # Mock configuration (endpoint building)
├── handlers/           # HTTP request handlers organized by resource
│   ├── index.ts        # Central export combining all handlers
│   ├── datasets.ts     # Dataset handlers
│   ├── trainings.ts    # Training handlers
│   └── ...
├── fixtures/           # Mock data fixtures organized by resource
│   ├── datasets.ts     # Dataset fixture store
│   ├── trainings.ts    # Training fixture store
│   ├── checkpoints.ts  # Checkpoint fixture store (shared resource)
│   └── ...
└── utils/              # Mock utilities
    ├── delay.ts        # Delay simulation
    └── space.ts        # Space helpers
```

## Handler Organization

### File Structure

One handler file per resource (e.g., `datasets.ts`, `trainings.ts`):

```typescript
import { HttpResponse, http } from 'msw'
import { buildEndpoint } from '@/mocks/config'
import { datasetStore } from '@/mocks/fixtures/datasets'
import { simulateDelay } from '@/mocks/utils/delay'
import { getSpaceIdFromRequest } from '@/mocks/utils/space'

const listDatasets = async ({ request }: { request: Request }) => {
  await simulateDelay()
  // ... handler logic
  return HttpResponse.json(data)
}

export const datasetHandlers: HttpHandler[] = [
  http.get(buildEndpoint('/datasets'), listDatasets),
  http.post(buildEndpoint('/datasets'), createDataset),
  // ...
]
```

### Handler Export Pattern

- Export handler array: `export const [resource]Handlers: HttpHandler[] = [...]`
- Use `http.get()`, `http.post()`, `http.put()`, `http.delete()` from MSW
- Central export in `handlers/index.ts` combining all handlers

### Central Handler Export

```typescript
// handlers/index.ts
import { apiKeyHandlers } from './api-keys'
import { datasetHandlers } from './datasets'
import { trainingHandlers } from './trainings'
// ... other handlers

export const handlers: HttpHandler[] = [
  ...apiKeyHandlers,
  ...datasetHandlers,
  ...trainingHandlers,
  // ... all handlers
]
```

## Handler Patterns

### Endpoint Building

Use `buildEndpoint()` from config for consistent URL building:

```typescript
import { buildEndpoint } from '@/mocks/config'

http.get(buildEndpoint('/datasets'), listDatasets)
http.post(buildEndpoint('/datasets/:datasetId/process'), processDataset)
```

### Delay Simulation

Use `simulateDelay()` for realistic async behavior:

```typescript
import { simulateDelay } from '@/mocks/utils/delay'

const listDatasets = async ({ request }: { request: Request }) => {
  await simulateDelay() // Simulates network delay
  // ... handler logic
}
```

### Space-Scoped Requests

Use `getSpaceIdFromRequest()` for space-scoped requests:

```typescript
import { getSpaceIdFromRequest } from '@/mocks/utils/space'

const listDatasets = async ({ request }: { request: Request }) => {
  const spaceId = getSpaceIdFromRequest(request)
  const datasets = datasetStore.all(spaceId)
  // ...
}
```

### Response Types

Return proper response types matching API contracts:

```typescript
import type { Dataset } from '@/types/Dataset'
import type { ListResponse } from '@/types/List'

const listDatasets = async ({ request }: { request: Request }) => {
  const data = datasetStore.all(spaceId)
  return HttpResponse.json<ListResponse<Dataset>>({
    data,
    count: data.length,
    has_more: false,
  })
}
```

### Query Parameters and Request Body

Handle query parameters and request body:

```typescript
const listDatasets = async ({ request }: { request: Request }) => {
  const url = new URL(request.url)
  const name = url.searchParams.get('name') || undefined
  const status = url.searchParams.getAll('status') as DatasetStatus[] | undefined
  
  const filtered = datasetStore.filter({ name, status, spaceId })
  // ...
}

const createDataset = async ({ request }: { request: Request }) => {
  const payload = (await request.json().catch(() => ({}))) as {
    name?: string
    description?: string
    // ...
  }
  // ...
}
```

## Fixture Patterns

### Fixture Store Structure

One fixture file per resource (e.g., `datasets.ts`, `trainings.ts`):

```typescript
// fixtures/datasets.ts
import { faker } from '@faker-js/faker'
import type { Dataset } from '@/types/Dataset'

const datasets: Dataset[] = [
  // ... mock data
]

export const datasetStore = {
  all: (spaceId?: string | null) => {
    // Return all datasets for space
  },
  findById: (id: string) => {
    // Find dataset by ID
  },
  filter: (filters: FilterOptions) => {
    // Filter datasets
  },
  add: (dataset: Dataset, spaceId?: string | null) => {
    // Add new dataset
  },
  update: (id: string, updates: Partial<Dataset>) => {
    // Update dataset
  },
  remove: (id: string) => {
    // Remove dataset
  },
}
```

### Faker.js Usage

Use Faker.js for generating realistic mock data:

```typescript
import { faker } from '@faker-js/faker'

const dataset: Dataset = {
  id: faker.string.uuid(),
  name: faker.system.fileName(),
  created_at: faker.date.recent({ days: 10 }).toISOString(),
  // ...
}
```

### CRUD Operations

Support filtering, CRUD operations on mock data:

```typescript
export const datasetStore = {
  all: (spaceId?: string | null) => {
    // Return all or filter by space
  },
  filter: (filters: {
    name?: string
    status?: DatasetStatus[]
    spaceId?: string | null
    // ...
  }) => {
    // Apply filters
  },
  add: (dataset: Dataset, spaceId?: string | null) => {
    datasets.unshift(dataset)
    return dataset
  },
  update: (id: string, updates: Partial<Dataset>) => {
    const index = datasets.findIndex(d => d.id === id)
    if (index === -1) return null
    datasets[index] = { ...datasets[index], ...updates }
    return datasets[index]
  },
  remove: (id: string) => {
    const index = datasets.findIndex(d => d.id === id)
    if (index === -1) return null
    return datasets.splice(index, 1)[0]
  },
}
```

## Data Consistency and Single Source of Truth

### Critical Principle: Single Source of Truth

**Shared resources must use a single fixture store**. When multiple endpoints interact with the same resource, they should all reference the same fixture store.

### Example: Checkpoints as Shared Resource

Checkpoints are used by:
- Training endpoints (training checkpoints)
- Benchmark endpoints (benchmark checkpoints)
- Leaderboard endpoints (leaderboard checkpoints)

**Solution**: Single `checkpointStore` fixture used by all handlers:

```typescript
// fixtures/checkpoints.ts
export const checkpointStore = {
  all: (spaceId?: string | null) => { /* ... */ },
  findByTrainingId: (trainingId: string) => { /* ... */ },
  findByBenchmarkId: (benchmarkId: string) => { /* ... */ },
  // ...
}

// handlers/trainings.ts
import { checkpointStore } from '@/mocks/fixtures/checkpoints'
// Use checkpointStore for training checkpoints

// handlers/benchmarks.ts
import { checkpointStore } from '@/mocks/fixtures/checkpoints'
// Use checkpointStore for benchmark checkpoints

// handlers/leaderboard.ts
import { checkpointStore } from '@/mocks/fixtures/checkpoints'
// Use checkpointStore for leaderboard checkpoints
```

### Schema Extension (Not Duplication)

**Do not duplicate schemas - extend them instead**.

Mock schemas can extend frontend types with additional fields. Think about how data is stored in the database - the backend/DB schema may be richer than what the frontend needs.

```typescript
// Frontend type (types/Checkpoint.ts)
export interface Checkpoint {
  id: string
  training_id: string
  epoch: number
  // ... frontend fields
}

// Mock fixture (fixtures/checkpoints.ts)
// Can have additional fields not in frontend type
interface CheckpointMock extends Checkpoint {
  // Additional DB/backend fields
  internal_metadata?: Record<string, unknown>
  processing_flags?: string[]
  // ... other backend-only fields
}

const checkpoints: CheckpointMock[] = [
  // Mock data with extended schema
]
```

**Key Points**:
- Mock data can have **more fields** than frontend types
- Frontend types represent what the API returns to the frontend
- Mock fixtures represent what might be stored in the database
- Extend, don't duplicate - maintain relationship between types

## Side Effects and Realistic Behavior

### Implement Realistic Side Effects

When an action triggers a process (e.g., starting a training run), simulate the **full lifecycle**, not just the initial state.

### Example: Training Lifecycle

Starting a training should:
1. Create training in "CREATING" or "RUNNING" state
2. Update to "RUNNING" state
3. Progress through epochs/steps
4. Update metrics over time
5. Create checkpoints periodically
6. Complete after 1-2 minutes

```typescript
// fixtures/trainings.ts
const trainingIntervals = new Map<string, NodeJS.Timeout>()

const startTrainingSimulation = (trainingId: string) => {
  const training = trainings.find(t => t.id === trainingId)
  if (!training) return

  let currentStep = 0
  let currentEpoch = 1
  let currentLoss = 0.9 + Math.random() * 0.2

  // Update every few seconds
  const updateInterval = setInterval(() => {
    const trainingIndex = trainings.findIndex(t => t.id === trainingId)
    if (trainingIndex === -1) {
      clearInterval(updateInterval)
      return
    }

    const currentTraining = trainings[trainingIndex]

    // Transition from CREATING to RUNNING
    if (currentTraining.status === TrainingStatus.CREATING) {
      trainings[trainingIndex] = {
        ...currentTraining,
        status: TrainingStatus.RUNNING,
        metrics: [],
        checkpoints: [],
        updated_at: new Date().toISOString(),
      }
      return
    }

    // Update progress
    currentStep += stepsPerUpdate
    if (currentStep > stepsPerEpoch * currentEpoch) {
      currentEpoch++
    }

    // Complete after reaching max epochs (1-2 minutes)
    if (currentEpoch > epochs) {
      trainings[trainingIndex] = {
        ...currentTraining,
        status: TrainingStatus.COMPLETED,
        updated_at: new Date().toISOString(),
      }
      clearInterval(updateInterval)
      trainingIntervals.delete(trainingId)
      return
    }

    // Update metrics, create checkpoints, etc.
    // ...
  }, 2000) // Update every 2 seconds

  trainingIntervals.set(trainingId, updateInterval)
}

export const trainingStore = {
  add: (training: Training, _spaceId?: string | null) => {
    trainings.unshift(training)
    if (!defaultTrainingIds.has(training.id)) {
      startTrainingSimulation(training.id) // Start simulation
    }
    return training
  },
  // ...
}
```

### Realistic Completion Times

Set realistic completion times:
- **Training runs**: 1-2 minutes to complete
- **Dataset processing**: 30-60 seconds
- **File analysis**: 10-30 seconds

Use timers/intervals to update state over time:

```typescript
// Update status every few seconds
const updateInterval = setInterval(() => {
  // Update state
}, 2000) // 2 seconds

// Complete after 1-2 minutes
setTimeout(() => {
  // Mark as completed
  clearInterval(updateInterval)
}, 60000 + Math.random() * 60000) // 1-2 minutes
```

### Emulate Entire Process

Don't just set initial state - emulate the entire process:
- Status transitions (CREATING → RUNNING → COMPLETED)
- Progress updates (epochs, steps, metrics)
- Side effects (checkpoint creation, log updates)
- Error scenarios (occasional failures)

## MSW Setup

### Worker Configuration

```typescript
// browser.ts
import { setupWorker } from 'msw/browser'
import { handlers } from './handlers'

export const worker = setupWorker(...handlers)

const defaultOptions: StartOptions = {
  onUnhandledRequest: 'bypass',
  quiet: true,
}
```

### Mocking Functions

Export `enableMocking()`, `disableMocking()`, `bootstrapMocking()` functions:

```typescript
export const enableMocking = (options?: StartOptions) => {
  if (typeof window === 'undefined') {
    return Promise.resolve(undefined)
  }

  window.__PORTAL_MSW_ENABLED = true
  writeStoredPreference(true)

  return worker.start({
    ...defaultOptions,
    ...options,
  })
}

export const disableMocking = async () => {
  if (typeof window === 'undefined') {
    return
  }

  window.__PORTAL_MSW_ENABLED = false
  writeStoredPreference(false)

  await worker.stop()
}

export const bootstrapMocking = (options?: StartOptions) => {
  if (!readStoredPreference()) {
    return Promise.resolve(undefined)
  }

  return enableMocking(options)
}
```

### LocalStorage Persistence

Use localStorage for persistence of mock state preference:

```typescript
const STORAGE_KEY = 'portal-web:msw-enabled'

const readStoredPreference = () => {
  if (typeof window === 'undefined') {
    return false
  }
  try {
    return window.localStorage.getItem(STORAGE_KEY) === 'true'
  } catch {
    return false
  }
}

const writeStoredPreference = (value: boolean) => {
  if (typeof window === 'undefined') {
    return
  }
  try {
    window.localStorage.setItem(STORAGE_KEY, value.toString())
  } catch {
    console.error('Error writing localStorage key: ', STORAGE_KEY)
  }
}
```

### Window Flag

Handle `window.__PORTAL_MSW_ENABLED` flag:

```typescript
declare global {
  interface Window {
    __PORTAL_MSW_ENABLED?: boolean
  }
}

export const isMockEnabled = () => {
  if (typeof window === 'undefined') {
    return false
  }
  return window.__PORTAL_MSW_ENABLED ?? readStoredPreference()
}
```

## Best Practices

### Keep Handlers in Sync

- Keep handlers in sync with API contracts
- Update handlers when API changes
- Test handlers with actual API responses

### Use Fixtures for Reusability

- Use fixtures for reusable mock data
- Don't hardcode data in handlers
- Share fixtures across related handlers

### Simulate Realistic Behavior

- Simulate realistic delays and errors
- Implement side effects and lifecycle behaviors
- Use timers/intervals for processes that take time

### Ensure Data Consistency

- **Single source of truth** for shared resources
- **Extend schemas** instead of duplicating
- Maintain relationships between related data

### Implement All Expected Side Effects

- When creating a training, start the simulation
- When processing a dataset, update status over time
- When creating checkpoints, update related resources
- Ensure all expected side effects are implemented

## Examples

- **Handler**: `src/mocks/handlers/datasets.ts` - Dataset CRUD handlers
- **Fixture**: `src/mocks/fixtures/datasets.ts` - Dataset fixture store
- **Shared resource**: `src/mocks/fixtures/checkpoints.ts` - Checkpoint fixture (used by training, benchmark, leaderboard handlers)
- **Side effects**: `src/mocks/fixtures/trainings.ts` - Training simulation with lifecycle
- **Setup**: `src/mocks/browser.ts` - MSW worker configuration


# Storybook Patterns

## File Organization

### Co-location with Components

Stories are co-located with components in the same directory:

```
components/
└── Button/
    ├── index.tsx              # Component
    ├── Button.stories.tsx     # Stories
    └── CloseButton.tsx        # Sub-components
```

### File Naming

- **File name**: `ComponentName.stories.tsx`
- **Same directory** as component (e.g., `components/Button/Button.stories.tsx`)
- Use `.stories.tsx` extension

## Story File Structure

### Basic Structure

```typescript
import type { Meta, StoryObj } from '@storybook/react-vite'
import { Button } from './index'

const meta = {
  title: 'Components/Button',
  component: Button,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    // ... argTypes configuration
  },
} satisfies Meta<typeof Button>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    children: 'Button',
  },
}
```

### Structure Guidelines

1. **Import component and dependencies** - Import the component and any needed dependencies
2. **Define `meta` object** - Component metadata and configuration
3. **Export `meta` as default** - Required by Storybook
4. **Define `Story` type** - `type Story = StoryObj<typeof meta>`
5. **Export individual stories** - Named exports for each story

## Meta Configuration

### Required Properties

```typescript
const meta = {
  title: 'Components/Button',        // Component path in Storybook
  component: Button,                 // Component reference
  parameters: {                      // Story parameters
    layout: 'centered',              // Layout option
  },
  tags: ['autodocs'],                // Enable automatic documentation
  argTypes: {                        // Control configuration
    variant: {
      control: 'select',
      options: ['default', 'secondary', 'ghost'],
      description: 'Button variant style',
    },
    // ... more argTypes
  },
} satisfies Meta<typeof Button>
```

### Title Pattern

Use hierarchical paths for organization:

- `'Components/Button'` - Component in Components category
- `'Pages/Auth/LoginForm'` - Page component
- `'Features/DatasetAnalysis'` - Feature component

### Parameters

Common parameters:

```typescript
parameters: {
  layout: 'centered',    // 'centered', 'padded', 'fullscreen'
  docs: {
    description: {
      component: 'Component description',
    },
  },
}
```

### Tags

Use `['autodocs']` for automatic documentation:

```typescript
tags: ['autodocs']
```

### ArgTypes

Configure controls for each prop:

```typescript
argTypes: {
  variant: {
    control: 'select',
    options: ['default', 'secondary', 'ghost', 'outline', 'dashed', 'link', 'danger', 'success'],
    description: 'Button variant style',
  },
  size: {
    control: 'select',
    options: ['sm', 'default', 'lg'],
    description: 'Button size',
  },
  disabled: {
    control: 'boolean',
    description: 'Disable the button',
  },
  children: {
    control: 'text',
    description: 'Button content',
  },
}
```

### Type Safety

Use `satisfies Meta<typeof Component>` for type safety:

```typescript
const meta = {
  // ... configuration
} satisfies Meta<typeof Button>
```

This ensures:
- Type checking for meta properties
- Autocomplete for component props
- Compile-time validation

## Story Patterns

### Default Story

Basic usage with minimal args:

```typescript
export const Default: Story = {
  args: {
    children: 'Button',
  },
}
```

### Variants Story

Show all variants:

```typescript
export const Variants: Story = {
  render: () => (
    <div className="flex flex-wrap gap-2">
      <Button variant="default">Default</Button>
      <Button variant="secondary">Secondary</Button>
      <Button variant="ghost">Ghost</Button>
      <Button variant="outline">Outline</Button>
      <Button variant="dashed">Dashed</Button>
      <Button variant="link">Link</Button>
      <Button variant="danger">Danger</Button>
      <Button variant="success">Success</Button>
    </div>
  ),
}
```

### Sizes Story

Show all sizes:

```typescript
export const Sizes: Story = {
  render: () => (
    <div className="flex items-center gap-2">
      <Button size="sm">Small</Button>
      <Button size="default">Default</Button>
      <Button size="lg">Large</Button>
    </div>
  ),
}
```

### Examples Story

Real-world usage examples:

```typescript
export const Examples: Story = {
  render: () => (
    <div className="space-y-4">
      <div className="flex flex-wrap gap-2">
        <Button>
          <Icon icon={PlusSignCircleIcon} />
          Add Item
        </Button>
        <Button variant="secondary">
          <Icon icon={Delete03Icon} />
          Delete
        </Button>
      </div>
      {/* More examples */}
    </div>
  ),
}
```

### Playground Story

Interactive story with all controls:

```typescript
export const Playground: Story = {
  args: {
    children: 'Button',
    variant: 'default',
    size: 'default',
    disabled: false,
    isIcon: false,
  },
}
```

### Using Args vs Render

- **Use `args`** for simple stories with prop variations
- **Use `render`** for complex compositions or custom layouts

```typescript
// Simple with args
export const Default: Story = {
  args: {
    children: 'Button',
    variant: 'default',
  },
}

// Complex with render
export const WithIcons: Story = {
  render: () => (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap gap-2">
        <Button>
          <Icon icon={PlusSignCircleIcon} />
          Add Item
        </Button>
        {/* More complex composition */}
      </div>
    </div>
  ),
}
```

## Story Naming

### Naming Conventions

- **PascalCase** for story names (e.g., `Default`, `WithIcons`, `IconOnly`)
- **Descriptive names** that explain the story's purpose
- **Common patterns**:
  - `Default` - Basic usage
  - `Variants` - All variants
  - `Sizes` - All sizes
  - `WithIcons` - With icons
  - `IconOnly` - Icon-only variant
  - `Disabled` - Disabled state
  - `Examples` - Real-world examples
  - `Playground` - Interactive playground

### Examples

```typescript
export const Default: Story = { ... }
export const Variants: Story = { ... }
export const Sizes: Story = { ... }
export const WithIcons: Story = { ... }
export const IconOnly: Story = { ... }
export const Disabled: Story = { ... }
export const Examples: Story = { ... }
export const Playground: Story = { ... }
```

## Best Practices

### Cover All Variants and States

- Cover all component variants (default, secondary, ghost, etc.)
- Include all sizes (sm, default, lg)
- Show disabled states
- Include loading states (if applicable)
- Show error states (if applicable)

### Show Real-World Usage

- Include examples of how the component is used in context
- Show combinations with other components
- Demonstrate common use cases

### Use Proper TypeScript Types

- Use `satisfies Meta<typeof Component>` for type safety
- Define `Story` type for consistency
- Use proper types for args and render functions

### Keep Stories Focused

- Each story should demonstrate one aspect
- Keep stories readable and maintainable
- Use descriptive names

### Component Structure

```typescript
import type { Meta, StoryObj } from '@storybook/react-vite'
// ... other imports
import { Component } from './index'

const meta = {
  title: 'Components/Component',
  component: Component,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    // ... argTypes
  },
} satisfies Meta<typeof Component>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    // ... default args
  },
}

export const Variants: Story = {
  render: () => (
    // ... variant examples
  ),
}

// ... more stories
```

## Complete Example

```typescript
import { Delete03Icon, PlusSignCircleIcon } from '@hugeicons-pro/core-stroke-rounded'
import { Tick02Icon } from '@hugeicons-pro/core-stroke-sharp'
import type { Meta, StoryObj } from '@storybook/react-vite'
import { Icon } from '@/components/Icon'
import { Button } from './index'

const meta = {
  title: 'Components/Button',
  component: Button,
  parameters: {
    layout: 'centered',
  },
  tags: ['autodocs'],
  argTypes: {
    variant: {
      control: 'select',
      options: ['default', 'secondary', 'ghost', 'outline', 'dashed', 'link', 'danger', 'success'],
      description: 'Button variant style',
    },
    size: {
      control: 'select',
      options: ['sm', 'default', 'lg'],
      description: 'Button size',
    },
    disabled: {
      control: 'boolean',
      description: 'Disable the button',
    },
    isIcon: {
      control: 'boolean',
      description: 'Render as icon-only button (square aspect ratio)',
    },
    asChild: {
      control: 'boolean',
      description: 'Render as child component using Radix Slot',
    },
    children: {
      control: 'text',
      description: 'Button content',
    },
  },
} satisfies Meta<typeof Button>

export default meta
type Story = StoryObj<typeof meta>

export const Default: Story = {
  args: {
    children: 'Button',
  },
}

export const Variants: Story = {
  render: () => (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap gap-2">
        <Button variant="default">Default</Button>
        <Button variant="secondary">Secondary</Button>
        <Button variant="ghost">Ghost</Button>
        <Button variant="outline">Outline</Button>
        <Button variant="dashed">Dashed</Button>
        <Button variant="link">Link</Button>
        <Button variant="danger">Danger</Button>
        <Button variant="success">Success</Button>
      </div>
    </div>
  ),
}

export const Sizes: Story = {
  render: () => (
    <div className="flex flex-col gap-4">
      <div className="flex items-center gap-2">
        <Button size="sm">Small</Button>
        <Button size="default">Default</Button>
        <Button size="lg">Large</Button>
      </div>
    </div>
  ),
}

export const WithIcons: Story = {
  render: () => (
    <div className="flex flex-col gap-4">
      <div className="flex flex-wrap gap-2">
        <Button>
          <Icon icon={PlusSignCircleIcon} />
          Add Item
        </Button>
        <Button variant="secondary">
          <Icon icon={Delete03Icon} />
          Delete
        </Button>
        <Button variant="success">
          <Icon icon={Tick02Icon} />
          Confirm
        </Button>
      </div>
    </div>
  ),
}

export const Playground: Story = {
  args: {
    children: 'Button',
    variant: 'default',
    size: 'default',
    disabled: false,
    isIcon: false,
  },
}
```

## Examples

- **Basic component**: `src/components/Button/Button.stories.tsx`
- **Simple component**: `src/components/Badge/Badge.stories.tsx`
- **Complex component**: `src/components/EmptyContent/EmptyContent.stories.tsx`

# style guide

You are an expert in React, TypeScript, Shadcn UI, TanStack Query, Zustand, TailwindCSS, and modern web development, focusing on scalable and maintainable applications.

## React Profile Context
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

## Style Guide (Important)
- Be **direct and concise**, no unnecessary explanations  
- **Simplicity first** — focus on clarity and consistency  
- Use **subtle micro-interactions** for interactive elements  
- **Respect the design system** and component patterns
- Prioritize **UX** — animations should enhance, not distract  
- Follow **React best practices** and modern patterns

## Code Commentary
- **Avoid obvious comments** that restate what the code already shows
- **Add comments for:**
  - Complex business logic or non-obvious decisions
  - Workarounds and temporary solutions (with context)
  - Complex state flows with non-linear dependencies
  - Performance optimizations that aren't self-evident
  - Edge cases and error handling rationale
  - Integration points with specific requirements
- **Prefer self-documenting code** through clear naming, small functions, and TypeScript types

## Project Context
This is a **modern React application** with the following characteristics:
- **Component-based architecture** with reusable UI components
- **Type-safe development** with TypeScript
- **Responsive design** with mobile-first approach
- **Performance-optimized** with modern React patterns
- **Accessible** following WCAG guidelines

## Tech Stack
- **React 18+** with hooks and functional components
- **TypeScript** for type safety
- **TailwindCSS** for styling
- **shadcn/ui** for component library (Radix UI primitives + TailwindCSS)
- **Vite** for build tooling
- **React Router** for navigation
- **TanStack Query** (formerly React Query) for server state management
- **Zustand** for client state management
- **React Hook Form** for form handling
- **Recharts** for all charts and data visualizations

## Code Conventions
- **File naming:** 
  - **PascalCase**: React components (folders with `index.tsx`), Pages, Types
    - Components: `Button/index.tsx`, `Dialog/index.tsx`
    - Pages: `Detail.tsx`, `Content.tsx`
    - Types: `Dataset.ts`, `Training.ts`, `User.ts`
  - **camelCase**: All hooks (API and non-API)
    - API hooks: `useGetAllDatasetsQuery.ts`, `usePostTrainingMutation.ts`
    - Shared hooks: `useAuth.ts`, `useLocalStorage.ts`, `useDebouncedValue.ts`
    - Page hooks: `useMetricsDashboardData.ts`, `useBenchmarkPreferences.ts`
  - **kebab-case**: Utils, stores, route files, query keys, constants
    - Utils: `format-number.ts`, `clear-empty-object-attributes.ts`
    - Stores: `auth.store.ts`, `user.store.ts`
    - Route files: `admin-route.tsx`, `neodata-route.tsx`
    - Query keys: `query-keys.ts`
    - Constants: `api-constants.ts`
  - **Index files**: Use `index.ts` or `index.tsx` for cleaner imports (e.g., `components/index.ts` enables `@/components` imports)
- '*.tsx' → React components  
- '*.ts' → utilities, types, and configs  
- **Export patterns:**
  - **Named exports** preferred for most components and utilities
  - **Default exports** required only for lazy-loaded pages/components (React.lazy() requires default exports)
- **Import order:** Handled automatically by Biome (lints and fixes on commit)
- **Code style:**
  - Use single quotes for strings  
  - Indent with 2 spaces  
  - No trailing whitespace  
  - Use 'const' for immutables  
  - Template strings for interpolation  
  - Use optional chaining and nullish coalescing

## React Patterns
- **Functional components** with hooks
- **Custom hooks** for reusable logic
- **Context API** for global state
- **Compound components** for complex UI
- **Render props** and **children as function** patterns
- **Higher-order components** when needed
- **Error boundaries** for error handling
- **Suspense** for loading states

## TypeScript Guidelines
- Define **interfaces** for component props and data structures
- Use **generic types** for reusable components
- **TypeScript any**: Acceptable when building custom solutions that cannot infer types (e.g., generic hooks, type utilities). Otherwise, avoid and use proper typing
- Use **union types** for component variants
- Implement **strict mode** configurations
- Use **utility types** (Pick, Omit, Partial, etc.)

## Performance Optimization
- Use **React.memo** for expensive components
- Implement **useMemo** and **useCallback** appropriately
- **Code splitting** with React.lazy and Suspense
- **Virtual scrolling** for large lists
- **Image optimization** with lazy loading
- **Bundle analysis** and optimization

## Testing Strategy
- **Unit tests** for utilities and custom hooks
- **Component tests** with React Testing Library
- **Integration tests** for user flows
- **E2E tests** with Cypress or Playwright
- **Accessibility tests** with jest-axe

## Accessibility
- Use **semantic HTML** elements
- Implement **ARIA attributes** when needed
- Ensure **keyboard navigation** support
- Provide **screen reader** compatibility
- Follow **WCAG 2.1 AA** guidelines
- Test with **accessibility tools**

## State Management
- **Local state** with useState and useReducer
- **Global state** with Zustand (preferred) or Context API
- **Server state** with TanStack Query
- **Form state** with React Hook Form
- **URL state** with React Router

## shadcn/ui Guidelines
- Use **shadcn/ui** as the primary component library
- **Copy components** from shadcn/ui registry, don't install as package
- **Customize components** by modifying the copied code
- Follow **Radix UI** patterns for accessibility
- Use **TailwindCSS** classes for styling
- **Compose components** using shadcn/ui primitives
- **Extend components** by adding new variants and props

## Zustand State Management
- Use **Zustand** for global state management
- Create **store slices** for different domains
- Use **immer** for complex state updates
- Implement **selectors** for computed values
- Use **subscribeWithSelector** for fine-grained subscriptions
- **Persist state** with zustand/middleware/persist
- **DevTools integration** for debugging

## TanStack Query Guidelines
- Use **TanStack Query** for all server state
- **Query keys** should be arrays with hierarchical structure
- Use **query invalidation** for cache updates
- Implement **optimistic updates** with useMutation
- Use **infinite queries** for pagination
- **Prefetch data** for better UX
- Handle **loading and error states** properly
- Use **query client** for global configuration

## Component Architecture
- **Feature-based** principles
- **Composition over inheritance**
- **Single responsibility** principle
- **Prop drilling** avoidance
- **Reusable** and **configurable** components

## Security Best Practices
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

## Error Handling
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

## Loading States
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

## Task Completion Checklist
Before completing a task, verify the following:

- [ ] **TypeScript errors** — Check for type errors using `read_lints` on modified files
- [ ] **Linter errors** — Verify no linting issues in changed files
- [ ] **MSW mocks** — Update MSW handlers if API contracts or data structures changed
- [ ] **Unused imports** — Remove any unused imports or variables
- [ ] **Console logs** — Remove development debug console.log statements (error logging with console.error is acceptable for user reporting and monitoring)
- [ ] **Import order** — Biome handles this automatically (no manual verification needed)
- [ ] **File naming** — Components use PascalCase folders with index.tsx
- [ ] **Export consistency** — Use named exports for most components, default exports only for lazy-loaded pages
- [ ] **Type safety** — Avoid 'any' types, use proper TypeScript types
- [ ] **Error handling** — Ensure proper error handling for async operations
- [ ] **Loading states** — Verify loading states are implemented where needed
- [ ] **Accessibility** — Check keyboard navigation and ARIA attributes if applicable

**Reference**
Refer to React official documentation and modern React patterns for best practices.


# Types Organization

## File Organization

All types are located in `src/types/` directory:

```
src/types/
├── Dataset.ts          # Dataset-related types
├── User.ts            # User-related types
├── Training.ts        # Training-related types
├── List.ts            # Generic list types
└── ...
```

### File Structure Principles

- **One file per domain concept** - Each file represents a domain entity or concept
- **PascalCase file names** - Match the primary type name (e.g., `Dataset.ts` contains `Dataset` interface)
- **Related types grouped together** - Types that belong to the same domain should be in the same file

## File Structure

### Basic Structure

Each type file should contain:

1. **Primary type** - The main type that matches the filename
2. **Related types** - Supporting types, interfaces, and enums for the domain
3. **Exports** - All types should be exported

### Example Structure

```typescript
// Dataset.ts
export interface DatasetFile {
  id: string
  paths: string[]
}

export enum DatasetStatusEnum {
  Running = 'RUNNING',
  Processing = 'PROCESSING',
  AwaitingModeling = 'AWAITING_MODELING',
  Ready = 'READY',
  Failed = 'FAILED',
}

export type DatasetStructure = 'EVENT_BASED' | 'FEATURE_BASED'

export type DatasetUsage = 'training' | 'benchmark' | 'partitioned'

export interface Dataset {
  id: string
  name: string
  description?: string
  structure: DatasetStructure
  status: DatasetStatusEnum
  // ... other fields
}
```

## Naming Conventions

### Interfaces

Use **PascalCase** for interfaces:

```typescript
export interface Dataset { ... }
export interface AuthUser { ... }
export interface Training { ... }
```

### Types

Use **PascalCase** for type aliases:

```typescript
export type DatasetStructure = 'EVENT_BASED' | 'FEATURE_BASED'
export type DatasetUsage = 'training' | 'benchmark' | 'partitioned'
```

### Enums

Use **PascalCase with "Enum" suffix** for enums:

```typescript
export enum DatasetStatusEnum {
  Running = 'RUNNING',
  Processing = 'PROCESSING',
  // ...
}

export enum DataModelTypeEnum {
  PreTraining = 'pre-training',
  PostTraining = 'post-training',
}
```

**Important**: Always suffix enum names with "Enum" to distinguish them from other types. This makes it clear when you're working with an enum versus an interface or type alias.

### Why "Enum" Suffix?

- **Clarity**: Distinguishes enums from interfaces and types
- **Consistency**: All enums follow the same naming pattern
- **Type safety**: Makes it easier to identify enum types in code

## Cross-References

### Importing Types

When types from one file need to reference types from another:

```typescript
// User.ts
import type { Space } from './Spaces'

export interface User {
  id: string
  email: string
  spaces?: Space[]  // Reference to Space type
  // ...
}
```

### Best Practices for Cross-References

- **Use `import type`** - For type-only imports (better for tree-shaking)
- **Avoid circular dependencies** - Structure types to prevent circular imports
- **Group related types** - Keep types that reference each other in the same file when possible

### Example: Avoiding Circular Dependencies

```typescript
// Good: User imports Space, but Space doesn't import User
// User.ts
import type { Space } from './Spaces'
export interface User {
  spaces?: Space[]
}

// Spaces.ts
export interface Space {
  id: string
  name: string
  // No User reference here
}
```

## When to Create New Type Files

### Create New File When:

1. **New domain concept** - A new entity or domain concept requires its own file
   - Example: New `Connector` entity → `Connector.ts`

2. **Distinct domain** - Types belong to a different domain
   - Example: `Dataset.ts` vs `Training.ts` (different domains)

3. **Generic/reusable types** - Types used across multiple domains
   - Example: `List.ts` for list response types, `Option.ts` for option types

### Group in Same File When:

1. **Related types** - Types that are closely related to the primary type
   - Example: `DatasetFile`, `DatasetStatusEnum` in `Dataset.ts`

2. **Supporting types** - Types that only make sense in context of the primary type
   - Example: `TrainingTimeRange` in `Training.ts`

## File Naming Examples

| Domain Concept | File Name | Primary Type |
|--------------|----------|--------------|
| Dataset | `Dataset.ts` | `Dataset` |
| User | `User.ts` | `User` |
| Training | `Training.ts` | `Training` |
| Generic List | `List.ts` | `ListResponse<T>` |
| Option | `Option.ts` | `Option` |

## Complete Example

```typescript
// Dataset.ts
export interface DatasetFile {
  id: string
  paths: string[]
}

export enum DatasetStatusEnum {
  Running = 'RUNNING',
  Processing = 'PROCESSING',
  AwaitingModeling = 'AWAITING_MODELING',
  Ready = 'READY',
  Failed = 'FAILED',
}

export type DatasetStructure = 'EVENT_BASED' | 'FEATURE_BASED'

export type DatasetUsage = 'training' | 'benchmark' | 'partitioned'

export interface Dataset {
  id: string
  name: string
  description?: string
  structure: DatasetStructure
  data_domain?: string[]
  created_at: string
  created_by?: string
  updated_at: string
  updated_by?: string
  files?: null | DatasetFile[]
  path: string
  connector_name: string
  size: number
  rows?: number
  features?: number
  status: DatasetStatusEnum
  usage?: DatasetUsage
  file_analysis_started_at?: string
  file_analysis_completed_at?: string
  processing_started_at?: string
  processing_completed_at?: string
}
```

## Best Practices

- **One domain per file** - Keep related types together
- **Match filename to primary type** - File name should match the main interface/type
- **Use Enum suffix** - Always suffix enums with "Enum"
- **Export all types** - Make types available for import
- **Use `import type`** - For type-only imports
- **Avoid circular dependencies** - Structure types carefully
- **Group related types** - Keep supporting types with their primary type

## Examples

- **Domain types**: `src/types/Dataset.ts`, `src/types/User.ts`, `src/types/Training.ts`
- **Generic types**: `src/types/List.ts`, `src/types/Option.ts`
- **Cross-references**: `User.ts` imports `Space` from `Spaces.ts`



# Zustand Store Patterns

## Store Organization

Organize stores based on scope and reusability:

### Global Stores
**Location**: `src/stores/[name].store.ts`
- Stores used across multiple pages/features
- App-wide state management (e.g., authentication, feature flags)
- Examples: `auth.store.ts`, `feature-flags.store.ts`

### Page-Scoped Stores
**Location**: `pages/[Feature]/hooks/[name].store.ts` or `pages/[Feature]/[name].store.ts`
- Stores used only within a specific page/feature
- Feature-specific state management
- Examples: `file-analysis.store.ts`, `model-selector.store.ts`

### Decision Criteria
- **Global stores**: State needed across multiple pages/features (e.g., user authentication, app settings)
- **Page-scoped stores**: State specific to one feature/page (e.g., file selection in dataset analysis, model selection in training)

## Naming Conventions

**All stores use the same pattern**: `[name].store.ts` (kebab-case filename)

- **File naming**: `[name].store.ts` (kebab-case)
  - Global: `auth.store.ts`, `feature-flags.store.ts`
  - Page-scoped: `file-analysis.store.ts`, `model-selector.store.ts`
- **Store hook**: `use[name]Store` (camelCase)
  - Examples: `useAuthStore`, `useFileAnalysisStore`, `useModelSelectorStore`
- **State interface**: `[Name]State` (PascalCase)
  - Examples: `AuthState`, `FileAnalysisState`, `ModelSelectorState`

**Important**: Use consistent `[name].store.ts` naming for both global and page-scoped stores. This avoids confusion and maintains consistency across the codebase.

## Store Structure

### Basic Structure

```typescript
import { create } from 'zustand'
import { persist } from 'zustand/middleware'

interface MyState {
  // State properties
  value: string
  count: number
  
  // Actions
  setValue: (value: string) => void
  increment: () => void
  reset: () => void
}

export const useMyStore = create<MyState>()(
  persist(
    (set, get) => ({
      // Initial state
      value: '',
      count: 0,
      
      // Actions
      setValue: (value) => set({ value }),
      increment: () => set(state => ({ count: state.count + 1 })),
      reset: () => set({ value: '', count: 0 }),
    }),
    {
      name: 'my-storage',
    },
  ),
)
```

### Structure Guidelines

1. **Define state interface first** - TypeScript interface for type safety
2. **Use `create<StateInterface>()`** - Generic type parameter for type safety
3. **Export store hook** - Named export (e.g., `useAuthStore`, `useFileAnalysisStore`)
4. **Export selectors separately** - For better performance and reusability

### Selectors Pattern

Export selectors separately for better performance and reusability:

```typescript
// In store file
export const useAuthStore = create<AuthState>()(...)

// Export selectors
export const selectUser = (state: AuthState) => state.user
export const selectIsAuthenticated = (state: AuthState) => state.isAuthenticated
export const selectIsLoading = (state: AuthState) => state.isLoading

// Usage in components
const user = useAuthStore(selectUser)
const isAuthenticated = useAuthStore(selectIsAuthenticated)
```

## Middleware Patterns

### Persist Middleware

Use `persist` for localStorage persistence:

```typescript
import { persist } from 'zustand/middleware'

export const useMyStore = create<MyState>()(
  persist(
    (set, get) => ({
      // ... state and actions
    }),
    {
      name: 'my-storage', // localStorage key
      partialize: (state) => ({
        // Only persist specific fields
        value: state.value,
        // Exclude sensitive or computed fields
      }),
    },
  ),
)
```

### SubscribeWithSelector Middleware

Use `subscribeWithSelector` for fine-grained subscriptions:

```typescript
import { subscribeWithSelector } from 'zustand/middleware'

export const useAuthStore = create<AuthState>()(
  subscribeWithSelector(
    persist(
      (set, get) => ({
        // ... state and actions
      }),
      { name: 'auth-storage' },
    ),
  ),
)
```

### Combining Middleware

Combine multiple middleware when needed:

```typescript
export const useAuthStore = create<AuthState>()(
  subscribeWithSelector(
    persist(
      (set, get) => ({
        // ... state and actions
      }),
      { name: 'auth-storage' },
    ),
  ),
)
```

## State Management Patterns

### Simple Updates

Use `set` for simple state updates:

```typescript
setValue: (value: string) => set({ value }),
```

### Accessing Current State

Use `get` for accessing current state in actions:

```typescript
increment: () => {
  const current = get()
  set({ count: current.count + 1 })
},
```

### Complex State Updates

Use functional updates for complex state:

```typescript
addItem: (item: Item) => {
  set(state => ({
    items: [...state.items, item],
  })),
},
```

### Reset Method

Implement `reset()` method for clearing state:

```typescript
reset: () => set({
  value: '',
  count: 0,
  items: [],
}),
```

### Getter Functions

Use getter functions for computed values:

```typescript
interface FileAnalysisState {
  selectedFileIds: Set<string>
  getFileById: (fileId: string) => FileTreeNode | null
  getFilesFromFolder: (folderId: string) => string[]
}

export const useFileAnalysisStore = create<FileAnalysisState>()(
  (set, get) => ({
    selectedFileIds: new Set<string>(),
    
    getFileById: (fileId: string) => {
      const { fileTree } = get()
      return findNodeById(fileTree, fileId)
    },
    
    getFilesFromFolder: (folderId: string) => {
      const { fileTree } = get()
      // ... computation logic
      return ids
    },
  }),
)
```

## Best Practices

- **Type safety**: Always define TypeScript interfaces for state
- **Selectors**: Export selectors separately for better performance
- **Persistence**: Use `persist` middleware only when state needs to survive page reloads
- **Middleware order**: Apply middleware in correct order (outermost first)
- **Reset methods**: Implement `reset()` for stores with complex state
- **Getter functions**: Use getter functions for computed values instead of storing them
- **Avoid over-fetching**: Use selectors to subscribe only to needed state slices

## Examples

- **Global store**: `src/stores/auth.store.ts` - Authentication state with persist middleware
- **Page-scoped store**: `src/pages/NeoData/Datasets/FileAnalysis/useFileAnalysisStore.ts` - File analysis state with complex getters

**Note**: Existing page-scoped stores may use `use*Store.ts` naming. They should be renamed to follow the `[name].store.ts` convention for consistency.