<script lang="ts">
    import { onMount } from "svelte";
    import type { User } from "./models/user";
    import ConfirmButton from "./common/ConfirmButton.svelte";

    type UserWithEdit = User & { edit: boolean };
    let users: UserWithEdit[] = [];
    let addUserData = {
        name: "",
        email: "",
        valid_until: new Date(
            Date.now() + 1000 * 60 * 60 * 24 * 365,
        ).toISOString(),
        is_admin: false,
    };

    onMount(async () => {
        try {
            const res = await fetch("/api/users", {
                credentials: "include",
            });

            if (res.ok) {
                users = await res.json();
            } else {
                users = [];
            }
        } catch (e) {
            console.error("Failed to fetch user:", e);
            users = [];
        }
    });

    const editUser = async (user: UserWithEdit) => {
        const res = await fetch("/api/users", {
            method: "PUT",
            credentials: "include",
            body: JSON.stringify(user as User),
        });

        if (!res.ok) {
            console.error("Failed to modify user:", await res.text());
            return;
        }

        user.edit = false;

        users = users.map((u) => {
            if (u.id === user.id) {
                return user;
            }
            return u;
        });
    };

    const addUser = async () => {
        const res = await fetch("/api/users", {
            method: "POST",
            credentials: "include",
            body: JSON.stringify(addUserData),
        });

        if (!res.ok) {
            console.error("Failed to modify user:", await res.text());
            return;
        }

        addUserData.name = "";
        addUserData.email = "";
        addUserData.valid_until = new Date(
            Date.now() + 1000 * 60 * 60 * 24 * 365,
        ).toISOString();
        addUserData.is_admin = false;

        let newUser = await res.json();
        newUser.edit = false;
        users = [...users, newUser];
    };

    const deleteUser = async (id: string) => {
        const res = await fetch(`/api/users/${id}`, {
            method: "DELETE",
            credentials: "include",
        });

        if (!res.ok) {
            console.error("Failed to delete user:", await res.text());
            return;
        }
        users = users.filter((u) => u.id !== id);
    };
</script>

<table>
    <thead>
        <tr>
            <th>ID</th>
            <th>Name</th>
            <th>Email</th>
            <th>ValidUntil</th>
            <th>Admin</th>
            <th></th>
            <th></th>
        </tr>
    </thead>
    <tbody>
        {#each users as user}
            <tr>
                <td>
                    {user.id}
                </td>
                <td>
                    {#if user.edit}
                        <input type="text" bind:value={user.name} />
                    {:else}
                        {user.name}
                    {/if}
                </td>
                <td>
                    {#if user.edit}
                        <input type="text" bind:value={user.email} />
                    {:else}
                        {user.email}
                    {/if}
                </td>
                <td>
                    {#if user.edit}
                        <input type="text" bind:value={user.valid_until} />
                    {:else}
                        {user.valid_until}
                    {/if}
                </td>
                <td>
                    {#if user.edit}
                        <input
                            type="checkbox"
                            bind:checked={user.is_admin}
                            readonly
                        />
                    {:else}
                        {user.is_admin ? "Yes" : "No"}
                    {/if}
                </td>
                <td>
                    <ConfirmButton
                        onConfirm={() => editUser(user)}
                        activateCallback={() => {
                            user.edit = !user.edit;
                        }}
                        dectivateCallback={() => {
                            user.edit = !user.edit;
                        }}
                        label="Edit"
                        timeout={9999999999}
                    />
                </td>
                <td>
                    <ConfirmButton
                        onConfirm={() => deleteUser(user.id)}
                        label="Remove"
                    />
                </td>
            </tr>
        {/each}
        <tr>
            <!-- todo: make visible on confirm -->
            <td></td>
            <td>
                <input type="text" bind:value={addUserData.name} />
            </td>
            <td>
                <input type="text" bind:value={addUserData.email} />
            </td>
            <td>
                <input type="text" bind:value={addUserData.valid_until} />
            </td>
            <td>
                <input type="checkbox" bind:checked={addUserData.is_admin} />
            </td>
            <td>
                <ConfirmButton onConfirm={() => addUser()} label="Add" />
            </td>
            <td></td>
        </tr>
    </tbody>
</table>
