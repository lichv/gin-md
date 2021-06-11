import request from '../utils/request'

let searchRequest = {};

searchRequest.send = function(data) {
    return request({
        url: `/api/search`,
        data,
        method: 'POST'
    })
}

export default searchRequest