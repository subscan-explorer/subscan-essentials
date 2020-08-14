import Vue from "vue";
import VueI18n from "vue-i18n";

import eleEnLocale from "element-ui/lib/locale/lang/en";
import eleZhLocale from "element-ui/lib/locale/lang/zh-CN";
import ElementLocale from "element-ui/lib/locale";

Vue.use(VueI18n)
export default ({ app, store }) => {
  app.i18n = new VueI18n({
    locale: store.state.locale,
    fallbackLocale: store.state.locale || 'en',
    messages: {
      'zh-CN': require('~/locales/zh-CN.json'),
      'en': require('~/locales/en.json')
    }
  })
  ElementLocale.i18n((key, value) => app.i18n.t(key, value));
}
