<script lang="ts">
    import { onMount } from "svelte";
    import type { Contact } from "./models/contact";
    import { PAGE_SIZE } from "./common/constants";
    onMount(async () => {
        getContacts();
    });

    type ContactWithEdit = Contact & { edit: boolean };
    let contacts: ContactWithEdit[] = [];

    let search = "";
    let all = false;
    let pages = 0;

    const getContacts = async (
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
                pages = Math.ceil(response.count / PAGE_SIZE);
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
    <button on:click={() => getContacts(search, all)}>Search</button>
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
{#if pages > 1}
    <div class="pagination">
        {#each Array(pages) as _, i}
            <button on:click={() => getContacts(search, all, i)}>{i}</button>
        {/each}
    </div>
{/if}
