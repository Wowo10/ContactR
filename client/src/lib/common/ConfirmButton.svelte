<script lang="ts">
    export let onConfirm: () => void;
    export let activateCallback: (() => void) | null = null;
    export let dectivateCallback: (() => void) | null = null;
    export let label = "Remove";
    export let confirmLabel = "Sure?";
    export let timeout = 3000; // auto-reset after 3s

    let confirming = false;
    let timer: any;

    function handleClick() {
        if (!confirming) {
            if (activateCallback) {
                activateCallback();
            }
            confirming = true;
            timer = setTimeout(() => {
                if (dectivateCallback) {
                    dectivateCallback();
                }
                confirming = false;
            }, timeout);
        } else {
            if (dectivateCallback) {
                dectivateCallback();
            }
            clearTimeout(timer);
            confirming = false;
            onConfirm();
        }
    }
</script>

<button on:click={handleClick} class:confirming>
    {confirming ? confirmLabel : label}
</button>

<style>
    button {
        padding: 0.5em 1em;
        border: none;
        border-radius: 6px;
        background-color: #eee;
        cursor: pointer;
        transition: background 0.2s ease;
    }

    button.confirming {
        background-color: #e74c3c;
        color: white;
    }
</style>
