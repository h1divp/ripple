import { PUBLIC_API_URL } from '$env/static/public';
import { storeMap } from '$lib/stores/user';

export async function getProfile() {
  const url = `${PUBLIC_API_URL}/profile`;
  try {
    const response = await fetch(url, {
      method: 'GET',
      credentials: 'include'
    });
    if (!response.ok) {
      console.error('Profile retrieval failed with status:', response.status);
      return;
    }
    const data = await response.json();
    Object.entries(storeMap).forEach(([key, store]) => {
      if (data[key] !== undefined) {
        store.set(data[key]);
      }
    });
  } catch (err) {
    console.error('Profile retrieval failed:', err);
  }
}

export async function getSession() {
  const url = `${PUBLIC_API_URL}/register`;
  try {
    const response = await fetch(url, {
      method: 'POST',
      credentials: 'include'
    });
    if (!response.ok) {
      console.error('Session registration failed with status:', response.status);
    }
  } catch (err) {
    console.error('Session registration failed:', err);
  }
}
