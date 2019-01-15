<template>
    <div class="astimodaler" id="astimodaler">
        <div class="astimodaler-background"></div>
        <div class="astimodaler-table">
            <div class="astimodaler-wrapper">
                <div class="astimodaler-body">
                    <i class="fa fa-close astimodaler-close" onclick="asticode.modaler.close()"></i>
                    <div id="astimodaler-content"></div>
                </div>
            </div>
        </div>
    </div>
</template>

<script>
if (typeof asticode === "undefined") {
    var asticode = {};
}
asticode.modaler = {
    close: function() {
        if (typeof asticode.modaler.onclose !== "undefined" && asticode.modaler.onclose !== null) {
            asticode.modaler.onclose();
        }
        asticode.modaler.hide();
    },
    hide: function() {
        document.getElementById("astimodaler").style.display = "none";
    },
    setContent: function(content) {
        document.getElementById("astimodaler-content").innerHTML = '';
        document.getElementById("astimodaler-content").appendChild(content);
    },
    show: function() {
        document.getElementById("astimodaler").style.display = "block";
    }
};

export default {
  name: 'Modal'
}

</script>

<style>
.astimodaler {
    color: #fff;
    display: none;
    height: 100%;
    left: 0;
    position: fixed;
    top: 0;
    width: 100%;
    z-index: 1;
}
.astimodaler-background {
    background-color: #000;
    height: 100%;
    left: 0;
    opacity: 0.7;
    position: absolute;
    top: 0;
    width: 100%;
}
.astimodaler-table {
    display: table;
    height: 100%;
    left: 0;
    position: absolute;
    top: 0;
    width: 100%;
}
.astimodaler-wrapper {
    display: table-cell;
    height: 100%;
    padding: 2em;
    text-align: center;
    vertical-align: middle;
    width: 100%;
}
.astimodaler-body {
    background-color: #f1f1f1;
    border-radius: 5px;
    color: #333;
    padding: 2em;
    position: relative;
}
.astimodaler-close {
    cursor: pointer;
    position: absolute;
    right: 0.65em;
    top: 0.65em;
}
</style>
