<template>
  <div class="connect">
    <h1>Welcome, {{username}}</h1>
    <p>Your ID is: <code>{{id}}</code></p>
    <p>Wait for a peer to connect to you or enter a peer's ID below.</p>
    <form @submit.prevent="onSubmit">
      <div class="field">
        <label for="peerID">Peer ID</label>
        <input id="peerID" type="text" placeholder="a37309618856280592f0a4342525e22ccdda656e0c509dd67f7394ec721df031" :value="peerID" @input="updatePeerID" />
      </div>
      <input type="submit" value="submit" />
    </form>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'connect',
  methods: {
    updatePeerID: function (e) {
      this.$store.dispatch('updatePeerID', e.target.value)
    },
    onSubmit: function (e) {
      if (this.peerID !== '') {
        this.$socket.emit('establish', this.peerID)
      }
    }
  },
  computed: {
    ...mapGetters({
      username: 'username',
      id: 'id',
      peerID: 'peerID'
    })
  }
}
</script>

<style lang="scss" scoped >
</style>
