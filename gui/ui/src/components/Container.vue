<template>
  <div id="container" class="cf" :class="{ entered: entered, connecting: connecting, connected: connected, stretch: stretch }">
    <div class="content" :class="{ focus: focus }">
      <transition name="fade" mode="out-in">
        <connect v-if="entered"></connect>
        <connecting v-else-if="connecting"></connecting>
        <div v-else-if="connected">
          <chat-bar :onFocusIn="onFocusIn" :onFocusOut="onFocusOut"></chat-bar>
        </div>
        <register v-else></register>
      </transition>
    </div>
  </div>
</template>

<script>
import Register from './Register'
import Connect from './Connect'
import Connecting from './Connecting'
import ChatBar from './ChatBar'
import { mapGetters } from 'vuex'

export default {
  name: 'container',
  data () {
    return {
      focus: false,
      stretch: false
    }
  },
  components: {
    Connect,
    Connecting,
    ChatBar,
    Register
  },
  methods: {
    onFocusIn () {
      this.focus = true
    },
    onFocusOut () {
      this.focus = false
    }
  },
  computed: {
    ...mapGetters({
      entered: 'entered',
      connecting: 'connecting',
      connected: 'connected'
    })
  },
  watch: {
    connected (val) {
      if (val) {
        setTimeout(() => {
          this.stretch = true
        }, 500)
      }
    }
  }
}
</script>

<style lang="scss" scoped>
#container {
  position: fixed;
  width: 100%;
  bottom: 50%;
  margin-bottom: -112px;
  transition: bottom 500ms, margin-bottom 500ms, padding 500ms;

  .content {
    height: 224px;
    width: 416px;
    padding: 1.5em;
    background-color: white;
    border-radius: 5px;
    position: relative;
    margin: 0 auto;
    box-shadow: 0 5px 10px rgba(0,0,0,0.1);
    white-space: nowrap;
    overflow: hidden;
    transition: height 500ms, width 500ms, padding 500ms, box-shadow 200ms, transform 200ms;
  }
}

#container.entered {
  margin-bottom: -150px;

  .content {
    width: 660px;
    height: 300px;
  }
}

#container.connecting {
  margin-bottom: -33px;

  .content {
    width: 500px;
    height: 67px;
  }
}

#container.connected {
  margin-bottom: 0px;
  bottom: 2em;
  padding: 0 1em;

  .content {
    padding: 1em 1em;
    width: 550px;
    height: 54px;
    box-shadow: 0 5px 10px rgba(0,0,0,0.15);

    &.focus {
      box-shadow: 0 10px 25px rgba(0,0,0,0.1);
      transform: scale(1.05);
    }
  }

  &.stretch .content {
    width: 100%;
    max-width: 550px;
    overflow: visible;
    height: auto;
  }
}

</style>
