import { themes as prismThemes } from 'prism-react-renderer';
import type { Config } from '@docusaurus/types';
import type * as Preset from '@docusaurus/preset-classic';

const config: Config = {
  title: 'LiteLog',
  tagline: 'Centralized logging without the infrastructure.',
  favicon: 'img/favicon.ico',

  url: 'https://yashnaiduu.github.io',
  baseUrl: '/Litelog/',

  organizationName: 'yashnaiduu',
  projectName: 'Litelog',
  trailingSlash: false,

  onBrokenLinks: 'warn',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale: 'en',
    locales: ['en'],
  },

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/yashnaiduu/Litelog/edit/main/website/',
        },
        blog: false,
        theme: {
          customCss: './src/css/custom.css',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    image: 'img/logo.png',
    colorMode: {
      defaultMode: 'dark',
      disableSwitch: true,
      respectPrefersColorScheme: false,
    },
    navbar: {
      title: 'LiteLog',
      logo: {
        alt: 'LiteLog',
        src: 'img/logo.png',
        height: 28,
        width: 'auto',
      },
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'docsSidebar',
          position: 'left',
          label: 'Docs',
        },
        {
          href: 'https://github.com/yashnaiduu/Litelog',
          label: 'GitHub',
          position: 'right',
        },
      ],
    },
    footer: {
      style: 'dark',
      links: [
        {
          title: 'Docs',
          items: [
            { label: 'Quick Start', to: '/docs/quick-start' },
            { label: 'Commands', to: '/docs/commands' },
            { label: 'Architecture', to: '/docs/architecture' },
          ],
        },
        {
          title: 'Community',
          items: [
            { label: 'GitHub', href: 'https://github.com/yashnaiduu/Litelog' },
            { label: 'Issues', href: 'https://github.com/yashnaiduu/Litelog/issues' },
            { label: 'Pull Requests', href: 'https://github.com/yashnaiduu/Litelog/pulls' },
          ],
        },
      ],
      copyright: `Copyright ${new Date().getFullYear()} LiteLog. Built with Go and SQLite.`,
    },
    prism: {
      theme: prismThemes.oneDark,
      darkTheme: prismThemes.oneDark,
      additionalLanguages: ['bash', 'json', 'go', 'sql'],
    },
  } satisfies Preset.ThemeConfig,
};

export default config;
