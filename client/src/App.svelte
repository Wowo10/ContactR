<script lang="ts">
  const handleLogin = () => {
    window.location.href = "http://localhost:3000/auth/google"
  }
</script>

<main> 

  <div>
    super Logo
  </div>

  <div>
    <button onclick={handleLogin}>
      Login With Google
    </button>
  </div>

  <div>
    <input type="text" placeholder="Search" />
    <label><input type="checkbox" />all</label> <!-- make nice style toggle all/any-->
  </div>

  <div id="results">
    results
  </div>

</main>

<style>

</style>
