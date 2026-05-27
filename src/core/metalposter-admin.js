/*
 *  web框架组
 *
 * */
// 加载网站配置文件夹
import { register } from './global'

export default {
  install: (app) => {
    register(app)
    console.log(`
      welcome to pixel admin system, have a good time!
    `)
  }
}
