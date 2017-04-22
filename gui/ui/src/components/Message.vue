<template>
  <div class="message" :class="{sent: message.sent}">
    <p>{{ message.text }}</p>
  </div>
</template>

<script>
import { mapGetters } from 'vuex'

export default {
  name: 'message',
  props: {
    message: Object
  },
  computed: {
    ...mapGetters({
    })
  }
}
</script>

<style lang="scss" scoped>
.message {
  box-shadow: 0 2px 6px rgba(0,0,0,0.1);
  border-radius: 5px;
  display: inline-block;
  right: 0px;
  max-width: 400px;
  background-color: white;
  padding: 1em;
  margin: 0.5em 0em;
  float: left;

  p {
    line-height: 1.4;
    margin: 0;
  }

  &.sent {
    background-color: #54BA75;
    color: white;
    float: right;
  }
}
</style>
