<script lang="ts">
    import Search from "./search.svelte";
    import UsersControl from "./usersControl.svelte";
    import type { User } from "./models/user";

    export let user: User | null = null;

    let route = "search";

    const routes: Record<string, typeof Search | typeof UsersControl> = {
        search: Search,
        usersControl: UsersControl,
    };

    let Current = routes[route];

    function navigate(to: string) {
        route = to;
        Current = routes[route];
    }
</script>

{#if user?.is_admin}
    <nav>
        <button
            on:click={() => navigate("search")}
            class:active={route === "search"}
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
