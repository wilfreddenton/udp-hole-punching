<template>
  <div class="connecting">
    <span>Connecting to {{peerUsername}} at <code>{{peerAddr}}</code></span>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'connecting',
  computed: {
    ...mapGetters({
      peerUsername: 'peerUsername',
      peerAddr: 'peerAddr',
      peerID: 'peerID'
    })
  }
}
</script>

<style lang="scss" scoped >
</style>
