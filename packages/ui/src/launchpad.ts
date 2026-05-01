export interface LaunchpadRoute {
  name: string;
  href: string;
  summary: string;
  actionLabel: string;
  phase: string;
  domain: string;
  status: string;
  featured?: boolean;
}

export const launchpadRoutes = [
  {
    name: 'Auth',
    href: '/module-auth',
    summary: 'Authentication and onboarding entry points.',
    actionLabel: 'Open auth',
    phase: 'Phase 6',
    domain: 'Identity',
    status: 'Scaffolded',
    featured: true,
  },
  {
    name: 'Billing',
    href: '/module-billing',
    summary: 'Checkout, plan management, and account billing flows.',
    actionLabel: 'Open billing',
    phase: 'Phase 6',
    domain: 'Revenue',
    status: 'Scaffolded',
    featured: true,
  },
  {
    name: 'Dashboard',
    href: '/module-dashboard',
    summary: 'Summary views and operational dashboards.',
    actionLabel: 'Open dashboard',
    phase: 'Phase 6',
    domain: 'Operations',
    status: 'Scaffolded',
  },
] as const satisfies readonly LaunchpadRoute[];

export const featuredLaunchpadRoutes = launchpadRoutes.filter((route) => route.featured);
