import Vue from 'vue'

/**
 * 按需引入
 */
import 'element-ui/lib/theme-chalk/index.css'
import {
  Button, Select, Menu, MenuItem, Table, TableColumn, Autocomplete, Divider, Option, Avatar, Tag,
  Dialog, Form, FormItem, Loading, Message
} from 'element-ui'

Vue.use(Button)
Vue.use(Select)
Vue.use(Menu)
Vue.use(MenuItem)
Vue.use(Table)
Vue.use(TableColumn)
Vue.use(Autocomplete)
Vue.use(Divider)
Vue.use(Option)
Vue.use(Avatar)
Vue.use(Tag)
Vue.use(Dialog)
Vue.use(Form)
Vue.use(FormItem)
Vue.use(Loading)
Vue.prototype.$message = Message // 避免浏览器触发默认消息