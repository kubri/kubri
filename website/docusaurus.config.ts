import { themes as prismThemes } from 'prism-react-renderer'
import type { Config } from '@docusaurus/types'
import type * as Preset from '@docusaurus/preset-classic'

const defaultLocale = 'en'

const config: Config = {
  title: 'Kubri',
  tagline: 'Sign and release software for common package managers and software update frameworks.',
  favicon: 'img/favicon.svg',

  url: 'https://kubri.dev',
  baseUrl: '/',

  organizationName: 'kubri',
  projectName: 'kubri',

  onBrokenLinks: 'throw',
  onBrokenMarkdownLinks: 'warn',

  i18n: {
    defaultLocale,
    locales: [defaultLocale],
  },

  plugins: [
    'docusaurus-plugin-sass',
    [
      './src/plugins/changelog/index.js',
      {
        blogTitle: 'Kubri changelog',
        blogDescription: 'Keep yourself up-to-date about new features in every release',
        blogSidebarCount: 'ALL',
        blogSidebarTitle: 'Changelog',
        routeBasePath: '/changelog',
        showReadingTime: false,
        postsPerPage: 20,
        archiveBasePath: null,
        authorsMapPath: 'authors.json',
        feedOptions: {
          type: 'all',
          title: 'Kubri changelog',
          description: 'Keep yourself up-to-date about new features in every release',
          language: defaultLocale,
        },
      },
    ],
  ],

  presets: [
    [
      'classic',
      {
        docs: {
          sidebarPath: './sidebars.ts',
          editUrl: 'https://github.com/kubri/kubri/tree/master/website/',
        },
        blog: {
          showReadingTime: true,
          editUrl: 'https://github.com/kubri/kubri/tree/master/website/',
        },
        theme: {
          customCss: './src/css/custom.scss',
        },
      } satisfies Preset.Options,
    ],
  ],

  themeConfig: {
    // Replace with your project's social card
    // image: 'img/docusaurus-social-card.jpg',
    navbar: {
      title: 'Kubri',
      logo: {
        alt: 'Kubri Logo',
        src: 'img/logo.svg',
      },
      style: 'dark',
      items: [
        {
          type: 'docSidebar',
          sidebarId: 'tutorialSidebar',
          position: 'left',
          label: 'Docs',
        },
        // { to: '/blog', label: 'Blog', position: 'left' },
        {
          href: 'https://github.com/kubri/kubri',
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
            {
              label: 'Tutorial',
              to: '/docs/intro',
            },
          ],
        },
        {
          title: 'Community',
          items: [
            {
              label: 'GitHub',
              href: 'https://github.com/kubri/kubri',
            },
            {
              label: 'Stack Overflow',
              href: 'https://stackoverflow.com/questions/tagged/kubri',
            },
          ],
        },
        {
          title: 'More',
          items: [
            // {
            //   label: 'Blog',
            //   to: '/blog',
            // },
            {
              label: 'Changelog',
              to: '/changelog',
            },
          ],
        },
      ],
      // copyright: `Copyright © ${new Date().getFullYear()} Adam Bouqdib.`,
    },
    prism: {
      theme: prismThemes.github,
      darkTheme: prismThemes.dracula,
      additionalLanguages: ['bash'],
    },
  } satisfies Preset.ThemeConfig,
}

export default config
