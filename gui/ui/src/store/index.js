import Vue from 'vue'
import Vuex from 'vuex'
import createLogger from 'vuex/dist/logger'
import messages from './modules/messages'
import * as types from './mutation-types'

Vue.use(Vuex)

const debug = process.env.NODE_ENV !== 'production'

export default new Vuex.Store({
  state: {
    socketConnected: false,
    entered: false,
    connecting: false,
    connected: false,
    message: null,
    username: '',
    protocol: 'TCP',
    id: '',
    peerID: '',
    peerUsername: '',
    peerAddr: ''
  },
  modules: {
    messages
  },
  mutations: {
    [types.SOCKET_CONNECT]: (state, status) => {
      state.socketConnected = true
    },
    [types.SOCKET_CONNECTING]: (state, objStr) => {
      console.log('connecting')
      const { username, addr } = JSON.parse(objStr)
      state.peerUsername = username
      state.peerAddr = addr
      state.entered = false
      state.connecting = true
    },
    [types.SOCKET_CONNECTED]: (state) => {
      console.log('connected')
      state.entered = false
      state.connecting = false
      state.connected = true
    },
    [types.SOCKET_ENTER]: (state, id) => {
      console.log('entered')
      state.id = id
      state.entered = true
    },
    [types.UPDATE_PROTOCOL]: (state, protocol) => {
      state.protocol = protocol
    },
    [types.UPDATE_USERNAME]: (state, username) => {
      state.username = username
    },
    [types.UPDATE_PEER_ID]: (state, peerID) => {
      state.peerID = peerID
    }
  },
  actions: {
    updateProtocol: ({ commit }, protocol) => {
      commit(types.UPDATE_PROTOCOL, protocol)
    },
    updateUsername: ({ commit }, username) => {
      commit(types.UPDATE_USERNAME, username)
    },
    updatePeerID: ({ commit }, peerID) => {
      commit(types.UPDATE_PEER_ID, peerID)
    },
    otherAction: (context, type) => {
      return true
    }
  },
  getters: {
    username: state => state.username,
    protocol: state => state.protocol,
    entered: state => state.entered,
    connecting: state => state.connecting,
    connected: state => state.connected,
    peerID: state => state.peerID,
    peerUsername: state => state.peerUsername,
    peerAddr: state => state.peerAddr,
    id: state => state.id
  },
  strict: debug,
  plugins: debug ? [createLogger()] : []
})
