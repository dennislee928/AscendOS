import { tokens } from '../tokens';

export type CardDensity = 'compact' | 'comfortable';
export type CardVariant = 'default' | 'raised' | 'subtle';

export interface CardOptions {
  density?: CardDensity;
  variant?: CardVariant;
}

export function cardClasses(options: CardOptions = {}): string {
  const classes = ['ui-card', tokens.color.surface, tokens.radius.md];

  classes.push(options.density === 'compact' ? tokens.spacing.sm : tokens.spacing.md);

  if (options.variant === 'raised') {
    classes.push('ui-card-raised');
  } else if (options.variant === 'subtle') {
    classes.push('ui-card-subtle');
  }

  return classes.join(' ');
}
