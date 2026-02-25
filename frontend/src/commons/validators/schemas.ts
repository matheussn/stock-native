import * as yup from 'yup'

export const namedEntitySchema = yup.object({
  name: yup.string().required('Nome é obrigatório').max(120, 'Nome muito longo'),
})

export const productSchema = yup.object({
  name: yup.string().required('Nome é obrigatório').max(120, 'Nome muito longo'),
  categoryId: yup.string().required('Categoria é obrigatória'),
  measureUnit: yup
    .mixed<'kg' | 'g' | 'l' | 'ml' | 'unidade'>()
    .oneOf(['kg', 'g', 'l', 'ml', 'unidade'], 'Unidade de medida inválida')
    .required('Unidade de medida é obrigatória'),
  packageAmount: yup
    .string()
    .required('Peso/volume é obrigatório')
    .matches(/^\d+(\.\d{1,3})?$/, 'Use número positivo com até 3 casas decimais'),
  lowStockThreshold: yup
    .string()
    .required('Limite é obrigatório')
    .matches(/^\d+(\.\d{1,3})?$/, 'Use número positivo com até 3 casas decimais'),
})

export const movementSchema = yup.object({
  productId: yup.string().required('Produto é obrigatório'),
  quantity: yup
    .string()
    .required('Quantidade é obrigatória')
    .matches(/^\d+(\.\d{1,3})?$/, 'Use número positivo com até 3 casas decimais'),
  sourceId: yup.string().required('Campo obrigatório'),
  movementDate: yup.string().required('Data é obrigatória'),
  note: yup.string().default('').max(255, 'Observação muito longa'),
})

export const reportFilterSchema = yup.object({
  dateFrom: yup.string().required('Data inicial é obrigatória'),
  dateTo: yup.string().required('Data final é obrigatória'),
})

export type NamedEntityForm = yup.InferType<typeof namedEntitySchema>
export type ProductForm = yup.InferType<typeof productSchema>
export type MovementForm = yup.InferType<typeof movementSchema>
export type ReportFilterForm = yup.InferType<typeof reportFilterSchema>
