import Vue from 'vue'
import IconSvg from '~/components/IconSvg'
// register globally

Vue.component('icon-svg', IconSvg)

const requireAll = requireContext => requireContext.keys().map(requireContext)
const req = require.context('~/assets/icons', false, /\.svg$/)
requireAll(req)
