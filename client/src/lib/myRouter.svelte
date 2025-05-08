<script lang="ts">
    import ContactSearch from "./contactSearch.svelte";
    import ContactCard from "./contactCard.svelte";
    import UsersControl from "./usersControl.svelte";
    import type { LoginData } from "./models/user";
    import { onMount } from "svelte";

    export let logged: LoginData | null = null;

    let route = window.location.pathname;

    const routes: Record<
        string,
        typeof ContactSearch | typeof UsersControl | typeof ContactCard
    > = {
        "/search": ContactSearch,
        "/users": UsersControl,
        "/contact": ContactCard,
    };

    let Current = routes[route];

    function navigate(to: string) {
        if (to.startsWith("/contact/")) {
            window.history.pushState(null, "", to);
            to = "/contact";
        } else {
            window.history.pushState(null, "", to);
        }

        route = to;
        Current = routes[to] || ContactSearch;
    }

    onMount(() => {
        const handlePopState = () => {
            route = window.location.pathname;
            Current = routes[route] || ContactSearch;
        };

        console.log(window.location.pathname);

        if (window.location.pathname === "/") {
            navigate("/search");
        } else {
            navigate(window.location.pathname);
        }

        window.addEventListener("popstate", handlePopState);

        return () => {
            window.removeEventListener("popstate", handlePopState);
        };
    });
</script>

{#if logged?.is_admin}
    <nav>
        <button
            on:click={() => navigate("/search")}
            class:active={route === "/search"}
        >
            Search
        </button>
        <button
            on:click={() => navigate("/contact")}
            class:active={route === "/contact"}
        >
            Add
        </button>
        <button
            on:click={() => navigate("/users")}
            class:active={route === "/users"}
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
