import { mount } from '@vue/test-utils'
import IconSvg from '@/components/IconSvg.vue'

describe('IconSvg', () => {
  test('is a Vue instance', () => {
    const wrapper = mount(IconSvg)
    expect(wrapper.isVueInstance()).toBeTruthy()
  })
})
