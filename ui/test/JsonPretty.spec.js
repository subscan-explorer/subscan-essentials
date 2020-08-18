import { mount } from '@vue/test-utils'
import JsonPretty from '@/components/JsonPretty.vue'

describe('JsonPretty', () => {
  test('is a Vue instance', () => {
    const wrapper = mount(JsonPretty)
    expect(wrapper.isVueInstance()).toBeTruthy()
  })
})
