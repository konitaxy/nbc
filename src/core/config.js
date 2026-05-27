/**
 * 网站配置文件
 */

const config = {
  appName: 'PixelCard.com',
  appLogo: '../asserts/logo.png',
  showViteLogo: false
}
export const viteLogo = (env) => {
  if (config.showViteLogo) {
    const chalk = require('chalk')
    console.log(
      chalk.green(
        `hello world`
      )
    )
    console.log('\n')
  }
}
export const mainLogo = (window) => {
  if (window.location.host.indexOf('ref') > -1) {
    config.appLogo = '../asserts/logo_ref.png'
  }
  return config.appLogo
}

export default config
