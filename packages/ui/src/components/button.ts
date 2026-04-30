export type ButtonVariant = 'primary' | 'secondary';

export function buttonClasses(variant: ButtonVariant = 'primary'): string {
  if (variant === 'secondary') {
    return 'ui-button ui-button-secondary';
  }

  return 'ui-button ui-button-primary';
}
