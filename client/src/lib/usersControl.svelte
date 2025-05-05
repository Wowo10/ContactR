<script lang="ts">
    import { onMount } from "svelte";
    import type { User } from "./models/user";
    import ConfirmButton from "./common/ConfirmButton.svelte";

    type UserWithEdit = User & { edit: boolean };
    let users: UserWithEdit[] = [];
    let addVisible = false;
    let addUserData = {
        name: "",
        email: "",
        valid_until: new Date(
            Date.now() + 1000 * 60 * 60 * 24 * 365,
        ).toISOString(),
        is_admin: false,
    };

    onMount(async () => {
        getUsers();
    });

    const getUsers = async (page: number = 0) => {
        try {
            const res = await fetch("/api/users?page=" + page, {
                credentials: "include",
            });

            if (res.ok) {
                const response = await res.json();
                users = response.users;
            } else {
                users = [];
            }
        } catch (e) {
            console.error("Failed to fetch user:", e);
            users = [];
        }
    };

    const editUser = async (user: UserWithEdit) => {
        if (!user.email || !user.valid_until) {
            getUsers();
            return;
        }

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
        if (!addUserData.email || !addUserData.valid_until) {
            return;
        }

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
                        timeout={0}
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
            <td></td>
            <td>
                <input
                    type="text"
                    bind:value={addUserData.name}
                    class:hide={!addVisible}
                />
            </td>
            <td>
                <input
                    type="text"
                    bind:value={addUserData.email}
                    class:hide={!addVisible}
                />
            </td>
            <td>
                <!-- todo: a datepicker would be nice -->
                <input
                    type="text"
                    bind:value={addUserData.valid_until}
                    class:hide={!addVisible}
                />
            </td>
            <td>
                <input
                    type="checkbox"
                    bind:checked={addUserData.is_admin}
                    class:hide={!addVisible}
                />
            </td>
            <td>
                <ConfirmButton
                    onConfirm={() => addUser()}
                    activateCallback={() => {
                        addVisible = true;
                    }}
                    dectivateCallback={() => {
                        addVisible = false;
                    }}
                    label="Add"
                    timeout={0}
                />
            </td>
            <td></td>
        </tr>
    </tbody>
</table>

<style>
    .hide {
        display: none;
    }
</style>
