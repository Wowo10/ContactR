<script lang="ts">
  import { onMount } from 'svelte';
  import type { Writable } from 'svelte/store';
  import { writable } from 'svelte/store';

  interface User {
    user_email: string;
  }

  export const user: Writable<User | null> = writable(null);

  onMount(async () => {
    try {
      const res = await fetch('/api/me', {
        credentials: 'include',
      });

      if (res.ok) {
        const data = await res.json();
        user.set(data);
      } else {
        user.set(null);
      }
    } catch (e) {
      console.error('Failed to fetch user:', e);
      user.set(null);
    }
  });

  const handleLogin = () => {
    window.location.href = 'auth/google';
  };

  async function logout() {
    await fetch('/auth/logout', {
      credentials: 'include',
    });
    user.set(null);
  }
</script>

<div>
  {#if $user}
    {$user.user_email}
    <button onclick={logout}>Logout</button>
  {:else}
    <button onclick={handleLogin}> Login With Google </button>
  {/if}
</div>
