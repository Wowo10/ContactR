<script>
    import { onMount } from 'svelte';
    import { writable } from 'svelte/store';

    // Store for auth state
    export const user = writable(null);

    onMount(async () => {
        try {
            const res = await fetch('http://localhost:3000/api/me', {
                credentials: 'include'  // Important for cookies/session
            });
            
            if (res.ok) {
                const data = await res.json();
                user.set(data);
            } else {
                user.set(null);  // Not logged in
            }
        } catch (e) {
            console.error('Failed to fetch user:', e);
            user.set(null);
        }
    });

    const handleLogin = () => {
        window.location.href = "http://localhost:3000/auth/google"
    }

    async function logout() {
        await fetch('/auth/logout', {
            method: 'POST',
            credentials: 'include'
        });
        user.set(null);
    }
</script>

<div>
    {#if $user}
    {$user.user_email}
    <button onclick={logout}>Logout</button>
    {:else}
    <button onclick={handleLogin}>
        Login With Google
    </button>
    {/if}
</div>