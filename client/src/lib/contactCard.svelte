<script lang="ts">
    const url = new URL(window.location.href);
    const lastSuffix = url.pathname.split("/").pop();

    let Id = "";
    if (lastSuffix !== "contact") {
        Id = lastSuffix!;
    }
</script>

<div>ContactCard {Id}</div>
