import { createApp } from 'vue'
import { createPinia } from 'pinia'
import ElementPlus from 'element-plus'
import 'element-plus/dist/index.css'
import router from './router' // 👈 新增这一行
import App from './App.vue'

const app = createApp(App)

app.use(createPinia())
app.use(ElementPlus)
app.use(router) // 👈 新增这一行

app.mount('#app')