<script lang="ts">
	import { enhance } from '$app/forms';
	import type { ActionData, PageData } from './$types';

	let { data, form }: { data: PageData; form: ActionData } = $props();

	let loading = $state(false);
</script>

<div class="min-h-screen bg-[var(--fc-page)] text-[var(--fc-fg)] flex flex-col justify-between selection:bg-[var(--fc-accent)] selection:text-[var(--fc-accent-fg)]">
	<header class="border-b border-[var(--fc-border)] px-6 py-4 flex items-center justify-between">
		<div class="flex items-center gap-3">
			<div class="w-8 h-8 rounded-md bg-[var(--fc-accent)] text-[var(--fc-accent-fg)] font-mono font-bold flex items-center justify-center text-sm shadow-xs">
				F
			</div>
			<div class="flex items-center gap-2">
				<span class="font-semibold text-sm tracking-tight">Facile Chat</span>
				<span class="text-xs px-2 py-0.5 rounded-full bg-[var(--fc-surface)] text-[var(--fc-fg-muted)] border border-[var(--fc-border)] font-mono">XMPP</span>
			</div>
		</div>
		<div class="flex items-center gap-2 text-xs text-[var(--fc-fg-muted)]">
			<span class="w-2 h-2 rounded-full bg-[var(--fc-success)] inline-block"></span>
			<span class="font-mono">{data.domain}</span>
		</div>
	</header>

	<main class="max-w-4xl mx-auto px-6 py-12 w-full grid grid-cols-1 md:grid-cols-12 gap-12 items-start">
		<div class="md:col-span-7 flex flex-col gap-6">
			<div>
				<div class="inline-flex items-center gap-2 px-2.5 py-1 rounded-md bg-[var(--fc-surface)] border border-[var(--fc-border)] text-xs text-[var(--fc-fg-muted)] mb-3">
					<span>Decentralized</span>
					<span>&bull;</span>
					<span>End-to-End Encrypted</span>
					<span>&bull;</span>
					<span>Self-Hosted</span>
				</div>
				<h1 class="text-3xl font-semibold tracking-tight text-[var(--fc-fg)]">
					Private messaging for the Facile team.
				</h1>
				<p class="text-[var(--fc-fg-muted)] mt-2 text-sm leading-relaxed">
					Self-hosted instant messaging, persistent group channels, file sharing, and encrypted WebRTC audio/video calls. No tracking, no metadata federation, and zero external dependencies.
				</p>
			</div>

			<div class="grid grid-cols-2 gap-3 pt-2">
				<div class="p-3.5 rounded-lg border border-[var(--fc-border)] bg-[var(--fc-component)]">
					<div class="text-xs font-semibold text-[var(--fc-fg)]">OMEMO E2EE</div>
					<div class="text-xs text-[var(--fc-fg-muted)] mt-1">Multi-device Double Ratchet encryption for 1:1 and private rooms.</div>
				</div>
				<div class="p-3.5 rounded-lg border border-[var(--fc-border)] bg-[var(--fc-component)]">
					<div class="text-xs font-semibold text-[var(--fc-fg)]">WebRTC Calling</div>
					<div class="text-xs text-[var(--fc-fg-muted)] mt-1">Low-latency peer-to-peer audio and video via dedicated Coturn relay.</div>
				</div>
				<div class="p-3.5 rounded-lg border border-[var(--fc-border)] bg-[var(--fc-component)]">
					<div class="text-xs font-semibold text-[var(--fc-fg)]">Persistent Rooms</div>
					<div class="text-xs text-[var(--fc-fg-muted)] mt-1">Private channels on <code class="text-xs font-mono">conference.{data.domain}</code> with history sync.</div>
				</div>
				<div class="p-3.5 rounded-lg border border-[var(--fc-border)] bg-[var(--fc-component)]">
					<div class="text-xs font-semibold text-[var(--fc-fg)]">Native Push</div>
					<div class="text-xs text-[var(--fc-fg-muted)] mt-1">XEP-0357 notifications wake iOS and Android clients for incoming calls.</div>
				</div>
			</div>

			<div class="pt-4 border-t border-[var(--fc-border)]">
				<div class="text-xs font-semibold uppercase tracking-wider text-[var(--fc-fg-muted)] mb-3">Recommended Clients</div>
				<div class="flex flex-wrap gap-2 text-xs">
					<a href="https://dino.im" target="_blank" rel="noreferrer" class="px-3 py-1.5 rounded-md border border-[var(--fc-border)] bg-[var(--fc-surface)] hover:bg-[var(--fc-component)] transition-colors flex items-center gap-1.5">
						<span class="font-medium">Dino</span>
						<span class="text-[var(--fc-fg-muted)] font-mono text-[10px]">Desktop</span>
					</a>
					<a href="https://conversations.im" target="_blank" rel="noreferrer" class="px-3 py-1.5 rounded-md border border-[var(--fc-border)] bg-[var(--fc-surface)] hover:bg-[var(--fc-component)] transition-colors flex items-center gap-1.5">
						<span class="font-medium">Conversations</span>
						<span class="text-[var(--fc-fg-muted)] font-mono text-[10px]">Android</span>
					</a>
					<a href="https://siskin.im" target="_blank" rel="noreferrer" class="px-3 py-1.5 rounded-md border border-[var(--fc-border)] bg-[var(--fc-surface)] hover:bg-[var(--fc-component)] transition-colors flex items-center gap-1.5">
						<span class="font-medium">Siskin IM</span>
						<span class="text-[var(--fc-fg-muted)] font-mono text-[10px]">iOS</span>
					</a>
					<a href="https://gajim.org" target="_blank" rel="noreferrer" class="px-3 py-1.5 rounded-md border border-[var(--fc-border)] bg-[var(--fc-surface)] hover:bg-[var(--fc-component)] transition-colors flex items-center gap-1.5">
						<span class="font-medium">Gajim</span>
						<span class="text-[var(--fc-fg-muted)] font-mono text-[10px]">Windows/Linux</span>
					</a>
				</div>
			</div>
		</div>

		<div class="md:col-span-5 w-full">
			<div class="p-6 rounded-xl border border-[var(--fc-border)] bg-[var(--fc-component)] shadow-xs">
				{#if form?.success}
					<div class="flex flex-col gap-4">
						<div class="w-10 h-10 rounded-full bg-[var(--fc-success)]/10 text-[var(--fc-success)] flex items-center justify-center font-bold text-lg">
							&check;
						</div>
						<div>
							<h2 class="text-lg font-semibold text-[var(--fc-fg)]">Account Created</h2>
							<p class="text-xs text-[var(--fc-fg-muted)] mt-1">Your personal XMPP handle is ready to use immediately.</p>
						</div>

						<div class="p-3.5 rounded-lg border border-[var(--fc-border)] bg-[var(--fc-surface)] flex flex-col gap-1">
							<span class="text-[10px] uppercase font-mono text-[var(--fc-fg-muted)] tracking-wider">Your Jabber ID (JID)</span>
							<span class="font-mono font-medium text-sm text-[var(--fc-fg)] select-all">{form.jid}</span>
						</div>

						<div class="p-3.5 rounded-lg border border-[var(--fc-border)] bg-[var(--fc-surface)] text-xs text-[var(--fc-fg-muted)] leading-relaxed flex flex-col gap-1.5">
							<div class="font-semibold text-[var(--fc-fg)]">How to connect:</div>
							<div>1. Open your client (Dino, Conversations, Siskin).</div>
							<div>2. Select <strong>Log In</strong> (not sign up).</div>
							<div>3. Enter your full JID <code class="font-mono text-[11px] bg-[var(--fc-component)] px-1 py-0.5 rounded border border-[var(--fc-border)]">{form.jid}</code> and password.</div>
						</div>

						<a href="/" class="mt-2 text-center text-xs text-[var(--fc-fg-muted)] hover:text-[var(--fc-fg)] underline decoration-[var(--fc-border)] underline-offset-4">
							Register another account
						</a>
					</div>
				{:else}
					<div class="flex flex-col gap-4">
						<div>
							<h2 class="text-lg font-semibold text-[var(--fc-fg)]">Create an Account</h2>
							<p class="text-xs text-[var(--fc-fg-muted)] mt-1">Register your handle to start chatting on {data.domain}.</p>
						</div>

						{#if form?.error}
							<div class="p-3 rounded-lg border border-[var(--fc-danger)]/30 bg-[var(--fc-danger)]/10 text-[var(--fc-danger)] text-xs leading-relaxed">
								{form.error}
							</div>
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
							class="flex flex-col gap-3.5"
						>
							<div class="flex flex-col gap-1.5">
								<label for="username" class="text-xs font-medium text-[var(--fc-fg)]">Username</label>
								<div class="flex items-center rounded-md border border-[var(--fc-border)] bg-[var(--fc-page)] focus-within:ring-2 focus-within:ring-[var(--fc-ring)] focus-within:border-transparent transition-all">
									<input
										id="username"
										name="username"
										type="text"
										required
										autocomplete="username"
										placeholder="username"
										value={form?.username || ''}
										class="w-full bg-transparent px-3 py-2 text-sm text-[var(--fc-fg)] placeholder:text-[var(--fc-fg-muted)]/50 focus:outline-hidden"
									/>
									<span class="pr-3 text-xs font-mono text-[var(--fc-fg-muted)] select-none">@{data.domain}</span>
								</div>
							</div>

							<div class="flex flex-col gap-1.5">
								<label for="password" class="text-xs font-medium text-[var(--fc-fg)]">Password</label>
								<input
									id="password"
									name="password"
									type="password"
									required
									autocomplete="new-password"
									placeholder="At least 8 characters"
									class="w-full rounded-md border border-[var(--fc-border)] bg-[var(--fc-page)] px-3 py-2 text-sm text-[var(--fc-fg)] placeholder:text-[var(--fc-fg-muted)]/50 focus:outline-hidden focus:ring-2 focus:ring-[var(--fc-ring)] focus:border-transparent transition-all"
								/>
							</div>

							<div class="flex flex-col gap-1.5">
								<label for="confirmPassword" class="text-xs font-medium text-[var(--fc-fg)]">Confirm Password</label>
								<input
									id="confirmPassword"
									name="confirmPassword"
									type="password"
									required
									autocomplete="new-password"
									placeholder="Repeat your password"
									class="w-full rounded-md border border-[var(--fc-border)] bg-[var(--fc-page)] px-3 py-2 text-sm text-[var(--fc-fg)] placeholder:text-[var(--fc-fg-muted)]/50 focus:outline-hidden focus:ring-2 focus:ring-[var(--fc-ring)] focus:border-transparent transition-all"
								/>
							</div>

							{#if data.inviteRequired}
								<div class="flex flex-col gap-1.5">
									<label for="inviteCode" class="text-xs font-medium text-[var(--fc-fg)]">Team Invite Code</label>
									<input
										id="inviteCode"
										name="inviteCode"
										type="text"
										required
										placeholder="Provided by administrator"
										class="w-full rounded-md border border-[var(--fc-border)] bg-[var(--fc-page)] px-3 py-2 text-sm text-[var(--fc-fg)] placeholder:text-[var(--fc-fg-muted)]/50 focus:outline-hidden focus:ring-2 focus:ring-[var(--fc-ring)] focus:border-transparent transition-all"
									/>
								</div>
							{/if}

							<button
								type="submit"
								disabled={loading}
								class="mt-2 w-full py-2.5 px-4 rounded-md bg-[var(--fc-accent)] text-[var(--fc-accent-fg)] text-xs font-medium tracking-tight hover:opacity-90 active:scale-[0.99] transition-all disabled:opacity-50 disabled:pointer-events-none cursor-pointer"
							>
								{loading ? 'Creating account...' : 'Create Account'}
							</button>
						</form>
					</div>
				{/if}
			</div>
		</div>
	</main>

	<footer class="border-t border-[var(--fc-border)] px-6 py-4 flex flex-col sm:flex-row items-center justify-between text-xs text-[var(--fc-fg-muted)] gap-2">
		<div>Facile Suite &middot; Privacy-first infrastructure</div>
		<div class="flex items-center gap-4">
			<a href="https://chat.facile.studio/admin/" class="hover:text-[var(--fc-fg)] transition-colors">Admin Console</a>
			<a href="https://compliance.conversations.im/server/{data.domain}/" target="_blank" rel="noreferrer" class="hover:text-[var(--fc-fg)] transition-colors">Compliance Score</a>
		</div>
	</footer>
</div>
