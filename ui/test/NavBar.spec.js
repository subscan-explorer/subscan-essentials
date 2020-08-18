import { mount, RouterLinkStub } from '@vue/test-utils'
import NavBar from '@/components/NavBar.vue'
const wrapper = mount(NavBar, {
  stubs: {
    RouterLink: RouterLinkStub
  }
})
describe('NavBar', () => {
  test('is a Vue instance', () => {
    expect(wrapper.isVueInstance()).toBeTruthy()
  })
  test('logo has link to home page', () => {
    expect(wrapper.find(RouterLinkStub).props().to).toBe('/')
  })
})
