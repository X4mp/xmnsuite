<template>
    <h1>{{ title }}</h1>
</template>

<script>
export default {
  name: 'Page',
  props: ['title']
}
</script>

<style>
h1 {
    border-bottom: 1px solid #000000;
}
</style>
