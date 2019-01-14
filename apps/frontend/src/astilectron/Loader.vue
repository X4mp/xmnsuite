<template>
    <div class="astiloader" id="astiloader">
        <div class="astiloader-background"></div>
        <div class="astiloader-table"><div class="astiloader-content"><i class="fa fa-spinner fa-spin fa-3x fa-fw"></i></div></div>
    </div>
</template>

<script>
if (typeof asticode === "undefined") {
    var asticode = {};
}
asticode.loader = {
    hide: function() {
        document.getElementById("astiloader").style.display = "none";
    },
    show: function() {
        document.getElementById("astiloader").style.display = "block";
    }
};


    export default {
      name: 'Loader'
    }
</script>

<style>
.astiloader {
    color: #fff;
    display: none;
    height: 100%;
    left: 0;
    position: absolute;
    top: 0;
    width: 100%;
    z-index: 1;
}
.astiloader-background {
    background-color: #000;
    height: 100%;
    left: 0;
    opacity: 0.7;
    position: absolute;
    top: 0;
    width: 100%;
}
.astiloader-table {
    display: table;
    height: 100%;
    left: 0;
    position: absolute;
    top: 0;
    width: 100%;
}
.astiloader-content {
    display: table-cell;
    height: 100%;
    text-align: center;
    vertical-align: middle;
    width: 100%;
}
</style>
