import clsx from 'clsx'
import Link from '@docusaurus/Link'
import useDocusaurusContext from '@docusaurus/useDocusaurusContext'
import Layout from '@theme/Layout'
import Heading from '@theme/Heading'
import { JSX } from 'react'

import styles from './index.module.css'
import Logo from '../../static/img/logo.svg'

function HomepageHeader() {
  const { siteConfig } = useDocusaurusContext()
  return (
    <header className={clsx('hero hero--dark', styles.heroBanner)}>
      <div className="container text--center">
        <Logo width={200} />
        <Heading as="h1" className="hero__title margin-vert--md">
          {siteConfig.title}
        </Heading>
        <p className="hero__subtitle">{siteConfig.tagline}</p>
        <div className={styles.buttons}>
          <Link className="button button--secondary button--lg" to="/docs/intro">
            Get Started
          </Link>
        </div>
      </div>
    </header>
  )
}

const FeatureList = [
  {
    title: 'Easy to Use',
    description: (
      <>
        Write a simple YAML file to define your release process. Kubri will take care of the rest.
      </>
    ),
  },
  {
    title: 'Works Everywhere',
    description: <>Kubri has zero dependencies and works on Linux, Windows, and MacOS.</>,
  },
  {
    title: 'Free and Open Source',
    description: (
      <>
        Kubri is free to use and open source. It is licensed under the{' '}
        <Link to="/license">MIT License</Link>.
      </>
    ),
  },
]

const Integrations = [
  {
    title: 'APT',
    description: 'Debian, Ubuntu etc.',
  },
  {
    title: 'YUM / DNF / Zypper',
    description: 'RHEL, Fedora, CentOS, OpenSUSE etc.',
  },
  {
    title: 'APK',
    description: 'Alpine Linux',
  },
  {
    title: 'App Installer',
    description: 'Windows',
  },
  {
    title: 'Sparkle / WinSparkle',
    description: 'MacOS, Windows',
  },
]

export default function Home(): JSX.Element {
  const { siteConfig } = useDocusaurusContext()
  return (
    <Layout
      title={`Hello from ${siteConfig.title}`}
      description="Description will go into a meta tag in <head />"
    >
      <HomepageHeader />
      <main>
        <section className={styles.section}>
          <div className="container padding-vert--xl">
            <div className="row">
              <div className={clsx('col col--6 margin-vert--md', styles.centered)}>
                <span className={styles.lead}>
                  Kubri takes the hassle out of releasing your software for multiple different
                  platforms. It takes your releases from GitHub, GitLab or cloud storage systems and
                  automatically creates repositories, releases, and artifacts for the platforms you
                  care about.
                </span>
              </div>
              <div className="col col--6">
                <div className="card">
                  <div className="card__header">
                    <Heading as="h3">Supported Platforms</Heading>
                  </div>
                  <div className="card__body">
                    <ul>
                      {Integrations.map((integration, idx) => (
                        <li key={idx}>
                          <strong>{integration.title}</strong> ({integration.description})
                        </li>
                      ))}
                    </ul>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </section>

        <section className={styles.section}>
          <div className="container padding-vert--xl text--center">
            <Heading as="h2" className="margin-bottom--lg">
              How does it work?
            </Heading>
            <div>
              <p>
                Kubri simplifies the release process by allowing you to define your entire release
                pipeline in a single YAML file. Whether your artifacts are stored on GitHub, GitLab,
                or other platforms, Kubri automatically handles the download, repository generation,
                package/metadata signing, and upload to your specified targets, allowing you to
                effortlessly distribute your software to a wide array of platforms.
              </p>
              <p>
                <Link to="/docs/intro#usage">Click here to see an example configuration.</Link>
              </p>
            </div>
          </div>
        </section>

        <section className={styles.section}>
          <div className="container padding-vert--xl text--center">
            <Heading as="h2" className="margin-bottom--lg">
              Why Kubri?
            </Heading>
            <div className="row">
              {FeatureList.map(({ title, description }, idx) => (
                <div className="col" key={idx}>
                  <Heading as="h3">{title}</Heading>
                  <p>{description}</p>
                </div>
              ))}
            </div>
          </div>
        </section>
      </main>
    </Layout>
  )
}
