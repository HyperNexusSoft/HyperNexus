/**
 * `borg billing` — Provider quotas, usage tracking, and model pricing
 */
import type { Command } from 'commander';

const TS_URL = 'http://127.0.0.1:4000/trpc';

export function registerBillingCommand(program: Command): void {
  const billing = program
    .command('billing')
    .description('Billing — provider quotas, usage, model pricing, and fallback chains');

  billing
    .command('status')
    .description('Show billing system status')
    .action(async () => {
      const chalk = (await import('chalk')).default;

      let status: any = {};
      try {
        const res = await fetch(`${TS_URL}/billing.getStatus`, { signal: AbortSignal.timeout(5000) });
        if (res.ok) status = (await res.json())?.result?.data ?? {};
      } catch {}

      console.log(chalk.bold.cyan('\n  Billing System\n'));
      console.log(chalk.dim('  Status:          ') + (status.status ?? 'active'));
      console.log(chalk.dim('  Active providers:') + ` ${status.activeProviders ?? 'see provider list'}`);
      console.log(chalk.dim('  Total models:    ') + `${status.totalModels ?? 'see provider list'}`);
      console.log(chalk.dim('  Depleted models: ') + `${status.depletedModels ?? 0}`);
      console.log('');
    });

  billing
    .command('quotas')
    .description('Show provider quota usage')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const chalk = (await import('chalk')).default;
      const Table = (await import('cli-table3')).default;

      let quotas: any[] = [];
      try {
        const res = await fetch(`${TS_URL}/billing.getProviderQuotas`, { signal: AbortSignal.timeout(5000) });
        if (res.ok) quotas = (await res.json())?.result?.data ?? [];
      } catch {}

      if (opts.json) {
        console.log(JSON.stringify({ quotas }, null, 2));
        return;
      }

      console.log(chalk.bold.cyan(`\n  Provider Quotas (${quotas.length})\n`));
      if (quotas.length === 0) {
        console.log(chalk.dim('  No quota data. Configure providers with `borg provider add`.\n'));
        return;
      }

      const table = new Table({
        head: ['Provider', 'Tier', 'Used', 'Limit', 'Remaining', 'Availability'],
        style: { head: ['cyan'] },
      });

      for (const q of quotas) {
        table.push([
          q.name ?? q.provider,
          q.tier ?? '-',
          q.used ?? 0,
          q.limit ?? '∞',
          q.remaining ?? '-',
          q.availability ?? '-',
        ]);
      }

      console.log(table.toString());
      console.log('');
    });

  billing
    .command('fallback')
    .description('Show or set model fallback chain')
    .option('--set <providers...>', 'Set fallback chain (e.g., openai anthropic google)')
    .option('--json', 'Output as JSON')
    .action(async (opts) => {
      const chalk = (await import('chalk')).default;

      if (opts.set) {
        try {
          const res = await fetch(`${TS_URL}/billing.setFallbackChain`, {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ json: { chain: opts.set } }),
            signal: AbortSignal.timeout(5000),
          });
          if (res.ok) {
            console.log(chalk.green(`  ✓ Fallback chain set: ${opts.set.join(' → ')}`));
          }
        } catch (e: any) {
          console.log(chalk.red(`  ✗ Error: ${e.message}`));
        }
        return;
      }

      let chain: any = {};
      try {
        const res = await fetch(`${TS_URL}/billing.getFallbackChain`, { signal: AbortSignal.timeout(5000) });
        if (res.ok) chain = (await res.json())?.result?.data ?? {};
      } catch {}

      if (opts.json) {
        console.log(JSON.stringify(chain, null, 2));
        return;
      }

      console.log(chalk.bold.cyan('\n  Fallback Chain\n'));
      const providers = chain.chain ?? chain.providers ?? chain ?? [];
      if (Array.isArray(providers) && providers.length > 0) {
        for (let i = 0; i < providers.length; i++) {
          const p = typeof providers[i] === 'string' ? providers[i] : providers[i].provider ?? providers[i].name;
          const prefix = i < providers.length - 1 ? '├─' : '└─';
          console.log(`  ${prefix} ${chalk.bold(p)}`);
        }
      } else {
        console.log(chalk.dim('  No fallback chain configured.'));
        console.log(chalk.dim('  Use: borg billing fallback --set openai anthropic google'));
      }
      console.log('');
    });

  billing
    .command('depleted')
    .description('Show depleted/rate-limited models')
    .action(async () => {
      const chalk = (await import('chalk')).default;

      let depleted: any[] = [];
      try {
        const res = await fetch(`${TS_URL}/billing.getDepletedModels`, { signal: AbortSignal.timeout(5000) });
        if (res.ok) depleted = (await res.json())?.result?.data ?? [];
      } catch {}

      console.log(chalk.bold.cyan(`\n  Depleted Models (${depleted.length})\n`));
      if (depleted.length === 0) {
        console.log(chalk.green('  No depleted models. All providers operational.\n'));
        return;
      }

      for (const d of depleted) {
        console.log(chalk.yellow(`  ⚠ ${d.model ?? d.name} (${d.provider})`) + ` — ${d.reason ?? 'rate limited'}`);
      }
      console.log('');
    });
}
