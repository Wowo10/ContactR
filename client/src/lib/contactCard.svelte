<script lang="ts">
    import { onMount } from "svelte";
    import { emptyContact, type Contact } from "./models/contact";
    import ConfirmButton from "./common/ConfirmButton.svelte";

    let edit = false;
    let id = "";
    let contact: Contact = emptyContact();
    let tagsStr = "";

    onMount(async () => {
        const url = new URL(window.location.href);
        const lastSuffix = url.pathname.split("/").pop();

        if (lastSuffix !== "contact") {
            id = lastSuffix!;
            edit = false;

            getContact(id);
        } else {
            edit = true;
            contact = emptyContact();
        }
    });

    const getContact = async (id: string) => {
        try {
            const res = await fetch("/api/contacts/" + id, {
                credentials: "include",
            });

            if (res.ok) {
                const response = await res.json();
                contact = response;

                tagsStr = contact.tags.join(", ");
            } else {
                contact = emptyContact();
            }
        } catch (e) {
            console.error("Failed to fetch user:", e);
            contact = emptyContact();
        }
    };

    const putContact = async () => {
        contact.tags = tagsStr.split(",").map((s) => s.trim());

        try {
            const res = await fetch("/api/contacts", {
                method: "PUT",
                credentials: "include",
                body: JSON.stringify(contact),
            });

            if (res.ok) {
                edit = false;
            } else {
                console.error("Failed to edit user:", await res.text());
            }
        } catch (e) {
            console.error("Failed to edit user:", e);
        }
    };

    const postContact = async () => {
        contact.tags = tagsStr.split(",").map((s) => s.trim());
        try {
            const res = await fetch("/api/contacts", {
                method: "POST",
                credentials: "include",
                body: JSON.stringify(contact),
            });

            if (res.ok) {
                const response = await res.json();
                contact = response;

                id = contact.id + "";
                tagsStr = contact.tags.join(", ");

                edit = false;

                window.history.pushState(null, "", `/contact/${contact.id}`);
            } else {
                console.error("Failed to edit user:", await res.text());
            }
        } catch (e) {
            console.error("Failed to edit user:", e);
        }
    };
</script>

{#if id !== ""}
    <ConfirmButton
        onConfirm={putContact}
        activateCallback={() => {
            edit = !edit;
        }}
        dectivateCallback={() => {
            edit = !edit;
        }}
        label="Edit"
        timeout={0}
        confirmLabel="Save"
    />
{:else}
    <ConfirmButton onConfirm={postContact} label="Create" confirmLabel="Save" />
{/if}
{#if edit}
    <div>
        <label>Name: <input type="text" bind:value={contact.name} /></label>
    </div>
    <div>
        <label>
            Linkedin:
            <input type="text" bind:value={contact.linkedInUrl} />
        </label>
    </div>
    <div>
        <label>
            Credly:
            <input type="text" bind:value={contact.credlyInUrl} />
        </label>
    </div>
    <div>
        <label>
            Contact:
            <input type="text" bind:value={contact.contact} />
        </label>
    </div>
    <div>
        <label>
            Tags:
            <input type="text" bind:value={tagsStr} />
        </label>
    </div>
{:else}
    <div>Name: {contact.name}</div>
    <div>Linkedin: {contact.linkedInUrl}</div>
    <div>Credly: {contact.credlyInUrl}</div>
    <div>Contact: {contact.contact}</div>
    <div>Tags: {contact.tags}</div>
{/if}
