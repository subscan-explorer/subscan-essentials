export default function ({ store, route, redirect }) {
  if (process.server) {
    store.commit('SET_PLUGINLIST', ['events', 'blocks'])
  }
}
