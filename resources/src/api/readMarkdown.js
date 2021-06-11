import request from '../utils/request'

let readMarkdown = {};

readMarkdown.send = function(data) {
    return request({
        url: `/api/markdown/read`,
        data,
        method: 'POST'
    })
}

export default readMarkdown