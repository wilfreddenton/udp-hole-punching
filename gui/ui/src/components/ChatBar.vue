<template>
  <div id="chat-bar">
    <form @submit.prevent="onSubmit">
      <textarea @focusin="onFocusIn" @focusout="onFocusOut" @keydown="onKeyDown" ref="textarea" rows="1" placeholder="message" type="text" v-model="text" ></textarea>
    </form>
  </div>
</template>

<script>
import Vue from 'vue'
import autosize from 'autosize'
import { mapGetters } from 'vuex'

export default {
  name: 'chat-bar',
  props: {
    onFocusIn: Function,
    onFocusOut: Function
  },
  data () {
    return {
      text: ''
    }
  },
  methods: {
    onSubmit () {
      if (this.text === '') {
        return
      }
      const text = this.text
      this.$store.dispatch('newMessage', { sent: true, text: text })
      this.$socket.emit('message', text)
      this.text = ''
      Vue.nextTick(() => {
        autosize.update(this.$refs.textarea)
      })
    },
    onKeyDown (e) {
      if (e.keyCode === 13) {
        if (!e.shiftKey) {
          e.preventDefault()
          this.onSubmit()
        }
      }
    }
  },
  computed: {
    ...mapGetters({
    })
  },
  mounted () {
    autosize(this.$refs.textarea)
    this.$refs.textarea.addEventListener('focusin', this.onFocusIn)
    this.$refs.textarea.addEventListener('focusout', this.onFocusOut)
  },
  beforeDestroy () {
    this.$refs.textarea.removeEventListener('focusin', this.onFocusIn)
    this.$refs.textarea.removeEventListener('focusout', this.onFocusOut)
  }
}
</script>

<style lang="scss" scoped>
#chat-bar {
  border-radius: 5px;
  transition: box-shadow 200ms, margin 200ms, transform 200ms;

  textarea {
    max-height: 50vh;
    border-radius: 5px;
    font-size: 16px;
    line-height: 1.4;
    display: block;
    width: 100%;
    border: none;
    position: relative;
    outline: none;
    resize: none;
    padding: 0;
  }
}
</style>
