<script lang="ts">
    import { onMount } from "svelte";
    import type { Contact } from "./models/contact";
    onMount(async () => {
        getUsers();
    });

    type ContactWithEdit = Contact & { edit: boolean };
    let contacts: ContactWithEdit[] = [];

    let search = "";
    let all = false;

    const getUsers = async (
        search: string = "",
        all: boolean = false,
        page: number = 0,
    ) => {
        try {
            const res = await fetch(
                "/api/contacts?search=" +
                    search +
                    "&matchAll=" +
                    all +
                    "&page=" +
                    page,
                {
                    credentials: "include",
                },
            );

            if (res.ok) {
                const response = await res.json();
                contacts = response.contacts;
            } else {
                contacts = [];
            }
        } catch (e) {
            console.error("Failed to fetch user:", e);
            contacts = [];
        }
    };
</script>

<div>
    <input type="text" placeholder="Search" bind:value={search} />
    <label><input type="checkbox" bind:checked={all} />all</label>
    <button on:click={() => getUsers(search, all)}>Search</button>
    <!-- make nice style toggle all/any-->
</div>

<table>
    <thead>
        <tr>
            <th>id</th>
            <th>name</th>
            <th>linkedinurl</th>
            <th>credlyurl</th>
            <th>tags</th>
            <th>contact</th>
        </tr>
    </thead>
    <tbody>
        {#each contacts as contact}
            <tr>
                <td>{contact.id}</td>
                <td>{contact.name}</td>
                <td>{contact.linkedInUrl}</td>
                <td>{contact.credlyInUrl}</td>
                <td>{contact.tags}</td>
                <td>{contact.contact}</td>
            </tr>
        {/each}
    </tbody>
</table>
