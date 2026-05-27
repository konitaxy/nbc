
import service from '@/utils/request'


export const getArtworkList = (data) => {
    return service({
      url: '/artor/artwork/list',
      method: 'post',
      data
    })
}

export const  addArtworkPreview = async(data) => {
  return service({
    url: '/artor/artwork/addPreview',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}
export const  addArtwork = async(data) => {
    return service({
      url: '/artor/artwork/add',
      method: 'post',
      headers: { 'Content-Type': 'multipart/form-data' },
      donNotShowLoading: true,
      data
    })
}

export const  previewAddArtwork = async(data) => {
  return service({
    url: '/artor/artwork/preUpload',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}

export const  getArtworkCount = async(params) => {
  return service({
    url: '/artor/artwork/count',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}



export const  getArtworkCatelogs = async(params) => {
  return service({
    url: '/artor/artwork/catelogs',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}

export const  listArtworkLabel = async(params) => {
  return service({
    url: '/artor/artwork/label/list',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}



export const  setUserArtworkLabel = async(data) => {
  return service({
    url: '/artor/artwork/setUserLabel',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}
export const  setArtworkLabel = async(data) => {
  return service({
    url: '/artor/artwork/setLabel',
    method: 'post',
    donNotShowLoading: true,
    data
  })
}
export const  listArtworkLabelByArtwork = async(params) => {
  return service({
    url: '/artor/artwork/label/listByArtwork',
    method: 'get',
    donNotShowLoading: true,
    params
  })
}
