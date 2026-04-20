/**
 * Copyright (c) Facebook, Inc. and its affiliates.
 *
 * This source code is licensed under the MIT license found in the
 * LICENSE file in the root directory of this source tree.
 */

import { HtmlClassNameProvider, PageMetadata, ThemeClassNames } from '@docusaurus/theme-common'
import BlogLayout from '@theme/BlogLayout'
import type { Props } from '@theme/BlogListPage'
import BlogListPaginator from '@theme/BlogListPaginator'
import BlogPostItems from '@theme/BlogPostItems'
import ChangelogItem from '@theme/ChangelogItem'
import ChangelogListHeader from '@theme/ChangelogList/Header'
import SearchMetadata from '@theme/SearchMetadata'
import clsx from 'clsx'
import type { JSX } from 'react'

function ChangelogListMetadata({ metadata }: Props): JSX.Element {
  const { blogTitle, blogDescription } = metadata
  return (
    <>
      <PageMetadata title={blogTitle} description={blogDescription} />
      <SearchMetadata tag="blog_posts_list" />
    </>
  )
}

function ChangelogListContent({ metadata, items, sidebar }: Props): JSX.Element {
  const { blogTitle } = metadata
  return (
    <BlogLayout sidebar={sidebar}>
      <ChangelogListHeader blogTitle={blogTitle} />
      <BlogPostItems items={items} component={ChangelogItem} />
      <BlogListPaginator metadata={metadata} />
    </BlogLayout>
  )
}

export default function ChangelogList(props: Props): JSX.Element {
  return (
    <HtmlClassNameProvider
      className={clsx(ThemeClassNames.wrapper.blogPages, ThemeClassNames.page.blogListPage)}
    >
      <ChangelogListMetadata {...props} />
      <ChangelogListContent {...props} />
    </HtmlClassNameProvider>
  )
}
