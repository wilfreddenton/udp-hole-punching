<template>
  <div id="messages">
    <transition-group name="fade" tag="ul">
      <li class="cf" v-for="(message, i) in messages" :key="i">
        <message :message="message"></message>
      </li>
    </transition-group>
    </ul>
  </div>
</template>

<script>
import Vue from 'vue'
import Message from './Message'
import { mapGetters } from 'vuex'

export default {
  name: 'messages',
  components: {
    Message
  },
  computed: {
    ...mapGetters({
      messages: 'messages'
    })
  },
  watch: {
    messages () {
      if (document.body.scrollTop + window.innerHeight !== document.body.scrollHeight) {
        return
      }
      Vue.nextTick(() => {
        window.scrollTo(0, document.body.scrollHeight)
      })
    }
  }
}
</script>

<style lang="scss" scoped>
#messages {
  position: relative;
  padding: 1em 1em 6em 1em;
  width: 100%;

  ul {
    max-width: 550px;
    margin: 0 auto;
    position: relative;
    display: block;
    padding: 0;
    list-style-type: none;

    li {
      display: block;
    }
  }
}
</style>
