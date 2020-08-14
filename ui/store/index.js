export const state = () => ({
  locales: ['en', 'zh-CN'],
  locale: 'en',
  metadata: {}
})

export const mutations = {
  SET_LANG(state, locale) {
    state.locale = locale
  },
  SET_METADATA(state, metadata) {
    state.metadata = metadata
  }
}
export const actions = {
  async SetMetaData({
    commit
  }) {
    const data = await this.$axios.$post(`/api/scan/metadata`)
    commit('SET_METADATA', data)
  },
  SetLanguage({
    commit
  }, language) {
    commit("SET_LANG", language);
    this.app.i18n.locale = language;
  }
}
