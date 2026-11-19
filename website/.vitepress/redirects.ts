/** Old URLs that send readers to a page's new home. Keys are directory paths ending in `/`. */
export const redirects: Record<string, string> = {
  '/changelog/': 'https://github.com/kubri/kubri/blob/master/CHANGELOG.md',
  '/docs/configuration/': '/docs/getting-started/',
  '/docs/guides/': '/docs/version-constraints/',
  '/docs/guides/version-constrains/': '/docs/version-constraints/',
}

/** The page written at an old URL. It keeps the URL's hash when JavaScript is available. */
export function redirectStub(to: string, base: string): string {
  return `<!doctype html>
<html lang="en">
<head>
<meta charset="utf-8">
<title>Redirecting</title>
<meta name="robots" content="noindex">
<link rel="canonical" href="${new URL(to, base).href}">
<script>location.replace(${JSON.stringify(to)} + location.hash)</script>
<meta http-equiv="refresh" content="0; url=${to}">
</head>
<body><a href="${to}">${to}</a></body>
</html>
`
}
