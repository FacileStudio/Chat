<script lang="ts">
	import { enhance } from '$app/forms';
	import { Alert, Badge, Button, Card, Field, Input, Spinner } from '@facile/muse';
	import type { ActionData, PageData } from './$types';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	let username = $state('');
	let password = $state('');
	let confirmPassword = $state('');
	let inviteCode = $state('');
	let loading = $state(false);
	let copied = $state(false);

	$effect(() => {
		if (form?.username) {
			username = form.username;
		}
	});

	async function copyJid(jid: string) {
		await navigator.clipboard.writeText(jid);
		copied = true;
		setTimeout(() => (copied = false), 2000);
	}
</script>

<svelte:head>
	<title>Chat — Private messaging for Facile</title>
	<meta name="description" content="Self-hosted, end-to-end encrypted messaging for the Facile team." />
</svelte:head>

<div class="grid gap-8 lg:grid-cols-12 lg:items-start lg:gap-12">
	<div class="flex flex-col gap-8 lg:col-span-7">
		<div class="flex flex-col gap-3">
			<div class="flex flex-wrap items-center gap-2">
				<Badge tone="neutral">Self-hosted</Badge>
				<Badge tone="neutral">End-to-end encrypted</Badge>
				<Badge tone="neutral">No telemetry</Badge>
			</div>
			<h1 class="text-fc-3xl font-semibold tracking-tight text-fc-fg sm:text-fc-4xl">
				Private messaging for the team.
			</h1>
			<p class="text-fc-md text-fc-fg-muted leading-relaxed">
				Encrypted direct messaging, persistent group channels, push notifications, and low-latency audio/video calls powered by modern XMPP and WebRTC.
			</p>
		</div>

		<div class="grid gap-3 sm:grid-cols-2">
			<Card class="flex flex-col gap-2 p-4">
				<div class="flex items-center gap-2 text-fc-fg">
					<iconify-icon icon="solar:shield-check-linear" width="18" height="18" class="block"></iconify-icon>
					<span class="text-fc-sm font-semibold">OMEMO E2EE</span>
				</div>
				<p class="text-fc-xs text-fc-fg-muted leading-relaxed">
					Multi-device Double Ratchet encryption for 1:1 chats and private group rooms.
				</p>
			</Card>

			<Card class="flex flex-col gap-2 p-4">
				<div class="flex items-center gap-2 text-fc-fg">
					<iconify-icon icon="solar:videocamera-record-linear" width="18" height="18" class="block"></iconify-icon>
					<span class="text-fc-sm font-semibold">A/V Calls</span>
				</div>
				<p class="text-fc-xs text-fc-fg-muted leading-relaxed">
					Peer-to-peer encrypted WebRTC audio and video calling with dedicated TURN relay.
				</p>
			</Card>

			<Card class="flex flex-col gap-2 p-4">
				<div class="flex items-center gap-2 text-fc-fg">
					<iconify-icon icon="solar:users-group-rounded-linear" width="18" height="18" class="block"></iconify-icon>
					<span class="text-fc-sm font-semibold">Team Channels</span>
				</div>
				<p class="text-fc-xs text-fc-fg-muted leading-relaxed">
					Members-only persistent channels with message archive sync across devices.
				</p>
			</Card>

			<Card class="flex flex-col gap-2 p-4">
				<div class="flex items-center gap-2 text-fc-fg">
					<iconify-icon icon="solar:bell-linear" width="18" height="18" class="block"></iconify-icon>
					<span class="text-fc-sm font-semibold">Mobile Push</span>
				</div>
				<p class="text-fc-xs text-fc-fg-muted leading-relaxed">
					XEP-0357 notifications wake mobile clients on iOS and Android for messages and calls.
				</p>
			</Card>
		</div>

		<div class="flex flex-col gap-3 pt-2">
			<span class="text-fc-xs font-medium tracking-wide uppercase text-fc-fg-muted">Recommended clients</span>
			<div class="flex flex-wrap gap-2">
				<Button
					href="https://dino.im"
					target="_blank"
					rel="noreferrer"
					variant="outline"
					size="sm"
					icon="solar:monitor-linear"
				>
					Dino (Linux, macOS)
				</Button>
				<Button
					href="https://conversations.im"
					target="_blank"
					rel="noreferrer"
					variant="outline"
					size="sm"
					icon="solar:smartphone-linear"
				>
					Conversations (Android)
				</Button>
				<Button
					href="https://siskin.im"
					target="_blank"
					rel="noreferrer"
					variant="outline"
					size="sm"
					icon="solar:smartphone-linear"
				>
					Siskin IM (iOS)
				</Button>
				<Button
					href="https://gajim.org"
					target="_blank"
					rel="noreferrer"
					variant="outline"
					size="sm"
					icon="solar:laptop-linear"
				>
					Gajim (Windows)
				</Button>
			</div>
		</div>
	</div>

	<div class="w-full lg:col-span-5">
		<Card class="flex flex-col gap-5 p-6 sm:p-7">
			{#if form?.success}
				<div class="flex flex-col gap-4">
					<div class="flex flex-col gap-1">
						<h2 class="text-fc-xl font-semibold tracking-tight text-fc-fg">Account created</h2>
						<p class="text-fc-sm text-fc-fg-muted">Your personal handle is ready to connect.</p>
					</div>

					<Alert tone="success" title="Success">
						Your handle has been provisioned on {data.domain}.
					</Alert>

					<div class="flex flex-col gap-1.5 rounded-fc-md bg-fc-surface p-3.5">
						<span class="text-fc-xs font-medium text-fc-fg-muted">Your Jabber ID (JID)</span>
						<div class="flex items-center justify-between gap-2">
							<span class="font-fc-mono text-fc-sm font-semibold text-fc-fg select-all">{form.jid}</span>
							<Button
								variant="ghost"
								size="sm"
								icon={copied ? 'solar:check-circle-linear' : 'solar:copy-linear'}
								onclick={() => copyJid(form.jid)}
							>
								{copied ? 'Copied' : 'Copy'}
							</Button>
						</div>
					</div>

					<div class="flex flex-col gap-2 rounded-fc-md bg-fc-surface p-3.5 text-fc-xs text-fc-fg-muted leading-relaxed">
						<div class="font-semibold text-fc-fg">Connecting from your client:</div>
						<p>1. Open Dino, Conversations, or Siskin.</p>
						<p>2. Select <strong>Log In</strong> (not sign up).</p>
						<p>3. Enter your JID <strong class="font-fc-mono text-fc-fg">{form.jid}</strong> and password.</p>
					</div>

					<Button href="/" variant="outline" size="md" class="w-full">
						Register another account
					</Button>
				</div>
			{:else}
				<div class="flex flex-col gap-1">
					<h2 class="text-fc-xl font-semibold tracking-tight text-fc-fg">Create an account</h2>
					<p class="text-fc-sm text-fc-fg-muted">Register your handle to join {data.domain}.</p>
				</div>

				{#if form?.error}
					<Alert tone="danger">{form.error}</Alert>
				{/if}

				<form
					method="POST"
					use:enhance={() => {
						loading = true;
						return async ({ update }) => {
							loading = false;
							await update();
						};
					}}
					class="flex flex-col gap-4"
				>
					<Field label="Username" helper="3 to 32 lowercase letters, numbers, dots, or hyphens">
						<Input
							bind:value={username}
							name="username"
							type="text"
							placeholder="username"
							autocomplete="username"
							required
							disabled={loading}
						/>
					</Field>

					<Field label="Password" helper="At least 8 characters">
						<Input
							bind:value={password}
							name="password"
							type="password"
							placeholder="••••••••"
							autocomplete="new-password"
							required
							disabled={loading}
						/>
					</Field>

					<Field label="Confirm password">
						<Input
							bind:value={confirmPassword}
							name="confirmPassword"
							type="password"
							placeholder="••••••••"
							autocomplete="new-password"
							required
							disabled={loading}
						/>
					</Field>

					{#if data.inviteRequired}
						<Field label="Invite code" helper="Required team access token">
							<Input
								bind:value={inviteCode}
								name="inviteCode"
								type="text"
								placeholder="Enter invite code"
								required
								disabled={loading}
							/>
						</Field>
					{/if}

					<Button
						type="submit"
						size="lg"
						icon={loading ? undefined : 'solar:user-plus-linear'}
						disabled={loading}
						class="w-full"
					>
						{#if loading}
							<Spinner size="sm" label="Creating account..." />
							Creating account...
						{:else}
							Create account
						{/if}
					</Button>
				</form>
			{/if}
		</Card>
	</div>
</div>
