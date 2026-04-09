import { writable } from 'svelte/store';

export const currentFilePath = writable<string | null>(null);
export const currentFileName = writable<string>('Untitled');
export const isDirty = writable<boolean>(false);
export const sidebarRootPath = writable<string | null>(null);
export const cursorPosition = writable<{ line: number; col: number }>({ line: 1, col: 1 });
export const currentContent = writable<string>('');
export const savedContent = writable<string>('');
export const lineEnding = writable<'CRLF' | 'LF'>('LF');
export const showWhitespace = writable<boolean>(false);
