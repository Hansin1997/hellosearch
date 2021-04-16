// @ts-ignore
import VueXss from 'vue-xss';
import { boot } from 'quasar/wrappers';


export default boot(({ Vue }) => {
  // eslint-disable-next-line @typescript-eslint/no-unsafe-member-access
  VueXss
  Vue.use(VueXss)
});
