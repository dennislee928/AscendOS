export interface LaunchpadRoute {
  name: string;
  href: string;
  summary: string;
  actionLabel: string;
  featured?: boolean;
}

export const launchpadRoutes = [
  {
    name: 'Auth',
    href: '/module-auth',
    summary: 'Authentication and onboarding entry points.',
    actionLabel: 'Open auth',
    featured: true,
  },
  {
    name: 'Billing',
    href: '/module-billing',
    summary: 'Checkout, plan management, and account billing flows.',
    actionLabel: 'Open billing',
    featured: true,
  },
  {
    name: 'Dashboard',
    href: '/module-dashboard',
    summary: 'Summary views and operational dashboards.',
    actionLabel: 'Open dashboard',
  },
] as const satisfies readonly LaunchpadRoute[];

export const featuredLaunchpadRoutes = launchpadRoutes.filter((route) => route.featured);
