<script lang="ts">
	import '../app.css';
	import { browser } from '$app/environment';
	import { setTheme, theme } from '$lib/theme.svelte';
	import ThemeToggle from '$lib/components/ThemeToggle.svelte';
	import { Badge, StatusDot } from '@facile/muse';
	import { page } from '$app/state';

	let { children } = $props();

	if (browser) {
		void import('iconify-icon');
		setTheme(theme.mode);
	}

	let domain = $derived(page.data?.domain || 'chat.facile.studio');
</script>

<div class="flex min-h-[100dvh] flex-col bg-fc-page text-fc-fg">
	<header class="border-b border-fc-border">
		<div class="mx-auto flex max-w-5xl items-center justify-between px-4 py-4 sm:px-6">
			<a href="/" class="flex items-center gap-3 transition-opacity hover:opacity-80">
				<iconify-icon icon="solar:chat-round-dots-bold-duotone" width="24" height="24" class="block text-fc-fg"></iconify-icon>
				<span class="text-fc-xl font-semibold tracking-tight">Chat</span>
				<Badge tone="neutral">XMPP</Badge>
			</a>
			<div class="flex items-center gap-4">
				<div class="hidden sm:flex items-center gap-2">
					<StatusDot tone="success" pulse />
					<span class="font-fc-mono text-fc-xs text-fc-fg-muted">{domain}</span>
				</div>
				<ThemeToggle />
			</div>
		</div>
	</header>

	<main class="mx-auto flex w-full max-w-5xl flex-1 flex-col px-4 py-8 sm:px-6 sm:py-12">
		{@render children()}
	</main>

	<footer class="border-t border-fc-border text-center">
		<div class="mx-auto max-w-5xl px-4 py-6 text-fc-sm text-fc-fg-muted sm:px-6">
			&copy; {new Date().getFullYear()} Chat by <a href="https://facile.studio" target="_blank" rel="noreferrer" class="font-medium underline hover:text-fc-fg">Facile.</a>
		</div>
	</footer>
</div>
