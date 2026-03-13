export function getFriendlyErrorMessage(err) {
  const raw = typeof err === 'string'
    ? err
    : (err && typeof err.message === 'string' ? err.message : '');

  const message = raw.trim();
  if (!message) {
    return 'Ocorreu um erro inesperado. Tente novamente.';
  }

  if (message.includes('UNIQUE constraint failed')) {
    return 'Já existe um registro com estes dados.';
  }
  if (message.includes('FOREIGN KEY constraint failed')) {
    return 'Não foi possível concluir a operação por referência inválida.';
  }
  if (message.includes('CHECK constraint failed')) {
    return 'Os dados informados violam uma regra de validação.';
  }
  if (message.toLowerCase().includes('database is locked') || message.toLowerCase().includes('database table is locked')) {
    return 'O banco de dados está ocupado no momento. Tente novamente em alguns segundos.';
  }
  if (message.toLowerCase().includes('insufficient stock')) {
    return 'Estoque insuficiente para concluir a movimentação.';
  }
  if (message.toLowerCase().includes('must be greater than zero')) {
    return 'Existem campos com valor inválido. Revise as quantidades.';
  }
  if (message.toLowerCase().includes('required')) {
    return 'Preencha os campos obrigatórios.';
  }

  return message;
}
