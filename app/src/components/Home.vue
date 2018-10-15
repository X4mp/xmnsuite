<template>
    <div id="home">
        <v-icon name="beer"/>
        <TitlePage
            v-bind:title="title"
        />

        <LeftMenu />
    </div>
</template>

<script>
import LeftMenu from './LeftMenu.vue'
import TitlePage from './titles/Page.vue'

export default {
  components: {
    LeftMenu,
    TitlePage
  },
  data () {
    return {
      'title': 'XMN Title'
    }
  },
  methods: {
  }
}
</script>

<style>
</style>
