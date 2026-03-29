import { writable } from 'svelte/store';

export const userDisplayName = writable("Anonymous Bat");
export const userAvatarUrl = writable("https://api.dicebear.com/9.x/icons/svg?backgroundColor=ffab91&seed=Adrian");

export const storeMap = {
  display_name: userDisplayName,
  avatar_url: userAvatarUrl,
};
