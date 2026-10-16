import { fail } from '@sveltejs/kit';
import type { Actions, PageServerLoad } from './$types';

export const load: PageServerLoad = async () => {
	const domain = process.env.DOMAIN || 'chat.facile.studio';
	const inviteRequired = Boolean(process.env.INVITE_CODE);
	return {
		domain,
		inviteRequired
	};
};

export const actions: Actions = {
	default: async ({ request }) => {
		const data = await request.formData();
		const username = data.get('username')?.toString().trim().toLowerCase() || '';
		const password = data.get('password')?.toString() || '';
		const confirmPassword = data.get('confirmPassword')?.toString() || '';
		const inviteCode = data.get('inviteCode')?.toString().trim() || '';

		const domain = process.env.DOMAIN || 'chat.facile.studio';
		const requiredInvite = process.env.INVITE_CODE;
		const adminPassword = process.env.ADMIN_PASSWORD || process.env.REGISTER_ADMIN_PASSWORD || '';
		const apiUrl = process.env.EJABBERD_API_URL || 'http://ejabberd:5443/api';

		if (!username) {
			return fail(400, { error: 'Username is required', username });
		}

		if (!/^[a-z0-9._-]+$/.test(username)) {
			return fail(400, {
				error: 'Username may only contain lowercase letters, numbers, dots, and hyphens',
				username
			});
		}

		if (username.length < 3 || username.length > 32) {
			return fail(400, {
				error: 'Username must be between 3 and 32 characters',
				username
			});
		}

		if (!password || password.length < 8) {
			return fail(400, {
				error: 'Password must be at least 8 characters long',
				username
			});
		}

		if (password !== confirmPassword) {
			return fail(400, {
				error: 'Passwords do not match',
				username
			});
		}

		if (requiredInvite && inviteCode !== requiredInvite) {
			return fail(403, {
				error: 'Invalid team invite code',
				username
			});
		}

		try {
			const auth = 'Basic ' + Buffer.from(`admin@${domain}:${adminPassword}`).toString('base64');
			const res = await fetch(`${apiUrl}/register`, {
				method: 'POST',
				headers: {
					'Content-Type': 'application/json',
					Authorization: auth
				},
				body: JSON.stringify({
					user: username,
					host: domain,
					password: password
				})
			});

			const body = await res.json().catch(() => null);

			if (res.status === 200 && (body === 0 || body?.status === 'ok' || body === 'ok')) {
				return {
					success: true,
					jid: `${username}@${domain}`
				};
			}

			if (body?.message?.includes('already registered') || body?.message?.includes('conflict') || body?.code === 409) {
				return fail(409, {
					error: 'This username is already taken. Please choose another.',
					username
				});
			}

			return fail(400, {
				error: body?.message || 'Registration failed on the server',
				username
			});
		} catch (err: any) {
			return fail(500, {
				error: 'Failed to connect to the authentication service: ' + (err.message || 'unknown error'),
				username
			});
		}
	}
};
