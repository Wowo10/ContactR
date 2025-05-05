<script lang="ts">
  import Login from "./lib/login.svelte";
  import Router from "./lib/router.svelte";

  import type { LoginData } from "./lib/models/user";

  let logged: LoginData | null = null;
</script>

<main>
  <div>super Logo</div>

  <Login bind:user={logged} />

  {JSON.stringify(logged)}

  {#if logged}
    {#if logged.is_valid}
      <Router {logged} />
    {:else}
      <p>Access denied</p>
    {/if}
  {:else}
    <p>Waiting for login...</p>
  {/if}
</main>
