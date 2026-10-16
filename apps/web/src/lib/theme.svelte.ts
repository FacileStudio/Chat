import { browser } from '$app/environment';

export type ThemeMode = 'system' | 'light' | 'dark';

const KEY = 'facile-chat-theme';

function stored(): ThemeMode {
	if (!browser) return 'system';
	const raw = localStorage.getItem(KEY);
	return raw === 'light' || raw === 'dark' || raw === 'system' ? raw : 'system';
}

export const theme = $state({ mode: stored() as ThemeMode });

export function setTheme(mode: ThemeMode) {
	theme.mode = mode;
	if (!browser) return;
	const root = document.documentElement;
	root.classList.toggle('dark', mode === 'dark');
	root.classList.toggle('light', mode === 'light');
	localStorage.setItem(KEY, mode);
}

export function isDark(): boolean {
	if (!browser) return false;
	if (theme.mode !== 'system') return theme.mode === 'dark';
	return window.matchMedia('(prefers-color-scheme: dark)').matches;
}
