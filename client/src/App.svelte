<script lang="ts">
  import Login from "./lib/login.svelte";
  import Router from "./lib/router.svelte";

  import type { User } from "./lib/models/user";

  let user: User | null = null;
</script>

<main>
  <div>super Logo</div>

  <Login bind:user />

  {JSON.stringify(user)}

  {#if user}
    {#if user.is_valid}
      <Router {user} />
    {:else}
      <p>Access denied</p>
    {/if}
  {:else}
    <p>Waiting for login...</p>
  {/if}
</main>
