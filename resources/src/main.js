import { createApp } from 'vue'
import App from './App.vue'
import './index.css'
import './base.css'
import './github.css'
import router from './router'
import ElementPlus from 'element-plus'
import 'element-plus/lib/theme-chalk/index.css'
import 'github-markdown-css/github-markdown.css'

const app = createApp(App)

app.use(ElementPlus, { size: 'small', zIndex: 3000, })

import ElMarkdownDisplay from './components/MarkdownDisplay.vue'

const components = [
ElMarkdownDisplay,
]
components.forEach(component => {
  app.component(component.name, component)
})

app.use(router)
app.mount('#app')
