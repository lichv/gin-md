import request from '../utils/request'

let getMarkdownFiles = {};

getMarkdownFiles.send = function(data) {
    return request({
        url: `/api/markdown/files`,
        data,
        method: 'GET'
    })
}

export default getMarkdownFiles