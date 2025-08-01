# Deployment Plan for Astro Web App to Cloudflare Workers

## Investigation Summary

The Astro app in packages/web is configured for Cloudflare Workers with SSR (output: 'server') using @astrojs/cloudflare adapter. It includes wrangler.jsonc for deployment, routes to opencode.j9xym.com/\*, and vars like VITE_API_URL.

## Build Steps

1. Install dependencies: `npm install`
2. Build the app: `npm run build` (runs `astro build`), generating dist/ with \_worker.js and static assets.

Reference: https://developers.cloudflare.com/workers/framework-guides/web-apps/astro/#deploy-an-existing-astro-project-on-workers (section on on-demand rendering).

## Deploy Steps

1. Ensure Wrangler is installed and logged in: `npx wrangler login`
2. Deploy: `npx wrangler deploy`
   This uploads the Worker and assets to Cloudflare, making it available on configured routes.

Reference: https://developers.cloudflare.com/workers/get-started/guide/ (deploy section) and https://developers.cloudflare.com/workers/framework-guides/web-apps/astro/.

For CI/CD, integrate with Workers Builds or external providers as per https://developers.cloudflare.com/workers/ci-cd/.
