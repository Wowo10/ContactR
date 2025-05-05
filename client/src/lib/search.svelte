<script lang="ts">
    import { onMount } from "svelte";
    import type { Contact } from "./models/contact";
    onMount(async () => {
        getUsers();
    });

    type ContactWithEdit = Contact & { edit: boolean };
    let contacts: ContactWithEdit[] = [];

    const getUsers = async () => {
        try {
            const res = await fetch("/api/contacts", {
                credentials: "include",
            });

            if (res.ok) {
                contacts = await res.json();
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
    <input type="text" placeholder="Search" />
    <label><input type="checkbox" />all</label>
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
