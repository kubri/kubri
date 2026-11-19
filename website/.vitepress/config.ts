import { existsSync, mkdirSync, writeFileSync } from 'node:fs'
import { join } from 'node:path'
import { defineConfig } from 'vitepress'
import { redirectStub, redirects } from './redirects'

const hostname = 'https://kubri.dev'

export default defineConfig({
  title: 'Kubri',
  description:
    'Kubri turns your releases into signed APT, YUM, APK and Pacman repositories, Windows App Installer files and Sparkle update feeds, from one YAML file.',
  cleanUrls: true,
  // GitHub Pages only serves trailing-slash URLs from <page>/index.html. Remove once on Cloudflare.
  rewrites: (id) => (/(^|\/)index\.md$/.test(id) ? id : id.replace(/\.md$/, '/index.md')),
  sitemap: { hostname },
  head: [
    ['link', { rel: 'icon', href: '/favicon.svg' }],
    ['style', {}, '.VPFooter.has-sidebar{display:block!important}'],
  ],
  transformPageData(pageData, { siteConfig: { site } }) {
    const url = new URL(pageData.relativePath.replace(/(?:(^|\/)index)?\.md$/, '$1'), hostname).href
    const title = pageData.title ? `${pageData.title} | ${site.title}` : site.title
    const description = pageData.description || site.description

    pageData.frontmatter.head ??= []
    pageData.frontmatter.head.push(
      ['meta', { property: 'og:url', content: url }],
      ['meta', { property: 'og:title', content: title }],
      ['meta', { property: 'og:description', content: description }],
      ['meta', { property: 'og:type', content: 'website' }],
      ['meta', { property: 'og:site_name', content: site.title }],
      ['link', { rel: 'canonical', href: url }],
    )
  },
  buildEnd({ outDir }) {
    for (const [from, to] of Object.entries(redirects)) {
      const file = join(outDir, from, 'index.html')
      if (existsSync(file)) {
        throw new Error(`${from} is both a page and a redirect`)
      }
      mkdirSync(join(outDir, from), { recursive: true })
      writeFileSync(file, redirectStub(to, hostname))
    }
  },
  themeConfig: {
    logo: { light: '/logo-light.svg', dark: '/logo-dark.svg', alt: 'Kubri' },
    search: { provider: 'local' },
    outline: [2, 3],
    nav: [
      {
        text: 'Guide',
        link: '/docs/intro/',
        activeMatch: '^/docs/(intro|installation|getting-started)/',
      },
      {
        text: 'Configuration',
        link: '/docs/configuration/',
        activeMatch: '^/docs/configuration/',
      },
      {
        text: 'Reference',
        link: '/docs/cli/',
        activeMatch: '^/docs/(cli|version-constraints)/',
      },
    ],
    socialLinks: [{ icon: 'github', link: 'https://github.com/kubri/kubri' }],
    sidebar: [
      {
        text: 'Guide',
        items: [
          { text: 'Introduction', link: '/docs/intro/' },
          { text: 'Installation', link: '/docs/installation/' },
          { text: 'Getting Started', link: '/docs/getting-started/' },
        ],
      },
      {
        text: 'Configuration',
        items: [
          { text: 'Overview', link: '/docs/configuration/' },
          {
            text: 'Source',
            link: '/docs/configuration/source/',
            collapsed: false,
            items: [
              { text: 'GitHub', link: '/docs/configuration/source/github/' },
              { text: 'GitLab', link: '/docs/configuration/source/gitlab/' },
              { text: 'Local Filesystem', link: '/docs/configuration/source/file/' },
            ],
          },
          {
            text: 'Target',
            link: '/docs/configuration/target/',
            collapsed: false,
            items: [
              { text: 'Azure Blob Storage', link: '/docs/configuration/target/azureblob/' },
              { text: 'Google Cloud Storage', link: '/docs/configuration/target/gcs/' },
              { text: 'GitHub', link: '/docs/configuration/target/github/' },
              { text: 'S3', link: '/docs/configuration/target/s3/' },
              { text: 'Local Filesystem', link: '/docs/configuration/target/file/' },
            ],
          },
          {
            text: 'Generators',
            link: '/docs/configuration/generators/',
            collapsed: false,
            items: [
              { text: 'APT Repository', link: '/docs/configuration/generators/apt/' },
              { text: 'Arch Linux Repository', link: '/docs/configuration/generators/arch/' },
              { text: 'YUM Repository', link: '/docs/configuration/generators/yum/' },
              { text: 'Alpine Linux Repository', link: '/docs/configuration/generators/apk/' },
              { text: 'Sparkle Framework', link: '/docs/configuration/generators/sparkle/' },
              {
                text: 'Windows App Installer',
                link: '/docs/configuration/generators/appinstaller/',
              },
            ],
          },
        ],
      },
      {
        text: 'Reference',
        items: [
          { text: 'CLI Reference', link: '/docs/cli/' },
          { text: 'Version Constraints', link: '/docs/version-constraints/' },
        ],
      },
    ],
    footer: {
      copyright: '&copy; 2026 Adam Bouqdib',
    },
    editLink: {
      pattern: 'https://github.com/kubri/kubri/edit/master/website/:path',
    },
  },
})
