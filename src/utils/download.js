

export const download = (url,filename)=>{
    fetch(url).then(res=>{
        res.blob().then(blob=>{
            let url = window.URL.createObjectURL(blob)
            let a = document.createElement('a')
            a.href = url
            a.download = filename
            a.click()
            a.remove()
        })
    })
}