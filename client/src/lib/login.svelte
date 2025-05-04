<script lang="ts">
  import { onMount } from "svelte";
  import type { User } from "./models/user";

  export let user: User | null = null;

  onMount(async () => {
    try {
      const res = await fetch("/api/me", {
        credentials: "include",
      });

      if (res.ok) {
        user = await res.json();
      } else {
        user = null;
      }
    } catch (e) {
      console.error("Failed to fetch user:", e);
      user = null;
    }
  });

  const handleLogin = () => {
    window.location.href = "auth/google";
  };

  const logout = async () => {
    await fetch("/auth/logout", {
      credentials: "include",
    });
    user = null;
  };
</script>

<div>
  {#if user}
    {user.user_email}
    <button onclick={logout}>Logout</button>
  {:else}
    <button onclick={handleLogin}> Login With Google </button>
  {/if}
</div>
