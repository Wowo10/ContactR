<script lang="ts">
  import Login from "./lib/login.svelte";
  import MyRouter from "./lib/myRouter.svelte";

  import type { LoginData } from "./lib/models/user";

  let logged: LoginData | null = null;
</script>

<main>
  <div><img src="/src/assets/logo.gif" alt="logo" /></div>

  <Login bind:user={logged} />

  {JSON.stringify(logged)}

  {#if logged}
    {#if logged.is_valid}
      <MyRouter {logged} />
    {:else}
      <p>Access denied</p>
    {/if}
  {:else}
    <p>Waiting for login...</p>
  {/if}
</main>
