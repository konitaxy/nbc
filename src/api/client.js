import service from '@/utils/request'
// @Tags authority
// @Summary 更改角色api权限
// @Security ApiKeyAuth
// @accept application/json
// @Produce application/json
// @Param data body api.CreateAuthorityPatams true "更改角色api权限"
// @Success 200 {string} json "{"success":true,"data":{},"msg":"获取成功"}"
// @Router /casbin/UpdateCasbin [post]

export const verifySetting = (data) => {
  return service({
    url: '/client/verifySetting',
    method: 'post',
    data
  })
}
export const genTOCPSecret = (params) => {
  return service({
    url: '/client/tocp',
    method: 'get',
    params
  })
}
export const confirmTOCPBind = (data) => {
  return service({
    url: '/client/tocp',
    method: 'post',
    data
  })
}
export const disableTOCPBind = (data) => {
  return service({
    url: '/client/tocp',
    method: 'delete',
    data
  })
}




  

