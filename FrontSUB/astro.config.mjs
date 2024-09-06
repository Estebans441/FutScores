import { defineConfig } from 'astro/config';
import { getMatches } from './src/services/matches';

import tailwind from '@astrojs/tailwind';


const matches = await getMatches();

// https://astro.build/config
export default defineConfig({
  integrations: [tailwind()],
  output: 'server'
});