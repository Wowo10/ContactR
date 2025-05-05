<script lang="ts">
    import ContactSearch from "./contactSearch.svelte";
    import UsersControl from "./usersControl.svelte";
    import type { LoginData } from "./models/user";

    export let logged: LoginData | null = null;

    let route = "contactSearch";

    const routes: Record<string, typeof ContactSearch | typeof UsersControl> = {
        contactSearch: ContactSearch,
        usersControl: UsersControl,
    };

    let Current = routes[route];

    function navigate(to: string) {
        route = to;
        Current = routes[route];
    }
</script>

{#if logged?.is_admin}
    <nav>
        <button
            on:click={() => navigate("contactSearch")}
            class:active={route === "contactSearch"}
        >
            Search
        </button>
        <button
            on:click={() => navigate("usersControl")}
            class:active={route === "usersControl"}
        >
            Users
        </button>
    </nav>
{/if}

<main>
    <svelte:component this={Current} />
</main>

<style>
    nav {
        display: flex;
        gap: 1rem;
        margin-bottom: 1rem;
    }

    button {
        padding: 0.5rem 1rem;
        cursor: pointer;
    }

    .active {
        background-color: #333;
        color: white;
    }
</style>
