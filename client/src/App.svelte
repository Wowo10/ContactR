<script lang="ts">
  import Login from './lib/login.svelte';
</script>

<main>
  <div>super Logo</div>

  <Login />

  <div>
    <input type="text" placeholder="Search" />
    <label><input type="checkbox" />all</label>
    <!-- make nice style toggle all/any-->
  </div>

  <div id="results">results</div>
</main>

<style>
</style>
