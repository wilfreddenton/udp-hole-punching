<template>
  <div id="app">
    <transition name="fade">
      <messages v-if="connected"></messages>
    </transition>
    <container></container>
  </div>
</template>

<script>
import Container from './components/Container'
import Messages from './components/Messages'
import { mapGetters } from 'vuex'

export default {
  name: 'app',
  components: {
    Container,
    Messages
  },
  computed: {
    ...mapGetters({
      connected: 'connected'
    })
  }
}
</script>

<style lang="scss">
* {
  box-sizing: border-box;
}

.cf:after {
  content: "";
  display: table;
  clear: both;
}

body {
  margin: 0;
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
  color: #444444;
  background-color: #efefef;
}

input {
  font-family: -apple-system, BlinkMacSystemFont, sans-serif;
}

#app {
  position: relative;
}

.fade-enter-active, .fade-leave-active {
  transition: opacity .1s
}
.fade-enter, .fade-leave-to /* .fade-leave-active in <2.1.8 */ {
  opacity: 0
}

h1 {
  margin-top: 0;
  margin-bottom: 1em;
}
.field {
  margin-top: 1.5em;

  > label:first-child {
    display: block;
    font-size: 12px;
  }

  input[type="text"] {
    outline: none;
    font-size: 16px;
    width: 100%;
    padding: 0.5em 0em;
    border-top: 0px;
    border-right: 0px;
    border-left: 0px;
    border-bottom: 1px solid #ccc;
  }

  .radios {
    padding: 0.5em 0em;
  }
}


input[type="submit"] {
  display: block;
  cursor: pointer;
  background-color: #54BA75;
  color: white;
  border: none;
  font-size: 16px;
  margin-top: 1.5em;
  padding: 0.33em 1.33em;
  outline: none;
  box-shadow: 0 1px 3px rgba(0,0,0,0.5);
  transition: transform 200ms, box-shadow 200ms;

  &:hover {
    transform: translateY(-1px);
    box-shadow: 0 2px 5px rgba(0,0,0,0.5);
  }

  &:active {
    transform: translateY(1px);
    box-shadow: none;
  }
}
</style>
