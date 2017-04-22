<template>
  <div class="register">
    <form @submit.prevent="onSubmit">
      <h1>👊 UDP Hole Punching 👊</h1>
      <div class="field">
        <label for="username" class="label">Username</label>
        <input id="username" type="text" placeholder="Yusuke" :value="username" @input="updateUsername"/>
      </div>
      <input type="submit" value="submit" />
    </form>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'register',
  methods: {
    onSubmit: function (e) {
      if (this.username !== '') {
        this.$socket.emit('enter', JSON.stringify({username: this.username, protocol: this.protocol}))
      }
    },
    updateProtocol: function (e) {
      const val = e.target.value
      if (val !== this.protocol) {
        this.$store.dispatch('updateProtocol', val)
      }
    },
    updateUsername: function (e) {
      this.$store.dispatch('updateUsername', e.target.value)
    }
  },
  computed: {
    ...mapGetters({
      username: 'username',
      protocol: 'protocol'
    })
  }
}
</script>

<style lang="scss" scoped>
</style>
