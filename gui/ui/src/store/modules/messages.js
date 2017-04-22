import * as types from '../mutation-types'

const state = {
  messages: []
}

const getters = {
  messages: state => state.messages
}

const actions = {
  newMessage ({ commit, state }, msg) {
    commit(types.NEW_MESSAGE, msg)
  }
}

const mutations = {
  [types.NEW_MESSAGE] (state, msg) {
    state.messages = state.messages.concat([msg])
  },
  [types.SOCKET_MESSAGE] (state, text) {
    state.messages = state.messages.concat([{ sent: false, text }])
  }
}

export default {
  state,
  getters,
  actions,
  mutations
}
