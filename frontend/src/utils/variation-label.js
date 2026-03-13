export function getVariationLabel({description, baseQuantity, baseUnit, fallback}) {
  const normalizedDescription = (description ?? '').trim();
  if (normalizedDescription) {
    return normalizedDescription;
  }

  if (Number.isFinite(baseQuantity) && baseQuantity > 0 && baseUnit) {
    return `${baseQuantity} ${baseUnit}`;
  }

  return fallback;
}
