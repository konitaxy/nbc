/**
 * 网站配置文件
 */

const config = {
  appName: 'Pixel',
  appLogo: 'https://metalposterpro.s3.us-east-1.amazonaws.com/static/display/admin/no_bug.png',
  showViteLogo: true
}

export const viteLogo = (env) => {
  if (config.showViteLogo) {
    const chalk = require('chalk')
    console.log(
      chalk.green(
        `> 欢迎来到pcard 后端`
      )
    )
    console.log('\n')
  }
}

export default config
